package test

import (
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestWorkerQueueDataSource tests the worker queue data source
func TestWorkerQueueDataSource(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Fatal("KUBIYA_CONTROL_PLANE_API_KEY environment variable is not set")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../../testdata/workers",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	// Verify data source outputs
	dataName := terraform.Output(t, terraformOptions, "data_minimal_name")
	assert.Equal(t, "test-worker-queue-minimal", dataName)

	dataEnvironmentID := terraform.Output(t, terraformOptions, "data_minimal_environment_id")
	assert.NotEmpty(t, dataEnvironmentID)

	dataDescription := terraform.Output(t, terraformOptions, "data_full_description")
	assert.Contains(t, dataDescription, "Comprehensive")

	dataMaxWorkers := terraform.Output(t, terraformOptions, "data_full_max_workers")
	assert.Equal(t, "10", dataMaxWorkers)

	dataTags := terraform.Output(t, terraformOptions, "data_full_tags")
	assert.NotEmpty(t, dataTags)

	dataStatus := terraform.Output(t, terraformOptions, "data_inactive_status")
	assert.Equal(t, "inactive", dataStatus)

	t.Logf("✓ Worker queue data source test passed")
}

// TestWorkerQueuesDataSource tests the worker_queues list data source
func TestWorkerQueuesDataSource(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Fatal("KUBIYA_CONTROL_PLANE_API_KEY environment variable is not set")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../../testdata/workers",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	// Verify list data source outputs
	queuesCount := terraform.Output(t, terraformOptions, "test_env_queues_count")
	require.NotEmpty(t, queuesCount)
	assert.NotEqual(t, "0", queuesCount, "Should have at least one queue in the environment")

	queuesList := terraform.Output(t, terraformOptions, "test_env_queues_list")
	assert.NotEmpty(t, queuesList)

	t.Logf("✓ Worker queues list data source test passed - found %s queues", queuesCount)
}
