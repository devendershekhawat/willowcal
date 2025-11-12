package config

import (
	"fmt"
	"strings"
	"testing"
)

func TestParseConfigValidSimple(t *testing.T) {
    yaml := `
version: "1.0"
workspace: "./workspace"
repositories:
  - name: backend
    url: https://github.com/test/backend.git
    path: ./backend
    setup_commands:
      - npm install
`
    
    config, err := ParseConfig([]byte(yaml))
    if err != nil {
        t.Fatalf("Expected no error, got: %v", err)
    }
    
    if config.Version != "1.0" {
        t.Errorf("Expected version 1.0, got: %s", config.Version)
    }
    
    if len(config.Repositories) != 1 {
        t.Errorf("Expected 1 repository, got: %d", len(config.Repositories))
    }
    
    if config.Repositories[0].Name != "backend" {
        t.Errorf("Expected repo name 'backend', got: %s", 
          config.Repositories[0].Name)
    }

		fmt.Println(config)
}

func TestParseConfigValidNoSetupCommands(t *testing.T) {
    yaml := `
version: "1.0"
workspace: "./workspace"
repositories:
  - name: docs
    url: https://github.com/test/docs.git
    path: ./docs
`
    
    config, err := ParseConfig([]byte(yaml))
    if err != nil {
        t.Fatalf("Expected no error, got: %v", err)
    }
    
    if len(config.Repositories[0].SetupCommands) != 0 {
        t.Errorf("Expected no setup commands, got: %d", 
          len(config.Repositories[0].SetupCommands))
    }

	fmt.Println(config)
}

func TestParseConfigInvalidVersion(t *testing.T) {
    yaml := `
version: "2.0"
workspace: "./workspace"
repositories: []
`
    
    _, err := ParseConfig([]byte(yaml))
    if err == nil {
        t.Fatal("Expected error for invalid version, got nil")
    }
    
    if !strings.Contains(err.Error(), "unsupported version") {
        t.Errorf("Expected 'unsupported version' error, got: %v", err)
    }
}

func TestParseConfigEmptyRepositories(t *testing.T) {
    yaml := `
version: "1.0"
workspace: "./workspace"
repositories: []
`
    
    _, err := ParseConfig([]byte(yaml))
    if err == nil {
        t.Fatal("Expected error for empty repositories, got nil")
    }
    
    if !strings.Contains(err.Error(), "at least one repository is required") {
        t.Errorf("Expected 'at least one repository' error, got: %v", err)
    }
}

func TestParseConfigDuplicateRepoNames(t *testing.T) {
    yaml := `
version: "1.0"
workspace: "./workspace"
repositories:
  - name: backend
    url: https://github.com/test/backend.git
    path: ./backend
  - name: backend
    url: https://github.com/test/backend2.git
    path: ./backend2
`
    
    _, err := ParseConfig([]byte(yaml))
    if err == nil {
        t.Fatal("Expected error for duplicate names, got nil")
    }
    
    if !strings.Contains(err.Error(), "duplicate repository name") {
        t.Errorf("Expected 'duplicate repository name' error, got: %v", err)
    }
}

func TestParseConfigMultipleErrors(t *testing.T) {
    yaml := `
version: "2.0"
workspace: ""
repositories:
  - name: ""
    url: "invalid-url"
    path: ""
`
    
    _, err := ParseConfig([]byte(yaml))
    if err == nil {
        t.Fatal("Expected multiple validation errors, got nil")
    }
    
    // Should contain multiple error messages
    errMsg := err.Error()
    if !strings.Contains(errMsg, "unsupported version") {
        t.Error("Expected version error in combined errors")
    }
    if !strings.Contains(errMsg, "workspace cannot be empty") {
        t.Error("Expected workspace error in combined errors")
    }
}