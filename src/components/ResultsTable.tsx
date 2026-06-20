import React from 'react';
import { CheckCircle, XCircle, Trash2, TrendingUp, TrendingDown } from 'lucide-react';
import { SimulationResult, POLICY_COLORS } from '../types';

interface ResultsTableProps {
  results: Record<string, SimulationResult> | null;
}

export function ResultsTable({ results }: ResultsTableProps) {
  if (!results || Object.keys(results).length === 0) {
    return (
      <div className="card p-6 text-center text-gray-500">
        <p>Run a simulation to see results</p>
      </div>
    );
  }

  const policies = Object.keys(results);
  const bestHitRatio = Math.max(
    ...policies.map((p) => results[p].metrics.cache_hits / results[p].metrics.total_requests)
  );
  const leastEvictions = Math.min(...policies.map((p) => results[p].metrics.evictions));

  return (
    <div className="card overflow-hidden animate-fade-in">
      <div className="bg-gradient-to-r from-slate-800 to-slate-900 px-6 py-4">
        <h3 className="text-lg font-semibold text-white">Simulation Results</h3>
        <p className="text-slate-400 text-sm">Performance comparison across all policies</p>
      </div>
      <div className="overflow-x-auto">
        <table className="w-full">
          <thead className="bg-gray-50 border-b border-gray-200">
            <tr>
              <th className="px-6 py-3 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider">
                Policy
              </th>
              <th className="px-6 py-3 text-center text-xs font-semibold text-gray-600 uppercase tracking-wider">
                Hits
              </th>
              <th className="px-6 py-3 text-center text-xs font-semibold text-gray-600 uppercase tracking-wider">
                Misses
              </th>
              <th className="px-6 py-3 text-center text-xs font-semibold text-gray-600 uppercase tracking-wider">
                Hit Ratio
              </th>
              <th className="px-6 py-3 text-center text-xs font-semibold text-gray-600 uppercase tracking-wider">
                Miss Ratio
              </th>
              <th className="px-6 py-3 text-center text-xs font-semibold text-gray-600 uppercase tracking-wider">
                Evictions
              </th>
              <th className="px-6 py-3 text-center text-xs font-semibold text-gray-600 uppercase tracking-wider">
                Exec Time
              </th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-200">
            {policies.map((policy) => {
              const result = results[policy];
              const hitRatio = result.metrics.cache_hits / result.metrics.total_requests;
              const missRatio = result.metrics.cache_misses / result.metrics.total_requests;
              const isBestHit = hitRatio === bestHitRatio;
              const isLeastEvict = result.metrics.evictions === leastEvictions;

              return (
                <tr key={policy} className="hover:bg-gray-50 transition-colors">
                  <td className="px-6 py-4 whitespace-nowrap">
                    <div className="flex items-center gap-3">
                      <div
                        className="w-3 h-3 rounded-full"
                        style={{ backgroundColor: POLICY_COLORS[policy] }}
                      />
                      <span className="font-semibold text-gray-900">{policy}</span>
                    </div>
                  </td>
                  <td className="px-6 py-4 text-center">
                    <div className="flex items-center justify-center gap-1">
                      <CheckCircle className="w-4 h-4 text-green-500" />
                      <span className="text-green-600 font-medium">
                        {result.metrics.cache_hits}
                      </span>
                    </div>
                  </td>
                  <td className="px-6 py-4 text-center">
                    <div className="flex items-center justify-center gap-1">
                      <XCircle className="w-4 h-4 text-red-500" />
                      <span className="text-red-600 font-medium">
                        {result.metrics.cache_misses}
                      </span>
                    </div>
                  </td>
                  <td className="px-6 py-4 text-center">
                    <div className="flex items-center justify-center gap-2">
                      <span
                        className={`text-lg font-bold ${
                          isBestHit ? 'text-green-600' : 'text-gray-700'
                        }`}
                      >
                        {(hitRatio * 100).toFixed(2)}%
                      </span>
                      {isBestHit && <TrendingUp className="w-4 h-4 text-green-500" />}
                    </div>
                  </td>
                  <td className="px-6 py-4 text-center">
                    <span className="text-red-500 font-medium">
                      {(missRatio * 100).toFixed(2)}%
                    </span>
                  </td>
                  <td className="px-6 py-4 text-center">
                    <div className="flex items-center justify-center gap-2">
                      <Trash2 className="w-4 h-4 text-gray-400" />
                      <span className={isLeastEvict ? 'text-green-600 font-medium' : ''}>
                        {result.metrics.evictions}
                      </span>
                    </div>
                  </td>
                  <td className="px-6 py-4 text-center">
                    <span className="text-gray-600 text-sm">
                      {(result.metrics.execution_time_nanos / 1000).toFixed(2)} μs
                    </span>
                  </td>
                </tr>
              );
            })}
          </tbody>
        </table>
      </div>
    </div>
  );
}
