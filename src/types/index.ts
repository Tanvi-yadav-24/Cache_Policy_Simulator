import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  BarElement,
  LineElement,
  PointElement,
  Title,
  Tooltip,
  Legend,
  Filler,
  ArcElement,
} from 'chart.js';

ChartJS.register(
  CategoryScale,
  LinearScale,
  BarElement,
  LineElement,
  PointElement,
  Title,
  Tooltip,
  Legend,
  Filler,
  ArcElement
);

// Re-export from lib for backward compatibility
export type { StepResult, SimulationResult, Metrics } from '../lib/simulator';
export type { BenchmarkResult, DistributionType } from '../lib/generator';

export const API_BASE = '/api';

export const POLICY_COLORS: Record<string, string> = {
  FIFO: '#3b82f6',    // Blue
  LRU: '#10b981',    // Green
  LFU: '#f59e0b',    // Amber
  Random: '#8b5cf6', // Purple
};

export const POLICY_DESCRIPTIONS: Record<string, string> = {
  FIFO: 'First In First Out - Removes the oldest element when cache is full',
  LRU: 'Least Recently Used - Removes the least recently accessed item',
  LFU: 'Least Frequently Used - Removes the least frequently used item (tie-breaker: oldest)',
  Random: 'Random Replacement - Randomly evicts an item when cache is full',
};
