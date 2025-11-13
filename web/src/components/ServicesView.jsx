import { useState } from 'react';
import { motion } from 'framer-motion';
import { RefreshCw } from 'lucide-react';
import { ServiceCard } from './ServiceCard';

export const ServicesView = ({ services, onStart, onStop, onViewLogs, onRefresh }) => {
  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold text-text-primary mb-2">Services</h2>
          <p className="text-text-tertiary">
            Manage and monitor your application services
          </p>
        </div>
        <button
          onClick={onRefresh}
          className="btn-ghost"
        >
          <RefreshCw className="w-4 h-4" />
          Refresh
        </button>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {services.map((service, index) => (
          <ServiceCard
            key={service.name}
            service={service}
            onStart={onStart}
            onStop={onStop}
            onViewLogs={onViewLogs}
            delay={index * 0.1}
          />
        ))}
      </div>

      {services.length === 0 && (
        <div className="text-center py-16 text-text-tertiary">
          <p>No services defined in configuration</p>
        </div>
      )}
    </div>
  );
};
