import { createCache, CacheInterface } from './cache';

export interface StepResult {
  request: number;
  cache_before: number[];
  cache_after: number[];
  evicted: number;
  hit: boolean;
}

export interface Metrics {
  total_requests: number;
  cache_hits: number;
  cache_misses: number;
  evictions: number;
  execution_time_nanos: number;
  memory_bytes: number;
}

export interface SimulationResult {
  policy: string;
  cache_size: number;
  steps: StepResult[];
  metrics: Metrics;
  request_trace: number[];
}

export class Simulator {
  private cache: CacheInterface;
  private policyName: string;
  private cap: number;

  constructor(policy: string, capacity: number) {
    this.policyName = policy;
    this.cap = capacity;
    this.cache = createCache(policy, capacity);
  }

  run(requests: number[]): SimulationResult {
    this.cache.clear();

    const steps: StepResult[] = [];
    let hits = 0;
    let misses = 0;
    let evictionCount = 0;

    const startTime = performance.now();

    for (const request of requests) {
      const cacheBefore = this.cache.keys();

      // Try to get from cache
      const result = this.cache.get(request);
      const hit = result !== undefined;

      let evicted = 0;

      if (!hit) {
        // Cache miss - add to cache
        const evictedKey = this.cache.put(request, request);
        if (evictedKey !== null) {
          evicted = evictedKey;
          evictionCount++;
        }
        misses++;
      } else {
        hits++;
      }

      const cacheAfter = this.cache.keys();

      steps.push({
        request,
        cache_before: cacheBefore,
        cache_after: cacheAfter,
        evicted,
        hit,
      });
    }

    const endTime = performance.now();
    const execTimeNanos = (endTime - startTime) * 1_000_000;

    return {
      policy: this.policyName,
      cache_size: this.cap,
      steps,
      metrics: {
        total_requests: requests.length,
        cache_hits: hits,
        cache_misses: misses,
        evictions: evictionCount,
        execution_time_nanos: execTimeNanos,
        memory_bytes: this.estimateMemory(),
      },
      request_trace: requests,
    };
  }

  private estimateMemory(): number {
    const baseOverhead = 100;
    const entrySize = 16;
    const pointerOverhead = 8;

    switch (this.policyName) {
      case 'FIFO':
        return baseOverhead + this.cap * entrySize + this.cap * pointerOverhead;
      case 'LRU':
        return baseOverhead + this.cap * entrySize + this.cap * pointerOverhead * 2;
      case 'LFU':
        return baseOverhead + this.cap * entrySize * 2 + this.cap * pointerOverhead * 2;
      case 'Random':
        return baseOverhead + this.cap * entrySize + this.cap * pointerOverhead;
      default:
        return baseOverhead + this.cap * entrySize;
    }
  }
}

export function runComparison(cacheSize: number, requests: number[]): Record<string, SimulationResult> {
  const policies = ['FIFO', 'LRU', 'LFU', 'Random'];
  const results: Record<string, SimulationResult> = {};

  for (const policy of policies) {
    const sim = new Simulator(policy, cacheSize);
    results[policy] = sim.run(requests);
  }

  return results;
}
