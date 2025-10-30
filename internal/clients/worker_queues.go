package clients

import (
	"fmt"
	"net/http"

	"kubiya-control-plane/internal/entities"
)

// CreateWorkerQueue creates a new worker queue
func (c *Client) CreateWorkerQueue(environmentID string, req *entities.WorkerQueueCreateRequest) (*entities.WorkerQueue, error) {
	resp, err := c.DoRequest(http.MethodPost, fmt.Sprintf("/api/v1/environments/%s/worker-queues", environmentID), req)
	if err != nil {
		return nil, err
	}

	var queue entities.WorkerQueue
	if err := ParseResponse(resp, &queue); err != nil {
		return nil, err
	}

	return &queue, nil
}

// GetWorkerQueue retrieves a worker queue by ID
func (c *Client) GetWorkerQueue(queueID string) (*entities.WorkerQueue, error) {
	resp, err := c.DoRequest(http.MethodGet, fmt.Sprintf("/api/v1/worker-queues/%s", queueID), nil)
	if err != nil {
		return nil, err
	}

	var queue entities.WorkerQueue
	if err := ParseResponse(resp, &queue); err != nil {
		return nil, err
	}

	return &queue, nil
}

// UpdateWorkerQueue updates a worker queue
func (c *Client) UpdateWorkerQueue(queueID string, req *entities.WorkerQueueUpdateRequest) (*entities.WorkerQueue, error) {
	resp, err := c.DoRequest(http.MethodPatch, fmt.Sprintf("/api/v1/worker-queues/%s", queueID), req)
	if err != nil {
		return nil, err
	}

	var queue entities.WorkerQueue
	if err := ParseResponse(resp, &queue); err != nil {
		return nil, err
	}

	return &queue, nil
}

// DeleteWorkerQueue deletes a worker queue
func (c *Client) DeleteWorkerQueue(queueID string) error {
	resp, err := c.DoRequest(http.MethodDelete, fmt.Sprintf("/api/v1/worker-queues/%s", queueID), nil)
	if err != nil {
		return err
	}

	return ParseResponse(resp, nil)
}

// ListWorkerQueues lists all worker queues in an environment
func (c *Client) ListWorkerQueues(environmentID string) ([]*entities.WorkerQueue, error) {
	resp, err := c.DoRequest(http.MethodGet, fmt.Sprintf("/api/v1/environments/%s/worker-queues", environmentID), nil)
	if err != nil {
		return nil, err
	}

	var queues []*entities.WorkerQueue
	if err := ParseResponse(resp, &queues); err != nil {
		return nil, err
	}

	return queues, nil
}
