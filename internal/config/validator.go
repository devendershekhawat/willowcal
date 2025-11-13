package config

import (
	"fmt"
	"strings"
)

// ValidateConfig validates the entire configuration
func ValidateConfig(config *Config) error {
    var errors []string
    
    // Validate version
    if config.Version != "1.0" {
        errors = append(errors, 
          fmt.Sprintf("unsupported version: %s (expected: 1.0)", 
            config.Version))
    }
    
    // Validate workspace
    if config.WorkspaceDir == "" {
        errors = append(errors, "workspace cannot be empty")
    }
    
    // Validate repositories
    if len(config.Repositories) == 0 {
        errors = append(errors, "at least one repository is required")
    }
    
    // Check for duplicate repository names
    repoNames := make(map[string]bool)
    for _, repo := range config.Repositories {
        if repoNames[repo.Name] {
            errors = append(errors,
              fmt.Sprintf("duplicate repository name: %s", repo.Name))
        }
        repoNames[repo.Name] = true

        // Validate individual repository
        if err := repo.Validate(); err != nil {
            errors = append(errors, err.Error())
        }
    }

    // Validate services (if any)
    if len(config.Services) > 0 {
        // Check for duplicate service names
        serviceNames := make(map[string]bool)
        for _, service := range config.Services {
            if serviceNames[service.Name] {
                errors = append(errors,
                  fmt.Sprintf("duplicate service name: %s", service.Name))
            }
            serviceNames[service.Name] = true

            // Validate individual service
            if err := service.Validate(); err != nil {
                errors = append(errors, err.Error())
            }
        }

        // Validate that services reference valid repositories
        if err := config.ValidateServices(); err != nil {
            errors = append(errors, err.Error())
        }
    }

    // Return all errors at once (not fail-fast)
    if len(errors) > 0 {
        return fmt.Errorf("config validation failed:\n  - %s",
          strings.Join(errors, "\n  - "))
    }

    return nil
}