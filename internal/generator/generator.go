package generator

import (
	"math"
	"math/rand"
	"time"
)

// DistributionType represents the type of trace distribution
type DistributionType string

const (
	Uniform DistributionType = "uniform"
	Zipf    DistributionType = "zipf"
)

// TraceGenerator generates request sequences with different distributions
type TraceGenerator struct {
	rng *rand.Rand
}

// NewTraceGenerator creates a new trace generator
func NewTraceGenerator() *TraceGenerator {
	return &TraceGenerator{
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// SetSeed sets the random seed for reproducible traces
func (g *TraceGenerator) SetSeed(seed int64) {
	g.rng = rand.New(rand.NewSource(seed))
}

// Generate creates a request sequence with the specified parameters
// numRequests: total number of requests to generate
// keyRange: range of possible key values (0 to keyRange-1)
// distType: distribution type (uniform or zipf)
func (g *TraceGenerator) Generate(numRequests, keyRange int, distType DistributionType) []int {
	switch distType {
	case Uniform:
		return g.generateUniform(numRequests, keyRange)
	case Zipf:
		return g.generateZipf(numRequests, keyRange)
	default:
		return g.generateUniform(numRequests, keyRange)
	}
}

// generateUniform generates a uniformly distributed request sequence
func (g *TraceGenerator) generateUniform(numRequests, keyRange int) []int {
	requests := make([]int, numRequests)
	for i := 0; i < numRequests; i++ {
		requests[i] = g.rng.Intn(keyRange)
	}
	return requests
}

// generateZipf generates a Zipf-distributed request sequence
// Zipf distribution models real-world access patterns where some items are more popular
func (g *TraceGenerator) generateZipf(numRequests, keyRange int) []int {
	// Zipf parameters
	s := 1.07 // skew parameter (higher = more skewed)

	// Calculate Zipf probabilities
	probabilities := make([]float64, keyRange)
	sum := 0.0
	for i := 1; i <= keyRange; i++ {
		prob := 1.0 / math.Pow(float64(i), s)
		probabilities[i-1] = prob
		sum += prob
	}

	// Normalize probabilities
	for i := range probabilities {
		probabilities[i] /= sum
	}

	// Generate requests based on cumulative distribution
	requests := make([]int, numRequests)
	for i := 0; i < numRequests; i++ {
		requests[i] = g.sampleFromDistribution(probabilities)
	}
	return requests
}

// sampleFromDistribution samples an index from a probability distribution
func (g *TraceGenerator) sampleFromDistribution(probabilities []float64) int {
	r := g.rng.Float64()
	cumSum := 0.0
	for i, prob := range probabilities {
		cumSum += prob
		if r <= cumSum {
			return i
		}
	}
	return len(probabilities) - 1 // Fallback
}

// GenerateLocalized generates a sequence with temporal locality
// Simulates working set patterns common in real applications
func (g *TraceGenerator) GenerateLocalized(numRequests, keyRange, workingSet int) []int {
	requests := make([]int, numRequests)

	// Create a smaller working set that changes occasionally
	workingSetKeys := make([]int, workingSet)
	for i := 0; i < workingSet; i++ {
		workingSetKeys[i] = g.rng.Intn(keyRange)
	}

	for i := 0; i < numRequests; i++ {
		// Occasionally shift the working set
		if g.rng.Float64() < 0.05 { // 5% chance to shift
			replaceIdx := g.rng.Intn(workingSet)
			workingSetKeys[replaceIdx] = g.rng.Intn(keyRange)
		}

		// Most requests (90%) come from working set
		if g.rng.Float64() < 0.9 {
			requests[i] = workingSetKeys[g.rng.Intn(workingSet)]
		} else {
			// 10% random requests outside working set
			requests[i] = g.rng.Intn(keyRange)
		}
	}
	return requests
}

// GenerateSequential generates a sequential access pattern with some random jumps
func (g *TraceGenerator) GenerateSequential(numRequests, keyRange int) []int {
	requests := make([]int, numRequests)
	current := g.rng.Intn(keyRange)

	for i := 0; i < numRequests; i++ {
		if g.rng.Float64() < 0.8 {
			// Sequential access (with wrap-around)
			requests[i] = current
			current = (current + 1) % keyRange
		} else {
			// Random jump
			requests[i] = g.rng.Intn(keyRange)
			current = requests[i]
		}
	}
	return requests
}

// GenerateLooping generates a looping pattern (common in programs)
func (g *TraceGenerator) GenerateLooping(numRequests, loopSize int) []int {
	requests := make([]int, numRequests)
	for i := 0; i < numRequests; i++ {
		requests[i] = i % loopSize
	}
	return requests
}
