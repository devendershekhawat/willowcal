package commands

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/devendershekhawat/teambiscuit/internal/api"
)

// ServerCommand starts the WebSocket server
func ServerCommand(port string, workspaceDir string, staticDir string) error {
	if workspaceDir == "" {
		workspaceDir = "./workspace"
	}

	// Create workspace directory
	if err := os.MkdirAll(workspaceDir, 0755); err != nil {
		return fmt.Errorf("failed to create workspace: %w", err)
	}

	// Create handler
	handler := api.NewHandler(workspaceDir)

	// Create server
	addr := fmt.Sprintf(":%s", port)
	server := api.NewServer(addr, handler)

	// Set broadcaster
	handler.SetBroadcaster(server.GetBroadcaster())

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("\n⚠️  Shutting down server...")
		os.Exit(0)
	}()

	// Start server
	log.Printf("📡 Starting willowcal server on port %s", port)
	log.Printf("📂 Workspace directory: %s", workspaceDir)
	if staticDir != "" {
		log.Printf("🌐 Web UI: http://localhost:%s", port)
	}
	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	if err := server.Start(staticDir); err != nil {
		return fmt.Errorf("server error: %w", err)
	}

	return nil
}
