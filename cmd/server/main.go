package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"cache-simulator/internal/cache"
	"cache-simulator/internal/db"
	"cache-simulator/internal/generator"
	"cache-simulator/internal/metrics"
	"cache-simulator/internal/simulator"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database connection
	if err := db.Init(); err != nil {
		// Log but don't fail - backend works without DB in local mode
		gin.DefaultWriter.Write([]byte("DB init warning: " + err.Error() + "\n"))
	}

	r := gin.Default()

	// Configure CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Client-Info", "Apikey"},
		AllowCredentials: true,
	}))

	api := r.Group("/api")
	{
		api.GET("/health", healthCheck)
		api.GET("/policies", getPolicies)

		api.POST("/simulate", runSimulation)
		api.POST("/generate", generateTrace)
		api.POST("/compare", runComparison)
		api.POST("/benchmark", runBenchmark)

		// History endpoints
		api.GET("/history/simulations", listSimulations)
		api.GET("/history/comparisons", listComparisons)
		api.GET("/history/benchmarks", listBenchmarks)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"timestamp": time.Now().Unix(),
		"version":   "1.0.0",
	})
}

// PolicyInfo contains information about a cache policy
type PolicyInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func getPolicies(c *gin.Context) {
	policies := []PolicyInfo{
		{"FIFO", "First In First Out - Removes the oldest element"},
		{"LRU", "Least Recently Used - Removes the least recently accessed item"},
		{"LFU", "Least Frequently Used - Removes the least frequently used item"},
		{"Random", "Random Replacement - Randomly evicts an item"},
	}
	c.JSON(http.StatusOK, policies)
}

// SimulationRequest represents a simulation request
type SimulationRequest struct {
	Policy       string `json:"policy"`
	CacheSize    int    `json:"cache_size"`
	Requests     []int  `json:"requests"`
	ShowSteps    bool   `json:"show_steps"`
	Distribution string `json:"distribution,omitempty"`
}

func runSimulation(c *gin.Context) {
	var req SimulationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	policy := cache.CachePolicy(req.Policy)
	sim := simulator.NewSimulator(policy, req.CacheSize)
	result := sim.Run(req.Requests)

	// If not showing steps, remove them from response
	if !req.ShowSteps {
		result.Steps = nil
	}

	// Persist to database (async, don't block response)
	go func() {
		if db.GetClient() != nil {
			record := map[string]interface{}{
				"policy":       req.Policy,
				"cache_size":   req.CacheSize,
				"requests":     req.Requests,
				"steps":        result.Steps,
				"metrics":      result.Metrics,
				"distribution": req.Distribution,
			}
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			db.InsertSimulation(ctx, record)
		}
	}()

	c.JSON(http.StatusOK, result)
}

// TraceRequest represents a trace generation request
type TraceRequest struct {
	NumRequests  int    `json:"num_requests"`
	KeyRange     int    `json:"key_range"`
	Distribution string `json:"distribution"`
	Seed         *int64 `json:"seed,omitempty"`
	WorkingSet   *int   `json:"working_set,omitempty"` // For localized pattern
}

func generateTrace(c *gin.Context) {
	var req TraceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	gen := generator.NewTraceGenerator()
	if req.Seed != nil {
		gen.SetSeed(*req.Seed)
	}

	distType := generator.DistributionType(req.Distribution)
	var requests []int

	if distType == "localized" && req.WorkingSet != nil {
		requests = gen.GenerateLocalized(req.NumRequests, req.KeyRange, *req.WorkingSet)
	} else if distType == "sequential" {
		requests = gen.GenerateSequential(req.NumRequests, req.KeyRange)
	} else if distType == "looping" && req.WorkingSet != nil {
		requests = gen.GenerateLooping(req.NumRequests, *req.WorkingSet)
	} else {
		requests = gen.Generate(req.NumRequests, req.KeyRange, distType)
	}

	// Persist to database (async)
	go func() {
		if db.GetClient() != nil {
			record := map[string]interface{}{
				"num_requests": req.NumRequests,
				"key_range":    req.KeyRange,
				"distribution": req.Distribution,
				"requests":     requests,
			}
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			db.InsertTrace(ctx, record)
		}
	}()

	c.JSON(http.StatusOK, gin.H{
		"requests": requests,
		"stats": gin.H{
			"num_requests": len(requests),
			"key_range":    req.KeyRange,
			"distribution": req.Distribution,
		},
	})
}

// ComparisonRequest represents a comparison request
type ComparisonRequest struct {
	CacheSize    int     `json:"cache_size"`
	Requests     []int   `json:"requests"`
	Distribution string  `json:"distribution,omitempty"`
	Seeds        []int64 `json:"seeds,omitempty"` // For random policy reproducibility
}

func runComparison(c *gin.Context) {
	var req ComparisonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	results := simulator.RunComparison(req.CacheSize, req.Requests)

	// Persist to database (async)
	go func() {
		if db.GetClient() != nil {
			record := map[string]interface{}{
				"cache_size":   req.CacheSize,
				"requests":     req.Requests,
				"results":      results,
				"distribution": req.Distribution,
			}
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			db.InsertComparison(ctx, record)
		}
	}()

	c.JSON(http.StatusOK, results)
}

// BenchmarkRequest represents a benchmark request
type BenchmarkRequest struct {
	CacheSizes   []int  `json:"cache_sizes"`
	NumRequests  int    `json:"num_requests"`
	KeyRange     int    `json:"key_range"`
	Distribution string `json:"distribution"`
	Seed         int64  `json:"seed"`
}

// BenchmarkResult represents benchmark results for a single configuration
type BenchmarkResult struct {
	Policy     string  `json:"policy"`
	CacheSize  int     `json:"cache_size"`
	HitRatio   float64 `json:"hit_ratio"`
	MissRatio  float64 `json:"miss_ratio"`
	Evictions  int     `json:"evictions"`
	ExecTimeNs int64   `json:"exec_time_ns"`
}

func runBenchmark(c *gin.Context) {
	var req BenchmarkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(req.CacheSizes) == 0 {
		req.CacheSizes = []int{10, 50, 100, 500, 1000}
	}
	if req.NumRequests == 0 {
		req.NumRequests = 10000
	}
	if req.KeyRange == 0 {
		req.KeyRange = 10000
	}
	if req.Distribution == "" {
		req.Distribution = "zipf"
	}

	gen := generator.NewTraceGenerator()
	gen.SetSeed(req.Seed)

	distType := generator.DistributionType(req.Distribution)
	requests := gen.Generate(req.NumRequests, req.KeyRange, distType)

	var results []BenchmarkResult
	policies := []cache.CachePolicy{cache.FIFO, cache.LRU, cache.LFU, cache.Random}

	for _, size := range req.CacheSizes {
		for _, policy := range policies {
			sim := simulator.NewSimulator(policy, size)
			result := sim.Run(requests)

			results = append(results, BenchmarkResult{
				Policy:     string(policy),
				CacheSize:  size,
				HitRatio:   result.Metrics.HitRatio(),
				MissRatio:  result.Metrics.MissRatio(),
				Evictions:  result.Metrics.Evictions,
				ExecTimeNs: result.Metrics.ExecutionTimeNanos,
			})
		}
	}

	// Persist to database (async)
	go func() {
		if db.GetClient() != nil {
			record := map[string]interface{}{
				"cache_sizes":  req.CacheSizes,
				"num_requests": req.NumRequests,
				"key_range":    req.KeyRange,
				"distribution": req.Distribution,
				"results":      results,
			}
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			db.InsertBenchmark(ctx, record)
		}
	}()

	c.JSON(http.StatusOK, gin.H{
		"configuration": gin.H{
			"cache_sizes":  req.CacheSizes,
			"num_requests": req.NumRequests,
			"key_range":    req.KeyRange,
			"distribution": req.Distribution,
			"seed":         req.Seed,
		},
		"results": results,
	})
}

// History endpoints
func listSimulations(c *gin.Context) {
	if db.GetClient() == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "database not configured"})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	results, err := db.ListSimulations(ctx, 20)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, results)
}

func listComparisons(c *gin.Context) {
	if db.GetClient() == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "database not configured"})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	results, err := db.ListComparisons(ctx, 20)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, results)
}

func listBenchmarks(c *gin.Context) {
	if db.GetClient() == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "database not configured"})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	results, err := db.ListBenchmarks(ctx, 20)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, results)
}
