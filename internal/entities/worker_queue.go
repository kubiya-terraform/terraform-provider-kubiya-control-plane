package entities

import "time"

// WorkerQueueStatus represents the status of a worker queue
type WorkerQueueStatus string

const (
	WorkerQueueStatusActive   WorkerQueueStatus = "active"
	WorkerQueueStatusInactive WorkerQueueStatus = "inactive"
	WorkerQueueStatusPaused   WorkerQueueStatus = "paused"
)

// WorkerQueue represents a worker queue in the control plane
type WorkerQueue struct {
	ID                string                 `json:"id,omitempty"`
	OrganizationID    string                 `json:"organization_id,omitempty"`
	EnvironmentID     string                 `json:"environment_id"`
	Name              string                 `json:"name"`
	DisplayName       *string                `json:"display_name,omitempty"`
	Description       *string                `json:"description,omitempty"`
	Status            WorkerQueueStatus      `json:"status,omitempty"`
	MaxWorkers        *int                   `json:"max_workers,omitempty"`
	HeartbeatInterval int                    `json:"heartbeat_interval"`
	Tags              []string               `json:"tags,omitempty"`
	Settings          map[string]interface{} `json:"settings,omitempty"`
	CreatedAt         *time.Time             `json:"created_at,omitempty"`
	UpdatedAt         *time.Time             `json:"updated_at,omitempty"`
	CreatedBy         *string                `json:"created_by,omitempty"`
	// Computed fields
	ActiveWorkers int    `json:"active_workers,omitempty"`
	TaskQueueName string `json:"task_queue_name,omitempty"`
}

// WorkerQueueCreateRequest represents the request to create a worker queue
type WorkerQueueCreateRequest struct {
	Name              string                 `json:"name"`
	DisplayName       *string                `json:"display_name,omitempty"`
	Description       *string                `json:"description,omitempty"`
	MaxWorkers        *int                   `json:"max_workers,omitempty"`
	HeartbeatInterval int                    `json:"heartbeat_interval"`
	Tags              []string               `json:"tags,omitempty"`
	Settings          map[string]interface{} `json:"settings,omitempty"`
}

// WorkerQueueUpdateRequest represents the request to update a worker queue
type WorkerQueueUpdateRequest struct {
	Name              *string                `json:"name,omitempty"`
	DisplayName       *string                `json:"display_name,omitempty"`
	Description       *string                `json:"description,omitempty"`
	Status            *string                `json:"status,omitempty"`
	MaxWorkers        *int                   `json:"max_workers,omitempty"`
	HeartbeatInterval *int                   `json:"heartbeat_interval,omitempty"`
	Tags              []string               `json:"tags,omitempty"`
	Settings          map[string]interface{} `json:"settings,omitempty"`
}
