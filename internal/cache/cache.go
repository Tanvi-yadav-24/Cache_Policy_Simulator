package cache

// Cache defines the interface for all cache replacement policies
type Cache interface {
	// Get retrieves a value from the cache if present
	// Returns the value and true if found, zero value and false otherwise
	Get(key int) (int, bool)

	// Put adds a key-value pair to the cache
	// Returns evicted key and true if eviction occurred, 0 and false otherwise
	Put(key, value int) (int, bool)

	// Size returns the current number of elements in the cache
	Size() int

	// Capacity returns the maximum number of elements the cache can hold
	Capacity() int

	// Clear removes all elements from the cache
	Clear()

	// Keys returns all keys currently in the cache
	Keys() []int

	// GetName returns the name of the cache policy
	GetName() string
}

// CachePolicy represents the type of cache replacement policy
type CachePolicy string

const (
	FIFO   CachePolicy = "FIFO"
	LRU    CachePolicy = "LRU"
	LFU    CachePolicy = "LFU"
	Random CachePolicy = "Random"
)

// Factory creates a new cache based on the specified policy
func Factory(policy CachePolicy, capacity int) Cache {
	switch policy {
	case FIFO:
		return NewFIFOCache(capacity)
	case LRU:
		return NewLRUCache(capacity)
	case LFU:
		return NewLFUCache(capacity)
	case Random:
		return NewRandomCache(capacity)
	default:
		return NewFIFOCache(capacity)
	}
}
