# Cache Policy Simulator

A comprehensive simulator for comparing cache replacement policies used in operating systems, databases, browsers, and distributed systems.

![TypeScript](https://img.shields.io/badge/TypeScript-5.0+-3178C6?style=flat&logo=typescript)
![React](https://img.shields.io/badge/React-18+-61DAFB?style=flat&logo=react)
![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/License-MIT-green)

## Overview

This project implements and compares four classic cache replacement algorithms:

- **FIFO** (First In First Out)
- **LRU** (Least Recently Used)
- **LFU** (Least Frequently Used)
- **Random Replacement**

Each algorithm is implemented from scratch with O(1) time complexity for both `Get()` and `Put()` operations.

**Two implementations available:**
- **TypeScript/React** - Runs entirely in the browser (no backend required)
- **Go** - CLI tool and HTTP API server (optional backend)

## Quick Start (Frontend Only)

The simulator works entirely in the browser. Just run:

```bash
npm install
npm run dev
```

Open `http://localhost:5173` and start simulating!

## Cache Algorithms

### FIFO (First In First Out)

**Strategy:** Removes the oldest element when the cache is full.

**Implementation:**
- Queue-based using `container/list` (doubly linked list)
- HashMap for O(1) lookups
- Evicts front of queue (oldest element)

**Use Cases:**
- Simple buffering scenarios
- When recency of access is not important
- Low overhead requirements

**Complexity:**
- Get: O(1)
- Put: O(1)

### LRU (Least Recently Used)

**Strategy:** Removes the least recently accessed item when the cache is full.

**Implementation:**
- **HashMap**: Maps keys to cache entries for O(1) lookup
- **Doubly Linked List**: Maintains access order (front = LRU, back = MRU)
- On `Get()`: Move accessed item to back (most recently used)
- On eviction: Remove front element

**Use Cases:**
- Web browser caches
- Database query caches
- When access patterns show temporal locality

**Complexity:**
- Get: O(1) - Hash lookup + list repositioning
- Put: O(1) - Hash insert + list append/remove

**Data Structure Diagram:**
```
HashMap              Doubly Linked List
------               -------------------
Key -> Node*   <-->  [Node(key)] <-> [Node(key)] <-> ...

Front of List = LRU (next to evict)
Back of List = MRU (most recently used)
```

### LFU (Least Frequently Used)

**Strategy:** Removes the least frequently accessed item. Tie-breaker: evicts the oldest among items with equal frequency.

**Implementation:**
- HashMap for O(1) key lookup
- Frequency-bucketed lists: `map[frequency]*list.List`
- Each frequency bucket contains keys with that frequency
- On `Get()`: Increment frequency, move to higher bucket
- On eviction: Remove oldest from minimum frequency bucket

**Use Cases:**
- Scenarios with steady access patterns
- When some items are consistently more popular
- Long-term caching decisions

**Complexity:**
- Get: O(1) - Hash lookup + frequency increment
- Put: O(1) - Hash insert + frequency bucket management

**Data Structures:**
```
HashMap          Frequency Buckets
-------          -----------------
Key -> Entry     Freq 1: [key, key, key] (oldest -> newest)
                 Freq 2: [key, key]
                 Freq 3: [key]
                 ...

minFreq tracks the minimum frequency for O(1) eviction
```

### Random Replacement

**Strategy:** Randomly evicts an item when the cache is full.

**Implementation:**
- HashMap for storage
- Slice of keys for random selection
- On eviction: Random index from key slice

**Use Cases:**
- When eviction speed is critical
- Unpredictable access patterns
- Low implementation overhead

**Complexity:**
- Get: O(1)
- Put: O(1)

## Project Structure

```
cache-simulator/
├── cmd/
│   ├── server/main.go     # HTTP API server
│   ├── cli/main.go        # Command-line interface
│   └── benchmark/         # Benchmarking utilities
├── internal/
│   ├── cache/
│   │   ├── cache.go       # Cache interface and factory
│   │   ├── fifo.go        # FIFO implementation
│   │   ├── lru.go         # LRU implementation
│   │   ├── lfu.go         # LFU implementation
│   │   └── random.go      # Random replacement
│   ├── simulator/
│   │   └── simulator.go   # Simulation engine
│   ├── metrics/
│   │   └── metrics.go     # Performance metrics
│   └── generator/
│       └── generator.go   # Request trace generation
├── tests/
│   ├── cache_test.go      # Unit tests for cache policies
│   ├── simulator_test.go  # Unit tests for simulator
│   ├── generator_test.go  # Unit tests for generator
│   └── benchmark_test.go  # Performance benchmarks
├── src/                   # React frontend
└── README.md
```

## Getting Started

### Prerequisites

- Go 1.21 or later
- Node.js 18+ (for frontend)

### Running the Backend

```bash
# Install Go dependencies
go mod tidy

# Run the HTTP API server
go run cmd/server/main.go

# Or use the CLI
go run cmd/cli/main.go simulate -policy LRU -size 3 -requests "1,2,3,1,4,5"

# Compare all policies
go run cmd/cli/main.go compare -size 3 -requests "1,2,3,1,4,5"

# Run benchmarks
go run cmd/cli/main.go benchmark -num 10000 -dist zipf
```

### Running the Frontend

```bash
# Install dependencies
npm install

# Run development server
npm run dev
```

The React app will be available at `http://localhost:5173` and the Go backend at `http://localhost:8080`.

## API Endpoints

### POST /api/simulate
Run a simulation with a single policy.

```json
{
  "policy": "LRU",
  "cache_size": 3,
  "requests": [1, 2, 3, 1, 4, 5],
  "show_steps": true
}
```

### POST /api/compare
Compare all policies on the same request sequence.

```json
{
  "cache_size": 3,
  "requests": [1, 2, 3, 1, 4, 5, 2, 1, 6, 7]
}
```

### POST /api/generate
Generate a random request trace.

```json
{
  "num_requests": 10000,
  "key_range": 1000,
  "distribution": "zipf"
}
```

### POST /api/benchmark
Run performance benchmarks across multiple cache sizes.

```json
{
  "cache_sizes": [10, 50, 100, 500, 1000],
  "num_requests": 10000,
  "key_range": 10000,
  "distribution": "zipf",
  "seed": 42
}
```

## Metrics

For each simulation, the following metrics are calculated:

| Metric | Description |
|--------|-------------|
| Cache Hits | Requests found in cache |
| Cache Misses | Requests not found in cache |
| Hit Ratio | hits / total_requests |
| Miss Ratio | misses / total_requests |
| Evictions | Number of elements removed |
| Execution Time | Time to process all requests |

## Trace Generation

The simulator can generate request sequences with different distributions:

| Distribution | Description |
|--------------|-------------|
| Uniform | Equal probability for all keys |
| Zipf | Power law distribution (some keys are more popular) |
| Localized | Working set pattern with temporal locality |
| Sequential | Mostly sequential access pattern |
| Looping | Repeated loop pattern |

## Performance Analysis

### Expected Behavior

1. **LRU** performs best when access patterns show temporal locality
2. **LFU** excels when some items are consistently more popular
3. **FIFO** is simple but doesn't adapt to access patterns
4. **Random** provides baseline comparison

### Benchmark Results

With Zipf distribution (10000 requests, key range 0-9999):

| Cache Size | FIFO Hit% | LRU Hit% | LFU Hit% | Random Hit% |
|------------|-----------|----------|----------|-------------|
| 10         | ~15%      | ~18%     | ~20%     | ~16%        |
| 50         | ~35%      | ~40%     | ~45%     | ~37%        |
| 100        | ~50%      | ~55%     | ~60%     | ~52%        |
| 500        | ~75%      | ~80%     | ~85%     | ~77%        |
| 1000       | ~85%      | ~90%     | ~92%     | ~87%        |

*Results vary based on distribution and seed*

## Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific tests
go test ./tests -run TestLRUCache

# Run benchmarks
go test -bench=. ./tests
```

## Complexity Analysis

| Operation | FIFO | LRU | LFU | Random |
|-----------|------|-----|-----|--------|
| Get       | O(1) | O(1)| O(1)| O(1)   |
| Put       | O(1) | O(1)| O(1)| O(1)   |

### Space Complexity

| Policy | Space |
|--------|-------|
| FIFO   | O(n) - Key to value map + order queue |
| LRU    | O(n) - Key to value map + access order list |
| LFU    | O(n) - Key to entry map + frequency buckets |
| Random | O(n) - Key to value map + key slice |

## Technologies Used

**Backend:**
- Go 1.21
- Gin web framework

**Frontend:**
- React 18
- Tailwind CSS
- Chart.js
- Lucide React icons

## License

MIT License
