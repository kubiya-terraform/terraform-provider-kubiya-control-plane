package clients

import (
	"fmt"
	"net/http"

	"terraform-provider-kubiya-control-plane/internal/entities"
)

// CreatePolicy creates a new policy
func (c *Client) CreatePolicy(req *entities.PolicyCreateRequest) (*entities.Policy, error) {
	resp, err := c.DoRequest(http.MethodPost, "/api/v1/policies", req)
	if err != nil {
		return nil, err
	}

	var policy entities.Policy
	if err := ParseResponse(resp, &policy); err != nil {
		return nil, err
	}

	return &policy, nil
}

// GetPolicy retrieves a policy by ID
func (c *Client) GetPolicy(id string) (*entities.Policy, error) {
	resp, err := c.DoRequest(http.MethodGet, fmt.Sprintf("/api/v1/policies/%s", id), nil)
	if err != nil {
		return nil, err
	}

	var policy entities.Policy
	if err := ParseResponse(resp, &policy); err != nil {
		return nil, err
	}

	return &policy, nil
}

// UpdatePolicy updates an existing policy
func (c *Client) UpdatePolicy(id string, req *entities.PolicyUpdateRequest) (*entities.Policy, error) {
	resp, err := c.DoRequest(http.MethodPut, fmt.Sprintf("/api/v1/policies/%s", id), req)
	if err != nil {
		return nil, err
	}

	var policy entities.Policy
	if err := ParseResponse(resp, &policy); err != nil {
		return nil, err
	}

	return &policy, nil
}

// DeletePolicy deletes a policy
func (c *Client) DeletePolicy(id string) error {
	resp, err := c.DoRequest(http.MethodDelete, fmt.Sprintf("/api/v1/policies/%s", id), nil)
	if err != nil {
		return err
	}

	return ParseResponse(resp, nil)
}

// ListPolicies lists all policies
func (c *Client) ListPolicies() ([]*entities.Policy, error) {
	resp, err := c.DoRequest(http.MethodGet, "/api/v1/policies", nil)
	if err != nil {
		return nil, err
	}

	var policies []*entities.Policy
	if err := ParseResponse(resp, &policies); err != nil {
		return nil, err
	}

	return policies, nil
}
