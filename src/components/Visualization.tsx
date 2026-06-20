import React from 'react';
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  BarElement,
  Title,
  Tooltip,
  Legend,
  LineElement,
  PointElement,
} from 'chart.js';
import { Bar, Line } from 'react-chartjs-2';
import { SimulationResult, BenchmarkResult, POLICY_COLORS } from '../types';

ChartJS.register(CategoryScale, LinearScale, BarElement, LineElement, PointElement, Title, Tooltip, Legend);

interface VisualizationProps {
  results: Record<string, SimulationResult> | null;
  benchmarkResults: BenchmarkResult[] | null;
}

const chartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: {
      position: 'top' as const,
    },
  },
  scales: {
    y: {
      beginAtZero: true,
    },
  },
};

export function Visualization({ results, benchmarkResults }: VisualizationProps) {
  if (!results && !benchmarkResults) {
    return (
      <div className="card p-6 text-center text-gray-500">
        <p>Run a simulation to see visualizations</p>
      </div>
    );
  }

  if (results) {
    const policies = Object.keys(results);
    const labels = policies;

    const hitRatioData = {
      labels,
      datasets: [
        {
          label: 'Hit Ratio',
          data: policies.map((p) => results[p].metrics.cache_hits / results[p].metrics.total_requests),
          backgroundColor: policies.map((p) => POLICY_COLORS[p] + '80'),
          borderColor: policies.map((p) => POLICY_COLORS[p]),
          borderWidth: 2,
        },
      ],
    };

    const missRatioData = {
      labels,
      datasets: [
        {
          label: 'Miss Ratio',
          data: policies.map((p) => results[p].metrics.cache_misses / results[p].metrics.total_requests),
          backgroundColor: policies.map((p) => POLICY_COLORS[p] + '60'),
          borderColor: policies.map((p) => POLICY_COLORS[p]),
          borderWidth: 2,
        },
      ],
    };

    const evictionData = {
      labels,
      datasets: [
        {
          label: 'Evictions',
          data: policies.map((p) => results[p].metrics.evictions),
          backgroundColor: policies.map((p) => POLICY_COLORS[p] + '70'),
          borderColor: policies.map((p) => POLICY_COLORS[p]),
          borderWidth: 2,
        },
      ],
    };

    const hitsMissesData = {
      labels,
      datasets: [
        {
          label: 'Hits',
          data: policies.map((p) => results[p].metrics.cache_hits),
          backgroundColor: '#10b981' + '80',
          borderColor: '#10b981',
          borderWidth: 2,
        },
        {
          label: 'Misses',
          data: policies.map((p) => results[p].metrics.cache_misses),
          backgroundColor: '#ef4444' + '80',
          borderColor: '#ef4444',
          borderWidth: 2,
        },
      ],
    };

    return (
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6 animate-fade-in">
        <div className="card p-6">
          <h3 className="text-lg font-semibold text-gray-800 mb-4">Hit Ratio Comparison</h3>
          <div className="h-64">
            <Bar data={hitRatioData} options={{ ...chartOptions, indexAxis: 'y' }} />
          </div>
        </div>

        <div className="card p-6">
          <h3 className="text-lg font-semibold text-gray-800 mb-4">Miss Ratio Comparison</h3>
          <div className="h-64">
            <Bar data={missRatioData} options={{ ...chartOptions, indexAxis: 'y' }} />
          </div>
        </div>

        <div className="card p-6">
          <h3 className="text-lg font-semibold text-gray-800 mb-4">Evictions Count</h3>
          <div className="h-64">
            <Bar data={evictionData} options={chartOptions} />
          </div>
        </div>

        <div className="card p-6">
          <h3 className="text-lg font-semibold text-gray-800 mb-4">Hits vs Misses</h3>
          <div className="h-64">
            <Bar data={hitsMissesData} options={chartOptions} />
          </div>
        </div>
      </div>
    );
  }

  if (benchmarkResults) {
    const cacheSizes = [...new Set(benchmarkResults.map((r) => r.cache_size))].sort((a, b) => a - b);
    const policies = [...new Set(benchmarkResults.map((r) => r.policy))];

    const hitRatioLineData = {
      labels: cacheSizes.map(String),
      datasets: policies.map((policy) => ({
        label: policy,
        data: cacheSizes.map(
          (size) =>
            benchmarkResults.find((r) => r.policy === policy && r.cache_size === size)
              ?.hit_ratio ?? 0
        ),
        borderColor: POLICY_COLORS[policy],
        backgroundColor: POLICY_COLORS[policy] + '20',
        tension: 0.3,
        pointRadius: 4,
        pointHoverRadius: 6,
      })),
    };

    const execTimeData = {
      labels: cacheSizes.map(String),
      datasets: policies.map((policy) => ({
        label: policy,
        data: cacheSizes.map(
          (size) =>
            (benchmarkResults.find((r) => r.policy === policy && r.cache_size === size)
              ?.exec_time_ns ?? 0) / 1000 // Convert to microseconds
        ),
        borderColor: POLICY_COLORS[policy],
        backgroundColor: POLICY_COLORS[policy] + '20',
        tension: 0.3,
        pointRadius: 4,
        pointHoverRadius: 6,
      })),
    };

    return (
      <div className="space-y-6 animate-fade-in">
        <div className="card p-6">
          <h3 className="text-lg font-semibold text-gray-800 mb-4">Hit Ratio vs Cache Size</h3>
          <div className="h-80">
            <Line
              data={hitRatioLineData}
              options={{
                ...chartOptions,
                plugins: {
                  ...chartOptions.plugins,
                  tooltip: {
                    callbacks: {
                      label: (context) => `${context.dataset.label}: ${(context.parsed.y * 100).toFixed(2)}%`,
                    },
                  },
                },
                scales: {
                  ...chartOptions.scales,
                  y: {
                    ...chartOptions.scales.y,
                    ticks: {
                      callback: (value: number) => `${(value * 100).toFixed(0)}%`,
                    },
                    title: {
                      display: true,
                      text: 'Hit Ratio',
                    },
                  },
                  x: {
                    title: {
                      display: true,
                      text: 'Cache Size',
                    },
                  },
                },
              }}
            />
          </div>
        </div>

        <div className="card p-6">
          <h3 className="text-lg font-semibold text-gray-800 mb-4">Execution Time vs Cache Size</h3>
          <div className="h-80">
            <Line
              data={execTimeData}
              options={{
                ...chartOptions,
                scales: {
                  ...chartOptions.scales,
                  y: {
                    ...chartOptions.scales.y,
                    title: {
                      display: true,
                      text: 'Execution Time (microseconds)',
                    },
                  },
                  x: {
                    title: {
                      display: true,
                      text: 'Cache Size',
                    },
                  },
                },
              }}
            />
          </div>
        </div>
      </div>
    );
  }

  return null;
}
