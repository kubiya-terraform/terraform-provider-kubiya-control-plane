package clients

import (
	"fmt"
	"net/http"

	"terraform-provider-kubiya-control-plane/internal/entities"
)

// CreateSkill creates a new skill (uses /api/v1/skills endpoint)
func (c *Client) CreateSkill(req *entities.SkillCreateRequest) (*entities.Skill, error) {
	resp, err := c.DoRequest(http.MethodPost, "/api/v1/skills", req)
	if err != nil {
		return nil, err
	}

	var skill entities.Skill
	if err := ParseResponse(resp, &skill); err != nil {
		return nil, err
	}

	return &skill, nil
}

// GetSkill retrieves a skill by ID (uses /api/v1/skills endpoint)
func (c *Client) GetSkill(id string) (*entities.Skill, error) {
	resp, err := c.DoRequest(http.MethodGet, fmt.Sprintf("/api/v1/skills/%s", id), nil)
	if err != nil {
		return nil, err
	}

	var skill entities.Skill
	if err := ParseResponse(resp, &skill); err != nil {
		return nil, err
	}

	return &skill, nil
}

// UpdateSkill updates an existing skill (uses /api/v1/skills endpoint)
func (c *Client) UpdateSkill(id string, req *entities.SkillUpdateRequest) (*entities.Skill, error) {
	resp, err := c.DoRequest(http.MethodPatch, fmt.Sprintf("/api/v1/skills/%s", id), req)
	if err != nil {
		return nil, err
	}

	var skill entities.Skill
	if err := ParseResponse(resp, &skill); err != nil {
		return nil, err
	}

	return &skill, nil
}

// DeleteSkill deletes a skill (uses /api/v1/skills endpoint)
func (c *Client) DeleteSkill(id string) error {
	resp, err := c.DoRequest(http.MethodDelete, fmt.Sprintf("/api/v1/skills/%s", id), nil)
	if err != nil {
		return err
	}

	return ParseResponse(resp, nil)
}

// ListSkills lists all skills (uses /api/v1/skills endpoint)
func (c *Client) ListSkills() ([]*entities.Skill, error) {
	resp, err := c.DoRequest(http.MethodGet, "/api/v1/skills", nil)
	if err != nil {
		return nil, err
	}

	var skills []*entities.Skill
	if err := ParseResponse(resp, &skills); err != nil {
		return nil, err
	}

	return skills, nil
}
