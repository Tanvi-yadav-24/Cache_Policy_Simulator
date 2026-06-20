package tests

import (
	"testing"

	"cache-simulator/internal/cache"
)

// TestFIFOCache tests the FIFO cache implementation
func TestFIFOCache(t *testing.T) {
	t.Run("Basic Operations", func(t *testing.T) {
		c := cache.NewFIFOCache(3)

		// Test initial state
		if c.Size() != 0 {
			t.Errorf("Expected size 0, got %d", c.Size())
		}
		if c.Capacity() != 3 {
			t.Errorf("Expected capacity 3, got %d", c.Capacity())
		}

		// Test Put and Get
		evictedKey, evicted := c.Put(1, 10)
		if evicted {
			t.Errorf("Expected no eviction on first put, got evicted key %d", evictedKey)
		}

		val, found := c.Get(1)
		if !found || val != 10 {
			t.Errorf("Expected to find key 1 with value 10, got found=%v, val=%d", found, val)
		}
	})

	t.Run("FIFO Eviction Order", func(t *testing.T) {
		c := cache.NewFIFOCache(3)

		// Fill cache: [1, 2, 3]
		c.Put(1, 10)
		c.Put(2, 20)
		c.Put(3, 30)

		// Add 4th element, should evict 1 (oldest)
		evictedKey, evicted := c.Put(4, 40)
		if !evicted || evictedKey != 1 {
			t.Errorf("Expected eviction of key 1, got evicted=%v, key=%d", evicted, evictedKey)
		}

		// Verify 1 is gone
		if _, found := c.Get(1); found {
			t.Error("Key 1 should have been evicted")
		}

		// Verify others remain: [2, 3, 4]
		if _, found := c.Get(2); !found {
			t.Error("Key 2 should still be in cache")
		}
		if _, found := c.Get(3); !found {
			t.Error("Key 3 should still be in cache")
		}
		if _, found := c.Get(4); !found {
			t.Error("Key 4 should be in cache")
		}
	})

	t.Run("Update Without Eviction", func(t *testing.T) {
		c := cache.NewFIFOCache(2)
		c.Put(1, 10)
		c.Put(2, 20)

		// Update existing key
		evictedKey, evicted := c.Put(1, 100)
		if evicted {
			t.Errorf("Update should not cause eviction, got evicted key %d", evictedKey)
		}

		val, _ := c.Get(1)
		if val != 100 {
			t.Errorf("Expected updated value 100, got %d", val)
		}
	})

	t.Run("Clear", func(t *testing.T) {
		c := cache.NewFIFOCache(3)
		c.Put(1, 10)
		c.Put(2, 20)

		c.Clear()
		if c.Size() != 0 {
			t.Errorf("Expected size 0 after clear, got %d", c.Size())
		}
	})

	t.Run("Keys Order", func(t *testing.T) {
		c := cache.NewFIFOCache(3)
		c.Put(1, 10)
		c.Put(2, 20)
		c.Put(3, 30)

		keys := c.Keys()
		expected := []int{1, 2, 3}
		for i, k := range keys {
			if k != expected[i] {
				t.Errorf("Expected key order %v, got %v", expected, keys)
				break
			}
		}
	})
}

// TestLRUCache tests the LRU cache implementation
func TestLRUCache(t *testing.T) {
	t.Run("Basic Operations", func(t *testing.T) {
		c := cache.NewLRUCache(3)

		c.Put(1, 10)
		val, found := c.Get(1)
		if !found || val != 10 {
			t.Errorf("Expected to find key 1, got found=%v, val=%d", found, val)
		}
	})

	t.Run("LRU Eviction", func(t *testing.T) {
		c := cache.NewLRUCache(3)

		// Fill cache: [1, 2, 3] (1 is LRU, 3 is MRU)
		c.Put(1, 10)
		c.Put(2, 20)
		c.Put(3, 30)

		// Access 1, making it MRU: [2, 3, 1]
		c.Get(1)

		// Add 4th element, should evict 2 (now LRU)
		evictedKey, evicted := c.Put(4, 40)
		if !evicted || evictedKey != 2 {
			t.Errorf("Expected eviction of key 2, got evicted=%v, key=%d", evicted, evictedKey)
		}

		// Verify 1 is still there (was accessed recently)
		if _, found := c.Get(1); !found {
			t.Error("Key 1 should still be in cache")
		}
	})

	t.Run("Get Updates Recency", func(t *testing.T) {
		c := cache.NewLRUCache(2)

		c.Put(1, 10)
		c.Put(2, 20)

		// Access 1, making it MRU
		c.Get(1)

		// Add 3, should evict 2 (LRU)
		evictedKey, _ := c.Put(3, 30)
		if evictedKey != 2 {
			t.Errorf("Expected eviction of key 2, got %d", evictedKey)
		}
	})

	t.Run("Keys Order LRU to MRU", func(t *testing.T) {
		c := cache.NewLRUCache(3)
		c.Put(1, 10)
		c.Put(2, 20)
		c.Put(3, 30)
		c.Get(1) // Move 1 to MRU

		keys := c.Keys()
		// Expected order: [2, 3, 1] (LRU -> MRU)
		if len(keys) != 3 || keys[2] != 1 {
			t.Errorf("Expected key 1 at end (MRU), got %v", keys)
		}
	})
}

// TestLFUCache tests the LFU cache implementation
func TestLFUCache(t *testing.T) {
	t.Run("Basic Operations", func(t *testing.T) {
		c := cache.NewLFUCache(3)

		c.Put(1, 10)
		val, found := c.Get(1)
		if !found || val != 10 {
			t.Errorf("Expected to find key 1, got found=%v, val=%d", found, val)
		}
	})

	t.Run("LFU Eviction", func(t *testing.T) {
		c := cache.NewLFUCache(3)

		// Fill cache: [1, 2, 3] all with freq=1
		c.Put(1, 10)
		c.Put(2, 20)
		c.Put(3, 30)

		// Access 1 and 2 multiple times
		c.Get(1) // freq(1) = 2
		c.Get(1) // freq(1) = 3
		c.Get(2) // freq(2) = 2

		// Add 4th element, should evict 3 (lowest freq)
		evictedKey, evicted := c.Put(4, 40)
		if !evicted || evictedKey != 3 {
			t.Errorf("Expected eviction of key 3 (lowest freq), got evicted=%v, key=%d", evicted, evictedKey)
		}
	})

	t.Run("Tie Breaker Uses FIFO", func(t *testing.T) {
		c := cache.NewLFUCache(3)

		// Fill cache: [1, 2, 3] all with freq=1
		c.Put(1, 10)
		c.Put(2, 20)
		c.Put(3, 30)

		// All have same frequency, should evict oldest (1)
		evictedKey, evicted := c.Put(4, 40)
		if !evicted || evictedKey != 1 {
			t.Errorf("Expected eviction of key 1 (oldest with same freq), got evicted=%v, key=%d", evicted, evictedKey)
		}
	})

	t.Run("Get Increments Frequency", func(t *testing.T) {
		c := cache.NewLFUCache(2)

		c.Put(1, 10)
		c.Put(2, 20)
		c.Get(1) // freq(1) = 2
		c.Get(1) // freq(1) = 3

		// Add 3, should evict 2 (lower freq)
		evictedKey, _ := c.Put(3, 30)
		if evictedKey != 2 {
			t.Errorf("Expected eviction of key 2, got %d", evictedKey)
		}
	})
}

// TestRandomCache tests the Random cache implementation
func TestRandomCache(t *testing.T) {
	t.Run("Basic Operations", func(t *testing.T) {
		c := cache.NewRandomCache(3)

		c.Put(1, 10)
		val, found := c.Get(1)
		if !found || val != 10 {
			t.Errorf("Expected to find key 1, got found=%v, val=%d", found, val)
		}
	})

	t.Run("Random Eviction", func(t *testing.T) {
		c := cache.NewRandomCache(3)
		c.SetSeed(42) // Reproducible randomness

		c.Put(1, 10)
		c.Put(2, 20)
		c.Put(3, 30)

		// Add 4th element, should evict randomly
		evictedKey, evicted := c.Put(4, 40)
		if !evicted {
			t.Error("Expected eviction when cache is full")
		}

		// Verify evicted key was valid
		if evictedKey < 1 || evictedKey > 3 {
			t.Errorf("Invalid evicted key %d, should be 1, 2, or 3", evictedKey)
		}

		// Verify cache size
		if c.Size() != 3 {
			t.Errorf("Expected size 3, got %d", c.Size())
		}
	})

	t.Run("Deterministic With Seed", func(t *testing.T) {
		c1 := cache.NewRandomCache(3)
		c1.SetSeed(123)
		c2 := cache.NewRandomCache(3)
		c2.SetSeed(123)

		for i := 1; i <= 10; i++ {
			ev1, _ := c1.Put(i, i*10)
			ev2, _ := c2.Put(i, i*10)

			if ev1 != ev2 {
				t.Errorf("Same seed should produce same evictions: %d vs %d", ev1, ev2)
			}
		}
	})
}

// TestCacheFactory tests the cache factory function
func TestCacheFactory(t *testing.T) {
	tests := []struct {
		policy cache.CachePolicy
		name   string
	}{
		{cache.FIFO, "FIFO"},
		{cache.LRU, "LRU"},
		{cache.LFU, "LFU"},
		{cache.Random, "Random"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := cache.Factory(tt.policy, 10)
			if c == nil {
				t.Errorf("Factory returned nil for policy %s", tt.policy)
			}
			if c.GetName() != tt.name {
				t.Errorf("Expected name %s, got %s", tt.name, c.GetName())
			}
		})
	}
}

// TestEdgeCases tests edge cases for all cache types
func TestEdgeCases(t *testing.T) {
	t.Run("Zero Capacity", func(t *testing.T) {
		// Zero capacity should be treated as 1
		for _, policy := range []cache.CachePolicy{cache.FIFO, cache.LRU, cache.LFU, cache.Random} {
			c := cache.Factory(policy, 0)
			if c.Capacity() != 1 {
				t.Errorf("Policy %s: zero capacity should become 1, got %d", policy, c.Capacity())
			}
		}
	})

	t.Run("Single Element Cache", func(t *testing.T) {
		for _, policy := range []cache.CachePolicy{cache.FIFO, cache.LRU, cache.LFU, cache.Random} {
			c := cache.Factory(policy, 1)

			c.Put(1, 10)
			evicted, _ := c.Put(2, 20)

			if evicted != 1 {
				t.Errorf("Policy %s: single element cache should evict 1, got %d", policy, evicted)
			}
		}
	})
}
