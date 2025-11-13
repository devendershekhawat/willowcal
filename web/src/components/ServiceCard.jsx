import { useState } from 'react';
import { Play, Square, Terminal, Activity } from 'lucide-react';
import { motion } from 'framer-motion';

export const ServiceCard = ({ service, onStart, onStop, onViewLogs, delay = 0 }) => {
  const [isHovered, setIsHovered] = useState(false);

  const getStatusColor = (status) => {
    switch (status) {
      case 'running':
        return 'text-green-400 bg-green-500/10 border-green-500/20';
      case 'starting':
        return 'text-yellow-400 bg-yellow-500/10 border-yellow-500/20';
      case 'stopped':
        return 'text-gray-400 bg-gray-500/10 border-gray-500/20';
      case 'failed':
        return 'text-red-400 bg-red-500/10 border-red-500/20';
      default:
        return 'text-gray-400 bg-gray-500/10 border-gray-500/20';
    }
  };

  const getStatusIcon = (status) => {
    if (status === 'running') {
      return <Activity className="w-3 h-3 animate-pulse" />;
    }
    return null;
  };

  const isRunning = service.status === 'running';
  const isStopped = service.status === 'stopped';

  return (
    <motion.div
      className={`card relative overflow-hidden transition-all duration-300 ${
        isHovered ? 'shadow-lg shadow-primary-500/10 border-primary-500/30' : ''
      }`}
      onMouseEnter={() => setIsHovered(true)}
      onMouseLeave={() => setIsHovered(false)}
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.3, delay }}
    >
      {/* Glow effect on hover */}
      {isHovered && (
        <div className="absolute inset-0 bg-gradient-to-r from-primary-500/5 to-transparent pointer-events-none" />
      )}

      <div className="relative">
        {/* Header */}
        <div className="flex items-start justify-between mb-4">
          <div className="flex-1">
            <h3 className="text-lg font-semibold text-text-primary mb-1">
              {service.name}
            </h3>
            <p className="text-sm text-text-tertiary font-mono">
              {service.repository}
            </p>
          </div>

          <div className={`badge ${getStatusColor(service.status)} flex items-center gap-1.5`}>
            {getStatusIcon(service.status)}
            <span className="capitalize">{service.status || 'stopped'}</span>
          </div>
        </div>

        {/* Command */}
        <div className="mb-4 p-3 bg-surface-base rounded-lg border border-surface-border">
          <p className="text-xs text-text-tertiary mb-1">Run Command</p>
          <p className="text-sm text-text-secondary font-mono break-all">
            {service.run_command}
          </p>
        </div>

        {/* Process Info */}
        {service.pid && isRunning && (
          <div className="mb-4 flex items-center gap-4 text-xs text-text-tertiary">
            <div className="flex items-center gap-2">
              <span className="text-text-secondary">PID:</span>
              <span className="font-mono text-text-primary">{service.pid}</span>
            </div>
          </div>
        )}

        {/* Actions */}
        <div className="flex items-center gap-2">
          {isStopped ? (
            <button
              onClick={() => onStart(service.name)}
              className="btn-primary flex-1"
            >
              <Play className="w-4 h-4" />
              Start Service
            </button>
          ) : (
            <button
              onClick={() => onStop(service.name)}
              className="btn-secondary flex-1"
              disabled={!isRunning}
            >
              <Square className="w-4 h-4" />
              Stop Service
            </button>
          )}

          <button
            onClick={() => onViewLogs(service.name)}
            className="btn-ghost"
          >
            <Terminal className="w-4 h-4" />
          </button>
        </div>
      </div>
    </motion.div>
  );
};
