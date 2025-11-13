import { useState, useEffect, useRef, useCallback } from 'react';

let messageId = 0;

export const useWebSocket = (url) => {
  const [isConnected, setIsConnected] = useState(false);
  const [messages, setMessages] = useState([]);
  const [services, setServices] = useState([]);
  const [config, setConfig] = useState(null);
  const ws = useRef(null);
  const reconnectTimeout = useRef(null);
  const messageHandlers = useRef(new Map());

  const connect = useCallback(() => {
    try {
      ws.current = new WebSocket(url);

      ws.current.onopen = () => {
        console.log('✅ Connected to willowcal');
        setIsConnected(true);
        setMessages(prev => [...prev, { type: 'system', text: 'Connected to willowcal server', timestamp: new Date() }]);
      };

      ws.current.onclose = () => {
        console.log('❌ Disconnected from willowcal');
        setIsConnected(false);
        setMessages(prev => [...prev, { type: 'system', text: 'Disconnected from server', timestamp: new Date() }]);

        // Attempt reconnection after 3 seconds
        reconnectTimeout.current = setTimeout(() => {
          console.log('🔄 Attempting to reconnect...');
          connect();
        }, 3000);
      };

      ws.current.onerror = (error) => {
        console.error('WebSocket error:', error);
        setMessages(prev => [...prev, { type: 'error', text: 'Connection error occurred', timestamp: new Date() }]);
      };

      ws.current.onmessage = (event) => {
        try {
          const message = JSON.parse(event.data);
          console.log('📥 Received:', message);

          // Call registered handler for this message ID
          if (message.id && messageHandlers.current.has(message.id)) {
            const handler = messageHandlers.current.get(message.id);
            handler(message);
            messageHandlers.current.delete(message.id);
          }

          // Handle specific message types
          switch (message.type) {
            case 'service.log':
              setMessages(prev => [...prev, {
                type: 'service-log',
                serviceName: message.payload.service_name,
                text: message.payload.line,
                stream: message.payload.stream,
                timestamp: new Date(message.payload.timestamp),
              }]);
              break;

            case 'service.started':
              setServices(prev => prev.map(s =>
                s.name === message.payload.service_name
                  ? { ...s, status: 'running' }
                  : s
              ));
              break;

            case 'service.stopped':
              setServices(prev => prev.map(s =>
                s.name === message.payload.service_name
                  ? { ...s, status: 'stopped' }
                  : s
              ));
              break;

            case 'init.progress':
              setMessages(prev => [...prev, {
                type: 'init-progress',
                repoName: message.payload.repo_name,
                text: message.payload.message,
                timestamp: new Date(),
              }]);
              break;

            case 'init.complete':
              setMessages(prev => [...prev, {
                type: 'init-complete',
                text: `Init complete: ${message.payload.success} succeeded, ${message.payload.failed} failed`,
                payload: message.payload,
                timestamp: new Date(),
              }]);
              break;

            case 'error':
              setMessages(prev => [...prev, {
                type: 'error',
                text: message.payload.message,
                timestamp: new Date(),
              }]);
              break;
          }
        } catch (error) {
          console.error('Error parsing message:', error);
        }
      };
    } catch (error) {
      console.error('Failed to connect:', error);
    }
  }, [url]);

  useEffect(() => {
    connect();

    return () => {
      if (reconnectTimeout.current) {
        clearTimeout(reconnectTimeout.current);
      }
      if (ws.current) {
        ws.current.close();
      }
    };
  }, [connect]);

  const sendMessage = useCallback((type, payload, onResponse) => {
    if (!ws.current || ws.current.readyState !== WebSocket.OPEN) {
      console.error('WebSocket is not connected');
      return;
    }

    const id = `req-${++messageId}`;
    const message = { type, id, payload };

    // Register response handler
    if (onResponse) {
      messageHandlers.current.set(id, onResponse);
    }

    console.log('📤 Sending:', message);
    ws.current.send(JSON.stringify(message));

    return id;
  }, []);

  const uploadConfig = useCallback((configYaml, onResponse) => {
    sendMessage('config.upload', { config_yaml: configYaml }, (response) => {
      if (response.type === 'success' && response.payload.valid) {
        setConfig(response.payload);
      }
      if (onResponse) onResponse(response);
    });
  }, [sendMessage]);

  const startInit = useCallback((onResponse) => {
    sendMessage('init.start', {}, onResponse);
  }, [sendMessage]);

  const listServices = useCallback((onResponse) => {
    sendMessage('service.list', {}, (response) => {
      if (response.type === 'success' && response.payload.services) {
        setServices(response.payload.services);
      }
      if (onResponse) onResponse(response);
    });
  }, [sendMessage]);

  const startService = useCallback((serviceName, onResponse) => {
    sendMessage('service.start', { service_name: serviceName }, onResponse);
  }, [sendMessage]);

  const stopService = useCallback((serviceName, onResponse) => {
    sendMessage('service.stop', { service_name: serviceName }, onResponse);
  }, [sendMessage]);

  const getServiceStatus = useCallback((onResponse) => {
    sendMessage('service.status', {}, (response) => {
      if (response.type === 'success' && response.payload.services) {
        setServices(prev => {
          const statusMap = new Map(response.payload.services.map(s => [s.name, s]));
          return prev.map(service => ({
            ...service,
            ...statusMap.get(service.name),
          }));
        });
      }
      if (onResponse) onResponse(response);
    });
  }, [sendMessage]);

  const clearMessages = useCallback(() => {
    setMessages([]);
  }, []);

  return {
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
  };
};
