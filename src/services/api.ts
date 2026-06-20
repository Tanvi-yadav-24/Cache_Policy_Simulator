import { Simulator, runComparison, SimulationResult } from '../lib/simulator';
import { TraceGenerator, runBenchmark, BenchmarkResult, DistributionType } from '../lib/generator';

export interface GenerateResponse {
  requests: number[];
  stats: {
    num_requests: number;
    key_range: number;
    distribution: string;
  };
}

const API_BASE = import.meta.env.VITE_API_URL || 'http://localhost:8080/api';

async function fetchBackend<T>(path: string, options?: RequestInit): Promise<T | null> {
  try {
    const res = await fetch(`${API_BASE}${path}`, {
      headers: { 'Content-Type': 'application/json' },
      ...options,
    });
    if (!res.ok) return null;
    return await res.json() as T;
  } catch {
    return null;
  }
}

export const api = {
  // Check if backend is available
  async isBackendAvailable(): Promise<boolean> {
    try {
      const res = await fetch(`${API_BASE}/health`, { method: 'GET' });
      return res.ok;
    } catch {
      return false;
    }
  },

  async simulate(
    policy: string,
    cacheSize: number,
    requests: number[],
    showSteps: boolean = true
  ): Promise<SimulationResult> {
    const backend = await fetchBackend<SimulationResult>('/simulate', {
      method: 'POST',
      body: JSON.stringify({
        policy,
        cache_size: cacheSize,
        requests,
        show_steps: showSteps,
      }),
    });
    if (backend) return backend;

    // Fallback to local
    const sim = new Simulator(policy, cacheSize);
    const result = sim.run(requests);
    if (!showSteps) {
      result.steps = [];
    }
    return result;
  },

  async compare(cacheSize: number, requests: number[]): Promise<Record<string, SimulationResult>> {
    const backend = await fetchBackend<Record<string, SimulationResult>>('/compare', {
      method: 'POST',
      body: JSON.stringify({ cache_size: cacheSize, requests }),
    });
    if (backend) return backend;

    // Fallback to local
    return runComparison(cacheSize, requests);
  },

  async generateTrace(
    numRequests: number,
    keyRange: number,
    distribution: DistributionType,
    seed?: number
  ): Promise<GenerateResponse> {
    const backend = await fetchBackend<GenerateResponse>('/generate', {
      method: 'POST',
      body: JSON.stringify({
        num_requests: numRequests,
        key_range: keyRange,
        distribution,
        seed,
      }),
    });
    if (backend) return backend;

    // Fallback to local
    const gen = new TraceGenerator(seed ?? Date.now());
    const requests = gen.generate(numRequests, keyRange, distribution);
    return {
      requests,
      stats: {
        num_requests: numRequests,
        key_range: keyRange,
        distribution,
      },
    };
  },

  async benchmark(
    cacheSizes: number[],
    numRequests: number,
    keyRange: number,
    distribution: DistributionType,
    _seed: number = 42
  ): Promise<{ configuration: { cache_sizes: number[]; num_requests: number; key_range: number; distribution: string; seed: number }; results: BenchmarkResult[] }> {
    const backend = await fetchBackend<any>('/benchmark', {
      method: 'POST',
      body: JSON.stringify({
        cache_sizes: cacheSizes,
        num_requests: numRequests,
        key_range: keyRange,
        distribution,
        seed: _seed,
      }),
    });
    if (backend) {
      return {
        configuration: backend.configuration,
        results: backend.results,
      };
    }

    // Fallback to local
    const results = runBenchmark(cacheSizes, numRequests, keyRange, distribution);
    return {
      configuration: {
        cache_sizes: cacheSizes,
        num_requests: numRequests,
        key_range: keyRange,
        distribution,
        seed: _seed,
      },
      results,
    };
  },

  // History endpoints
  async getSimulationHistory(): Promise<any[]> {
    const data = await fetchBackend<any[]>('/history/simulations');
    return data ?? [];
  },

  async getComparisonHistory(): Promise<any[]> {
    const data = await fetchBackend<any[]>('/history/comparisons');
    return data ?? [];
  },

  async getBenchmarkHistory(): Promise<any[]> {
    const data = await fetchBackend<any[]>('/history/benchmarks');
    return data ?? [];
  },
};
