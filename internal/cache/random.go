package cache

import (
	"math/rand"
	"time"
)

// RandomCache implements Random Replacement cache policy
// Randomly evicts an element when cache is full
type RandomCache struct {
	capacity int
	elements map[int]int // key -> value
	keys     []int       // maintains list of keys for random selection
	rng      *rand.Rand
}

// NewRandomCache creates a new Random cache with the given capacity
func NewRandomCache(capacity int) *RandomCache {
	if capacity <= 0 {
		capacity = 1
	}
	return &RandomCache{
		capacity: capacity,
		elements: make(map[int]int),
		keys:     make([]int, 0),
		rng:      rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Get retrieves a value from the cache if present
func (c *RandomCache) Get(key int) (int, bool) {
	val, found := c.elements[key]
	return val, found
}

// Put adds a key-value pair to the cache
// Returns evicted key and true if eviction occurred
func (c *RandomCache) Put(key, value int) (int, bool) {
	// If key already exists, update value (no eviction)
	if _, exists := c.elements[key]; exists {
		c.elements[key] = value
		return 0, false
	}

	var evictedKey int
	evicted := false

	// If cache is full, randomly evict an element
	if len(c.elements) >= c.capacity {
		// Select random index
		randomIndex := c.rng.Intn(len(c.keys))
		evictedKey = c.keys[randomIndex]

		// Remove evicted element
		delete(c.elements, evictedKey)

		// Remove from keys slice (swap with last and truncate)
		c.keys[randomIndex] = c.keys[len(c.keys)-1]
		c.keys = c.keys[:len(c.keys)-1]

		evicted = true
	}

	// Add new element
	c.elements[key] = value
	c.keys = append(c.keys, key)

	if evicted {
		return evictedKey, true
	}
	return 0, false
}

// Size returns the current number of elements in the cache
func (c *RandomCache) Size() int {
	return len(c.elements)
}

// Capacity returns the maximum number of elements the cache can hold
func (c *RandomCache) Capacity() int {
	return c.capacity
}

// Clear removes all elements from the cache
func (c *RandomCache) Clear() {
	c.elements = make(map[int]int)
	c.keys = make([]int, 0)
}

// Keys returns all keys currently in the cache
func (c *RandomCache) Keys() []int {
	keys := make([]int, len(c.keys))
	copy(keys, c.keys)
	return keys
}

// GetName returns the name of the cache policy
func (c *RandomCache) GetName() string {
	return string(Random)
}

// SetSeed sets the random seed for reproducible results (useful for testing)
func (c *RandomCache) SetSeed(seed int64) {
	c.rng = rand.New(rand.NewSource(seed))
}
