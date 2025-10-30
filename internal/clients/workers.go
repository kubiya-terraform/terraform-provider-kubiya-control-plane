package clients

import (
	"encoding/json"
	"fmt"
	"io"
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
	defer func() { _ = resp.Body.Close() }()

	// Read the body once
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(bodyBytes))
	}

	// Try parsing as a single worker object first
	var worker entities.Worker
	if err := json.Unmarshal(bodyBytes, &worker); err == nil {
		// Successfully parsed as single object
		if worker.ID != "" || worker.WorkerID != "" {
			return &worker, nil
		}
	}

	// Try parsing as an array of workers
	var workers []entities.Worker
	if err := json.Unmarshal(bodyBytes, &workers); err == nil {
		if len(workers) == 0 {
			return nil, fmt.Errorf("worker not found")
		}
		return &workers[0], nil
	}

	// Try parsing as a runner/queue response with nested workers array
	var queueResp struct {
		Workers []entities.Worker `json:"workers"`
	}
	if err := json.Unmarshal(bodyBytes, &queueResp); err == nil {
		if len(queueResp.Workers) == 0 {
			return nil, fmt.Errorf("worker not found")
		}
		return &queueResp.Workers[0], nil
	}

	return nil, fmt.Errorf("failed to parse worker response. Body: %s", string(bodyBytes))
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
