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
			"TF_CLI_CONFIG_FILE":           os.Getenv("TF_CLI_CONFIG_FILE"),
			"HOME":                    os.Getenv("HOME"),
			"TF_SKIP_PROVIDER_VERIFY": "1",
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
			"TF_CLI_CONFIG_FILE":           os.Getenv("TF_CLI_CONFIG_FILE"),
			"HOME":                    os.Getenv("HOME"),
			"TF_SKIP_PROVIDER_VERIFY": "1",
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
			"TF_CLI_CONFIG_FILE":           os.Getenv("TF_CLI_CONFIG_FILE"),
			"HOME":                    os.Getenv("HOME"),
			"TF_SKIP_PROVIDER_VERIFY": "1",
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
			"TF_CLI_CONFIG_FILE":           os.Getenv("TF_CLI_CONFIG_FILE"),
			"HOME":                    os.Getenv("HOME"),
			"TF_SKIP_PROVIDER_VERIFY": "1",
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

// ============================================================================
// STATE MANAGEMENT TESTS - Update Lifecycle & Import
// ============================================================================

// TestWorkerQueueUpdate_Fields tests updating worker queue fields
func TestWorkerQueueUpdate_Fields(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Skip("KUBIYA_CONTROL_PLANE_API_KEY not set, skipping test")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../../testdata/worker_queues/update_fields",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
			"TF_CLI_CONFIG_FILE":           os.Getenv("TF_CLI_CONFIG_FILE"),
			"HOME":                    os.Getenv("HOME"),
			"TF_SKIP_PROVIDER_VERIFY": "1",
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	workerQueueID := terraform.Output(t, terraformOptions, "worker_queue_id")
	require.NotEmpty(t, workerQueueID)

	// Update description
	terraformOptions.Vars = map[string]interface{}{
		"description": "Updated worker queue description",
	}
	terraform.Apply(t, terraformOptions)

	// Verify in-place update
	updatedWorkerQueueID := terraform.Output(t, terraformOptions, "worker_queue_id")
	assert.Equal(t, workerQueueID, updatedWorkerQueueID)

	updatedDescription := terraform.Output(t, terraformOptions, "worker_queue_description")
	assert.Equal(t, "Updated worker queue description", updatedDescription)

	t.Logf("✓ Worker queue update test passed: ID=%s remained stable", workerQueueID)
}

// TestWorkerQueueImport tests importing an existing worker queue
func TestWorkerQueueImport(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Skip("KUBIYA_CONTROL_PLANE_API_KEY not set, skipping test")
	}

	createOptions := &terraform.Options{
		TerraformDir: "../../testdata/workers/minimal",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
			"TF_CLI_CONFIG_FILE":           os.Getenv("TF_CLI_CONFIG_FILE"),
			"HOME":                    os.Getenv("HOME"),
			"TF_SKIP_PROVIDER_VERIFY": "1",
		},
	}

	terraform.InitAndApply(t, createOptions)
	workerQueueID := terraform.Output(t, createOptions, "worker_queue_id")
	workerQueueName := terraform.Output(t, createOptions, "worker_queue_name")
	require.NotEmpty(t, workerQueueID)

	terraform.RunTerraformCommand(t, createOptions, "state", "rm", "controlplane_worker_queue.minimal")

	importOptions := &terraform.Options{
		TerraformDir: "../../testdata/worker_queues/import",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
			"TF_CLI_CONFIG_FILE":           os.Getenv("TF_CLI_CONFIG_FILE"),
			"HOME":                    os.Getenv("HOME"),
			"TF_SKIP_PROVIDER_VERIFY": "1",
		},
		Vars: map[string]interface{}{
			"worker_queue_id":   workerQueueID,
			"worker_queue_name": workerQueueName,
		},
	}

	defer terraform.Destroy(t, importOptions)

	terraform.Init(t, importOptions)
	terraform.RunTerraformCommand(t, importOptions, "import", "controlplane_worker_queue.imported", workerQueueID)

	importedID := terraform.Output(t, importOptions, "imported_worker_queue_id")
	assert.Equal(t, workerQueueID, importedID)

	importedName := terraform.Output(t, importOptions, "imported_worker_queue_name")
	assert.Equal(t, workerQueueName, importedName)

	t.Logf("✓ Worker queue import test passed: Successfully imported worker queue %s", workerQueueID)
}

// TestWorkerQueueStateRefresh tests terraform refresh
func TestWorkerQueueStateRefresh(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Skip("KUBIYA_CONTROL_PLANE_API_KEY not set, skipping test")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../../testdata/workers/minimal",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
			"TF_CLI_CONFIG_FILE":           os.Getenv("TF_CLI_CONFIG_FILE"),
			"HOME":                    os.Getenv("HOME"),
			"TF_SKIP_PROVIDER_VERIFY": "1",
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	workerQueueID := terraform.Output(t, terraformOptions, "worker_queue_id")
	require.NotEmpty(t, workerQueueID)

	terraform.RunTerraformCommand(t, terraformOptions, "refresh")

	refreshedID := terraform.Output(t, terraformOptions, "worker_queue_id")
	assert.Equal(t, workerQueueID, refreshedID)

	t.Logf("✓ State refresh test passed for worker queue %s", workerQueueID)
}
