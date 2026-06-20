package simulator

import (
	"time"

	"cache-simulator/internal/cache"
	"cache-simulator/internal/metrics"
)

// StepResult represents the result of processing a single request
type StepResult struct {
	Request     int   `json:"request"`      // Current request value
	CacheBefore []int `json:"cache_before"` // Cache state before processing
	CacheAfter  []int `json:"cache_after"`  // Cache state after processing
	Evicted     int   `json:"evicted"`      // Evicted key (0 if none)
	Hit         bool  `json:"hit"`          // True if cache hit, false if miss
}

// SimulationResult represents the complete result of a simulation
type SimulationResult struct {
	Policy       string        `json:"policy"`       // Cache policy name
	CacheSize    int           `json:"cache_size"`    // Cache capacity
	Steps        []StepResult  `json:"steps"`        // Step-by-step results
	Metrics      metrics.Metrics ` `json:"metrics"`    // Final metrics
	RequestTrace []int         `json:"request_trace"` // Original request sequence
}

// Simulator runs cache simulations
type Simulator struct {
	cache    cache.Cache
	policy   cache.CachePolicy
	capacity int
	metrics  *metrics.Metrics
}

// NewSimulator creates a new simulator with the given policy and capacity
func NewSimulator(policy cache.CachePolicy, capacity int) *Simulator {
	return &Simulator{
		cache:    cache.Factory(policy, capacity),
		policy:   policy,
		capacity: capacity,
		metrics:  metrics.NewMetrics(),
	}
}

// Run executes a simulation on the given request sequence
// Returns step-by-step results
func (s *Simulator) Run(requests []int) SimulationResult {
	s.cache.Clear()
	s.metrics = metrics.NewMetrics()

	steps := make([]StepResult, 0, len(requests))
	startTime := time.Now()

	for _, req := range requests {
		step := s.processRequest(req)
		steps = append(steps, step)
	}

	executionTime := time.Since(startTime)

	// Approximate memory usage
	memUsage := s.estimateMemory()

	result := SimulationResult{
		Policy:       string(s.policy),
		CacheSize:    s.capacity,
		Steps:        steps,
		RequestTrace: requests,
	}
	result.Metrics = *s.metrics
	result.Metrics.ExecutionTimeNanos = executionTime.Nanoseconds()
	result.Metrics.MemoryBytes = memUsage

	return result
}

// processRequest handles a single request and returns the step result
func (s *Simulator) processRequest(request int) StepResult {
	cacheBefore := make([]int, len(s.cache.Keys()))
	copy(cacheBefore, s.cache.Keys())

	// Try to get from cache
	_, hit := s.cache.Get(request)

	var evicted int
	evictionOccurred := false

	if !hit {
		// Cache miss - add to cache (may cause eviction)
		evicted, evictionOccurred = s.cache.Put(request, request)
		s.metrics.AddMiss()
		if evictionOccurred {
			s.metrics.AddEviction()
		}
	} else {
		// For LRU/LFU, Get already updates access info
		// For FIFO, we need to check if key was already in cache
		s.metrics.AddHit()
	}

	cacheAfter := s.cache.Keys()

	result := StepResult{
		Request:     request,
		CacheBefore: cacheBefore,
		CacheAfter:  cacheAfter,
		Evicted:     evicted,
		Hit:         hit,
	}

	return result
}

// estimateMemory returns an approximate memory usage in bytes
func (s *Simulator) estimateMemory() uint64 {
	// Rough estimation based on cache size and policy
	// Each entry: key (int=8 bytes) + value (int=8 bytes) + overhead
	baseOverhead := uint64(100) // struct overhead
	entrySize := uint64(16)      // key + value
	pointerOverhead := uint64(8) // pointers for LRU/LFU lists

	switch s.policy {
	case cache.FIFO:
		return baseOverhead + uint64(s.capacity)*entrySize + uint64(s.capacity)*pointerOverhead
	case cache.LRU:
		return baseOverhead + uint64(s.capacity)*entrySize + uint64(s.capacity)*pointerOverhead*2
	case cache.LFU:
		// LFU has additional frequency tracking
		return baseOverhead + uint64(s.capacity)*entrySize*2 + uint64(s.capacity)*pointerOverhead*2
	case cache.Random:
		return baseOverhead + uint64(s.capacity)*entrySize + uint64(s.capacity)*pointerOverhead
	default:
		return baseOverhead + uint64(s.capacity)*entrySize
	}
}

// RunComparison runs simulations for all policies and returns comparison results
func RunComparison(capacity int, requests []int) map[string]SimulationResult {
	policies := []cache.CachePolicy{cache.FIFO, cache.LRU, cache.LFU, cache.Random}
	results := make(map[string]SimulationResult)

	for _, policy := range policies {
		sim := NewSimulator(policy, capacity)
		result := sim.Run(requests)
		results[string(policy)] = result
	}

	return results
}
