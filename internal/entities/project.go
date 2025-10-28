package entities

import "time"

// ProjectStatus represents the status of a project
type ProjectStatus string

const (
	ProjectStatusActive   ProjectStatus = "active"
	ProjectStatusArchived ProjectStatus = "archived"
	ProjectStatusPaused   ProjectStatus = "paused"
)

// Project represents a project in the control plane
type Project struct {
	ID                    string                 `json:"id,omitempty"`
	OrganizationID        string                 `json:"organization_id,omitempty"`
	Name                  string                 `json:"name"`
	Key                   string                 `json:"key"`
	Description           *string                `json:"description,omitempty"`
	Goals                 *string                `json:"goals,omitempty"`
	Settings              map[string]interface{} `json:"settings,omitempty"`
	Status                ProjectStatus          `json:"status,omitempty"`
	Visibility            string                 `json:"visibility,omitempty"`
	OwnerID               *string                `json:"owner_id,omitempty"`
	OwnerEmail            *string                `json:"owner_email,omitempty"`
	RestrictToEnvironment bool                   `json:"restrict_to_environment"`
	PolicyIDs             []string               `json:"policy_ids,omitempty"`
	DefaultModel          *string                `json:"default_model,omitempty"`
	CreatedAt             *time.Time             `json:"created_at,omitempty"`
	UpdatedAt             *time.Time             `json:"updated_at,omitempty"`
	ArchivedAt            *time.Time             `json:"archived_at,omitempty"`
	AgentCount            int                    `json:"agent_count,omitempty"`
	TeamCount             int                    `json:"team_count,omitempty"`
}

// ProjectCreateRequest represents the request to create a project
type ProjectCreateRequest struct {
	Name                  string                 `json:"name"`
	Key                   string                 `json:"key"`
	Description           *string                `json:"description,omitempty"`
	Goals                 *string                `json:"goals,omitempty"`
	Settings              map[string]interface{} `json:"settings,omitempty"`
	Visibility            string                 `json:"visibility,omitempty"`
	RestrictToEnvironment bool                   `json:"restrict_to_environment"`
	PolicyIDs             []string               `json:"policy_ids,omitempty"`
	DefaultModel          *string                `json:"default_model,omitempty"`
}

// ProjectUpdateRequest represents the request to update a project
type ProjectUpdateRequest struct {
	Name                  *string                `json:"name,omitempty"`
	Key                   *string                `json:"key,omitempty"`
	Description           *string                `json:"description,omitempty"`
	Goals                 *string                `json:"goals,omitempty"`
	Settings              map[string]interface{} `json:"settings,omitempty"`
	Status                *ProjectStatus         `json:"status,omitempty"`
	Visibility            *string                `json:"visibility,omitempty"`
	RestrictToEnvironment *bool                  `json:"restrict_to_environment,omitempty"`
	PolicyIDs             []string               `json:"policy_ids,omitempty"`
	DefaultModel          *string                `json:"default_model,omitempty"`
}
