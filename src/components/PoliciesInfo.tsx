import React from 'react';
import { Clock, TrendingDown, BarChart2, Shuffle } from 'lucide-react';
import { POLICY_COLORS } from '../types';

const policyIcons: Record<string, React.ReactNode> = {
  FIFO: <Clock className="w-5 h-5" />,
  LRU: <TrendingDown className="w-5 h-5" />,
  LFU: <BarChart2 className="w-5 h-5" />,
  Random: <Shuffle className="w-5 h-5" />,
};

const policyDescriptions: Record<string, { title: string; desc: string; complexity: string }> = {
  FIFO: {
    title: 'First In First Out',
    desc: 'Removes the oldest element when cache is full. Simple and fast, but ignores access patterns.',
    complexity: 'Get: O(1), Put: O(1)',
  },
  LRU: {
    title: 'Least Recently Used',
    desc: 'Removes the least recently accessed item. Adapts to access patterns, good for temporal locality.',
    complexity: 'Get: O(1), Put: O(1)',
  },
  LFU: {
    title: 'Least Frequently Used',
    desc: 'Removes the least frequently used item (tie-breaker: oldest). Good for steady access patterns.',
    complexity: 'Get: O(1), Put: O(1)',
  },
  Random: {
    title: 'Random Replacement',
    desc: 'Randomly evicts an item when cache is full. No overhead for tracking, but unpredictable.',
    complexity: 'Get: O(1), Put: O(1)',
  },
};

export function PoliciesInfo() {
  return (
    <div className="card p-6">
      <h3 className="text-lg font-semibold text-gray-800 mb-4">Cache Policies</h3>
      <div className="space-y-4">
        {Object.entries(policyDescriptions).map(([key, info]) => (
          <div
            key={key}
            className="flex gap-4 p-3 rounded-lg hover:bg-gray-50 transition-colors"
          >
            <div
              className="flex-shrink-0 w-10 h-10 rounded-lg flex items-center justify-center text-white"
              style={{ backgroundColor: POLICY_COLORS[key] }}
            >
              {policyIcons[key]}
            </div>
            <div className="flex-1">
              <div className="flex items-center gap-2">
                <span className="font-semibold text-gray-900">{key}</span>
                <span className="text-xs text-gray-500 bg-gray-100 px-2 py-0.5 rounded">
                  {info.complexity}
                </span>
              </div>
              <p className="text-sm text-gray-600 mt-1">{info.desc}</p>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
