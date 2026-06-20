package tests

import (
	"fmt"
	"testing"

	"cache-simulator/internal/cache"
	"cache-simulator/internal/generator"
	"cache-simulator/internal/simulator"
)

// BenchmarkCachePolicies benchmarks different cache policies
func BenchmarkCachePolicies(b *testing.B) {
	policies := []cache.CachePolicy{cache.FIFO, cache.LRU, cache.LFU, cache.Random}
	cacheSizes := []int{10, 50, 100, 500, 1000}

	gen := generator.NewTraceGenerator()
	gen.SetSeed(42)
	requests := gen.Generate(10000, 10000, generator.Zipf)

	for _, policy := range policies {
		for _, size := range cacheSizes {
			name := fmt.Sprintf("%s_Size%d", policy, size)
			b.Run(name, func(b *testing.B) {
				sim := simulator.NewSimulator(policy, size)
				for i := 0; i < b.N; i++ {
					sim.Run(requests)
				}
			})
		}
	}
}

// BenchmarkCacheOperations benchmarks individual cache operations
func BenchmarkCacheOperations(b *testing.B) {
	policies := []cache.CachePolicy{cache.FIFO, cache.LRU, cache.LFU, cache.Random}

	for _, policy := range policies {
		b.Run(string(policy)+"_Put", func(b *testing.B) {
			c := cache.Factory(policy, 10000)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				c.Put(i, i*10)
			}
		})

		b.Run(string(policy)+"_Get", func(b *testing.B) {
			c := cache.Factory(policy, 10000)
			// Pre-fill cache
			for i := 0; i < 10000; i++ {
				c.Put(i, i*10)
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				c.Get(i % 10000)
			}
		})

		b.Run(string(policy)+"_Mixed", func(b *testing.B) {
			c := cache.Factory(policy, 10000)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if i%2 == 0 {
					c.Put(i, i*10)
				} else {
					c.Get(i % 10000)
				}
			}
		})
	}
}

// BenchmarkGenerator benchmarks trace generation
func BenchmarkGenerator(b *testing.B) {
	gen := generator.NewTraceGenerator()
	gen.SetSeed(42)

	b.Run("Uniform_1000", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			gen.Generate(1000, 1000, generator.Uniform)
		}
	})

	b.Run("Zipf_1000", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			gen.Generate(1000, 1000, generator.Zipf)
		}
	})

	b.Run("Localized_1000", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			gen.GenerateLocalized(1000, 1000, 100)
		}
	})
}

// BenchmarkSimulation benchmarks full simulation runs
func BenchmarkSimulation(b *testing.B) {
	gen := generator.NewTraceGenerator()
	gen.SetSeed(42)

	sizes := []int{100, 1000, 10000, 100000}

	for _, size := range sizes {
		requests := gen.Generate(size, size/2, generator.Uniform)

		b.Run(fmt.Sprintf("Full_%d_requests", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				simulator.RunComparison(100, requests)
			}
		})
	}
}
