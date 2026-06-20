package tests

import (
	"testing"

	"cache-simulator/internal/cache"
	"cache-simulator/internal/metrics"
	"cache-simulator/internal/simulator"
)

// TestSimulator tests the simulation engine
func TestSimulator(t *testing.T) {
	t.Run("Basic Simulation", func(t *testing.T) {
		requests := []int{1, 2, 3, 1, 4, 5, 2, 1, 6, 7}
		sim := simulator.NewSimulator(cache.FIFO, 3)

		result := sim.Run(requests)

		if len(result.Steps) != 10 {
			t.Errorf("Expected 10 steps, got %d", len(result.Steps))
		}
		if result.Policy != "FIFO" {
			t.Errorf("Expected policy FIFO, got %s", result.Policy)
		}
		if result.CacheSize != 3 {
			t.Errorf("Expected cache size 3, got %d", result.CacheSize)
		}
	})

	t.Run("Hit Tracking", func(t *testing.T) {
		// Pattern: 1, 2, 3, 1 -> should hit on second 1
		requests := []int{1, 2, 3, 1}
		sim := simulator.NewSimulator(cache.FIFO, 3)

		result := sim.Run(requests)

		// First 3 are misses, 4th is a hit
		expectedHits := 1
		if result.Metrics.CacheHits != expectedHits {
			t.Errorf("Expected %d hits, got %d", expectedHits, result.Metrics.CacheHits)
		}
	})

	t.Run("Miss Tracking", func(t *testing.T) {
		requests := []int{1, 2, 3, 1}
		sim := simulator.NewSimulator(cache.FIFO, 3)

		result := sim.Run(requests)

		// First 3 are misses
		expectedMisses := 3
		if result.Metrics.CacheMisses != expectedMisses {
			t.Errorf("Expected %d misses, got %d", expectedMisses, result.Metrics.CacheMisses)
		}
	})

	t.Run("Eviction Tracking", func(t *testing.T) {
		requests := []int{1, 2, 3, 4, 5, 6}
		sim := simulator.NewSimulator(cache.FIFO, 3)

		result := sim.Run(requests)

		// When cache is full of 3, 4, 5, 6 cause evictions
		// Requests 4, 5, 6 each cause one eviction
		if result.Metrics.Evictions != 3 {
			t.Errorf("Expected 3 evictions, got %d", result.Metrics.Evictions)
		}
	})

	t.Run("Step Result", func(t *testing.T) {
		requests := []int{1, 2, 3, 4}
		sim := simulator.NewSimulator(cache.FIFO, 3)

		result := sim.Run(requests)

		// 4th request should cause eviction of 1
		step4 := result.Steps[3]
		if step4.Request != 4 {
			t.Errorf("Expected request 4, got %d", step4.Request)
		}
		if step4.Hit {
			t.Error("Request 4 should be a miss")
		}
		if step4.Evicted != 1 {
			t.Errorf("Expected eviction of 1, got %d", step4.Evicted)
		}
	})

	t.Run("LRU Behavior", func(t *testing.T) {
		// With LRU, accessing an item makes it MRU
		requests := []int{1, 2, 3, 1, 4} // After filling, access 1, then add 4
		sim := simulator.NewSimulator(cache.LRU, 3)

		result := sim.Run(requests)

		// When 4 is added, 2 should be evicted (LRU), not 1
		step5 := result.Steps[4]
		if step5.Evicted != 2 {
			t.Errorf("LRU: Expected eviction of 2, got %d", step5.Evicted)
		}
	})

	t.Run("Clear Between Runs", func(t *testing.T) {
		sim := simulator.NewSimulator(cache.FIFO, 3)

		// First run
		_ = sim.Run([]int{1, 2, 3})

		// Second run should start fresh
		result := sim.Run([]int{1, 2, 3})

		// All should be misses, no previous state
		if result.Metrics.CacheHits != 0 {
			t.Errorf("Expected 0 hits in fresh run, got %d", result.Metrics.CacheHits)
		}
	})
}

// TestRunComparison tests the comparison function
func TestRunComparison(t *testing.T) {
	requests := []int{1, 2, 3, 1, 4, 5, 2, 1, 6, 7}
	results := simulator.RunComparison(3, requests)

	// Should have results for all 4 policies
	if len(results) != 4 {
		t.Errorf("Expected 4 policy results, got %d", len(results))
	}

	for _, policy := range []string{"FIFO", "LRU", "LFU", "Random"} {
		if _, exists := results[policy]; !exists {
			t.Errorf("Missing results for policy %s", policy)
		}
	}
}

// TestMetrics tests the metrics calculations
func TestMetrics(t *testing.T) {
	t.Run("Hit Ratio", func(t *testing.T) {
		m := metrics.NewMetrics()
		m.TotalRequests = 10
		m.CacheHits = 7

		ratio := m.HitRatio()
		if ratio != 0.7 {
			t.Errorf("Expected hit ratio 0.7, got %f", ratio)
		}
	})

	t.Run("Miss Ratio", func(t *testing.T) {
		m := metrics.NewMetrics()
		m.TotalRequests = 10
		m.CacheMisses = 3

		ratio := m.MissRatio()
		if ratio != 0.3 {
			t.Errorf("Expected miss ratio 0.3, got %f", ratio)
		}
	})

	t.Run("Zero Requests", func(t *testing.T) {
		m := metrics.NewMetrics()

		if m.HitRatio() != 0.0 {
			t.Errorf("Hit ratio with 0 requests should be 0, got %f", m.HitRatio())
		}
		if m.MissRatio() != 0.0 {
			t.Errorf("Miss ratio with 0 requests should be 0, got %f", m.MissRatio())
		}
	})

	t.Run("Counters", func(t *testing.T) {
		m := metrics.NewMetrics()

		m.AddHit()
		m.AddHit()
		m.AddMiss()
		m.AddEviction()

		if m.CacheHits != 2 {
			t.Errorf("Expected 2 hits, got %d", m.CacheHits)
		}
		if m.CacheMisses != 1 {
			t.Errorf("Expected 1 miss, got %d", m.CacheMisses)
		}
		if m.TotalRequests != 3 {
			t.Errorf("Expected 3 total requests, got %d", m.TotalRequests)
		}
		if m.Evictions != 1 {
			t.Errorf("Expected 1 eviction, got %d", m.Evictions)
		}
	})
}
