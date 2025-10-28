package clients

import (
	"fmt"
	"net/http"

	"kubiya-control-plane/internal/entities"
)

// CreateWorker creates a new worker
func (c *Client) CreateWorker(req *entities.WorkerCreateRequest) (*entities.Worker, error) {
	resp, err := c.DoRequest(http.MethodPost, "/api/v1/workers/register", req)
	if err != nil {
		return nil, err
	}

	var worker entities.Worker
	if err := ParseResponse(resp, &worker); err != nil {
		return nil, err
	}

	return &worker, nil
}

// GetWorker retrieves a worker by ID
func (c *Client) GetWorker(id string) (*entities.Worker, error) {
	resp, err := c.DoRequest(http.MethodGet, fmt.Sprintf("/api/v1/workers/%s", id), nil)
	if err != nil {
		return nil, err
	}

	var worker entities.Worker
	if err := ParseResponse(resp, &worker); err != nil {
		return nil, err
	}

	return &worker, nil
}

// DeleteWorker deletes a worker
func (c *Client) DeleteWorker(id string) error {
	resp, err := c.DoRequest(http.MethodDelete, fmt.Sprintf("/api/v1/workers/%s", id), nil)
	if err != nil {
		return err
	}

	return ParseResponse(resp, nil)
}

// ListWorkers lists all workers
func (c *Client) ListWorkers() ([]*entities.Worker, error) {
	resp, err := c.DoRequest(http.MethodGet, "/api/v1/workers", nil)
	if err != nil {
		return nil, err
	}

	var workers []*entities.Worker
	if err := ParseResponse(resp, &workers); err != nil {
		return nil, err
	}

	return workers, nil
}
