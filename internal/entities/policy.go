package entities

import "time"

// PolicyType represents the type of policy
type PolicyType string

const (
	PolicyTypeRego PolicyType = "rego"
	PolicyTypeJSON PolicyType = "json"
)

// Policy represents an OPA policy in the control plane
type Policy struct {
	ID             string     `json:"id,omitempty"`
	OrganizationID string     `json:"organization_id,omitempty"`
	Name           string     `json:"name"`
	Description    *string    `json:"description,omitempty"`
	PolicyContent  string     `json:"policy_content"`
	PolicyType     PolicyType `json:"policy_type,omitempty"`
	Enabled        bool       `json:"enabled"`
	Tags           []string   `json:"tags,omitempty"`
	Version        int64      `json:"version,omitempty"`
	CreatedAt      *time.Time `json:"created_at,omitempty"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty"`
}

// PolicyCreateRequest represents the request to create a policy
type PolicyCreateRequest struct {
	Name          string     `json:"name"`
	Description   *string    `json:"description,omitempty"`
	PolicyContent string     `json:"policy_content"`
	PolicyType    PolicyType `json:"policy_type,omitempty"`
	Enabled       bool       `json:"enabled"`
	Tags          []string   `json:"tags,omitempty"`
}

// PolicyUpdateRequest represents the request to update a policy
type PolicyUpdateRequest struct {
	Name          *string    `json:"name,omitempty"`
	Description   *string    `json:"description,omitempty"`
	PolicyContent *string    `json:"policy_content,omitempty"`
	Enabled       *bool      `json:"enabled,omitempty"`
	Tags          []string   `json:"tags,omitempty"`
}
