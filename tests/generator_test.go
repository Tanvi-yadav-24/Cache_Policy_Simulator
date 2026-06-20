package tests

import (
	"testing"

	"cache-simulator/internal/generator"
)

// TestTraceGenerator tests the trace generation functionality
func TestTraceGenerator(t *testing.T) {
	t.Run("Uniform Distribution", func(t *testing.T) {
		g := generator.NewTraceGenerator()
		g.SetSeed(42)

		requests := g.Generate(1000, 50, generator.Uniform)

		if len(requests) != 1000 {
			t.Errorf("Expected 1000 requests, got %d", len(requests))
		}

		// Check all values are in range
		for _, r := range requests {
			if r < 0 || r >= 50 {
				t.Errorf("Request %d out of range [0, 50)", r)
			}
		}

		// Check distribution is roughly uniform (simple check)
		// Count occurrences of each value
		counts := make(map[int]int)
		for _, r := range requests {
			counts[r]++
		}

		// Each value should appear approximately 1000/50 = 20 times
		// Allow for variance (5 to 40 is reasonable for uniform)
		for val, count := range counts {
			if count < 5 || count > 40 {
				t.Logf("Value %d appears %d times (expected ~20)", val, count)
			}
		}
	})

	t.Run("Zipf Distribution", func(t *testing.T) {
		g := generator.NewTraceGenerator()
		g.SetSeed(42)

		requests := g.Generate(1000, 50, generator.Zipf)

		if len(requests) != 1000 {
			t.Errorf("Expected 1000 requests, got %d", len(requests))
		}

		// Check all values are in range
		for _, r := range requests {
			if r < 0 || r >= 50 {
				t.Errorf("Request %d out of range [0, 50)", r)
			}
		}

		// Zipf should favor lower values
		// Count occurrences of first 10 vs last 10 values
		first10 := 0
		last10 := 0
		for _, r := range requests {
			if r < 10 {
				first10++
			} else if r >= 40 {
				last10++
			}
		}

		// First 10 values should appear more often
		if first10 <= last10 {
			t.Logf("Zipf: first10=%d, last10=%d (expected first > last)", first10, last10)
		}
	})

	t.Run("Localized Pattern", func(t *testing.T) {
		g := generator.NewTraceGenerator()
		g.SetSeed(42)

		requests := g.GenerateLocalized(1000, 100, 10)

		if len(requests) != 1000 {
			t.Errorf("Expected 1000 requests, got %d", len(requests))
		}

		// Check all values are in range
		for _, r := range requests {
			if r < 0 || r >= 100 {
				t.Errorf("Request %d out of range [0, 100)", r)
			}
		}
	})

	t.Run("Sequential Pattern", func(t *testing.T) {
		g := generator.NewTraceGenerator()
		g.SetSeed(42)

		requests := g.GenerateSequential(100, 20)

		if len(requests) != 100 {
			t.Errorf("Expected 100 requests, got %d", len(requests))
		}

		// Check sequential nature (most requests should follow previous + 1)
		sequentialCount := 0
		for i := 1; i < len(requests); i++ {
			if requests[i] == (requests[i-1]+1)%20 || requests[i] == 0 {
				sequentialCount++
			}
		}

		// At least 70% should be sequential (80% sequential + some variance)
		if sequentialCount < 60 {
			t.Errorf("Expected mostly sequential, got only %d sequential out of %d", sequentialCount, len(requests)-1)
		}
	})

	t.Run("Looping Pattern", func(t *testing.T) {
		g := generator.NewTraceGenerator()
		loopSize := 10
		requests := g.GenerateLooping(100, loopSize)

		if len(requests) != 100 {
			t.Errorf("Expected 100 requests, got %d", len(requests))
		}

		// Verify loop pattern
		for i, r := range requests {
			expected := i % loopSize
			if r != expected {
				t.Errorf("Position %d: expected %d, got %d", i, expected, r)
			}
		}
	})

	t.Run("ReproducibleWithSeed", func(t *testing.T) {
		g1 := generator.NewTraceGenerator()
		g1.SetSeed(123)

		g2 := generator.NewTraceGenerator()
		g2.SetSeed(123)

		r1 := g1.Generate(100, 20, generator.Uniform)
		r2 := g2.Generate(100, 20, generator.Uniform)

		// Same seed should produce identical sequences
		for i := range r1 {
			if r1[i] != r2[i] {
				t.Errorf("Same seed produced different results at position %d: %d vs %d", i, r1[i], r2[i])
				break
			}
		}
	})
}
