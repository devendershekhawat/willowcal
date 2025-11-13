import { useState, useEffect } from 'react';
import { Terminal as TerminalIcon, Wifi, WifiOff, FileCode, FolderGit2, Layers } from 'lucide-react';
import { AnimatePresence, motion } from 'framer-motion';
import { useWebSocket } from './hooks/useWebSocket';
import { ConfigEditor } from './components/ConfigEditor';
import { InitSection } from './components/InitSection';
import { RepositoriesView } from './components/RepositoriesView';
import { ServicesView } from './components/ServicesView';
import { Terminal } from './components/Terminal';

function App() {
  const [showTerminal, setShowTerminal] = useState(false);
  const [selectedService, setSelectedService] = useState(null);
  const [activeTab, setActiveTab] = useState('setup');
  const [initResult, setInitResult] = useState(null);

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

  // Poll service status every 3 seconds when on services tab
  useEffect(() => {
    if (!isConnected || services.length === 0 || activeTab !== 'services') return;

    const interval = setInterval(() => {
      getServiceStatus();
    }, 3000);

    return () => clearInterval(interval);
  }, [isConnected, services.length, activeTab, getServiceStatus]);

  // Fetch services list when config is uploaded
  useEffect(() => {
    if (config && isConnected) {
      listServices();
    }
  }, [config, isConnected, listServices]);

  // Auto-switch to repositories tab after successful init
  useEffect(() => {
    const initCompleteMsg = messages.find(m => m.type === 'init-complete');
    if (initCompleteMsg && initCompleteMsg.payload) {
      setInitResult(initCompleteMsg.payload);
      // Auto-switch to repositories tab after init
      setTimeout(() => {
        setActiveTab('repositories');
      }, 1000);
    }
  }, [messages]);

  const handleStartInit = (callback) => {
    setInitResult(null);
    setShowTerminal(true);
    startInit(callback);
  };

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

  // Build repositories data from config and init result
  const repositories = config?.repositories?.map(repo => {
    const result = initResult?.repositories?.find(r => r.name === repo.name);
    return {
      ...repo,
      status: result?.status || 'pending',
      duration_seconds: result?.duration_seconds,
      error: result?.error,
    };
  }) || [];

  const tabs = [
    { id: 'setup', label: 'Setup', icon: FileCode },
    { id: 'repositories', label: 'Repositories', icon: FolderGit2, disabled: !config },
    { id: 'services', label: 'Services', icon: Layers, disabled: !config || services.length === 0 },
  ];

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

          {/* Tabs */}
          <div className="mt-6 flex items-center gap-2 border-b border-surface-border">
            {tabs.map((tab) => {
              const Icon = tab.icon;
              const isActive = activeTab === tab.id;
              const isDisabled = tab.disabled;

              return (
                <button
                  key={tab.id}
                  onClick={() => !isDisabled && setActiveTab(tab.id)}
                  disabled={isDisabled}
                  className={`relative px-4 py-3 flex items-center gap-2 font-medium transition-colors ${
                    isActive
                      ? 'text-primary-400'
                      : isDisabled
                      ? 'text-text-tertiary opacity-50 cursor-not-allowed'
                      : 'text-text-secondary hover:text-text-primary'
                  }`}
                >
                  <Icon className="w-4 h-4" />
                  {tab.label}

                  {/* Active indicator */}
                  {isActive && (
                    <motion.div
                      layoutId="activeTab"
                      className="absolute bottom-0 left-0 right-0 h-0.5 bg-primary-500"
                      transition={{ type: "spring", stiffness: 500, damping: 30 }}
                    />
                  )}
                </button>
              );
            })}
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-7xl mx-auto px-6 py-8 pb-32">
        <AnimatePresence mode="wait">
          {activeTab === 'setup' && (
            <motion.div
              key="setup"
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0, y: -20 }}
              transition={{ duration: 0.3 }}
              className="space-y-6"
            >
              <ConfigEditor onUpload={uploadConfig} config={config} />
              <InitSection onStartInit={handleStartInit} config={config} disabled={!isConnected} />
            </motion.div>
          )}

          {activeTab === 'repositories' && (
            <motion.div
              key="repositories"
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0, y: -20 }}
              transition={{ duration: 0.3 }}
            >
              <RepositoriesView repositories={repositories} />
            </motion.div>
          )}

          {activeTab === 'services' && (
            <motion.div
              key="services"
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0, y: -20 }}
              transition={{ duration: 0.3 }}
            >
              <ServicesView
                services={services}
                onStart={handleStartService}
                onStop={handleStopService}
                onViewLogs={handleViewLogs}
                onRefresh={getServiceStatus}
              />
            </motion.div>
          )}
        </AnimatePresence>
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
