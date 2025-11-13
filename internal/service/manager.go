package service

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"time"

	"github.com/devendershekhawat/teambiscuit/internal/config"
	"github.com/devendershekhawat/teambiscuit/internal/models"
)

// ServiceState represents the current state of a service
type ServiceState string

const (
	StateStopped  ServiceState = "stopped"
	StateStarting ServiceState = "starting"
	StateRunning  ServiceState = "running"
	StateFailed   ServiceState = "failed"
)

// ServiceInstance represents a running service
type ServiceInstance struct {
	Name       string
	Service    models.Service
	State      ServiceState
	Process    *exec.Cmd
	StartTime  time.Time
	Error      string
	ctx        context.Context
	cancel     context.CancelFunc
	logChan    chan LogEntry
	mu         sync.RWMutex
}

// LogEntry represents a single log entry from a service
type LogEntry struct {
	Timestamp   time.Time
	ServiceName string
	Line        string
	Stream      string // "stdout" or "stderr"
}

// Manager manages all services
type Manager struct {
	config       *config.Config
	workspaceDir string
	services     map[string]*ServiceInstance
	mu           sync.RWMutex
	logBroadcast chan LogEntry
}

// NewManager creates a new service manager
func NewManager(cfg *config.Config, workspaceDir string) *Manager {
	return &Manager{
		config:       cfg,
		workspaceDir: workspaceDir,
		services:     make(map[string]*ServiceInstance),
		logBroadcast: make(chan LogEntry, 1000),
	}
}

// GetLogChannel returns the channel for receiving all service logs
func (m *Manager) GetLogChannel() <-chan LogEntry {
	return m.logBroadcast
}

// Start starts a specific service
func (m *Manager) Start(serviceName string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if service exists in config
	svc, err := m.config.GetServiceByName(serviceName)
	if err != nil {
		return fmt.Errorf("service not found: %s", serviceName)
	}

	// Check if already running
	if instance, exists := m.services[serviceName]; exists {
		if instance.State == StateRunning || instance.State == StateStarting {
			return fmt.Errorf("service already running")
		}
	}

	// Get repository path
	repo, err := m.config.GetRepositoryByName(svc.Repository)
	if err != nil {
		return fmt.Errorf("repository not found: %s", svc.Repository)
	}

	servicePath := repo.GetFullPath(m.workspaceDir)

	// Create context for cancellation
	ctx, cancel := context.WithCancel(context.Background())

	// Create service instance
	instance := &ServiceInstance{
		Name:      serviceName,
		Service:   *svc,
		State:     StateStarting,
		ctx:       ctx,
		cancel:    cancel,
		logChan:   make(chan LogEntry, 100),
		StartTime: time.Now(),
	}

	// Create command
	cmd := exec.CommandContext(ctx, "sh", "-c", svc.RunCommand)
	cmd.Dir = servicePath

	// Setup pipes
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		cancel()
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		cancel()
		return fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		cancel()
		return fmt.Errorf("failed to start service: %w", err)
	}

	instance.Process = cmd
	instance.State = StateRunning
	m.services[serviceName] = instance

	// Start log streaming goroutines
	go m.streamOutput(instance, stdout, "stdout")
	go m.streamOutput(instance, stderr, "stderr")

	// Monitor process
	go m.monitorProcess(instance)

	return nil
}

// Stop stops a specific service
func (m *Manager) Stop(serviceName string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	instance, exists := m.services[serviceName]
	if !exists {
		return fmt.Errorf("service not found or not running: %s", serviceName)
	}

	if instance.State != StateRunning {
		return fmt.Errorf("service not running: %s", serviceName)
	}

	// Cancel context to stop the process
	instance.cancel()

	// Try graceful shutdown first
	if instance.Process != nil && instance.Process.Process != nil {
		instance.Process.Process.Signal(syscall.SIGTERM)

		// Wait for graceful shutdown with timeout
		done := make(chan error, 1)
		go func() {
			done <- instance.Process.Wait()
		}()

		select {
		case <-time.After(5 * time.Second):
			// Force kill if not stopped gracefully
			instance.Process.Process.Kill()
		case <-done:
			// Process stopped gracefully
		}
	}

	instance.State = StateStopped
	return nil
}

// StopAll stops all running services
func (m *Manager) StopAll() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for name := range m.services {
		if instance := m.services[name]; instance.State == StateRunning {
			instance.cancel()
			if instance.Process != nil && instance.Process.Process != nil {
				instance.Process.Process.Signal(syscall.SIGTERM)
			}
		}
	}

	// Wait a bit for graceful shutdown
	time.Sleep(2 * time.Second)

	// Force kill any remaining
	for _, instance := range m.services {
		if instance.Process != nil && instance.Process.Process != nil {
			instance.Process.Process.Kill()
		}
		instance.State = StateStopped
	}
}

// GetStatus returns the status of a specific service
func (m *Manager) GetStatus(serviceName string) (*ServiceInstance, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	instance, exists := m.services[serviceName]
	if !exists {
		// Service exists in config but not started
		if _, err := m.config.GetServiceByName(serviceName); err == nil {
			return &ServiceInstance{
				Name:  serviceName,
				State: StateStopped,
			}, nil
		}
		return nil, fmt.Errorf("service not found: %s", serviceName)
	}

	// Return a copy to avoid race conditions
	return m.copyInstance(instance), nil
}

// GetAllStatuses returns status of all services
func (m *Manager) GetAllStatuses() []*ServiceInstance {
	m.mu.RLock()
	defer m.mu.RUnlock()

	statuses := make([]*ServiceInstance, 0, len(m.config.Services))

	for _, svc := range m.config.Services {
		if instance, exists := m.services[svc.Name]; exists {
			statuses = append(statuses, m.copyInstance(instance))
		} else {
			// Service exists in config but not started
			statuses = append(statuses, &ServiceInstance{
				Name:    svc.Name,
				Service: svc,
				State:   StateStopped,
			})
		}
	}

	return statuses
}

// UpdateConfig updates the manager's configuration
func (m *Manager) UpdateConfig(cfg *config.Config) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.config = cfg
}

// streamOutput reads from a pipe and broadcasts logs
func (m *Manager) streamOutput(instance *ServiceInstance, reader io.Reader, stream string) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		entry := LogEntry{
			Timestamp:   time.Now(),
			ServiceName: instance.Name,
			Line:        scanner.Text(),
			Stream:      stream,
		}

		// Send to instance channel
		select {
		case instance.logChan <- entry:
		default:
			// Channel full, skip
		}

		// Broadcast to global channel
		select {
		case m.logBroadcast <- entry:
		default:
			// Channel full, skip
		}
	}
}

// monitorProcess monitors the process and updates state
func (m *Manager) monitorProcess(instance *ServiceInstance) {
	err := instance.Process.Wait()

	m.mu.Lock()
	defer m.mu.Unlock()

	if err != nil && instance.ctx.Err() == nil {
		// Process failed (not cancelled by us)
		instance.State = StateFailed
		instance.Error = err.Error()
	} else {
		instance.State = StateStopped
	}

	// Close log channel
	close(instance.logChan)
}

// copyInstance creates a copy of a service instance for safe reading
func (m *Manager) copyInstance(instance *ServiceInstance) *ServiceInstance {
	instance.mu.RLock()
	defer instance.mu.RUnlock()

	var pid int
	if instance.Process != nil && instance.Process.Process != nil {
		pid = instance.Process.Process.Pid
	}

	return &ServiceInstance{
		Name:      instance.Name,
		Service:   instance.Service,
		State:     instance.State,
		StartTime: instance.StartTime,
		Error:     instance.Error,
		Process: &exec.Cmd{
			Process: &os.Process{Pid: pid},
		},
	}
}
