package clients

import (
	"fmt"
	"net/http"

	"kubiya-control-plane/internal/entities"
)

// CreateProject creates a new project
func (c *Client) CreateProject(req *entities.ProjectCreateRequest) (*entities.Project, error) {
	resp, err := c.DoRequest(http.MethodPost, "/api/v1/projects", req)
	if err != nil {
		return nil, err
	}

	var project entities.Project
	if err := ParseResponse(resp, &project); err != nil {
		return nil, err
	}

	return &project, nil
}

// GetProject retrieves a project by ID
func (c *Client) GetProject(id string) (*entities.Project, error) {
	resp, err := c.DoRequest(http.MethodGet, fmt.Sprintf("/api/v1/projects/%s", id), nil)
	if err != nil {
		return nil, err
	}

	var project entities.Project
	if err := ParseResponse(resp, &project); err != nil {
		return nil, err
	}

	return &project, nil
}

// UpdateProject updates an existing project
func (c *Client) UpdateProject(id string, req *entities.ProjectUpdateRequest) (*entities.Project, error) {
	resp, err := c.DoRequest(http.MethodPatch, fmt.Sprintf("/api/v1/projects/%s", id), req)
	if err != nil {
		return nil, err
	}

	var project entities.Project
	if err := ParseResponse(resp, &project); err != nil {
		return nil, err
	}

	return &project, nil
}

// DeleteProject deletes a project
func (c *Client) DeleteProject(id string) error {
	resp, err := c.DoRequest(http.MethodDelete, fmt.Sprintf("/api/v1/projects/%s", id), nil)
	if err != nil {
		return err
	}

	return ParseResponse(resp, nil)
}

// ListProjects lists all projects
func (c *Client) ListProjects() ([]*entities.Project, error) {
	resp, err := c.DoRequest(http.MethodGet, "/api/v1/projects", nil)
	if err != nil {
		return nil, err
	}

	var projects []*entities.Project
	if err := ParseResponse(resp, &projects); err != nil {
		return nil, err
	}

	return projects, nil
}
