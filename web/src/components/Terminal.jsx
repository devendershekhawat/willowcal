import { useEffect, useRef } from 'react';
import { X, Trash2 } from 'lucide-react';
import { motion, AnimatePresence } from 'framer-motion';

export const Terminal = ({ messages, onClear, onClose, title = 'Terminal' }) => {
  const terminalRef = useRef(null);

  useEffect(() => {
    if (terminalRef.current) {
      terminalRef.current.scrollTop = terminalRef.current.scrollHeight;
    }
  }, [messages]);

  const getMessageColor = (type) => {
    switch (type) {
      case 'error':
        return 'text-red-400';
      case 'system':
        return 'text-blue-400';
      case 'init-progress':
        return 'text-yellow-400';
      case 'init-complete':
        return 'text-green-400';
      case 'service-log':
        return 'text-text-primary';
      default:
        return 'text-text-secondary';
    }
  };

  const formatTimestamp = (timestamp) => {
    if (!(timestamp instanceof Date)) {
      timestamp = new Date(timestamp);
    }
    return timestamp.toLocaleTimeString('en-US', {
      hour12: false,
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit'
    });
  };

  return (
    <motion.div
      className="fixed inset-x-0 bottom-0 z-50 max-h-[60vh] flex flex-col bg-surface-elevated border-t border-surface-border shadow-2xl"
      initial={{ y: '100%' }}
      animate={{ y: 0 }}
      exit={{ y: '100%' }}
      transition={{ type: 'spring', damping: 25, stiffness: 300 }}
    >
      {/* Header */}
      <div className="flex items-center justify-between px-4 py-3 border-b border-surface-border bg-surface-base">
        <div className="flex items-center gap-3">
          <div className="flex items-center gap-2">
            <div className="w-3 h-3 rounded-full bg-red-500" />
            <div className="w-3 h-3 rounded-full bg-yellow-500" />
            <div className="w-3 h-3 rounded-full bg-green-500" />
          </div>
          <h3 className="text-sm font-semibold text-text-primary">{title}</h3>
          <span className="text-xs text-text-tertiary">
            {messages.length} {messages.length === 1 ? 'line' : 'lines'}
          </span>
        </div>

        <div className="flex items-center gap-2">
          <button
            onClick={onClear}
            className="p-2 hover:bg-surface-hover rounded-lg transition-colors"
            title="Clear terminal"
          >
            <Trash2 className="w-4 h-4 text-text-tertiary" />
          </button>
          <button
            onClick={onClose}
            className="p-2 hover:bg-surface-hover rounded-lg transition-colors"
            title="Close terminal"
          >
            <X className="w-4 h-4 text-text-tertiary" />
          </button>
        </div>
      </div>

      {/* Terminal Content */}
      <div
        ref={terminalRef}
        className="flex-1 overflow-y-auto p-4 bg-black/50 font-mono text-sm"
      >
        {messages.length === 0 ? (
          <div className="text-text-tertiary italic">No messages yet...</div>
        ) : (
          <AnimatePresence initial={false}>
            {messages.map((message, index) => (
              <motion.div
                key={index}
                initial={{ opacity: 0, x: -10 }}
                animate={{ opacity: 1, x: 0 }}
                transition={{ duration: 0.2 }}
                className={`mb-1 ${getMessageColor(message.type)}`}
              >
                <span className="text-text-tertiary">
                  [{formatTimestamp(message.timestamp)}]
                </span>
                {message.serviceName && (
                  <span className="text-primary-400 ml-2">
                    [{message.serviceName}]
                  </span>
                )}
                {message.repoName && (
                  <span className="text-yellow-400 ml-2">
                    [{message.repoName}]
                  </span>
                )}
                <span className="ml-2">{message.text}</span>
              </motion.div>
            ))}
          </AnimatePresence>
        )}
      </div>
    </motion.div>
  );
};
