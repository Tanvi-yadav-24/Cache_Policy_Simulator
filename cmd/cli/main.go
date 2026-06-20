package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"cache-simulator/internal/cache"
	"cache-simulator/internal/generator"
	"cache-simulator/internal/simulator"
)

func main() {
	// Subcommands
	simCmd := flag.NewFlagSet("simulate", flag.ExitOnError)
	compareCmd := flag.NewFlagSet("compare", flag.ExitOnError)
	genCmd := flag.NewFlagSet("generate", flag.ExitOnError)
	benchCmd := flag.NewFlagSet("benchmark", flag.ExitOnError)

	// Simulation flags
	policy := simCmd.String("policy", "FIFO", "Cache policy: FIFO, LRU, LFU, Random")
	cacheSize := simCmd.Int("size", 3, "Cache size/capacity")
	requests := simCmd.String("requests", "", "Comma-separated request sequence (e.g., '1,2,3,4,5')")
	showSteps := simCmd.Bool("steps", true, "Show step-by-step results")
	outputFormat := simCmd.String("output", "text", "Output format: text or json")

	// Compare flags
	compareCacheSize := compareCmd.Int("size", 3, "Cache size/capacity")
	compareRequests := compareCmd.String("requests", "", "Comma-separated request sequence")

	// Generate flags
	genNum := genCmd.Int("num", 100, "Number of requests to generate")
	genRange := genCmd.Int("range", 100, "Key range (0 to range-1)")
	genDist := genCmd.String("dist", "uniform", "Distribution: uniform, zipf, localized, sequential, looping")
	genSeed := genCmd.Int64("seed", 0, "Random seed (0 for time-based)")

	// Benchmark flags
	benchNum := benchCmd.Int("num", 10000, "Number of requests")
	benchRange := benchCmd.Int("range", 10000, "Key range")
	benchDist := benchCmd.String("dist", "zipf", "Distribution type")
	benchSeed := benchCmd.Int64("seed", 42, "Random seed")
	benchSizes := benchCmd.String("sizes", "10,50,100,500,1000", "Comma-separated cache sizes")

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "simulate":
		simCmd.Parse(os.Args[2:])
	 runSimulation(*policy, *cacheSize, *requests, *showSteps, *outputFormat)
	case "compare":
		compareCmd.Parse(os.Args[2:])
		runComparison(*compareCacheSize, *compareRequests, *outputFormat)
	case "generate":
		genCmd.Parse(os.Args[2:])
		generateTrace(*genNum, *genRange, *genDist, *genSeed, *outputFormat)
	case "benchmark":
		benchCmd.Parse(os.Args[2:])
		runBenchmark(*benchNum, *benchRange, *benchDist, *benchSeed, *benchSizes, *outputFormat)
	default:
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Cache Policy Simulator")
	fmt.Println("")
	fmt.Println("Commands:")
	fmt.Println("  simulate   - Run a simulation with a single policy")
	fmt.Println("  compare    - Compare all policies on the same input")
	fmt.Println("  generate   - Generate a request trace")
	fmt.Println("  benchmark  - Run performance benchmarks")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  cache-simulator simulate -policy LRU -size 3 -requests \"1,2,3,1,4,5\"")
	fmt.Println("  cache-simulator compare -size 3 -requests \"1,2,3,1,4,5\"")
	fmt.Println("  cache-simulator generate -num 100 -range 50 -dist zipf")
	fmt.Println("  cache-simulator benchmark -num 10000 -dist zipf")
}

func parseRequests(reqStr string) []int {
	if reqStr == "" {
		return []int{1, 2, 3, 1, 4, 5, 2, 1, 6, 7}
	}

	parts := strings.Split(reqStr, ",")
	requests := make([]int, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		n, err := strconv.Atoi(p)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing request: %s\n", p)
			os.Exit(1)
		}
		requests = append(requests, n)
	}
	return requests
}

func runSimulation(policyStr string, size int, reqStr string, showSteps bool, format string) {
	reqs := parseRequests(reqStr)
	policy := cache.CachePolicy(policyStr)

	sim := simulator.NewSimulator(policy, size)
	result := sim.Run(reqs)

	if format == "json" {
		if !showSteps {
			result.Steps = nil
		}
		data, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(data))
		return
	}

	// Text output
	fmt.Printf("\n=== %s Cache Simulation ===\n", policyStr)
	fmt.Printf("Cache Size: %d\n", size)
	fmt.Printf("Requests: %v\n\n", reqs)

	if showSteps {
		fmt.Println("--- Step by Step ---")
		for i, step := range result.Steps {
			fmt.Printf("Step %d: Request %d\n", i+1, step.Request)
			fmt.Printf("  Cache Before: %v\n", step.CacheBefore)
			fmt.Printf("  Result: %s\n", map[bool]string{true: "HIT", false: "MISS"}[step.Hit])
			if step.Evicted != 0 {
				fmt.Printf("  Evicted: %d\n", step.Evicted)
			}
			fmt.Printf("  Cache After:  %v\n\n", step.CacheAfter)
		}
	}

	fmt.Println("--- Summary ---")
	fmt.Printf("Total Requests: %d\n", result.Metrics.TotalRequests)
	fmt.Printf("Cache Hits: %d\n", result.Metrics.CacheHits)
	fmt.Printf("Cache Misses: %d\n", result.Metrics.CacheMisses)
	fmt.Printf("Evictions: %d\n", result.Metrics.Evictions)
	fmt.Printf("Hit Ratio: %.4f (%.2f%%)\n", result.Metrics.HitRatio(), result.Metrics.HitRatio()*100)
	fmt.Printf("Miss Ratio: %.4f (%.2f%%)\n", result.Metrics.MissRatio(), result.Metrics.MissRatio()*100)
	fmt.Printf("Execution Time: %d ns\n", result.Metrics.ExecutionTimeNanos)
}

func runComparison(size int, reqStr string, format string) {
	reqs := parseRequests(reqStr)
	results := simulator.RunComparison(size, reqs)

	if format == "json" {
		data, _ := json.MarshalIndent(results, "", "  ")
		fmt.Println(string(data))
		return
	}

	// Text output
	fmt.Printf("\n=== Policy Comparison ===\n")
	fmt.Printf("Cache Size: %d\n", size)
	fmt.Printf("Requests: %v\n\n", reqs)

	fmt.Printf("%-10s %6s %6s %6s %10s %10s\n", "Policy", "Hits", "Misses", "Evicts", "Hit%", "Miss%")
	fmt.Println(strings.Repeat("-", 55))

	for _, policy := range []string{"FIFO", "LRU", "LFU", "Random"} {
		r := results[policy]
		fmt.Printf("%-10s %6d %6d %6d %9.2f%% %9.2f%%\n",
			policy,
			r.Metrics.CacheHits,
			r.Metrics.CacheMisses,
			r.Metrics.Evictions,
			r.Metrics.HitRatio()*100,
			r.Metrics.MissRatio()*100,
		)
	}
}

func generateTrace(num, keyRange int, dist string, seed int64, format string) {
	gen := generator.NewTraceGenerator()
	if seed != 0 {
		gen.SetSeed(seed)
	}

	var requests []int
	distType := generator.DistributionType(dist)

	switch dist {
	case "localized":
		requests = gen.GenerateLocalized(num, keyRange, keyRange/10)
	case "sequential":
		requests = gen.GenerateSequential(num, keyRange)
	case "looping":
		requests = gen.GenerateLooping(num, keyRange/10)
	default:
		requests = gen.Generate(num, keyRange, distType)
	}

	if format == "json" {
		data, _ := json.MarshalIndent(map[string]interface{}{
			"requests": requests,
			"stats": map[string]interface{}{
				"num_requests": num,
				"key_range":    keyRange,
				"distribution": dist,
			},
		}, "", "  ")
		fmt.Println(string(data))
		return
	}

	fmt.Printf("Generated %d requests with %s distribution (range: 0-%d)\n", num, dist, keyRange-1)
	fmt.Printf("First 20 requests: %v\n", requests[:min(20, len(requests))])
	fmt.Printf("Use -output json for full trace\n")
}

func runBenchmark(num, keyRange int, dist string, seed int64, sizesStr, format string) {
	sizes := parseSizes(sizesStr)

	gen := generator.NewTraceGenerator()
	gen.SetSeed(seed)

	distType := generator.DistributionType(dist)
	var requests []int
	switch dist {
	case "localized":
		requests = gen.GenerateLocalized(num, keyRange, keyRange/10)
	case "sequential":
		requests = gen.GenerateSequential(num, keyRange)
	case "looping":
		requests = gen.GenerateLooping(num, keyRange/10)
	default:
		requests = gen.Generate(num, keyRange, distType)
	}

	if format == "json" {
		results := make(map[string]interface{})
		results["configuration"] = map[string]interface{}{
			"num_requests": num,
			"key_range":    keyRange,
			"distribution": dist,
			"cache_sizes":  sizes,
		}

		var benchResults []map[string]interface{}
		for _, size := range sizes {
			for _, policy := range []cache.CachePolicy{cache.FIFO, cache.LRU, cache.LFU, cache.Random} {
				sim := simulator.NewSimulator(policy, size)
				result := sim.Run(requests)
				benchResults = append(benchResults, map[string]interface{}{
					"policy":      string(policy),
					"cache_size":  size,
					"hit_ratio":   result.Metrics.HitRatio(),
					"miss_ratio":  result.Metrics.MissRatio(),
					"evictions":   result.Metrics.Evictions,
					"exec_time_nanos": result.Metrics.ExecutionTimeNanos,
				})
			}
		}
		results["results"] = benchResults

		data, _ := json.MarshalIndent(results, "", "  ")
		fmt.Println(string(data))
		return
	}

	// Text output
	fmt.Printf("\n=== Performance Benchmark ===\n")
	fmt.Printf("Requests: %d | Range: %d | Distribution: %s\n", num, keyRange, dist)
	fmt.Println()

	for _, size := range sizes {
		fmt.Printf("\n--- Cache Size: %d ---\n", size)
		fmt.Printf("%-10s %10s %10s %10s %15s\n", "Policy", "Hit%", "Miss%", "Evictions", "Exec Time (ns)")
		fmt.Println(strings.Repeat("-", 60))

		for _, policy := range []cache.CachePolicy{cache.FIFO, cache.LRU, cache.LFU, cache.Random} {
			sim := simulator.NewSimulator(policy, size)
			result := sim.Run(requests)
			fmt.Printf("%-10s %9.2f%% %9.2f%% %10d %15d\n",
				policy,
				result.Metrics.HitRatio()*100,
				result.Metrics.MissRatio()*100,
				result.Metrics.Evictions,
				result.Metrics.ExecutionTimeNanos,
			)
		}
	}
}

func parseSizes(sizesStr string) []int {
	parts := strings.Split(sizesStr, ",")
	sizes := make([]int, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		n, err := strconv.Atoi(p)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing size: %s\n", p)
			os.Exit(1)
		}
		sizes = append(sizes, n)
	}
	return sizes
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
