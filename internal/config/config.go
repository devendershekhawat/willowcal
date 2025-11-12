package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/devendershekhawat/teambiscuit/internal/models"
)

type Config struct {
	Version string `yaml:"version"`
	WorkspaceDir string `yaml:"workspace_dir"`
	Repositories []models.Repository `yaml:"repositories"`
}

func (c *Config) GetRepositoryByName(name string) (*models.Repository, error) {
    for i := range c.Repositories {
        if c.Repositories[i].Name == name {
            return &c.Repositories[i], nil
        }
    }
    return nil, fmt.Errorf("repository not found: %s", name)
}

func (c *Config) GetAbsoluteWorkspace() (string, error) {
    if filepath.IsAbs(c.WorkspaceDir) {
        return c.WorkspaceDir, nil
    }
    
    cwd, err := os.Getwd()
    if err != nil {
        return "", fmt.Errorf("failed to get working directory: %w", err)
    }
    
    return filepath.Join(cwd, c.WorkspaceDir), nil
}

func (c *Config) RepositoryCount() int {
    return len(c.Repositories)
}

func (c *Config) String() string {
	out := ""
	out += fmt.Sprintf("Version: %s\n", c.Version)
	out += fmt.Sprintf("WorkspaceDir: %s\n", c.WorkspaceDir)
	out += fmt.Sprintf("Repositories: %d\n", c.RepositoryCount())
	for _, repo := range c.Repositories {
		out += fmt.Sprintf("  - Name: %s\n", repo.Name)
		out += fmt.Sprintf("    URL: %s\n", repo.URL)
		out += fmt.Sprintf("    Path: %s\n", repo.Path)
		out += fmt.Sprintf("    SetupCommands: %v\n", repo.SetupCommands)
	}
	return out
}