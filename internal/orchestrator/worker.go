package orchestrator

import (
	"github.com/devendershekhawat/teambiscuit/internal/executor"
	"github.com/devendershekhawat/teambiscuit/internal/git"
	"github.com/devendershekhawat/teambiscuit/internal/models"
)

// Worker processes repositories from job queue
func Worker(
    id int,
    jobs <-chan models.Repository,
    results chan<- *models.RepoState,
    gitService *git.GitService,
    execService *executor.Service,
) {
    for repo := range jobs {
        // Process repository
        state := ProcessRepository(repo, gitService, execService)
        
        // Send result
        results <- state
    }
}