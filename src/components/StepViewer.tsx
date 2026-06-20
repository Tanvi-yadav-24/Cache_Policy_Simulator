import React, { useState } from 'react';
import { ChevronLeft, ChevronRight, CheckCircle, XCircle, Trash2 } from 'lucide-react';
import { StepResult, POLICY_COLORS } from '../types';

interface StepViewerProps {
  steps: StepResult[];
  policy: string;
  visible?: boolean;
}

export function StepViewer({ steps, policy, visible = true }: StepViewerProps) {
  const [currentStep, setCurrentStep] = useState(0);

  if (!visible || steps.length === 0) {
    return null;
  }

  const step = steps[currentStep];

  return (
    <div className="card p-6 animate-fade-in">
      <div className="flex items-center justify-between mb-4">
        <h3 className="text-lg font-semibold text-gray-800 flex items-center gap-2">
          <div
            className="w-3 h-3 rounded-full"
            style={{ backgroundColor: POLICY_COLORS[policy] }}
          />
          Step-by-Step Execution
          <span className="badge badge-info">{policy}</span>
        </h3>
        <div className="flex items-center gap-2">
          <button
            onClick={() => setCurrentStep(Math.max(0, currentStep - 1))}
            disabled={currentStep === 0}
            className="p-2 rounded-lg hover:bg-gray-100 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            <ChevronLeft className="w-5 h-5" />
          </button>
          <span className="text-sm text-gray-600 min-w-[80px] text-center">
            {currentStep + 1} / {steps.length}
          </span>
          <button
            onClick={() => setCurrentStep(Math.min(steps.length - 1, currentStep + 1))}
            disabled={currentStep === steps.length - 1}
            className="p-2 rounded-lg hover:bg-gray-100 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            <ChevronRight className="w-5 h-5" />
          </button>
        </div>
      </div>

      <div className="bg-gray-50 rounded-lg p-4 space-y-4">
        <div className="grid grid-cols-2 gap-4">
          <div>
            <span className="text-sm font-medium text-gray-600">Request Value:</span>
            <div className="text-2xl font-bold text-blue-600 mt-1">{step.request}</div>
          </div>
          <div>
            <span className="text-sm font-medium text-gray-600">Result:</span>
            <div className="mt-1">
              {step.hit ? (
                <span className="badge badge-success flex items-center gap-1 w-fit">
                  <CheckCircle className="w-4 h-4" /> HIT
                </span>
              ) : (
                <span className="badge badge-error flex items-center gap-1 w-fit">
                  <XCircle className="w-4 h-4" /> MISS
                </span>
              )}
            </div>
          </div>
        </div>

        <div className="grid grid-cols-2 gap-4">
          <div>
            <span className="text-sm font-medium text-gray-600">Cache Before:</span>
            <div className="flex flex-wrap gap-1 mt-2">
              {step.cache_before.length === 0 ? (
                <span className="text-gray-400 text-sm">[empty]</span>
              ) : (
                step.cache_before.map((key) => (
                  <span
                    key={key}
                    className="inline-flex items-center px-2 py-1 rounded bg-gray-200 text-gray-800 text-sm font-mono"
                  >
                    {key}
                  </span>
                ))
              )}
            </div>
          </div>
          <div>
            <span className="text-sm font-medium text-gray-600">Cache After:</span>
            <div className="flex flex-wrap gap-1 mt-2">
              {step.cache_after.length === 0 ? (
                <span className="text-gray-400 text-sm">[empty]</span>
              ) : (
                step.cache_after.map((key) => (
                  <span
                    key={key}
                    className={`inline-flex items-center px-2 py-1 rounded text-sm font-mono ${
                      key === step.request && step.hit === false
                        ? 'bg-blue-100 text-blue-800 ring-2 ring-blue-400'
                        : 'bg-gray-200 text-gray-800'
                    }`}
                  >
                    {key}
                  </span>
                ))
              )}
            </div>
          </div>
        </div>

        {step.evicted !== 0 && (
          <div className="bg-red-50 border border-red-200 rounded-lg p-3 flex items-center gap-3">
            <Trash2 className="w-5 h-5 text-red-500" />
            <div>
              <span className="text-sm font-medium text-red-700">Evicted: </span>
              <span className="font-bold text-red-800">{step.evicted}</span>
            </div>
          </div>
        )}
      </div>

      {/* Progress bar */}
      <div className="mt-4">
        <input
          type="range"
          min={0}
          max={steps.length - 1}
          value={currentStep}
          onChange={(e) => setCurrentStep(parseInt(e.target.value, 10))}
          className="w-full h-2 bg-gray-200 rounded-lg appearance-none cursor-pointer"
        />
      </div>
    </div>
  );
}
