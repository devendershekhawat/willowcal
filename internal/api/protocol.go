package api

// MessageType defines the type of WebSocket message
type MessageType string

const (
	// Client -> Server messages
	TypeConfigUpload    MessageType = "config.upload"
	TypeConfigParse     MessageType = "config.parse"
	TypeInitStart       MessageType = "init.start"
	TypeServiceList     MessageType = "service.list"
	TypeServiceStart    MessageType = "service.start"
	TypeServiceStop     MessageType = "service.stop"
	TypeServiceStatus   MessageType = "service.status"
	TypeServiceLogs     MessageType = "service.logs"
	TypeConfigUpdate    MessageType = "config.update"
	TypeConfigDiff      MessageType = "config.diff"

	// Server -> Client messages (Events)
	TypeInitProgress    MessageType = "init.progress"
	TypeInitComplete    MessageType = "init.complete"
	TypeInitError       MessageType = "init.error"
	TypeServiceLog      MessageType = "service.log"
	TypeServiceStarted  MessageType = "service.started"
	TypeServiceStopped  MessageType = "service.stopped"
	TypeServiceError    MessageType = "service.error"
	TypeError           MessageType = "error"
	TypeSuccess         MessageType = "success"
)

// Message represents a WebSocket message
type Message struct {
	Type    MessageType `json:"type"`
	ID      string      `json:"id,omitempty"`      // Request ID for correlation
	Payload interface{} `json:"payload,omitempty"`
}

// ConfigUploadPayload is sent when uploading a config
type ConfigUploadPayload struct {
	ConfigYAML string `json:"config_yaml"`
	Name       string `json:"name,omitempty"` // Optional config name/identifier
}

// ConfigParseResponse returns parsed config info
type ConfigParseResponse struct {
	Valid        bool     `json:"valid"`
	Repositories int      `json:"repositories"`
	Services     int      `json:"services"`
	WorkspaceDir string   `json:"workspace_dir"`
	Errors       []string `json:"errors,omitempty"`
}

// InitStartPayload starts the initialization process
type InitStartPayload struct {
	ConfigYAML string `json:"config_yaml,omitempty"`
}

// InitProgressPayload streams init progress
type InitProgressPayload struct {
	RepoName string `json:"repo_name"`
	Status   string `json:"status"`
	Message  string `json:"message"`
	LogLine  string `json:"log_line,omitempty"`
}

// InitCompletePayload is sent when init completes
type InitCompletePayload struct {
	Success      int     `json:"success"`
	Failed       int     `json:"failed"`
	TotalTime    float64 `json:"total_time_seconds"`
	Repositories []RepoSummary `json:"repositories"`
}

// RepoSummary provides summary of a repository operation
type RepoSummary struct {
	Name     string  `json:"name"`
	Status   string  `json:"status"`
	Duration float64 `json:"duration_seconds"`
	Error    string  `json:"error,omitempty"`
}

// ServiceListResponse returns list of services
type ServiceListResponse struct {
	Services []ServiceInfo `json:"services"`
}

// ServiceInfo provides information about a service
type ServiceInfo struct {
	Name       string `json:"name"`
	Repository string `json:"repository"`
	RunCommand string `json:"run_command"`
	Status     string `json:"status"` // "stopped", "starting", "running", "failed"
}

// ServiceStartPayload starts a service
type ServiceStartPayload struct {
	ServiceName string `json:"service_name"`
}

// ServiceStopPayload stops a service
type ServiceStopPayload struct {
	ServiceName string `json:"service_name"`
}

// ServiceStatusPayload requests service status
type ServiceStatusPayload struct {
	ServiceName string `json:"service_name,omitempty"` // Empty means all services
}

// ServiceStatusResponse returns service status
type ServiceStatusResponse struct {
	Services []ServiceStatus `json:"services"`
}

// ServiceStatus represents the current status of a service
type ServiceStatus struct {
	Name      string  `json:"name"`
	Status    string  `json:"status"`
	PID       int     `json:"pid,omitempty"`
	Uptime    float64 `json:"uptime_seconds,omitempty"`
	Error     string  `json:"error,omitempty"`
}

// ServiceLogPayload is sent when streaming service logs
type ServiceLogPayload struct {
	ServiceName string `json:"service_name"`
	Timestamp   string `json:"timestamp"`
	Line        string `json:"line"`
	Stream      string `json:"stream"` // "stdout" or "stderr"
}

// ServiceLogsPayload requests service logs
type ServiceLogsPayload struct {
	ServiceName string `json:"service_name"`
	Follow      bool   `json:"follow"`      // Stream logs in real-time
	Tail        int    `json:"tail"`        // Number of recent lines to return
}

// ConfigUpdatePayload updates the config
type ConfigUpdatePayload struct {
	ConfigYAML string `json:"config_yaml"`
}

// ConfigDiffPayload requests diff between current and new config
type ConfigDiffPayload struct {
	NewConfigYAML string `json:"new_config_yaml"`
}

// ConfigDiffResponse returns the differences
type ConfigDiffResponse struct {
	HasChanges       bool     `json:"has_changes"`
	AddedRepos       []string `json:"added_repos"`
	RemovedRepos     []string `json:"removed_repos"`
	ModifiedRepos    []string `json:"modified_repos"`
	AddedServices    []string `json:"added_services"`
	RemovedServices  []string `json:"removed_services"`
	ModifiedServices []string `json:"modified_services"`
}

// ErrorPayload represents an error response
type ErrorPayload struct {
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
}

// SuccessPayload represents a success response
type SuccessPayload struct {
	Message string `json:"message"`
}
