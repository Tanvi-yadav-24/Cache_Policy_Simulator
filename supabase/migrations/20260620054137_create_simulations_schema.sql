/*
# Cache Simulator Database Schema

1. New Tables
- `simulations` - Stores individual simulation runs
  - `id` (uuid, primary key)
  - `policy` (text, not null) - FIFO, LRU, LFU, Random
  - `cache_size` (integer, not null)
  - `requests` (jsonb, not null) - request sequence
  - `steps` (jsonb) - step-by-step results
  - `metrics` (jsonb) - hit/miss/eviction metrics
  - `distribution` (text) - distribution type used
  - `created_at` (timestamptz)

- `comparisons` - Stores comparison runs (all 4 policies at once)
  - `id` (uuid, primary key)
  - `cache_size` (integer, not null)
  - `requests` (jsonb, not null)
  - `results` (jsonb, not null) - map of policy -> result
  - `distribution` (text)
  - `created_at` (timestamptz)

- `benchmarks` - Stores benchmark runs
  - `id` (uuid, primary key)
  - `cache_sizes` (jsonb, not null)
  - `num_requests` (integer, not null)
  - `key_range` (integer, not null)
  - `distribution` (text)
  - `results` (jsonb, not null) - array of benchmark results
  - `created_at` (timestamptz)

- `traces` - Stores generated request traces
  - `id` (uuid, primary key)
  - `num_requests` (integer, not null)
  - `key_range` (integer, not null)
  - `distribution` (text, not null)
  - `requests` (jsonb, not null)
  - `created_at` (timestamptz)

2. Security
- Enable RLS on all tables.
- Allow anon + authenticated CRUD because this is a single-tenant demo app.
*/

CREATE TABLE IF NOT EXISTS simulations (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  policy text NOT NULL,
  cache_size integer NOT NULL,
  requests jsonb NOT NULL,
  steps jsonb,
  metrics jsonb,
  distribution text,
  created_at timestamptz DEFAULT now()
);

CREATE TABLE IF NOT EXISTS comparisons (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  cache_size integer NOT NULL,
  requests jsonb NOT NULL,
  results jsonb NOT NULL,
  distribution text,
  created_at timestamptz DEFAULT now()
);

CREATE TABLE IF NOT EXISTS benchmarks (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  cache_sizes jsonb NOT NULL,
  num_requests integer NOT NULL,
  key_range integer NOT NULL,
  distribution text,
  results jsonb NOT NULL,
  created_at timestamptz DEFAULT now()
);

CREATE TABLE IF NOT EXISTS traces (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  num_requests integer NOT NULL,
  key_range integer NOT NULL,
  distribution text NOT NULL,
  requests jsonb NOT NULL,
  created_at timestamptz DEFAULT now()
);

ALTER TABLE simulations ENABLE ROW LEVEL SECURITY;
ALTER TABLE comparisons ENABLE ROW LEVEL SECURITY;
ALTER TABLE benchmarks ENABLE ROW LEVEL SECURITY;
ALTER TABLE traces ENABLE ROW LEVEL SECURITY;

-- Simulations policies
DROP POLICY IF EXISTS "anon_select_simulations" ON simulations;
CREATE POLICY "anon_select_simulations" ON simulations FOR SELECT
TO anon, authenticated USING (true);

DROP POLICY IF EXISTS "anon_insert_simulations" ON simulations;
CREATE POLICY "anon_insert_simulations" ON simulations FOR INSERT
TO anon, authenticated WITH CHECK (true);

DROP POLICY IF EXISTS "anon_update_simulations" ON simulations;
CREATE POLICY "anon_update_simulations" ON simulations FOR UPDATE
TO anon, authenticated USING (true) WITH CHECK (true);

DROP POLICY IF EXISTS "anon_delete_simulations" ON simulations;
CREATE POLICY "anon_delete_simulations" ON simulations FOR DELETE
TO anon, authenticated USING (true);

-- Comparisons policies
DROP POLICY IF EXISTS "anon_select_comparisons" ON comparisons;
CREATE POLICY "anon_select_comparisons" ON comparisons FOR SELECT
TO anon, authenticated USING (true);

DROP POLICY IF EXISTS "anon_insert_comparisons" ON comparisons;
CREATE POLICY "anon_insert_comparisons" ON comparisons FOR INSERT
TO anon, authenticated WITH CHECK (true);

DROP POLICY IF EXISTS "anon_update_comparisons" ON comparisons;
CREATE POLICY "anon_update_comparisons" ON comparisons FOR UPDATE
TO anon, authenticated USING (true) WITH CHECK (true);

DROP POLICY IF EXISTS "anon_delete_comparisons" ON comparisons;
CREATE POLICY "anon_delete_comparisons" ON comparisons FOR DELETE
TO anon, authenticated USING (true);

-- Benchmarks policies
DROP POLICY IF EXISTS "anon_select_benchmarks" ON benchmarks;
CREATE POLICY "anon_select_benchmarks" ON benchmarks FOR SELECT
TO anon, authenticated USING (true);

DROP POLICY IF EXISTS "anon_insert_benchmarks" ON benchmarks;
CREATE POLICY "anon_insert_benchmarks" ON benchmarks FOR INSERT
TO anon, authenticated WITH CHECK (true);

DROP POLICY IF EXISTS "anon_update_benchmarks" ON benchmarks;
CREATE POLICY "anon_update_benchmarks" ON benchmarks FOR UPDATE
TO anon, authenticated USING (true) WITH CHECK (true);

DROP POLICY IF EXISTS "anon_delete_benchmarks" ON benchmarks;
CREATE POLICY "anon_delete_benchmarks" ON benchmarks FOR DELETE
TO anon, authenticated USING (true);

-- Traces policies
DROP POLICY IF EXISTS "anon_select_traces" ON traces;
CREATE POLICY "anon_select_traces" ON traces FOR SELECT
TO anon, authenticated USING (true);

DROP POLICY IF EXISTS "anon_insert_traces" ON traces;
CREATE POLICY "anon_insert_traces" ON traces FOR INSERT
TO anon, authenticated WITH CHECK (true);

DROP POLICY IF EXISTS "anon_update_traces" ON traces;
CREATE POLICY "anon_update_traces" ON traces FOR UPDATE
TO anon, authenticated USING (true) WITH CHECK (true);

DROP POLICY IF EXISTS "anon_delete_traces" ON traces;
CREATE POLICY "anon_delete_traces" ON traces FOR DELETE
TO anon, authenticated USING (true);
