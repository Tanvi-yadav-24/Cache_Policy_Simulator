import React, { useState } from 'react';
import { Play, Shuffle, BarChart3, Loader2 } from 'lucide-react';
import { DistributionType } from '../types';

interface SimulatorFormProps {
  mode: 'single' | 'compare' | 'benchmark';
  onSimulate: (policy: string, cacheSize: number, requests: number[]) => void;
  onCompare: (cacheSize: number, requests: number[]) => void;
  onGenerate: (
    numRequests: number,
    keyRange: number,
    distribution: DistributionType
  ) => Promise<number[]>;
  onBenchmark: (
    cacheSizes: number[],
    numRequests: number,
    keyRange: number,
    distribution: DistributionType
  ) => void;
  loading: boolean;
}

export function SimulatorForm({
  mode,
  onSimulate,
  onCompare,
  onGenerate,
  onBenchmark,
  loading,
}: SimulatorFormProps) {
  const [policy, setPolicy] = useState<string>('FIFO');
  const [cacheSize, setCacheSize] = useState<number>(3);
  const [requestsInput, setRequestsInput] = useState<string>('1, 2, 3, 1, 4, 5, 2, 1, 6, 7');
  const [numRequests, setNumRequests] = useState<number>(10000);
  const [keyRange, setKeyRange] = useState<number>(10000);
  const [distribution, setDistribution] = useState<DistributionType>('zipf');
  const [benchmarkSizes, setBenchmarkSizes] = useState<string>('10, 50, 100, 500, 1000');

  const policies = ['FIFO', 'LRU', 'LFU', 'Random'];

  const parseRequests = (input: string): number[] => {
    return input
      .split(/[,\s]+/)
      .map((s) => parseInt(s.trim(), 10))
      .filter((n) => !isNaN(n));
  };

  const handleGenerate = async () => {
    const reqs = await onGenerate(numRequests, keyRange, distribution);
    setRequestsInput(reqs.slice(0, 100).join(', ') + (reqs.length > 100 ? '...' : ''));
  };

  const handleSubmit = () => {
    const requests = parseRequests(requestsInput);

    if (mode === 'single') {
      onSimulate(policy, cacheSize, requests);
    } else if (mode === 'compare') {
      onCompare(cacheSize, requests);
    } else if (mode === 'benchmark') {
      const sizes = benchmarkSizes
        .split(',')
        .map((s) => parseInt(s.trim(), 10))
        .filter((n) => !isNaN(n));
      onBenchmark(sizes, numRequests, keyRange, distribution);
    }
  };

  return (
    <div className="card p-6 space-y-6">
      <h2 className="text-xl font-bold text-gray-800 flex items-center gap-2">
        <BarChart3 className="w-6 h-6 text-blue-600" />
        Simulator Configuration
      </h2>

      {mode === 'single' && (
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">Cache Policy</label>
          <select
            value={policy}
            onChange={(e) => setPolicy(e.target.value)}
            className="select-field"
          >
            {policies.map((p) => (
              <option key={p} value={p}>
                {p}
              </option>
            ))}
          </select>
          <p className="mt-2 text-sm text-gray-500">
            {policy === 'FIFO' && 'First In First Out - removes oldest element'}
            {policy === 'LRU' && 'Least Recently Used - removes least recently accessed'}
            {policy === 'LFU' && 'Least Frequently Used - removes least frequently used'}
            {policy === 'Random' && 'Random Replacement - randomly evicts an item'}
          </p>
        </div>
      )}

      {mode !== 'benchmark' && (
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">Cache Size</label>
          <input
            type="number"
            value={cacheSize}
            onChange={(e) => setCacheSize(parseInt(e.target.value, 10) || 1)}
            min={1}
            className="input-field"
          />
        </div>
      )}

      {mode !== 'benchmark' && (
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">
            Request Sequence
          </label>
          <textarea
            value={requestsInput}
            onChange={(e) => setRequestsInput(e.target.value)}
            placeholder="e.g., 1, 2, 3, 1, 4, 5"
            className="input-field h-24 font-mono text-sm"
          />
          <button
            onClick={handleGenerate}
            disabled={loading}
            className="mt-2 btn-secondary flex items-center gap-2 text-sm"
          >
            <Shuffle className="w-4 h-4" />
            Generate Random Trace
          </button>
        </div>
      )}

      {mode === 'benchmark' && (
        <>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Number of Requests
              </label>
              <input
                type="number"
                value={numRequests}
                onChange={(e) => setNumRequests(parseInt(e.target.value, 10) || 1000)}
                min={100}
                className="input-field"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">Key Range</label>
              <input
                type="number"
                value={keyRange}
                onChange={(e) => setKeyRange(parseInt(e.target.value, 10) || 100)}
                min={10}
                className="input-field"
              />
            </div>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Cache Sizes (comma-separated)
            </label>
            <input
              type="text"
              value={benchmarkSizes}
              onChange={(e) => setBenchmarkSizes(e.target.value)}
              placeholder="e.g., 10, 50, 100, 500"
              className="input-field"
            />
          </div>
        </>
      )}

      <div>
        <label className="block text-sm font-medium text-gray-700 mb-2">Distribution</label>
        <select
          value={distribution}
          onChange={(e) => setDistribution(e.target.value as DistributionType)}
          className="select-field"
        >
          <option value="uniform">Uniform - Equal probability for all keys</option>
          <option value="zipf">Zipf - Power law (some keys more popular)</option>
          <option value="localized">Localized - Working set pattern</option>
          <option value="sequential">Sequential - Mostly sequential access</option>
          <option value="looping">Looping - Repeated loop pattern</option>
        </select>
      </div>

      <button
        onClick={handleSubmit}
        disabled={loading}
        className="btn-primary w-full flex items-center justify-center gap-2"
      >
        {loading ? (
          <>
            <Loader2 className="w-5 h-5 animate-spin" />
            Running...
          </>
        ) : (
          <>
            <Play className="w-5 h-5" />
            {mode === 'single' && 'Run Simulation'}
            {mode === 'compare' && 'Compare All Policies'}
            {mode === 'benchmark' && 'Run Benchmark'}
          </>
        )}
      </button>
    </div>
  );
}
