package entities

import "time"

// AgentStatus represents the status of an agent
type AgentStatus string

const (
	AgentStatusIdle      AgentStatus = "idle"
	AgentStatusRunning   AgentStatus = "running"
	AgentStatusPaused    AgentStatus = "paused"
	AgentStatusCompleted AgentStatus = "completed"
	AgentStatusFailed    AgentStatus = "failed"
	AgentStatusStopped   AgentStatus = "stopped"
)

// RuntimeType represents the agent runtime type
type RuntimeType string

const (
	RuntimeDefault    RuntimeType = "default"
	RuntimeClaudeCode RuntimeType = "claude_code"
)

// Agent represents an agent in the control plane
type Agent struct {
	ID            string                 `json:"id,omitempty"`
	Name          string                 `json:"name"`
	Description   *string                `json:"description,omitempty"`
	Status        AgentStatus            `json:"status,omitempty"`
	Capabilities  []string               `json:"capabilities,omitempty"`
	Configuration map[string]interface{} `json:"configuration,omitempty"`
	ModelID       *string                `json:"model_id,omitempty"`
	LLMConfig     map[string]interface{} `json:"llm_config,omitempty"`
	Runtime       RuntimeType            `json:"runtime,omitempty"`
	TeamID        *string                `json:"team_id,omitempty"`
	CreatedAt     *time.Time             `json:"created_at,omitempty"`
	UpdatedAt     *time.Time             `json:"updated_at,omitempty"`
	LastActiveAt  *time.Time             `json:"last_active_at,omitempty"`
	State         map[string]interface{} `json:"state,omitempty"`
	ErrorMessage  *string                `json:"error_message,omitempty"`
}

// AgentCreateRequest represents the request to create an agent
type AgentCreateRequest struct {
	Name          string                 `json:"name"`
	Description   *string                `json:"description,omitempty"`
	Capabilities  []string               `json:"capabilities,omitempty"`
	Configuration map[string]interface{} `json:"configuration,omitempty"`
	ModelID       *string                `json:"model_id,omitempty"`
	LLMConfig     map[string]interface{} `json:"llm_config,omitempty"`
	Runtime       *RuntimeType           `json:"runtime,omitempty"`
	TeamID        *string                `json:"team_id,omitempty"`
}

// AgentUpdateRequest represents the request to update an agent
type AgentUpdateRequest struct {
	Name          *string                `json:"name,omitempty"`
	Description   *string                `json:"description,omitempty"`
	Status        *AgentStatus           `json:"status,omitempty"`
	Capabilities  []string               `json:"capabilities,omitempty"`
	Configuration map[string]interface{} `json:"configuration,omitempty"`
	State         map[string]interface{} `json:"state,omitempty"`
	ModelID       *string                `json:"model_id,omitempty"`
	LLMConfig     map[string]interface{} `json:"llm_config,omitempty"`
	Runtime       *RuntimeType           `json:"runtime,omitempty"`
	TeamID        *string                `json:"team_id,omitempty"`
}
