package reporter

import (
	"fmt"
	"time"

	"github.com/devendershekhawat/teambiscuit/internal/models"
)

// PrintProgress prints live progress (call this from orchestrator during execution)
func PrintProgress(repoName string, status models.RepoStatus, message string) {
    icon := "⏳"
    switch status {
    case models.RepoStatusCloning:
        icon = "📥"
    case models.RepoStatusSetupRunning:
        icon = "⚙️ "
    case models.RepoStatusSuccess:
        icon = "✅"
    case models.RepoStatusFailed:
        icon = "❌"
    }
    
    fmt.Printf("%s [%s] %s\n", icon, repoName, message)
}

// PrintFinalSummary prints execution summary
func PrintFinalSummary(state *models.ExecutionState) {
    duration := state.EndTime.Sub(state.StartTime)
    
    fmt.Println("\n" + "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
    fmt.Println("📊 Execution Summary")
    fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
    fmt.Printf("Total Duration: %v\n", duration.Round(time.Millisecond))
    fmt.Printf("Total Repositories: %d\n", state.TotalRepos)
    fmt.Printf("✅ Successful: %d\n", state.SuccessCount)
    fmt.Printf("❌ Failed: %d\n", state.FailureCount)
    fmt.Printf("🔄 Total Retries: %d\n\n", state.RetryCount)
    
    // Print failed repositories
    if state.FailureCount > 0 {
        fmt.Println("Failed Repositories:")
        for name, repoState := range state.RepoStates {
            if repoState.Status == models.RepoStatusFailed {
                fmt.Printf("  • %s\n", name)
                fmt.Printf("    Error: %s\n", repoState.Error)
                if repoState.CurrentRetry > 0 {
                    fmt.Printf("    Retries: %d/%d\n", repoState.CurrentRetry, 3)
                }
            }
        }
        fmt.Println()
    }
    
    // Print successful repositories
    if state.SuccessCount > 0 {
        fmt.Println("Successful Repositories:")
        for name, repoState := range state.RepoStates {
            if repoState.Status == models.RepoStatusSuccess {
                duration := repoState.EndTime.Sub(repoState.StartTime)
                fmt.Printf("  • %s (%.1fs)\n", name, duration.Seconds())
            }
        }
    }
}