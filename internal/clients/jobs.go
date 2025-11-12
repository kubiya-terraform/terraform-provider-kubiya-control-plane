package clients

import (
	"fmt"
	"net/http"

	"kubiya-control-plane/internal/entities"
)

// CreateJob creates a new job
func (c *Client) CreateJob(req *entities.JobCreateRequest) (*entities.Job, error) {
	resp, err := c.DoRequest(http.MethodPost, "/api/v1/jobs", req)
	if err != nil {
		return nil, err
	}

	var job entities.Job
	if err := ParseResponse(resp, &job); err != nil {
		return nil, err
	}

	return &job, nil
}

// GetJob retrieves a job by ID
func (c *Client) GetJob(id string) (*entities.Job, error) {
	resp, err := c.DoRequest(http.MethodGet, fmt.Sprintf("/api/v1/jobs/%s", id), nil)
	if err != nil {
		return nil, err
	}

	var job entities.Job
	if err := ParseResponse(resp, &job); err != nil {
		return nil, err
	}

	return &job, nil
}

// UpdateJob updates an existing job
func (c *Client) UpdateJob(id string, req *entities.JobUpdateRequest) (*entities.Job, error) {
	resp, err := c.DoRequest(http.MethodPatch, fmt.Sprintf("/api/v1/jobs/%s", id), req)
	if err != nil {
		return nil, err
	}

	var job entities.Job
	if err := ParseResponse(resp, &job); err != nil {
		return nil, err
	}

	return &job, nil
}

// DeleteJob deletes a job
func (c *Client) DeleteJob(id string) error {
	resp, err := c.DoRequest(http.MethodDelete, fmt.Sprintf("/api/v1/jobs/%s", id), nil)
	if err != nil {
		return err
	}

	return ParseResponse(resp, nil)
}

// ListJobs lists all jobs
func (c *Client) ListJobs() ([]*entities.Job, error) {
	resp, err := c.DoRequest(http.MethodGet, "/api/v1/jobs", nil)
	if err != nil {
		return nil, err
	}

	var jobs []*entities.Job
	if err := ParseResponse(resp, &jobs); err != nil {
		return nil, err
	}

	return jobs, nil
}

// EnableJob enables a job
func (c *Client) EnableJob(id string) (*entities.Job, error) {
	resp, err := c.DoRequest(http.MethodPost, fmt.Sprintf("/api/v1/jobs/%s/enable", id), nil)
	if err != nil {
		return nil, err
	}

	var job entities.Job
	if err := ParseResponse(resp, &job); err != nil {
		return nil, err
	}

	return &job, nil
}

// DisableJob disables a job
func (c *Client) DisableJob(id string) (*entities.Job, error) {
	resp, err := c.DoRequest(http.MethodPost, fmt.Sprintf("/api/v1/jobs/%s/disable", id), nil)
	if err != nil {
		return nil, err
	}

	var job entities.Job
	if err := ParseResponse(resp, &job); err != nil {
		return nil, err
	}

	return &job, nil
}
