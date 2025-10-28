package entities

import "time"

// WorkerStatus represents the status of a worker
type WorkerStatus string

const (
	WorkerStatusActive       WorkerStatus = "active"
	WorkerStatusInactive     WorkerStatus = "inactive"
	WorkerStatusDisconnected WorkerStatus = "disconnected"
)

// Worker represents a worker in the control plane
type Worker struct {
	ID              string                 `json:"id,omitempty"`
	WorkerID        string                 `json:"worker_id,omitempty"`
	OrganizationID  string                 `json:"organization_id,omitempty"`
	EnvironmentName string                 `json:"environment_name"`
	Hostname        *string                `json:"hostname,omitempty"`
	Status          WorkerStatus           `json:"status,omitempty"`
	WorkerMetadata  map[string]interface{} `json:"worker_metadata,omitempty"`
	RegisteredAt    *time.Time             `json:"registered_at,omitempty"`
	LastHeartbeat   *time.Time             `json:"last_heartbeat,omitempty"`
	UpdatedAt       *time.Time             `json:"updated_at,omitempty"`
}

// WorkerCreateRequest represents the request to create/register a worker
type WorkerCreateRequest struct {
	EnvironmentName string                 `json:"environment_name"`
	Hostname        *string                `json:"hostname,omitempty"`
	WorkerMetadata  map[string]interface{} `json:"worker_metadata,omitempty"`
}
