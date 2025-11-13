package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/devendershekhawat/teambiscuit/internal/models"
)

type Config struct {
	Version      string               `yaml:"version"`
	WorkspaceDir string               `yaml:"workspace_dir"`
	Repositories []models.Repository  `yaml:"repositories"`
	Services     []models.Service     `yaml:"services"`
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

func (c *Config) GetServiceByName(name string) (*models.Service, error) {
    for i := range c.Services {
        if c.Services[i].Name == name {
            return &c.Services[i], nil
        }
    }
    return nil, fmt.Errorf("service not found: %s", name)
}

func (c *Config) ServiceCount() int {
    return len(c.Services)
}

// ValidateServices checks that all services reference valid repositories
func (c *Config) ValidateServices() error {
    for _, service := range c.Services {
        if err := service.Validate(); err != nil {
            return err
        }

        // Check that referenced repository exists
        if _, err := c.GetRepositoryByName(service.Repository); err != nil {
            return fmt.Errorf("service '%s' references non-existent repository '%s'",
                service.Name, service.Repository)
        }
    }
    return nil
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
	out += fmt.Sprintf("Services: %d\n", c.ServiceCount())
	for _, service := range c.Services {
		out += fmt.Sprintf("  - Name: %s\n", service.Name)
		out += fmt.Sprintf("    Repository: %s\n", service.Repository)
		out += fmt.Sprintf("    RunCommand: %s\n", service.RunCommand)
	}
	return out
}