package git

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/devendershekhawat/teambiscuit/internal/models"
)

type GitService struct {
    workspaceDir string
}

func NewGitService(workspaceDir string) *GitService {
    return &GitService{
        workspaceDir: workspaceDir,
    }
}

func (s *GitService) Clone(repoURL, relativePath string) *models.CloneResult {
    start := time.Now()
    result := &models.CloneResult{}
    
    // Calculate full path
    fullPath := filepath.Join(s.workspaceDir, relativePath)
    
    // Check if already exists - if so, treat as already cloned (success)
    if s.RepositoryExists(relativePath) {
        result.Success = true
        result.Output = fmt.Sprintf("repository already exists at %s, skipping clone", fullPath)
        result.Duration = time.Since(start)
        return result
    }
    
    // Create parent directories
    parentDir := filepath.Dir(fullPath)
    if err := os.MkdirAll(parentDir, 0755); err != nil {
        result.Success = false
        result.Error = fmt.Sprintf("failed to create parent directory: %v", err)
        result.Duration = time.Since(start)
        return result
    }
    
    // Execute git clone
    cmd := exec.Command("git", "clone", repoURL, fullPath)
    
    var stdout, stderr bytes.Buffer
    cmd.Stdout = &stdout
    cmd.Stderr = &stderr
    
    err := cmd.Run()
    
    result.Duration = time.Since(start)
    result.Output = stdout.String() + stderr.String()
    
    if err != nil {
        result.Success = false
        result.Error = fmt.Sprintf("git clone failed: %v", err)
        return result
    }
    
    result.Success = true
    return result
}

func (s *GitService) RepositoryExists(relativePath string) bool {
    fullPath := filepath.Join(s.workspaceDir, relativePath)
    gitDir := filepath.Join(fullPath, ".git")
    
    _, err := os.Stat(gitDir)
    return err == nil
}