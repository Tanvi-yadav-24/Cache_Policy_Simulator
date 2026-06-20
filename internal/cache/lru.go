package cache

import (
	"container/list"
)

// LRUCache implements Least Recently Used cache replacement policy
// Time Complexity: Get O(1), Put O(1)
// Data Structures: HashMap + Doubly Linked List
type LRUCache struct {
	capacity int
	elements map[int]int          // key -> value
	order    *list.List           // doubly linked list (front = LRU, back = MRU)
	keyPos   map[int]*list.Element // key -> position in list
}

// NewLRUCache creates a new LRU cache with the given capacity
func NewLRUCache(capacity int) *LRUCache {
	if capacity <= 0 {
		capacity = 1
	}
	return &LRUCache{
		capacity: capacity,
		elements: make(map[int]int),
		order:    list.New(),
		keyPos:   make(map[int]*list.Element),
	}
}

// Get retrieves a value from the cache if present and moves it to most recently used
func (c *LRUCache) Get(key int) (int, bool) {
	val, found := c.elements[key]
	if !found {
		return 0, false
	}

	// Move to back (most recently used)
	c.moveToBack(key)
	return val, true
}

// Put adds a key-value pair to the cache
// Returns evicted key and true if eviction occurred
func (c *LRUCache) Put(key, value int) (int, bool) {
	// If key exists, update value and move to most recently used
	if _, exists := c.elements[key]; exists {
		c.elements[key] = value
		c.moveToBack(key)
		return 0, false
	}

	var evictedKey int
	evicted := false

	// If cache is full, remove the least recently used (front)
	if c.order.Len() >= c.capacity {
		lru := c.order.Front()
		evictedKey = lru.Value.(int)
		c.order.Remove(lru)
		delete(c.elements, evictedKey)
		delete(c.keyPos, evictedKey)
		evicted = true
	}

	// Add new element at the back (most recently used)
	c.elements[key] = value
	element := c.order.PushBack(key)
	c.keyPos[key] = element

	if evicted {
		return evictedKey, true
	}
	return 0, false
}

// moveToBack moves the key to the back of the list (most recently used)
func (c *LRUCache) moveToBack(key int) {
	if element, exists := c.keyPos[key]; exists {
		c.order.MoveToBack(element)
	}
}

// Size returns the current number of elements in the cache
func (c *LRUCache) Size() int {
	return len(c.elements)
}

// Capacity returns the maximum number of elements the cache can hold
func (c *LRUCache) Capacity() int {
	return c.capacity
}

// Clear removes all elements from the cache
func (c *LRUCache) Clear() {
	c.elements = make(map[int]int)
	c.order = list.New()
	c.keyPos = make(map[int]*list.Element)
}

// Keys returns all keys currently in the cache (ordered from LRU to MRU)
func (c *LRUCache) Keys() []int {
	keys := make([]int, 0, len(c.elements))
	for e := c.order.Front(); e != nil; e = e.Next() {
		keys = append(keys, e.Value.(int))
	}
	return keys
}

// GetName returns the name of the cache policy
func (c *LRUCache) GetName() string {
	return string(LRU)
}
