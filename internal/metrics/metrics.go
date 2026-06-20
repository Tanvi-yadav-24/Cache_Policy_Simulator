package metrics

import "fmt"

// Metrics stores performance metrics for a cache policy simulation
type Metrics struct {
	TotalRequests  int     // Total number of requests processed
	CacheHits      int     // Number of cache hits
	CacheMisses    int     // Number of cache misses
	Evictions      int     // Number of elements evicted
	ExecutionTimeNanos int64 // Execution time in nanoseconds
	MemoryBytes    uint64  // Memory usage in bytes (approximate)
}

// NewMetrics creates a new Metrics instance
func NewMetrics() *Metrics {
	return &Metrics{}
}

// HitRatio returns the cache hit ratio (hits / total_requests)
func (m *Metrics) HitRatio() float64 {
	if m.TotalRequests == 0 {
		return 0.0
	}
	return float64(m.CacheHits) / float64(m.TotalRequests)
}

// MissRatio returns the cache miss ratio (misses / total_requests)
func (m *Metrics) MissRatio() float64 {
	if m.TotalRequests == 0 {
		return 0.0
	}
	return float64(m.CacheMisses) / float64(m.TotalRequests)
}

// String returns a formatted string representation of the metrics
func (m *Metrics) String() string {
	return fmt.Sprintf(
		"Total Requests: %d\nCache Hits: %d\nCache Misses: %d\nEvictions: %d\nHit Ratio: %.4f (%.2f%%)\nMiss Ratio: %.4f (%.2f%%)",
		m.TotalRequests, m.CacheHits, m.CacheMisses, m.Evictions,
		m.HitRatio(), m.HitRatio()*100,
		m.MissRatio(), m.MissRatio()*100,
	)
}

// AddHit increments the cache hit counter
func (m *Metrics) AddHit() {
	m.CacheHits++
	m.TotalRequests++
}

// AddMiss increments the cache miss counter
func (m *Metrics) AddMiss() {
	m.CacheMisses++
	m.TotalRequests++
}

// AddEviction increments the eviction counter
func (m *Metrics) AddEviction() {
	m.Evictions++
}
