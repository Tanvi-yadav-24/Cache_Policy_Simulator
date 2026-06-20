import React from 'react';
import { Server, Database, Cpu, Activity } from 'lucide-react';

export function Header() {
  return (
    <header className="bg-gradient-to-r from-slate-900 via-blue-900 to-slate-900 text-white py-8 px-6 shadow-xl">
      <div className="max-w-7xl mx-auto">
        <div className="flex items-center gap-4 mb-2">
          <div className="p-3 bg-blue-500 rounded-xl shadow-lg">
            <Database className="w-8 h-8" />
          </div>
          <div>
            <h1 className="text-3xl font-bold tracking-tight">Cache Policy Simulator</h1>
            <p className="text-blue-200 mt-1">
              Compare cache replacement algorithms: FIFO, LRU, LFU, and Random
            </p>
          </div>
        </div>
        <div className="flex gap-6 mt-6 text-sm">
          <div className="flex items-center gap-2">
            <Server className="w-4 h-4 text-blue-400" />
            <span className="text-slate-300">Go Backend</span>
          </div>
          <div className="flex items-center gap-2">
            <Activity className="w-4 h-4 text-green-400" />
            <span className="text-slate-300">Real-time Simulation</span>
          </div>
          <div className="flex items-center gap-2">
            <Cpu className="w-4 h-4 text-amber-400" />
            <span className="text-slate-300">Performance Benchmarking</span>
          </div>
        </div>
      </div>
    </header>
  );
}
