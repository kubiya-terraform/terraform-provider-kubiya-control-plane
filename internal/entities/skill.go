package entities

import (
	"strings"
	"time"
)

// FlexibleTime is a custom time type that can parse timestamps with or without timezone
type FlexibleTime struct {
	time.Time
}

// UnmarshalJSON handles parsing timestamps in multiple formats
func (ft *FlexibleTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "null" || s == "" {
		return nil
	}

	// Try parsing with timezone first (RFC3339)
	t, err := time.Parse(time.RFC3339, s)
	if err == nil {
		ft.Time = t
		return nil
	}

	// Try parsing without timezone, assume UTC
	layouts := []string{
		"2006-01-02T15:04:05.999999",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05.999999",
		"2006-01-02 15:04:05",
	}

	for _, layout := range layouts {
		t, err := time.Parse(layout, s)
		if err == nil {
			ft.Time = t.UTC()
			return nil
		}
	}

	return err
}

// SkillType represents the type of skill
type SkillType string

const (
	SkillTypeFileSystem     SkillType = "file_system"
	SkillTypeShell          SkillType = "shell"
	SkillTypeDocker         SkillType = "docker"
	SkillTypePython         SkillType = "python"
	SkillTypeFileGeneration SkillType = "file_generation"
	SkillTypeCustom         SkillType = "custom"
)

// Skill represents a skill in the control plane
type Skill struct {
	ID             string                 `json:"id,omitempty"`
	OrganizationID string                 `json:"organization_id,omitempty"`
	Name           string                 `json:"name"`
	Type           SkillType              `json:"type"`
	Description    *string                `json:"description,omitempty"`
	Icon           string                 `json:"icon,omitempty"`
	Enabled        bool                   `json:"enabled"`
	Configuration  map[string]interface{} `json:"configuration,omitempty"`
	CreatedAt      *FlexibleTime          `json:"created_at,omitempty"`
	UpdatedAt      *FlexibleTime          `json:"updated_at,omitempty"`
}

// SkillCreateRequest represents the request to create a skill
type SkillCreateRequest struct {
	Name          string                 `json:"name"`
	Type          SkillType              `json:"type"`
	Description   *string                `json:"description,omitempty"`
	Icon          string                 `json:"icon,omitempty"`
	Enabled       bool                   `json:"enabled"`
	Configuration map[string]interface{} `json:"configuration,omitempty"`
}

// SkillUpdateRequest represents the request to update a skill
type SkillUpdateRequest struct {
	Name          *string                `json:"name,omitempty"`
	Description   *string                `json:"description,omitempty"`
	Icon          *string                `json:"icon,omitempty"`
	Enabled       *bool                  `json:"enabled,omitempty"`
	Configuration map[string]interface{} `json:"configuration,omitempty"`
}
