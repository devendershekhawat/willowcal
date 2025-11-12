package executor

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/devendershekhawat/teambiscuit/internal/models"
)

const DefaultTimeout = 5 * time.Minute

type Service struct {
    workspaceDir string
    timeout      time.Duration
}

func NewService(workspaceDir string) *Service {
    return &Service{
        workspaceDir: workspaceDir,
        timeout:      DefaultTimeout,
    }
}

// ExecuteCommand runs a command in the specified directory
func (s *Service) ExecuteCommand(command, relativePath string) *models.CommandResult {
    start := time.Now()
    result := &models.CommandResult{
        Command: command,
    }
    
    // Create context with timeout
    ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
    defer cancel()
    
    // Determine working directory
    workingDir := s.workspaceDir
    if relativePath != "" {
        workingDir = workingDir + "/" + relativePath
    }
    
    // Parse and create command
    var cmd *exec.Cmd
    if s.isComplexCommand(command) {
        // Use shell for complex commands
        cmd = exec.CommandContext(ctx, "sh", "-c", command)
    } else {
        // Split simple commands
        parts := strings.Fields(command)
        if len(parts) == 0 {
            result.Success = false
            result.Error = "empty command"
            result.ExitCode = -1
            result.Duration = time.Since(start)
            return result
        }
        cmd = exec.CommandContext(ctx, parts[0], parts[1:]...)
    }
    
    cmd.Dir = workingDir
    
    // Capture output
    var stdout, stderr bytes.Buffer
    cmd.Stdout = &stdout
    cmd.Stderr = &stderr
    
    // Execute
    err := cmd.Run()
    
    result.Duration = time.Since(start)
    result.Stdout = stdout.String()
    result.Stderr = stderr.String()
    
    if err != nil {
        result.Success = false
        result.Error = err.Error()
        
        if exitErr, ok := err.(*exec.ExitError); ok {
            result.ExitCode = exitErr.ExitCode()
        } else if ctx.Err() == context.DeadlineExceeded {
            result.Error = fmt.Sprintf("command timeout after %v", s.timeout)
            result.ExitCode = -1
        } else {
            result.ExitCode = -1
        }
        
        return result
    }
    
    result.Success = true
    result.ExitCode = 0
    return result
}

// isComplexCommand checks if command needs shell
func (s *Service) isComplexCommand(cmd string) bool {
    return strings.Contains(cmd, "&&") ||
           strings.Contains(cmd, "||") ||
           strings.Contains(cmd, "|") ||
           strings.Contains(cmd, ";") ||
           strings.Contains(cmd, ">") ||
           strings.Contains(cmd, "<")
}