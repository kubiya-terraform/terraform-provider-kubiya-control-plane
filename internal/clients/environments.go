package clients

import (
	"fmt"
	"net/http"

	"terraform-provider-kubiya-control-plane/internal/entities"
)

// CreateEnvironment creates a new environment
func (c *Client) CreateEnvironment(req *entities.EnvironmentCreateRequest) (*entities.Environment, error) {
	resp, err := c.DoRequest(http.MethodPost, "/api/v1/environments", req)
	if err != nil {
		return nil, err
	}

	var environment entities.Environment
	if err := ParseResponse(resp, &environment); err != nil {
		return nil, err
	}

	return &environment, nil
}

// GetEnvironment retrieves an environment by ID
func (c *Client) GetEnvironment(id string) (*entities.Environment, error) {
	resp, err := c.DoRequest(http.MethodGet, fmt.Sprintf("/api/v1/environments/%s", id), nil)
	if err != nil {
		return nil, err
	}

	var environment entities.Environment
	if err := ParseResponse(resp, &environment); err != nil {
		return nil, err
	}

	return &environment, nil
}

// UpdateEnvironment updates an existing environment
func (c *Client) UpdateEnvironment(id string, req *entities.EnvironmentUpdateRequest) (*entities.Environment, error) {
	resp, err := c.DoRequest(http.MethodPatch, fmt.Sprintf("/api/v1/environments/%s", id), req)
	if err != nil {
		return nil, err
	}

	var environment entities.Environment
	if err := ParseResponse(resp, &environment); err != nil {
		return nil, err
	}

	return &environment, nil
}

// DeleteEnvironment deletes an environment
func (c *Client) DeleteEnvironment(id string) error {
	resp, err := c.DoRequest(http.MethodDelete, fmt.Sprintf("/api/v1/environments/%s", id), nil)
	if err != nil {
		return err
	}

	return ParseResponse(resp, nil)
}

// ListEnvironments lists all environments
func (c *Client) ListEnvironments() ([]*entities.Environment, error) {
	resp, err := c.DoRequest(http.MethodGet, "/api/v1/environments", nil)
	if err != nil {
		return nil, err
	}

	var environments []*entities.Environment
	if err := ParseResponse(resp, &environments); err != nil {
		return nil, err
	}

	return environments, nil
}
