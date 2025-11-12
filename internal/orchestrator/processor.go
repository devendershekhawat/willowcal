package orchestrator

import (
	"fmt"
	"strings"
	"time"

	"github.com/devendershekhawat/teambiscuit/internal/executor"
	"github.com/devendershekhawat/teambiscuit/internal/git"
	"github.com/devendershekhawat/teambiscuit/internal/models"
	"github.com/devendershekhawat/teambiscuit/internal/reporter"
)

// ProcessRepository handles cloning and setup for a single repository
func ProcessRepository(
    repo models.Repository,
    gitService *git.GitService,
    execService *executor.Service,
) *models.RepoState {
    
    state := models.NewRepoState(repo.Name)
    
    // Step 1: Clone repository
    state.Status = models.RepoStatusCloning
    reporter.PrintProgress(repo.Name, state.Status, "Cloning repository...")
    cloneResult := gitService.Clone(repo.URL, repo.Path)
    state.CloneResult = cloneResult
    
    // Show message if repository already exists
    if cloneResult.Success && cloneResult.Output != "" {
        if strings.Contains(cloneResult.Output, "already exists") || strings.Contains(cloneResult.Output, "skipping clone") {
            reporter.PrintProgress(repo.Name, state.Status, "Repository already exists, skipping clone")
        }
    }
    
    if !cloneResult.Success {
        state.Status = models.RepoStatusFailed
        state.Error = cloneResult.Error
        reporter.PrintProgress(repo.Name, state.Status, cloneResult.Error)
        state.EndTime = time.Now()
        return state
    }
    
    // Step 2: Run setup commands sequentially
    if len(repo.SetupCommands) > 0 {
        state.Status = models.RepoStatusSetupRunning
        reporter.PrintProgress(repo.Name, state.Status, fmt.Sprintf("Running %d setup command(s)...", len(repo.SetupCommands)))
        
        for i, cmd := range repo.SetupCommands {
            reporter.PrintProgress(repo.Name, state.Status, fmt.Sprintf("Executing command %d/%d: %s", i+1, len(repo.SetupCommands), cmd))
            cmdResult := execService.ExecuteCommand(cmd, repo.Path)
            state.SetupResults = append(state.SetupResults, cmdResult)
            
            // Stop on first command failure
            if !cmdResult.Success {
                state.Status = models.RepoStatusFailed
                state.Error = fmt.Sprintf("command '%s' failed: %s", cmd, cmdResult.Error)
                reporter.PrintProgress(repo.Name, state.Status, state.Error)
                state.EndTime = time.Now()
                return state
            }
        }
    }
    
    // Success!
    state.Status = models.RepoStatusSuccess
    reporter.PrintProgress(repo.Name, state.Status, "Repository initialized successfully")
    state.EndTime = time.Now()
    return state
}