import { Simulator, SimulationResult } from './simulator';

export type DistributionType = 'uniform' | 'zipf' | 'localized' | 'sequential' | 'looping';

export class TraceGenerator {
  private seed: number;

  constructor(seed?: number) {
    this.seed = seed ?? Date.now();
  }

  setSeed(seed: number): void {
    this.seed = seed;
  }

  // Simple seeded random number generator (Linear Congruential Generator)
  private random(): number {
    this.seed = (this.seed * 1103515245 + 12345) & 0x7fffffff;
    return this.seed / 0x7fffffff;
  }

  generate(numRequests: number, keyRange: number, distribution: DistributionType): number[] {
    switch (distribution) {
      case 'uniform':
        return this.generateUniform(numRequests, keyRange);
      case 'zipf':
        return this.generateZipf(numRequests, keyRange);
      case 'localized':
        return this.generateLocalized(numRequests, keyRange);
      case 'sequential':
        return this.generateSequential(numRequests, keyRange);
      case 'looping':
        return this.generateLooping(numRequests, keyRange);
      default:
        return this.generateUniform(numRequests, keyRange);
    }
  }

  private generateUniform(numRequests: number, keyRange: number): number[] {
    const requests: number[] = [];
    for (let i = 0; i < numRequests; i++) {
      requests.push(Math.floor(this.random() * keyRange));
    }
    return requests;
  }

  private generateZipf(numRequests: number, keyRange: number): number[] {
    const s = 1.07; // Skew parameter

    // Calculate Zipf probabilities
    const probabilities: number[] = [];
    let sum = 0;
    for (let i = 1; i <= keyRange; i++) {
      const prob = 1 / Math.pow(i, s);
      probabilities.push(prob);
      sum += prob;
    }

    // Normalize
    for (let i = 0; i < probabilities.length; i++) {
      probabilities[i] /= sum;
    }

    // Generate requests
    const requests: number[] = [];
    for (let i = 0; i < numRequests; i++) {
      const r = this.random();
      let cumSum = 0;
      for (let j = 0; j < probabilities.length; j++) {
        cumSum += probabilities[j];
        if (r <= cumSum) {
          requests.push(j);
          break;
        }
      }
    }

    return requests;
  }

  private generateLocalized(numRequests: number, keyRange: number): number[] {
    const workingSet = Math.max(1, Math.min(Math.floor(keyRange / 10), 100));
    const workingSetKeys: number[] = [];

    // Initialize working set
    for (let i = 0; i < workingSet; i++) {
      workingSetKeys.push(Math.floor(this.random() * keyRange));
    }

    const requests: number[] = [];
    for (let i = 0; i < numRequests; i++) {
      // Occasionally shift working set
      if (this.random() < 0.05) {
        const replaceIdx = Math.floor(this.random() * workingSet);
        workingSetKeys[replaceIdx] = Math.floor(this.random() * keyRange);
      }

      // 90% from working set, 10% random
      if (this.random() < 0.9) {
        requests.push(workingSetKeys[Math.floor(this.random() * workingSet)]);
      } else {
        requests.push(Math.floor(this.random() * keyRange));
      }
    }

    return requests;
  }

  private generateSequential(numRequests: number, keyRange: number): number[] {
    const requests: number[] = [];
    let current = Math.floor(this.random() * keyRange);

    for (let i = 0; i < numRequests; i++) {
      if (this.random() < 0.8) {
        // Sequential
        requests.push(current);
        current = (current + 1) % keyRange;
      } else {
        // Random jump
        current = Math.floor(this.random() * keyRange);
        requests.push(current);
      }
    }

    return requests;
  }

  private generateLooping(numRequests: number, keyRange: number): number[] {
    const loopSize = Math.max(1, Math.min(keyRange, Math.floor(keyRange / 2)));
    const requests: number[] = [];

    for (let i = 0; i < numRequests; i++) {
      requests.push(i % loopSize);
    }

    return requests;
  }
}

export interface BenchmarkResult {
  policy: string;
  cache_size: number;
  hit_ratio: number;
  miss_ratio: number;
  evictions: number;
  exec_time_ns: number;
}

export function runBenchmark(
  cacheSizes: number[],
  numRequests: number,
  keyRange: number,
  distribution: DistributionType
): BenchmarkResult[] {
  const gen = new TraceGenerator(42);
  const requests = gen.generate(numRequests, keyRange, distribution);
  const results: BenchmarkResult[] = [];

  for (const size of cacheSizes) {
    for (const policy of ['FIFO', 'LRU', 'LFU', 'Random']) {
      const sim = new Simulator(policy, size);
      const result = sim.run(requests);

      results.push({
        policy,
        cache_size: size,
        hit_ratio: result.metrics.cache_hits / result.metrics.total_requests,
        miss_ratio: result.metrics.cache_misses / result.metrics.total_requests,
        evictions: result.metrics.evictions,
        exec_time_ns: result.metrics.execution_time_nanos,
      });
    }
  }

  return results;
}
