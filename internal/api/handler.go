package api

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/devendershekhawat/teambiscuit/internal/config"
	"github.com/devendershekhawat/teambiscuit/internal/orchestrator"
	"github.com/devendershekhawat/teambiscuit/internal/service"
	"gopkg.in/yaml.v3"
)

// Handler handles WebSocket messages
type Handler struct {
	config        *config.Config
	workspaceDir  string
	serviceManager *service.Manager
	broadcaster   func(Message)
}

// NewHandler creates a new message handler
func NewHandler(workspaceDir string) *Handler {
	return &Handler{
		workspaceDir: workspaceDir,
	}
}

// SetBroadcaster sets the broadcast function
func (h *Handler) SetBroadcaster(broadcaster func(Message)) {
	h.broadcaster = broadcaster
}

// HandleMessage processes incoming WebSocket messages
func (h *Handler) HandleMessage(msg Message) *Message {
	log.Printf("📨 Received message: type=%s id=%s", msg.Type, msg.ID)

	switch msg.Type {
	case TypeConfigUpload:
		return h.handleConfigUpload(msg)
	case TypeConfigParse:
		return h.handleConfigParse(msg)
	case TypeInitStart:
		return h.handleInitStart(msg)
	case TypeServiceList:
		return h.handleServiceList(msg)
	case TypeServiceStart:
		return h.handleServiceStart(msg)
	case TypeServiceStop:
		return h.handleServiceStop(msg)
	case TypeServiceStatus:
		return h.handleServiceStatus(msg)
	case TypeConfigDiff:
		return h.handleConfigDiff(msg)
	case TypeConfigUpdate:
		return h.handleConfigUpdate(msg)
	default:
		return &Message{
			Type: TypeError,
			ID:   msg.ID,
			Payload: ErrorPayload{
				Message: fmt.Sprintf("Unknown message type: %s", msg.Type),
			},
		}
	}
}

// handleConfigUpload parses and validates a config
func (h *Handler) handleConfigUpload(msg Message) *Message {
	payload, ok := msg.Payload.(map[string]interface{})
	if !ok {
		return h.errorResponse(msg.ID, "Invalid payload format")
	}

	configYAML, ok := payload["config_yaml"].(string)
	if !ok || configYAML == "" {
		return h.errorResponse(msg.ID, "Missing config_yaml field")
	}

	// Parse config
	var cfg config.Config
	if err := yaml.Unmarshal([]byte(configYAML), &cfg); err != nil {
		return h.errorResponse(msg.ID, fmt.Sprintf("Failed to parse YAML: %v", err))
	}

	// Validate config
	if err := config.ValidateConfig(&cfg); err != nil {
		return &Message{
			Type: TypeSuccess,
			ID:   msg.ID,
			Payload: ConfigParseResponse{
				Valid:  false,
				Errors: strings.Split(err.Error(), "\n"),
			},
		}
	}

	// Store config
	h.config = &cfg

	// Get absolute workspace
	workspaceDir, err := cfg.GetAbsoluteWorkspace()
	if err != nil {
		return h.errorResponse(msg.ID, fmt.Sprintf("Failed to resolve workspace: %v", err))
	}
	h.workspaceDir = workspaceDir

	// Create workspace directory
	if err := os.MkdirAll(workspaceDir, 0755); err != nil {
		return h.errorResponse(msg.ID, fmt.Sprintf("Failed to create workspace: %v", err))
	}

	// Initialize service manager
	h.serviceManager = service.NewManager(&cfg, workspaceDir)

	// Start log broadcaster
	go h.broadcastServiceLogs()

	return &Message{
		Type: TypeSuccess,
		ID:   msg.ID,
		Payload: ConfigParseResponse{
			Valid:        true,
			Repositories: cfg.RepositoryCount(),
			Services:     cfg.ServiceCount(),
			WorkspaceDir: cfg.WorkspaceDir,
		},
	}
}

// handleConfigParse just parses without storing
func (h *Handler) handleConfigParse(msg Message) *Message {
	payload, ok := msg.Payload.(map[string]interface{})
	if !ok {
		return h.errorResponse(msg.ID, "Invalid payload format")
	}

	configYAML, ok := payload["config_yaml"].(string)
	if !ok || configYAML == "" {
		return h.errorResponse(msg.ID, "Missing config_yaml field")
	}

	// Parse config
	var cfg config.Config
	if err := yaml.Unmarshal([]byte(configYAML), &cfg); err != nil {
		return h.errorResponse(msg.ID, fmt.Sprintf("Failed to parse YAML: %v", err))
	}

	// Validate config
	if err := config.ValidateConfig(&cfg); err != nil {
		return &Message{
			Type: TypeSuccess,
			ID:   msg.ID,
			Payload: ConfigParseResponse{
				Valid:  false,
				Errors: strings.Split(err.Error(), "\n"),
			},
		}
	}

	return &Message{
		Type: TypeSuccess,
		ID:   msg.ID,
		Payload: ConfigParseResponse{
			Valid:        true,
			Repositories: cfg.RepositoryCount(),
			Services:     cfg.ServiceCount(),
			WorkspaceDir: cfg.WorkspaceDir,
		},
	}
}

// handleInitStart starts the initialization process
func (h *Handler) handleInitStart(msg Message) *Message {
	if h.config == nil {
		return h.errorResponse(msg.ID, "No config uploaded")
	}

	// Start init in background
	go h.runInit(msg.ID)

	return &Message{
		Type: TypeSuccess,
		ID:   msg.ID,
		Payload: SuccessPayload{
			Message: "Initialization started",
		},
	}
}

// runInit runs the initialization process
func (h *Handler) runInit(requestID string) {
	log.Printf("🚀 Starting initialization...")

	orch := orchestrator.NewOrchestrator(h.config, h.workspaceDir)
	state := orch.Execute()

	// Convert state to response
	repos := make([]RepoSummary, 0, len(state.RepoStates))
	for _, repoState := range state.RepoStates {
		duration := repoState.EndTime.Sub(repoState.StartTime).Seconds()
		repos = append(repos, RepoSummary{
			Name:     repoState.Name,
			Status:   string(repoState.Status),
			Duration: duration,
			Error:    repoState.Error,
		})
	}

	totalTime := state.EndTime.Sub(state.StartTime).Seconds()

	// Broadcast completion
	if h.broadcaster != nil {
		h.broadcaster(Message{
			Type: TypeInitComplete,
			ID:   requestID,
			Payload: InitCompletePayload{
				Success:      state.SuccessCount,
				Failed:       state.FailureCount,
				TotalTime:    totalTime,
				Repositories: repos,
			},
		})
	}

	log.Printf("✅ Initialization complete: %d succeeded, %d failed", state.SuccessCount, state.FailureCount)
}

// handleServiceList returns list of services
func (h *Handler) handleServiceList(msg Message) *Message {
	if h.config == nil {
		return h.errorResponse(msg.ID, "No config uploaded")
	}

	services := make([]ServiceInfo, 0, len(h.config.Services))
	statuses := h.serviceManager.GetAllStatuses()

	for _, status := range statuses {
		services = append(services, ServiceInfo{
			Name:       status.Name,
			Repository: status.Service.Repository,
			RunCommand: status.Service.RunCommand,
			Status:     string(status.State),
		})
	}

	return &Message{
		Type: TypeSuccess,
		ID:   msg.ID,
		Payload: ServiceListResponse{
			Services: services,
		},
	}
}

// handleServiceStart starts a service
func (h *Handler) handleServiceStart(msg Message) *Message {
	if h.serviceManager == nil {
		return h.errorResponse(msg.ID, "Service manager not initialized")
	}

	payload, ok := msg.Payload.(map[string]interface{})
	if !ok {
		return h.errorResponse(msg.ID, "Invalid payload format")
	}

	serviceName, ok := payload["service_name"].(string)
	if !ok || serviceName == "" {
		return h.errorResponse(msg.ID, "Missing service_name field")
	}

	if err := h.serviceManager.Start(serviceName); err != nil {
		return h.errorResponse(msg.ID, fmt.Sprintf("Failed to start service: %v", err))
	}

	// Broadcast service started event
	if h.broadcaster != nil {
		h.broadcaster(Message{
			Type: TypeServiceStarted,
			Payload: map[string]string{
				"service_name": serviceName,
			},
		})
	}

	return &Message{
		Type: TypeSuccess,
		ID:   msg.ID,
		Payload: SuccessPayload{
			Message: fmt.Sprintf("Service '%s' started", serviceName),
		},
	}
}

// handleServiceStop stops a service
func (h *Handler) handleServiceStop(msg Message) *Message {
	if h.serviceManager == nil {
		return h.errorResponse(msg.ID, "Service manager not initialized")
	}

	payload, ok := msg.Payload.(map[string]interface{})
	if !ok {
		return h.errorResponse(msg.ID, "Invalid payload format")
	}

	serviceName, ok := payload["service_name"].(string)
	if !ok || serviceName == "" {
		return h.errorResponse(msg.ID, "Missing service_name field")
	}

	if err := h.serviceManager.Stop(serviceName); err != nil {
		return h.errorResponse(msg.ID, fmt.Sprintf("Failed to stop service: %v", err))
	}

	// Broadcast service stopped event
	if h.broadcaster != nil {
		h.broadcaster(Message{
			Type: TypeServiceStopped,
			Payload: map[string]string{
				"service_name": serviceName,
			},
		})
	}

	return &Message{
		Type: TypeSuccess,
		ID:   msg.ID,
		Payload: SuccessPayload{
			Message: fmt.Sprintf("Service '%s' stopped", serviceName),
		},
	}
}

// handleServiceStatus returns service status
func (h *Handler) handleServiceStatus(msg Message) *Message {
	if h.serviceManager == nil {
		return h.errorResponse(msg.ID, "Service manager not initialized")
	}

	statuses := h.serviceManager.GetAllStatuses()
	serviceStatuses := make([]ServiceStatus, 0, len(statuses))

	for _, status := range statuses {
		var pid int
		var uptime float64

		if status.Process != nil && status.Process.Process != nil {
			pid = status.Process.Process.Pid
		}

		if status.State == service.StateRunning {
			uptime = 0 // Will be calculated on client side or we can add it
		}

		serviceStatuses = append(serviceStatuses, ServiceStatus{
			Name:   status.Name,
			Status: string(status.State),
			PID:    pid,
			Uptime: uptime,
			Error:  status.Error,
		})
	}

	return &Message{
		Type: TypeSuccess,
		ID:   msg.ID,
		Payload: ServiceStatusResponse{
			Services: serviceStatuses,
		},
	}
}

// handleConfigDiff computes diff between current and new config
func (h *Handler) handleConfigDiff(msg Message) *Message {
	if h.config == nil {
		return h.errorResponse(msg.ID, "No current config to compare against")
	}

	payload, ok := msg.Payload.(map[string]interface{})
	if !ok {
		return h.errorResponse(msg.ID, "Invalid payload format")
	}

	newConfigYAML, ok := payload["new_config_yaml"].(string)
	if !ok || newConfigYAML == "" {
		return h.errorResponse(msg.ID, "Missing new_config_yaml field")
	}

	// Parse new config
	var newCfg config.Config
	if err := yaml.Unmarshal([]byte(newConfigYAML), &newCfg); err != nil {
		return h.errorResponse(msg.ID, fmt.Sprintf("Failed to parse YAML: %v", err))
	}

	// Compute diff
	diff := h.computeConfigDiff(h.config, &newCfg)

	return &Message{
		Type: TypeSuccess,
		ID:   msg.ID,
		Payload: diff,
	}
}

// handleConfigUpdate updates the config
func (h *Handler) handleConfigUpdate(msg Message) *Message {
	payload, ok := msg.Payload.(map[string]interface{})
	if !ok {
		return h.errorResponse(msg.ID, "Invalid payload format")
	}

	configYAML, ok := payload["config_yaml"].(string)
	if !ok || configYAML == "" {
		return h.errorResponse(msg.ID, "Missing config_yaml field")
	}

	// Parse config
	var cfg config.Config
	if err := yaml.Unmarshal([]byte(configYAML), &cfg); err != nil {
		return h.errorResponse(msg.ID, fmt.Sprintf("Failed to parse YAML: %v", err))
	}

	// Validate config
	if err := config.ValidateConfig(&cfg); err != nil {
		return h.errorResponse(msg.ID, fmt.Sprintf("Invalid config: %v", err))
	}

	// Update config
	h.config = &cfg

	// Update service manager
	if h.serviceManager != nil {
		h.serviceManager.UpdateConfig(&cfg)
	}

	return &Message{
		Type: TypeSuccess,
		ID:   msg.ID,
		Payload: SuccessPayload{
			Message: "Config updated successfully",
		},
	}
}

// computeConfigDiff compares two configs and returns differences
func (h *Handler) computeConfigDiff(oldCfg, newCfg *config.Config) ConfigDiffResponse {
	diff := ConfigDiffResponse{}

	// Create maps for quick lookup
	oldRepos := make(map[string]bool)
	newRepos := make(map[string]bool)
	oldServices := make(map[string]bool)
	newServices := make(map[string]bool)

	for _, repo := range oldCfg.Repositories {
		oldRepos[repo.Name] = true
	}
	for _, repo := range newCfg.Repositories {
		newRepos[repo.Name] = true
	}
	for _, svc := range oldCfg.Services {
		oldServices[svc.Name] = true
	}
	for _, svc := range newCfg.Services {
		newServices[svc.Name] = true
	}

	// Find added/removed repos
	for name := range newRepos {
		if !oldRepos[name] {
			diff.AddedRepos = append(diff.AddedRepos, name)
			diff.HasChanges = true
		}
	}
	for name := range oldRepos {
		if !newRepos[name] {
			diff.RemovedRepos = append(diff.RemovedRepos, name)
			diff.HasChanges = true
		}
	}

	// Find added/removed services
	for name := range newServices {
		if !oldServices[name] {
			diff.AddedServices = append(diff.AddedServices, name)
			diff.HasChanges = true
		}
	}
	for name := range oldServices {
		if !newServices[name] {
			diff.RemovedServices = append(diff.RemovedServices, name)
			diff.HasChanges = true
		}
	}

	// TODO: Check for modified repos/services (URL changes, command changes, etc.)

	return diff
}

// broadcastServiceLogs broadcasts service logs to all clients
func (h *Handler) broadcastServiceLogs() {
	if h.serviceManager == nil || h.broadcaster == nil {
		return
	}

	logChan := h.serviceManager.GetLogChannel()
	for entry := range logChan {
		h.broadcaster(Message{
			Type: TypeServiceLog,
			Payload: ServiceLogPayload{
				ServiceName: entry.ServiceName,
				Timestamp:   entry.Timestamp.Format("15:04:05"),
				Line:        entry.Line,
				Stream:      entry.Stream,
			},
		})
	}
}

// errorResponse creates an error response message
func (h *Handler) errorResponse(requestID string, message string) *Message {
	return &Message{
		Type: TypeError,
		ID:   requestID,
		Payload: ErrorPayload{
			Message: message,
		},
	}
}
