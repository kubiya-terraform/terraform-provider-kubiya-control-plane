package test

import (
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestWorkerQueueBasic tests the basic worker queue resource lifecycle using the example
func TestWorkerQueueBasic(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Fatal("KUBIYA_CONTROL_PLANE_API_KEY environment variable is not set")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../../examples/worker_queue",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	queueID := terraform.Output(t, terraformOptions, "queue_id")
	t.Logf("Created worker queue with ID: %s", queueID)
}

// TestWorkerQueueMinimal tests minimal worker queue configuration
func TestWorkerQueueMinimal(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Skip("KUBIYA_CONTROL_PLANE_API_KEY not set, skipping testdata test")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../../testdata/workers/minimal",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	workerID := terraform.Output(t, terraformOptions, "worker_id")
	require.NotEmpty(t, workerID)

	workerName := terraform.Output(t, terraformOptions, "worker_name")
	assert.Equal(t, "test-worker-queue-minimal", workerName)

	workerEnvironmentID := terraform.Output(t, terraformOptions, "worker_environment_id")
	assert.NotEmpty(t, workerEnvironmentID)

	// Verify data source
	dataWorkerName := terraform.Output(t, terraformOptions, "data_worker_name")
	assert.Equal(t, workerName, dataWorkerName)

	t.Logf("✓ Minimal worker queue test passed: ID=%s, Name=%s", workerID, workerName)
}

// TestWorkerQueueFull tests worker queue with all optional fields
func TestWorkerQueueFull(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Skip("KUBIYA_CONTROL_PLANE_API_KEY not set, skipping testdata test")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../../testdata/workers/full",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	workerID := terraform.Output(t, terraformOptions, "worker_id")
	require.NotEmpty(t, workerID)

	workerName := terraform.Output(t, terraformOptions, "worker_name")
	assert.Equal(t, "test-worker-queue-full", workerName)

	workerDescription := terraform.Output(t, terraformOptions, "worker_description")
	assert.Contains(t, workerDescription, "Comprehensive")

	maxWorkers := terraform.Output(t, terraformOptions, "worker_max_workers")
	assert.Equal(t, "10", maxWorkers)

	heartbeat := terraform.Output(t, terraformOptions, "worker_heartbeat_interval")
	assert.Equal(t, "60", heartbeat)

	// Verify data source
	dataDescription := terraform.Output(t, terraformOptions, "data_worker_description")
	assert.Equal(t, workerDescription, dataDescription)

	t.Logf("✓ Full worker queue test passed: ID=%s", workerID)
}

// TestWorkerQueueComprehensive tests all worker queue resource scenarios together
func TestWorkerQueueComprehensive(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Fatal("KUBIYA_CONTROL_PLANE_API_KEY environment variable is not set")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../../testdata/workers/comprehensive",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	// Test minimal worker queue
	minimalID := terraform.Output(t, terraformOptions, "minimal_worker_queue_id")
	require.NotEmpty(t, minimalID)
	minimalName := terraform.Output(t, terraformOptions, "minimal_worker_queue_name")
	assert.Equal(t, "test-worker-queue-minimal", minimalName)

	// Test full worker queue
	fullID := terraform.Output(t, terraformOptions, "full_worker_queue_id")
	require.NotEmpty(t, fullID)
	fullMaxWorkers := terraform.Output(t, terraformOptions, "full_worker_queue_max_workers")
	assert.Equal(t, "10", fullMaxWorkers)

	t.Logf("✓ Comprehensive worker queue tests passed")
}
