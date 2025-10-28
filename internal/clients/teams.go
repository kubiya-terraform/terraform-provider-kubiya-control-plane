package clients

import (
	"fmt"
	"net/http"

	"kubiya-control-plane/internal/entities"
)

// CreateTeam creates a new team
func (c *Client) CreateTeam(req *entities.TeamCreateRequest) (*entities.Team, error) {
	resp, err := c.DoRequest(http.MethodPost, "/api/v1/teams", req)
	if err != nil {
		return nil, err
	}

	var team entities.Team
	if err := ParseResponse(resp, &team); err != nil {
		return nil, err
	}

	return &team, nil
}

// GetTeam retrieves a team by ID
func (c *Client) GetTeam(id string) (*entities.Team, error) {
	resp, err := c.DoRequest(http.MethodGet, fmt.Sprintf("/api/v1/teams/%s", id), nil)
	if err != nil {
		return nil, err
	}

	var team entities.Team
	if err := ParseResponse(resp, &team); err != nil {
		return nil, err
	}

	return &team, nil
}

// UpdateTeam updates an existing team
func (c *Client) UpdateTeam(id string, req *entities.TeamUpdateRequest) (*entities.Team, error) {
	resp, err := c.DoRequest(http.MethodPut, fmt.Sprintf("/api/v1/teams/%s", id), req)
	if err != nil {
		return nil, err
	}

	var team entities.Team
	if err := ParseResponse(resp, &team); err != nil {
		return nil, err
	}

	return &team, nil
}

// DeleteTeam deletes a team
func (c *Client) DeleteTeam(id string) error {
	resp, err := c.DoRequest(http.MethodDelete, fmt.Sprintf("/api/v1/teams/%s", id), nil)
	if err != nil {
		return err
	}

	return ParseResponse(resp, nil)
}

// ListTeams lists all teams
func (c *Client) ListTeams() ([]*entities.Team, error) {
	resp, err := c.DoRequest(http.MethodGet, "/api/v1/teams", nil)
	if err != nil {
		return nil, err
	}

	var teams []*entities.Team
	if err := ParseResponse(resp, &teams); err != nil {
		return nil, err
	}

	return teams, nil
}
