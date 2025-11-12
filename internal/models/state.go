package models

import "time"

type ExecutionState struct {
    StartTime     time.Time
    EndTime       time.Time
    TotalRepos    int
    SuccessCount  int
    FailureCount  int
    RetryCount    int
    RepoStates    map[string]*RepoState
    Status        ExecutionStatus
}

type ExecutionStatus string

const (
    ExecutionStatusRunning   ExecutionStatus = "running"
    ExecutionStatusCompleted ExecutionStatus = "completed"
    ExecutionStatusFailed    ExecutionStatus = "failed"
)

type RepoState struct {
    Name             string
    Status           RepoStatus
    CloneResult      *CloneResult
    SetupResults     []*CommandResult
    CurrentRetry     int
    Error            string
    StartTime        time.Time
    EndTime          time.Time
}

type RepoStatus string

const (
    RepoStatusPending      RepoStatus = "pending"
    RepoStatusCloning      RepoStatus = "cloning"
    RepoStatusSetupRunning RepoStatus = "setup_running"
    RepoStatusSuccess      RepoStatus = "success"
    RepoStatusFailed       RepoStatus = "failed"
)

type CloneResult struct {
    Success  bool
    Error    string
    Duration time.Duration
    Output   string
}



type CommandResult struct {
    Command  string
    Success  bool
    ExitCode int
    Error    string
    Duration time.Duration
    Stdout   string
    Stderr   string
}

func NewExecutionState(totalRepos int) *ExecutionState {
    return &ExecutionState{
        StartTime:    time.Now(),
        TotalRepos:   totalRepos,
        RepoStates:   make(map[string]*RepoState),
        Status:       ExecutionStatusRunning,
    }
}

// NewRepoState creates initial repo state
func NewRepoState(name string) *RepoState {
    return &RepoState{
        Name:         name,
        Status:       RepoStatusPending,
        SetupResults: make([]*CommandResult, 0),
        StartTime:    time.Now(),
    }
}