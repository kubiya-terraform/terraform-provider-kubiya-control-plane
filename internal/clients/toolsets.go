package clients

import (
	"fmt"
	"net/http"

	"kubiya-control-plane/internal/entities"
)

// CreateToolSet creates a new toolset
func (c *Client) CreateToolSet(req *entities.ToolSetCreateRequest) (*entities.ToolSet, error) {
	resp, err := c.DoRequest(http.MethodPost, "/api/v1/toolsets", req)
	if err != nil {
		return nil, err
	}

	var toolset entities.ToolSet
	if err := ParseResponse(resp, &toolset); err != nil {
		return nil, err
	}

	return &toolset, nil
}

// GetToolSet retrieves a toolset by ID
func (c *Client) GetToolSet(id string) (*entities.ToolSet, error) {
	resp, err := c.DoRequest(http.MethodGet, fmt.Sprintf("/api/v1/toolsets/%s", id), nil)
	if err != nil {
		return nil, err
	}

	var toolset entities.ToolSet
	if err := ParseResponse(resp, &toolset); err != nil {
		return nil, err
	}

	return &toolset, nil
}

// UpdateToolSet updates an existing toolset
func (c *Client) UpdateToolSet(id string, req *entities.ToolSetUpdateRequest) (*entities.ToolSet, error) {
	resp, err := c.DoRequest(http.MethodPatch, fmt.Sprintf("/api/v1/toolsets/%s", id), req)
	if err != nil {
		return nil, err
	}

	var toolset entities.ToolSet
	if err := ParseResponse(resp, &toolset); err != nil {
		return nil, err
	}

	return &toolset, nil
}

// DeleteToolSet deletes a toolset
func (c *Client) DeleteToolSet(id string) error {
	resp, err := c.DoRequest(http.MethodDelete, fmt.Sprintf("/api/v1/toolsets/%s", id), nil)
	if err != nil {
		return err
	}

	return ParseResponse(resp, nil)
}

// ListToolSets lists all toolsets
func (c *Client) ListToolSets() ([]*entities.ToolSet, error) {
	resp, err := c.DoRequest(http.MethodGet, "/api/v1/toolsets", nil)
	if err != nil {
		return nil, err
	}

	var toolsets []*entities.ToolSet
	if err := ParseResponse(resp, &toolsets); err != nil {
		return nil, err
	}

	return toolsets, nil
}
