package commands

import (
	"fmt"
	"log"
	"os"

	"github.com/devendershekhawat/teambiscuit/internal/config"
	"github.com/devendershekhawat/teambiscuit/internal/orchestrator"
	"github.com/devendershekhawat/teambiscuit/internal/reporter"
)

// InitCommand handles the 'init' command
func InitCommand(configPath string) error {
	// Parse config
	fmt.Println("📖 Parsing configuration...")
	cfg, err := config.ParseConfigFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	fmt.Printf("✅ Config parsed successfully\n")
	fmt.Printf("   Workspace: %s\n", cfg.WorkspaceDir)
	fmt.Printf("   Repositories: %d\n\n", len(cfg.Repositories))

	// Get absolute workspace path
	workspaceDir, err := cfg.GetAbsoluteWorkspace()
	if err != nil {
		return fmt.Errorf("failed to resolve workspace: %w", err)
	}

	// Create workspace directory
	if err := os.MkdirAll(workspaceDir, 0755); err != nil {
		return fmt.Errorf("failed to create workspace: %w", err)
	}

	// Execute
	fmt.Println("🚀 Starting parallel initialization...")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	orch := orchestrator.NewOrchestrator(cfg, workspaceDir)
	state := orch.Execute()

	// Print summary
	reporter.PrintFinalSummary(state)

	// Exit with appropriate code
	if state.FailureCount > 0 {
		log.Fatal("❌ Initialization failed")
	}

	return nil
}
