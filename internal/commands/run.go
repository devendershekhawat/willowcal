package commands

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/devendershekhawat/teambiscuit/internal/config"
	"github.com/devendershekhawat/teambiscuit/internal/models"
	"github.com/devendershekhawat/teambiscuit/internal/orchestrator"
	"github.com/devendershekhawat/teambiscuit/internal/reporter"
	"github.com/devendershekhawat/teambiscuit/internal/runner"
)

// RunCommand handles the 'run' command
func RunCommand(configPath string) error {
	// Parse config
	fmt.Println("📖 Parsing configuration...")
	cfg, err := config.ParseConfigFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	fmt.Printf("✅ Config parsed successfully\n")
	fmt.Printf("   Workspace: %s\n", cfg.WorkspaceDir)
	fmt.Printf("   Repositories: %d\n", len(cfg.Repositories))
	fmt.Printf("   Services: %d\n\n", len(cfg.Services))

	// Validate services
	if len(cfg.Services) == 0 {
		return fmt.Errorf("no services defined in config")
	}

	// Get absolute workspace path
	workspaceDir, err := cfg.GetAbsoluteWorkspace()
	if err != nil {
		return fmt.Errorf("failed to resolve workspace: %w", err)
	}

	// Create workspace directory if it doesn't exist
	if err := os.MkdirAll(workspaceDir, 0755); err != nil {
		return fmt.Errorf("failed to create workspace: %w", err)
	}

	// Check which repositories need to be cloned
	fmt.Println("🔍 Checking repository status...")
	missingRepos, err := checkMissingRepositories(cfg, workspaceDir)
	if err != nil {
		return fmt.Errorf("failed to check repositories: %w", err)
	}

	// Clone missing repositories
	if len(missingRepos) > 0 {
		fmt.Printf("\n📦 Found %d missing repositories, cloning...\n", len(missingRepos))
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

		if err := cloneMissingRepositories(missingRepos, cfg, workspaceDir); err != nil {
			return fmt.Errorf("failed to clone repositories: %w", err)
		}

		fmt.Println("\n✅ All repositories ready")
	} else {
		fmt.Println("✅ All repositories already cloned")
	}

	// Run services
	fmt.Println()
	serviceRunner := runner.NewServiceRunner(cfg, workspaceDir)
	if err := serviceRunner.Run(); err != nil {
		return fmt.Errorf("failed to run services: %w", err)
	}

	return nil
}

// checkMissingRepositories returns a list of repositories that need to be cloned
// for the services defined in the config
func checkMissingRepositories(cfg *config.Config, workspaceDir string) ([]models.Repository, error) {
	missingRepos := make([]models.Repository, 0)
	checkedRepos := make(map[string]bool)

	// Check each service's repository
	for _, service := range cfg.Services {
		// Skip if we've already checked this repo
		if checkedRepos[service.Repository] {
			continue
		}
		checkedRepos[service.Repository] = true

		// Get repository config
		repo, err := cfg.GetRepositoryByName(service.Repository)
		if err != nil {
			return nil, err
		}

		// Check if repository is cloned
		repoPath := repo.GetFullPath(workspaceDir)
		if !isRepositoryCloned(repoPath) {
			fmt.Printf("   ⚠️  Repository '%s' not found at %s\n", repo.Name, repoPath)
			missingRepos = append(missingRepos, *repo)
		} else {
			fmt.Printf("   ✅ Repository '%s' found\n", repo.Name)
		}
	}

	return missingRepos, nil
}

// isRepositoryCloned checks if a git repository exists at the given path
func isRepositoryCloned(path string) bool {
	gitDir := filepath.Join(path, ".git")
	info, err := os.Stat(gitDir)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// cloneMissingRepositories clones the missing repositories and runs setup commands
func cloneMissingRepositories(repos []models.Repository, cfg *config.Config, workspaceDir string) error {
	// Use orchestrator for parallel cloning
	tempConfig := &config.Config{
		Version:      cfg.Version,
		WorkspaceDir: cfg.WorkspaceDir,
		Repositories: repos,
	}

	orch := orchestrator.NewOrchestrator(tempConfig, workspaceDir)
	state := orch.Execute()

	// Print summary
	reporter.PrintFinalSummary(state)

	// Check for failures
	if state.FailureCount > 0 {
		log.Fatal("❌ Failed to clone all required repositories")
	}

	return nil
}
