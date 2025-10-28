package clients

import (
	"fmt"
	"net/http"

	"kubiya-control-plane/internal/entities"
)

// CreateAgent creates a new agent
func (c *Client) CreateAgent(req *entities.AgentCreateRequest) (*entities.Agent, error) {
	resp, err := c.DoRequest(http.MethodPost, "/api/v1/agents", req)
	if err != nil {
		return nil, err
	}

	var agent entities.Agent
	if err := ParseResponse(resp, &agent); err != nil {
		return nil, err
	}

	return &agent, nil
}

// GetAgent retrieves an agent by ID
func (c *Client) GetAgent(id string) (*entities.Agent, error) {
	resp, err := c.DoRequest(http.MethodGet, fmt.Sprintf("/api/v1/agents/%s", id), nil)
	if err != nil {
		return nil, err
	}

	var agent entities.Agent
	if err := ParseResponse(resp, &agent); err != nil {
		return nil, err
	}

	return &agent, nil
}

// UpdateAgent updates an existing agent
func (c *Client) UpdateAgent(id string, req *entities.AgentUpdateRequest) (*entities.Agent, error) {
	resp, err := c.DoRequest(http.MethodPut, fmt.Sprintf("/api/v1/agents/%s", id), req)
	if err != nil {
		return nil, err
	}

	var agent entities.Agent
	if err := ParseResponse(resp, &agent); err != nil {
		return nil, err
	}

	return &agent, nil
}

// DeleteAgent deletes an agent
func (c *Client) DeleteAgent(id string) error {
	resp, err := c.DoRequest(http.MethodDelete, fmt.Sprintf("/api/v1/agents/%s", id), nil)
	if err != nil {
		return err
	}

	return ParseResponse(resp, nil)
}

// ListAgents lists all agents
func (c *Client) ListAgents() ([]*entities.Agent, error) {
	resp, err := c.DoRequest(http.MethodGet, "/api/v1/agents", nil)
	if err != nil {
		return nil, err
	}

	var agents []*entities.Agent
	if err := ParseResponse(resp, &agents); err != nil {
		return nil, err
	}

	return agents, nil
}
