package db

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/supabase-community/supabase-go"
)

var client *supabase.Client

// Init initializes the Supabase client
func Init() error {
	url := os.Getenv("SUPABASE_URL")
	if url == "" {
		url = os.Getenv("VITE_SUPABASE_URL")
	}
	key := os.Getenv("SUPABASE_SERVICE_ROLE_KEY")
	if key == "" {
		key = os.Getenv("SUPABASE_ANON_KEY")
	}

	if url == "" || key == "" {
		return fmt.Errorf("SUPABASE_URL and SUPABASE_SERVICE_ROLE_KEY must be set")
	}

	var err error
	client, err = supabase.NewClient(url, key, nil)
	if err != nil {
		return fmt.Errorf("failed to create supabase client: %w", err)
	}

	return nil
}

// GetClient returns the Supabase client
func GetClient() *supabase.Client {
	return client
}

// SimulationRecord for database insert
func InsertSimulation(ctx context.Context, record map[string]interface{}) error {
	_, _, err := client.From("simulations").Insert(record, false, "", "representation", "minimal").Execute()
	return err
}

// InsertComparison saves a comparison run
func InsertComparison(ctx context.Context, record map[string]interface{}) error {
	_, _, err := client.From("comparisons").Insert(record, false, "", "representation", "minimal").Execute()
	return err
}

// InsertBenchmark saves a benchmark run
func InsertBenchmark(ctx context.Context, record map[string]interface{}) error {
	_, _, err := client.From("benchmarks").Insert(record, false, "", "representation", "minimal").Execute()
	return err
}

// InsertTrace saves a generated trace
func InsertTrace(ctx context.Context, record map[string]interface{}) error {
	_, _, err := client.From("traces").Insert(record, false, "", "representation", "minimal").Execute()
	return err
}

// ListSimulations returns recent simulations
func ListSimulations(ctx context.Context, limit int) ([]map[string]interface{}, error) {
	if limit <= 0 {
		limit = 10
	}
	data, _, err := client.From("simulations").
		Select("*", "exact", false).
		Order("created_at", &supabase.OrderOptions{Ascending: false}).
		Limit(limit, "").
		Execute()
	if err != nil {
		return nil, err
	}

	var results []map[string]interface{}
	if err := json.Unmarshal(data, &results); err != nil {
		return nil, err
	}
	return results, nil
}

// ListComparisons returns recent comparisons
func ListComparisons(ctx context.Context, limit int) ([]map[string]interface{}, error) {
	if limit <= 0 {
		limit = 10
	}
	data, _, err := client.From("comparisons").
		Select("*", "exact", false).
		Order("created_at", &supabase.OrderOptions{Ascending: false}).
		Limit(limit, "").
		Execute()
	if err != nil {
		return nil, err
	}

	var results []map[string]interface{}
	if err := json.Unmarshal(data, &results); err != nil {
		return nil, err
	}
	return results, nil
}

// ListBenchmarks returns recent benchmarks
func ListBenchmarks(ctx context.Context, limit int) ([]map[string]interface{}, error) {
	if limit <= 0 {
		limit = 10
	}
	data, _, err := client.From("benchmarks").
		Select("*", "exact", false).
		Order("created_at", &supabase.OrderOptions{Ascending: false}).
		Limit(limit, "").
		Execute()
	if err != nil {
		return nil, err
	}

	var results []map[string]interface{}
	if err := json.Unmarshal(data, &results); err != nil {
		return nil, err
	}
	return results, nil
}
