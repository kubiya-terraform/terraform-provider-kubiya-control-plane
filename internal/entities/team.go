package entities

import "time"

// TeamStatus represents the status of a team
type TeamStatus string

const (
	TeamStatusActive   TeamStatus = "active"
	TeamStatusInactive TeamStatus = "inactive"
	TeamStatusArchived TeamStatus = "archived"
)

// Team represents a team in the control plane
type Team struct {
	ID                   string                 `json:"id,omitempty"`
	OrganizationID       string                 `json:"organization_id,omitempty"`
	Name                 string                 `json:"name"`
	Description          *string                `json:"description,omitempty"`
	Status               TeamStatus             `json:"status,omitempty"`
	Configuration        map[string]interface{} `json:"configuration,omitempty"`
	SkillIDs             []string               `json:"skill_ids,omitempty"`
	ExecutionEnvironment map[string]interface{} `json:"execution_environment,omitempty"`
	CreatedAt            *time.Time             `json:"created_at,omitempty"`
	UpdatedAt            *time.Time             `json:"updated_at,omitempty"`
}

// TeamCreateRequest represents the request to create a team
type TeamCreateRequest struct {
	Name                 string                 `json:"name"`
	Description          *string                `json:"description,omitempty"`
	Configuration        map[string]interface{} `json:"configuration,omitempty"`
	SkillIDs             []string               `json:"skill_ids,omitempty"`
	ExecutionEnvironment map[string]interface{} `json:"execution_environment,omitempty"`
}

// TeamUpdateRequest represents the request to update a team
type TeamUpdateRequest struct {
	Name                 *string                `json:"name,omitempty"`
	Description          *string                `json:"description,omitempty"`
	Status               *TeamStatus            `json:"status,omitempty"`
	Configuration        map[string]interface{} `json:"configuration,omitempty"`
	SkillIDs             []string               `json:"skill_ids,omitempty"`
	ExecutionEnvironment map[string]interface{} `json:"execution_environment,omitempty"`
}
