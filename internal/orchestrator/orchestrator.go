package orchestrator

import (
	"sync"
	"time"

	"github.com/devendershekhawat/teambiscuit/internal/config"
	"github.com/devendershekhawat/teambiscuit/internal/executor"
	"github.com/devendershekhawat/teambiscuit/internal/git"
	"github.com/devendershekhawat/teambiscuit/internal/models"
)

const (
    MaxParallelRepos = 5
    MaxRetries       = 3
)

type Orchestrator struct {
    config      *config.Config
    gitService  *git.GitService
    execService *executor.Service
    state       *models.ExecutionState
    mu          sync.Mutex
}

func NewOrchestrator(cfg *config.Config, workspaceDir string) *Orchestrator {
    return &Orchestrator{
        config:      cfg,
        gitService:  git.NewGitService(workspaceDir),
        execService: executor.NewService(workspaceDir),
        state:       models.NewExecutionState(len(cfg.Repositories)),
    }
}

// Execute runs the orchestration
func (o *Orchestrator) Execute() *models.ExecutionState {
    repos := o.config.Repositories
    totalRepos := len(repos)
    
    // Create channels
    jobs := make(chan models.Repository, totalRepos)
    results := make(chan *models.RepoState, totalRepos)
    
    // Start worker pool
    numWorkers := min(MaxParallelRepos, totalRepos)
    var wg sync.WaitGroup
    
    for i := 0; i < numWorkers; i++ {
        wg.Add(1)
        go func(workerID int) {
            defer wg.Done()
            Worker(workerID, jobs, results, o.gitService, o.execService)
        }(i)
    }
    
    // Send jobs
    for _, repo := range repos {
        jobs <- repo
    }
    close(jobs)
    
    // Collect results in separate goroutine
    go func() {
        wg.Wait()
        close(results)
    }()
    
    // Process results and handle retries
    processed := 0
    retryQueue := make([]models.Repository, 0)
    retryCounts := make(map[string]int)
    
    for state := range results {
        processed++
        
        o.mu.Lock()
        o.state.RepoStates[state.Name] = state
        o.mu.Unlock()
        
        if state.Status == models.RepoStatusFailed {
            // Find original repo
            var repo *models.Repository
            for i := range repos {
                if repos[i].Name == state.Name {
                    repo = &repos[i]
                    break
                }
            }
            
            if repo != nil {
                retryCount := retryCounts[repo.Name]
                if retryCount < MaxRetries {
                    retryQueue = append(retryQueue, *repo)
                    retryCounts[repo.Name] = retryCount + 1
                    o.state.RetryCount++
                }
            }
        }
    }
    
    // Handle retries
    if len(retryQueue) > 0 {
        o.executeRetries(retryQueue, retryCounts)
    }
    
    // Finalize state
    o.finalizeState()
    
    return o.state
}

// executeRetries handles retry logic
func (o *Orchestrator) executeRetries(retryQueue []models.Repository, retryCounts map[string]int) {
    for len(retryQueue) > 0 {
        repo := retryQueue[0]
        retryQueue = retryQueue[1:]
        
        // Process retry
        state := ProcessRepository(repo, o.gitService, o.execService)
        state.CurrentRetry = retryCounts[repo.Name]
        
        o.mu.Lock()
        o.state.RepoStates[repo.Name] = state
        o.mu.Unlock()
        
        // If still failed and retries remaining, add back to queue
        if state.Status == models.RepoStatusFailed {
            retryCount := retryCounts[repo.Name]
            if retryCount < MaxRetries {
                retryQueue = append(retryQueue, repo)
                retryCounts[repo.Name] = retryCount + 1
                o.state.RetryCount++
            }
        }
    }
}

// finalizeState calculates final statistics
func (o *Orchestrator) finalizeState() {
    o.mu.Lock()
    defer o.mu.Unlock()
    
    for _, state := range o.state.RepoStates {
        if state.Status == models.RepoStatusSuccess {
            o.state.SuccessCount++
        } else {
            o.state.FailureCount++
        }
    }
    
    o.state.EndTime = time.Now()
    if o.state.FailureCount > 0 {
        o.state.Status = models.ExecutionStatusFailed
    } else {
        o.state.Status = models.ExecutionStatusCompleted
    }
}

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}