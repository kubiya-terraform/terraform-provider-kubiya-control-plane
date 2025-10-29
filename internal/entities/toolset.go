package entities

import "time"

// ToolSetType represents the type of toolset
type ToolSetType string

const (
	ToolSetTypeFileSystem     ToolSetType = "file_system"
	ToolSetTypeShell          ToolSetType = "shell"
	ToolSetTypeDocker         ToolSetType = "docker"
	ToolSetTypePython         ToolSetType = "python"
	ToolSetTypeFileGeneration ToolSetType = "file_generation"
	ToolSetTypeCustom         ToolSetType = "custom"
)

// ToolSet represents a toolset in the control plane
type ToolSet struct {
	ID             string                 `json:"id,omitempty"`
	OrganizationID string                 `json:"organization_id,omitempty"`
	Name           string                 `json:"name"`
	Type           ToolSetType            `json:"type"`
	Description    *string                `json:"description,omitempty"`
	Icon           string                 `json:"icon,omitempty"`
	Enabled        bool                   `json:"enabled"`
	Configuration  map[string]interface{} `json:"configuration,omitempty"`
	CreatedAt      *time.Time             `json:"created_at,omitempty"`
	UpdatedAt      *time.Time             `json:"updated_at,omitempty"`
}

// ToolSetCreateRequest represents the request to create a toolset
type ToolSetCreateRequest struct {
	Name          string                 `json:"name"`
	Type          ToolSetType            `json:"type"`
	Description   *string                `json:"description,omitempty"`
	Icon          string                 `json:"icon,omitempty"`
	Enabled       bool                   `json:"enabled"`
	Configuration map[string]interface{} `json:"configuration,omitempty"`
}

// ToolSetUpdateRequest represents the request to update a toolset
type ToolSetUpdateRequest struct {
	Name          *string                `json:"name,omitempty"`
	Description   *string                `json:"description,omitempty"`
	Icon          *string                `json:"icon,omitempty"`
	Enabled       *bool                  `json:"enabled,omitempty"`
	Configuration map[string]interface{} `json:"configuration,omitempty"`
}
