package cache

import (
	"container/list"
)

// LFUCache implements Least Frequently Used cache replacement policy
// When frequencies are equal, evicts the oldest (FIFO within same frequency)
// Time Complexity: Get O(1), Put O(1)
type LFUCache struct {
	capacity  int
	elements  map[int]*lfuEntry          // key -> entry
	freqLists map[int]*list.List          // frequency -> list of keys
	keyPos    map[int]*list.Element       // key -> position in frequency list
	minFreq   int                         // current minimum frequency
}

type lfuEntry struct {
	key   int
	value int
	freq  int
}

// NewLFUCache creates a new LFU cache with the given capacity
func NewLFUCache(capacity int) *LFUCache {
	if capacity <= 0 {
		capacity = 1
	}
	return &LFUCache{
		capacity:  capacity,
		elements:  make(map[int]*lfuEntry),
		freqLists: make(map[int]*list.List),
		keyPos:    make(map[int]*list.Element),
		minFreq:   0,
	}
}

// Get retrieves a value from the cache if present and increments its frequency
func (c *LFUCache) Get(key int) (int, bool) {
	entry, found := c.elements[key]
	if !found {
		return 0, false
	}

	// Increment frequency
	c.incrementFreq(key)
	return entry.value, true
}

// Put adds a key-value pair to the cache
// Returns evicted key and true if eviction occurred
func (c *LFUCache) Put(key, value int) (int, bool) {
	// If key exists, update value and increment frequency
	if entry, exists := c.elements[key]; exists {
		entry.value = value
		c.incrementFreq(key)
		return 0, false
	}

	var evictedKey int
	evicted := false

	// If cache is full, remove the least frequently used (and oldest if tied)
	if len(c.elements) >= c.capacity {
		evictedKey, evicted = c.evictLFU()
	}

	// Add new element with frequency 1
	c.addElement(key, value)

	if evicted {
		return evictedKey, true
	}
	return 0, false
}

// incrementFreq increments the frequency of the key and updates data structures
func (c *LFUCache) incrementFreq(key int) {
	entry := c.elements[key]
	oldFreq := entry.freq
	newFreq := oldFreq + 1

	// Remove from old frequency list
	if list, ok := c.freqLists[oldFreq]; ok {
		if element, exists := c.keyPos[key]; exists {
			list.Remove(element)
			delete(c.keyPos, key)
		}
		// Clean up empty frequency lists
		if list.Len() == 0 {
			delete(c.freqLists, oldFreq)
			// Update minFreq if needed
			if c.minFreq == oldFreq {
				c.minFreq = newFreq
			}
		}
	}

	// Add to new frequency list
	if _, ok := c.freqLists[newFreq]; !ok {
		c.freqLists[newFreq] = list.New()
	}
	element := c.freqLists[newFreq].PushBack(key)
	c.keyPos[key] = element
	entry.freq = newFreq
}

// addElement adds a new element with frequency 1
func (c *LFUCache) addElement(key, value int) {
	entry := &lfuEntry{key: key, value: value, freq: 1}
	c.elements[key] = entry

	// Initialize frequency 1 list if needed
	if _, ok := c.freqLists[1]; !ok {
		c.freqLists[1] = list.New()
	}

	element := c.freqLists[1].PushBack(key)
	c.keyPos[key] = element
	c.minFreq = 1 // New elements always start with frequency 1
}

// evictLFU removes the least frequently used element
// If multiple elements have the same minimum frequency, removes the oldest (front of list)
func (c *LFUCache) evictLFU() (int, bool) {
	minList, ok := c.freqLists[c.minFreq]
	if !ok || minList.Len() == 0 {
		return 0, false
	}

	// Remove oldest element from minimum frequency list
	front := minList.Front()
	evictedKey := front.Value.(int)

	minList.Remove(front)
	delete(c.keyPos, evictedKey)
	delete(c.elements, evictedKey)

	// Clean up empty frequency list
	if minList.Len() == 0 {
		delete(c.freqLists, c.minFreq)
	}

	return evictedKey, true
}

// Size returns the current number of elements in the cache
func (c *LFUCache) Size() int {
	return len(c.elements)
}

// Capacity returns the maximum number of elements the cache can hold
func (c *LFUCache) Capacity() int {
	return c.capacity
}

// Clear removes all elements from the cache
func (c *LFUCache) Clear() {
	c.elements = make(map[int]*lfuEntry)
	c.freqLists = make(map[int]*list.List)
	c.keyPos = make(map[int]*list.Element)
	c.minFreq = 0
}

// Keys returns all keys currently in the cache
func (c *LFUCache) Keys() []int {
	keys := make([]int, 0, len(c.elements))
	for key := range c.elements {
		keys = append(keys, key)
	}
	return keys
}

// GetName returns the name of the cache policy
func (c *LFUCache) GetName() string {
	return string(LFU)
}
