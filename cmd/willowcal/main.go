package main

import (
	"fmt"
	"log"
	"os"

	"github.com/devendershekhawat/teambiscuit/internal/config"
	"github.com/devendershekhawat/teambiscuit/internal/orchestrator"
	"github.com/devendershekhawat/teambiscuit/internal/reporter"
)

func main() {
    // Parse CLI args
    if len(os.Args) < 2 {
        log.Fatal("Usage: willowcal init <config.yaml>")
    }
    
    configPath := os.Args[1]
    
    // Parse config
    fmt.Println("📖 Parsing configuration...")
    cfg, err := config.ParseConfigFile(configPath)
    if err != nil {
        log.Fatalf("❌ Failed to parse config: %v", err)
    }
    
    fmt.Printf("✅ Config parsed successfully\n")
    fmt.Printf("   Workspace: %s\n", cfg.WorkspaceDir)
    fmt.Printf("   Repositories: %d\n\n", len(cfg.Repositories))
    
    // Get absolute workspace path
    workspaceDir, err := cfg.GetAbsoluteWorkspace()
    if err != nil {
        log.Fatalf("❌ Failed to resolve workspace: %v", err)
    }
    
    // Create workspace directory
    if err := os.MkdirAll(workspaceDir, 0755); err != nil {
        log.Fatalf("❌ Failed to create workspace: %v", err)
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
        os.Exit(1)
    }
}