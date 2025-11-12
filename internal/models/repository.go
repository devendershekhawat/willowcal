package models

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

type Repository struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
	Path string `yaml:"path"`
	SetupCommands []string `yaml:"setup_commands"`
}

func (r *Repository) GetFullPath(workspaceDir string) string {
	return filepath.Join(workspaceDir, r.Path)
}

func (r *Repository) HasSetupCommands() bool {
	return len(r.SetupCommands) > 0
}

func (r *Repository) Validate() error {
	if r.Name == "" {
        return fmt.Errorf("repository name cannot be empty")
    }
    
    if r.URL == "" {
        return fmt.Errorf("repository '%s' must have a URL", r.Name)
    }
    
    if r.Path == "" {
        return fmt.Errorf("repository '%s' must have a path", r.Name)
    }
    
    // Validate URL format
    if !isValidGitURL(r.URL) {
        return fmt.Errorf("repository '%s' has invalid git URL: %s", 
          r.Name, r.URL)
    }
    
    return nil
}

func isValidGitURL(url string) bool {
    return strings.HasPrefix(url, "http://") ||
           strings.HasPrefix(url, "https://") ||
           strings.HasPrefix(url, "git@") ||
           strings.HasPrefix(url, "ssh://")
}

func (r *Repository) Ping() error {
	gitCmd := exec.Command("git", "ls-remote", r.URL)
	if err := gitCmd.Run(); err != nil {
		return fmt.Errorf("failed to ping repository: %w", err)
	}
	return nil
}