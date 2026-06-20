import React, { useState } from 'react';
import { Header } from './components/Header';
import { SimulatorForm } from './components/SimulatorForm';
import { ResultsTable } from './components/ResultsTable';
import { Visualization } from './components/Visualization';
import { StepViewer } from './components/StepViewer';
import { PoliciesInfo } from './components/PoliciesInfo';
import { api } from './services/api';
import {
  SimulationResult,
  BenchmarkResult,
  DistributionType,
} from './types';
import { runBenchmark } from './lib/generator';

type Mode = 'single' | 'compare' | 'benchmark';

function App() {
  const [mode, setMode] = useState<Mode>('compare');
  const [results, setResults] = useState<Record<string, SimulationResult> | null>(null);
  const [benchmarkResults, setBenchmarkResults] = useState<BenchmarkResult[] | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleSimulate = async (
    policy: string,
    cacheSize: number,
    requests: number[]
  ) => {
    setLoading(true);
    setError(null);
    setBenchmarkResults(null);
    try {
      const result = await api.simulate(policy, cacheSize, requests, true);
      setResults({ [policy]: result });
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Simulation failed');
    } finally {
      setLoading(false);
    }
  };

  const handleCompare = async (cacheSize: number, requests: number[]) => {
    setLoading(true);
    setError(null);
    setBenchmarkResults(null);
    try {
      const comparisonResults = await api.compare(cacheSize, requests);
      setResults(comparisonResults);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Comparison failed');
    } finally {
      setLoading(false);
    }
  };

  const handleGenerate = async (
    numRequests: number,
    keyRange: number,
    distribution: DistributionType
  ): Promise<number[]> => {
    try {
      const response = await api.generateTrace(numRequests, keyRange, distribution);
      return response.requests;
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Generation failed');
      return [];
    }
  };

  const handleBenchmark = (
    cacheSizes: number[],
    numRequests: number,
    keyRange: number,
    distribution: DistributionType
  ) => {
    setLoading(true);
    setError(null);
    setResults(null);
    try {
      const benchResults = runBenchmark(cacheSizes, numRequests, keyRange, distribution);
      setBenchmarkResults(benchResults);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Benchmark failed');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-100 via-gray-50 to-slate-100">
      <Header />

      <main className="max-w-7xl mx-auto px-6 py-8">
        <div className="bg-green-50 border border-green-200 rounded-lg p-4 mb-6 flex items-start gap-3">
          <div className="text-green-600 font-medium">Ready:</div>
          <div className="text-green-800">
            Simulator is running locally with pure TypeScript implementations. All cache policies (FIFO, LRU, LFU, Random) work without a backend server.
          </div>
        </div>

        {error && (
          <div className="bg-red-50 border border-red-200 rounded-lg p-4 mb-6 flex items-start gap-3">
            <div className="text-red-600 font-medium">Error:</div>
            <div className="text-red-800">{error}</div>
          </div>
        )}

        <div className="mb-6">
          <div className="flex gap-2 justify-center">
            <button
              onClick={() => setMode('single')}
              className={`px-6 py-2.5 rounded-lg text-sm font-medium transition-all ${
                mode === 'single'
                  ? 'bg-blue-600 text-white shadow-md'
                  : 'bg-white text-gray-600 hover:bg-gray-50 border border-gray-200'
              }`}
            >
              Single
            </button>
            <button
              onClick={() => setMode('compare')}
              className={`px-6 py-2.5 rounded-lg text-sm font-medium transition-all ${
                mode === 'compare'
                  ? 'bg-blue-600 text-white shadow-md'
                  : 'bg-white text-gray-600 hover:bg-gray-50 border border-gray-200'
              }`}
            >
              Compare
            </button>
            <button
              onClick={() => setMode('benchmark')}
              className={`px-6 py-2.5 rounded-lg text-sm font-medium transition-all ${
                mode === 'benchmark'
                  ? 'bg-blue-600 text-white shadow-md'
                  : 'bg-white text-gray-600 hover:bg-gray-50 border border-gray-200'
              }`}
            >
              Benchmark
            </button>
          </div>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <div className="lg:col-span-1 space-y-6">
            <SimulatorForm
              mode={mode}
              onSimulate={handleSimulate}
              onCompare={handleCompare}
              onGenerate={handleGenerate}
              onBenchmark={handleBenchmark}
              loading={loading}
            />
            <PoliciesInfo />
          </div>

          <div className="lg:col-span-2 space-y-6">
            {results && <ResultsTable results={results} />}

            {results && Object.values(results)[0]?.steps && (
              <StepViewer
                steps={Object.values(results)[0].steps}
                policy={Object.keys(results)[0]}
                visible={true}
              />
            )}

            {(results || benchmarkResults) && (
              <Visualization results={results} benchmarkResults={benchmarkResults} />
            )}
          </div>
        </div>
      </main>

      <footer className="border-t border-gray-200 bg-white mt-12">
        <div className="max-w-7xl mx-auto px-6 py-6 text-center text-gray-500 text-sm">
          <p>
            Cache Policy Simulator - Compare FIFO, LRU, LFU, and Random cache replacement
            algorithms
          </p>
          <p className="mt-1">
            Pure TypeScript implementation with O(1) Get/Put operations | HashMap + Doubly Linked
            List
          </p>
        </div>
      </footer>
    </div>
  );
}

export default App;
