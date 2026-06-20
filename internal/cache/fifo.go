package cache

import (
	"container/list"
)

// FIFOCache implements First In First Out cache replacement policy
// Time Complexity: Get O(1), Put O(1)
type FIFOCache struct {
	capacity int
	elements map[int]int          // key -> value
	order    *list.List           // maintains insertion order
	keyPos   map[int]*list.Element // key -> position in order list
}

// NewFIFOCache creates a new FIFO cache with the given capacity
func NewFIFOCache(capacity int) *FIFOCache {
	if capacity <= 0 {
		capacity = 1
	}
	return &FIFOCache{
		capacity: capacity,
		elements: make(map[int]int),
		order:    list.New(),
		keyPos:   make(map[int]*list.Element),
	}
}

// Get retrieves a value from the cache if present
func (c *FIFOCache) Get(key int) (int, bool) {
	val, found := c.elements[key]
	return val, found
}

// Put adds a key-value pair to the cache
// Returns evicted key and true if eviction occurred
func (c *FIFOCache) Put(key, value int) (int, bool) {
	// If key already exists, update value (no eviction)
	if _, exists := c.elements[key]; exists {
		c.elements[key] = value
		return 0, false
	}

	var evictedKey int
	evicted := false

	// If cache is full, remove the oldest element
	if c.order.Len() >= c.capacity {
		oldest := c.order.Front()
		evictedKey = oldest.Value.(int)
		c.order.Remove(oldest)
		delete(c.elements, evictedKey)
		delete(c.keyPos, evictedKey)
		evicted = true
	}

	// Add new element
	c.elements[key] = value
	element := c.order.PushBack(key)
	c.keyPos[key] = element

	if evicted {
		return evictedKey, true
	}
	return 0, false
}

// Size returns the current number of elements in the cache
func (c *FIFOCache) Size() int {
	return len(c.elements)
}

// Capacity returns the maximum number of elements the cache can hold
func (c *FIFOCache) Capacity() int {
	return c.capacity
}

// Clear removes all elements from the cache
func (c *FIFOCache) Clear() {
	c.elements = make(map[int]int)
	c.order = list.New()
	c.keyPos = make(map[int]*list.Element)
}

// Keys returns all keys currently in the cache
func (c *FIFOCache) Keys() []int {
	keys := make([]int, 0, len(c.elements))
	for e := c.order.Front(); e != nil; e = e.Next() {
		keys = append(keys, e.Value.(int))
	}
	return keys
}

// GetName returns the name of the cache policy
func (c *FIFOCache) GetName() string {
	return string(FIFO)
}
