import { useState } from 'react';
import { Rocket, Loader2, CheckCircle2 } from 'lucide-react';
import { motion } from 'framer-motion';

export const InitSection = ({ onStartInit, config, disabled }) => {
  const [isInitializing, setIsInitializing] = useState(false);
  const [initResult, setInitResult] = useState(null);

  const handleInit = () => {
    setIsInitializing(true);
    setInitResult(null);

    onStartInit((response) => {
      if (response.type === 'init.complete') {
        setIsInitializing(false);
        setInitResult(response.payload);
      }
    });
  };

  return (
    <motion.div
      className="card"
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.3, delay: 0.1 }}
    >
      <div className="flex items-start justify-between">
        <div className="flex-1">
          <h2 className="text-xl font-bold text-text-primary flex items-center gap-2 mb-2">
            <Rocket className="w-5 h-5 text-primary-500" />
            Initialize Repositories
          </h2>
          <p className="text-sm text-text-tertiary mb-4">
            Clone repositories and run setup commands defined in your configuration
          </p>

          {!config && (
            <div className="p-3 bg-yellow-500/10 border border-yellow-500/20 rounded-lg text-yellow-400 text-sm mb-4">
              Please upload a configuration first
            </div>
          )}

          {initResult && (
            <motion.div
              initial={{ opacity: 0, scale: 0.95 }}
              animate={{ opacity: 1, scale: 1 }}
              className="p-4 bg-green-500/10 border border-green-500/20 rounded-lg mb-4"
            >
              <div className="flex items-center gap-2 text-green-400 font-medium mb-2">
                <CheckCircle2 className="w-5 h-5" />
                Initialization Complete
              </div>
              <div className="grid grid-cols-3 gap-4 text-sm">
                <div>
                  <div className="text-text-tertiary">Success</div>
                  <div className="text-green-400 font-semibold text-lg">
                    {initResult.success}
                  </div>
                </div>
                <div>
                  <div className="text-text-tertiary">Failed</div>
                  <div className="text-red-400 font-semibold text-lg">
                    {initResult.failed}
                  </div>
                </div>
                <div>
                  <div className="text-text-tertiary">Duration</div>
                  <div className="text-text-primary font-semibold text-lg">
                    {initResult.total_time_seconds?.toFixed(1)}s
                  </div>
                </div>
              </div>
            </motion.div>
          )}

          <button
            onClick={handleInit}
            disabled={disabled || !config || isInitializing}
            className="btn-primary"
          >
            {isInitializing ? (
              <>
                <Loader2 className="w-4 h-4 animate-spin" />
                Initializing...
              </>
            ) : (
              <>
                <Rocket className="w-4 h-4" />
                Start Initialization
              </>
            )}
          </button>
        </div>
      </div>
    </motion.div>
  );
};
