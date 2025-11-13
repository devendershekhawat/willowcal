package runner

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/devendershekhawat/teambiscuit/internal/config"
	"github.com/devendershekhawat/teambiscuit/internal/models"
)

const (
	// Color codes for different services
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[37m"
)

var colors = []string{
	colorGreen,
	colorYellow,
	colorBlue,
	colorPurple,
	colorCyan,
	colorWhite,
}

type ServiceRunner struct {
	config       *config.Config
	workspaceDir string
	processes    []*exec.Cmd
	mu           sync.Mutex
}

func NewServiceRunner(cfg *config.Config, workspaceDir string) *ServiceRunner {
	return &ServiceRunner{
		config:       cfg,
		workspaceDir: workspaceDir,
		processes:    make([]*exec.Cmd, 0),
	}
}

// Run starts all services and waits for them to complete or for Ctrl+C
func (sr *ServiceRunner) Run() error {
	if len(sr.config.Services) == 0 {
		return fmt.Errorf("no services defined in config")
	}

	fmt.Println("🚀 Starting services...")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup signal handling for Ctrl+C
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Start all services
	var wg sync.WaitGroup
	errChan := make(chan error, len(sr.config.Services))

	for i, service := range sr.config.Services {
		wg.Add(1)
		color := colors[i%len(colors)]
		go func(svc models.Service, serviceColor string) {
			defer wg.Done()
			if err := sr.runService(ctx, svc, serviceColor); err != nil {
				errChan <- fmt.Errorf("service '%s' failed: %w", svc.Name, err)
			}
		}(service, color)
	}

	// Wait for either completion, error, or Ctrl+C
	go func() {
		wg.Wait()
		cancel()
	}()

	select {
	case <-ctx.Done():
		// Services completed
	case err := <-errChan:
		fmt.Printf("\n❌ %v\n", err)
		cancel()
	case <-sigChan:
		fmt.Println("\n\n⚠️  Received interrupt signal, shutting down services...")
		cancel()
	}

	// Give processes time to shut down gracefully
	sr.stopAllServices()

	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("✅ All services stopped")

	return nil
}

// runService runs a single service and streams its logs with a prefix
func (sr *ServiceRunner) runService(ctx context.Context, service models.Service, color string) error {
	// Get repository to find the path
	repo, err := sr.config.GetRepositoryByName(service.Repository)
	if err != nil {
		return err
	}

	servicePath := repo.GetFullPath(sr.workspaceDir)

	// Log service start
	prefix := fmt.Sprintf("%s[%s]%s", color, service.Name, colorReset)
	fmt.Printf("%s Starting service: %s\n", prefix, service.RunCommand)

	// Create command - use shell to support complex commands
	cmd := exec.CommandContext(ctx, "sh", "-c", service.RunCommand)
	cmd.Dir = servicePath

	// Get pipes for stdout and stderr
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start command: %w", err)
	}

	// Track the process
	sr.mu.Lock()
	sr.processes = append(sr.processes, cmd)
	sr.mu.Unlock()

	// Stream output with service prefix
	var wg sync.WaitGroup
	wg.Add(2)

	go sr.streamOutput(&wg, stdout, prefix, service.Name)
	go sr.streamOutput(&wg, stderr, prefix, service.Name)

	// Wait for output to finish
	wg.Wait()

	// Wait for command to finish
	if err := cmd.Wait(); err != nil {
		if ctx.Err() == context.Canceled {
			// Context was cancelled (Ctrl+C or other service failed)
			return nil
		}
		return fmt.Errorf("command exited with error: %w", err)
	}

	fmt.Printf("%s Service exited\n", prefix)
	return nil
}

// streamOutput reads from a pipe and writes to stdout with a prefix
func (sr *ServiceRunner) streamOutput(wg *sync.WaitGroup, reader io.Reader, prefix, serviceName string) {
	defer wg.Done()

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		timestamp := time.Now().Format("15:04:05")
		fmt.Printf("%s [%s] %s\n", prefix, timestamp, scanner.Text())
	}
}

// stopAllServices sends kill signals to all running processes
func (sr *ServiceRunner) stopAllServices() {
	sr.mu.Lock()
	defer sr.mu.Unlock()

	for _, cmd := range sr.processes {
		if cmd.Process != nil {
			// Try graceful shutdown first
			cmd.Process.Signal(syscall.SIGTERM)
		}
	}

	// Wait a bit for graceful shutdown
	time.Sleep(2 * time.Second)

	// Force kill if still running
	for _, cmd := range sr.processes {
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
	}
}
