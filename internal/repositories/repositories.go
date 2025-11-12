package repositories

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type Repository struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
	Path string `yaml:"path"`
	SetupCommands []string `yaml:"setup_commands"`
	cloned bool
}

type RepositoryService struct {
	BaseDir string;
}

func NewRepositoryService(basePath string) (*RepositoryService, error) {
	// Check if the directory exists
	_, err := os.Stat(basePath)
	if os.IsNotExist(err) {
		// Directory doesn't exist, create it
		if err := os.MkdirAll(basePath, 0755); err != nil {
			return nil, fmt.Errorf("failed to create base directory: %w", err)
		}
	} else if err != nil {
		// Some other error occurred (permissions, etc.)
		return nil, fmt.Errorf("failed to access base directory: %w", err)
	}
	// If err == nil, directory exists, we're good!

	return &RepositoryService{
		BaseDir: basePath,
	}, nil
}


func (r *RepositoryService) CloneAndSetupRepository(repo Repository) error {
	if repo.cloned {
		return nil
	}

	// Clone the repository
	repoPath := filepath.Join(r.BaseDir, repo.Name)
	if err := os.MkdirAll(repoPath, 0755); err != nil {
		return fmt.Errorf("failed to create repository directory: %w", err)
	}

	// Clone the repository using git
	gitCmd := exec.Command("git", "clone", repo.URL, repoPath)
	if err := gitCmd.Run(); err != nil {
		return fmt.Errorf("failed to clone repository: %w", err)
	}

	repo.cloned = true
	return nil
}