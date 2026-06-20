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

export const api = {
  // Always available - uses local simulation
  isLocal(): boolean {
    return true;
  },

  async simulate(
    policy: string,
    cacheSize: number,
    requests: number[],
    showSteps: boolean = true
  ): Promise<SimulationResult> {
    const sim = new Simulator(policy, cacheSize);
    const result = sim.run(requests);
    if (!showSteps) {
      result.steps = [];
    }
    return result;
  },

  async compare(cacheSize: number, requests: number[]): Promise<Record<string, SimulationResult>> {
    return runComparison(cacheSize, requests);
  },

  async generateTrace(
    numRequests: number,
    keyRange: number,
    distribution: DistributionType,
    seed?: number
  ): Promise<GenerateResponse> {
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
};
