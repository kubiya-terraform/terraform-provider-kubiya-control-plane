package entities

import "time"

// EnvironmentStatus represents the status of an environment
type EnvironmentStatus string

const (
	EnvironmentStatusActive   EnvironmentStatus = "active"
	EnvironmentStatusInactive EnvironmentStatus = "inactive"
)

// ExecutionEnvironment represents execution environment configuration
type ExecutionEnvironment struct {
	EnvVars        map[string]string `json:"env_vars,omitempty"`
	Secrets        []string          `json:"secrets,omitempty"`
	IntegrationIDs []string          `json:"integration_ids,omitempty"`
}

// Environment represents an environment in the control plane
type Environment struct {
	ID                     string                   `json:"id,omitempty"`
	OrganizationID         string                   `json:"organization_id,omitempty"`
	Name                   string                   `json:"name"`
	DisplayName            *string                  `json:"display_name,omitempty"`
	Description            *string                  `json:"description,omitempty"`
	Tags                   []string                 `json:"tags,omitempty"`
	Settings               map[string]interface{}   `json:"settings,omitempty"`
	Status                 EnvironmentStatus        `json:"status,omitempty"`
	CreatedAt              *time.Time               `json:"created_at,omitempty"`
	UpdatedAt              *time.Time               `json:"updated_at,omitempty"`
	CreatedBy              *string                  `json:"created_by,omitempty"`
	WorkerToken            *string                  `json:"worker_token,omitempty"`
	ProvisioningWorkflowID *string                  `json:"provisioning_workflow_id,omitempty"`
	ProvisionedAt          *time.Time               `json:"provisioned_at,omitempty"`
	ErrorMessage           *string                  `json:"error_message,omitempty"`
	TemporalNamespaceID    *string                  `json:"temporal_namespace_id,omitempty"`
	ActiveWorkers          int                      `json:"active_workers,omitempty"`
	IdleWorkers            int                      `json:"idle_workers,omitempty"`
	BusyWorkers            int                      `json:"busy_workers,omitempty"`
	ToolsetIDs             []string                 `json:"toolset_ids,omitempty"`
	Toolsets               []map[string]interface{} `json:"toolsets,omitempty"`
	ExecutionEnvironment   map[string]interface{}   `json:"execution_environment,omitempty"`
}

// EnvironmentCreateRequest represents the request to create an environment
type EnvironmentCreateRequest struct {
	Name                 string                 `json:"name"`
	DisplayName          *string                `json:"display_name,omitempty"`
	Description          *string                `json:"description,omitempty"`
	Tags                 []string               `json:"tags,omitempty"`
	Settings             map[string]interface{} `json:"settings,omitempty"`
	ExecutionEnvironment map[string]interface{} `json:"execution_environment,omitempty"`
}

// EnvironmentUpdateRequest represents the request to update an environment
type EnvironmentUpdateRequest struct {
	Name                 *string                `json:"name,omitempty"`
	DisplayName          *string                `json:"display_name,omitempty"`
	Description          *string                `json:"description,omitempty"`
	Tags                 []string               `json:"tags,omitempty"`
	Settings             map[string]interface{} `json:"settings,omitempty"`
	Status               *EnvironmentStatus     `json:"status,omitempty"`
	ExecutionEnvironment map[string]interface{} `json:"execution_environment,omitempty"`
}
