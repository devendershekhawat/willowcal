import { useState, useEffect } from 'react';
import { Terminal as TerminalIcon, Wifi, WifiOff, RefreshCw } from 'lucide-react';
import { AnimatePresence, motion } from 'framer-motion';
import { useWebSocket } from './hooks/useWebSocket';
import { ConfigEditor } from './components/ConfigEditor';
import { InitSection } from './components/InitSection';
import { ServiceCard } from './components/ServiceCard';
import { Terminal } from './components/Terminal';

function App() {
  const [showTerminal, setShowTerminal] = useState(false);
  const [selectedService, setSelectedService] = useState(null);

  const {
    isConnected,
    messages,
    services,
    config,
    uploadConfig,
    startInit,
    listServices,
    startService,
    stopService,
    getServiceStatus,
    clearMessages,
  } = useWebSocket('ws://localhost:8080/ws');

  // Poll service status every 3 seconds
  useEffect(() => {
    if (!isConnected || services.length === 0) return;

    const interval = setInterval(() => {
      getServiceStatus();
    }, 3000);

    return () => clearInterval(interval);
  }, [isConnected, services.length, getServiceStatus]);

  // Fetch services list when config is uploaded
  useEffect(() => {
    if (config && isConnected) {
      listServices();
    }
  }, [config, isConnected, listServices]);

  const handleStartService = (serviceName) => {
    startService(serviceName, (response) => {
      if (response.type === 'success') {
        setShowTerminal(true);
        getServiceStatus();
      }
    });
  };

  const handleStopService = (serviceName) => {
    stopService(serviceName, (response) => {
      if (response.type === 'success') {
        getServiceStatus();
      }
    });
  };

  const handleViewLogs = (serviceName) => {
    setSelectedService(serviceName);
    setShowTerminal(true);
  };

  const filteredMessages = selectedService
    ? messages.filter(m => !m.serviceName || m.serviceName === selectedService)
    : messages;

  return (
    <div className="min-h-screen bg-surface-base">
      {/* Header */}
      <header className="sticky top-0 z-40 backdrop-blur-xl bg-surface-base/80 border-b border-surface-border">
        <div className="max-w-7xl mx-auto px-6 py-4">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-3">
              <div className="w-10 h-10 bg-gradient-to-br from-primary-500 to-primary-700 rounded-xl flex items-center justify-center shadow-lg shadow-primary-500/20">
                <span className="text-xl font-bold text-white">W</span>
              </div>
              <div>
                <h1 className="text-2xl font-bold text-gradient">willowcal</h1>
                <p className="text-xs text-text-tertiary">Service Orchestration</p>
              </div>
            </div>

            <div className="flex items-center gap-4">
              <button
                onClick={() => setShowTerminal(!showTerminal)}
                className={`btn-ghost ${showTerminal ? 'bg-surface-hover' : ''}`}
              >
                <TerminalIcon className="w-4 h-4" />
                Terminal
                {messages.length > 0 && (
                  <span className="badge-info text-xs">
                    {messages.length}
                  </span>
                )}
              </button>

              <div className="flex items-center gap-2 px-3 py-1.5 rounded-lg bg-surface-elevated border border-surface-border">
                {isConnected ? (
                  <>
                    <Wifi className="w-4 h-4 text-green-400" />
                    <span className="text-sm text-green-400 font-medium">Connected</span>
                  </>
                ) : (
                  <>
                    <WifiOff className="w-4 h-4 text-red-400 animate-pulse" />
                    <span className="text-sm text-red-400 font-medium">Disconnected</span>
                  </>
                )}
              </div>
            </div>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-7xl mx-auto px-6 py-8 pb-32">
        <div className="space-y-6">
          {/* Config Editor */}
          <ConfigEditor
            onUpload={uploadConfig}
            config={config}
          />

          {/* Init Section */}
          <InitSection
            onStartInit={startInit}
            config={config}
            disabled={!isConnected}
          />

          {/* Services Section */}
          {services.length > 0 && (
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.3, delay: 0.2 }}
            >
              <div className="flex items-center justify-between mb-6">
                <div>
                  <h2 className="text-xl font-bold text-text-primary">Services</h2>
                  <p className="text-sm text-text-tertiary mt-1">
                    Manage your application services
                  </p>
                </div>
                <button
                  onClick={getServiceStatus}
                  className="btn-ghost"
                >
                  <RefreshCw className="w-4 h-4" />
                  Refresh
                </button>
              </div>

              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                {services.map((service) => (
                  <ServiceCard
                    key={service.name}
                    service={service}
                    onStart={handleStartService}
                    onStop={handleStopService}
                    onViewLogs={handleViewLogs}
                  />
                ))}
              </div>
            </motion.div>
          )}

          {/* Empty State */}
          {!config && (
            <motion.div
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              className="text-center py-16"
            >
              <div className="w-16 h-16 bg-surface-elevated rounded-full flex items-center justify-center mx-auto mb-4">
                <TerminalIcon className="w-8 h-8 text-text-tertiary" />
              </div>
              <h3 className="text-lg font-semibold text-text-primary mb-2">
                No Configuration Loaded
              </h3>
              <p className="text-text-tertiary">
                Upload a YAML configuration to get started
              </p>
            </motion.div>
          )}
        </div>
      </main>

      {/* Terminal Overlay */}
      <AnimatePresence>
        {showTerminal && (
          <Terminal
            messages={filteredMessages}
            onClear={clearMessages}
            onClose={() => {
              setShowTerminal(false);
              setSelectedService(null);
            }}
            title={selectedService ? `${selectedService} logs` : 'All Logs'}
          />
        )}
      </AnimatePresence>
    </div>
  );
}

export default App;
