package test

import (
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestJobDataSource tests the job data source
func TestJobDataSource(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Fatal("KUBIYA_CONTROL_PLANE_API_KEY environment variable is not set")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../../testdata/jobs",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
			"TF_CLI_CONFIG_FILE":           os.Getenv("TF_CLI_CONFIG_FILE"),
			"HOME":                    os.Getenv("HOME"),
			"TF_SKIP_PROVIDER_VERIFY": "1",
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	// Verify data source outputs
	dataName := terraform.Output(t, terraformOptions, "data_minimal_cron_name")
	assert.Equal(t, "test-job-minimal-cron", dataName)

	dataTriggerType := terraform.Output(t, terraformOptions, "data_minimal_cron_trigger_type")
	assert.Equal(t, "cron", dataTriggerType)

	dataDescription := terraform.Output(t, terraformOptions, "data_full_cron_description")
	assert.Contains(t, dataDescription, "Comprehensive")

	dataCronSchedule := terraform.Output(t, terraformOptions, "data_full_cron_cron_schedule")
	assert.Equal(t, "0 17 * * 1-5", dataCronSchedule)

	dataPlanningMode := terraform.Output(t, terraformOptions, "data_full_cron_planning_mode")
	assert.Equal(t, "predefined_agent", dataPlanningMode)

	dataWebhookTriggerType := terraform.Output(t, terraformOptions, "data_minimal_webhook_trigger_type")
	assert.Equal(t, "webhook", dataWebhookTriggerType)

	dataWebhookURL := terraform.Output(t, terraformOptions, "data_minimal_webhook_webhook_url")
	assert.NotEmpty(t, dataWebhookURL, "Webhook URL should be populated")

	dataDisabledEnabled := terraform.Output(t, terraformOptions, "data_disabled_enabled")
	assert.Equal(t, "false", dataDisabledEnabled)

	t.Logf("✓ Job data source test passed")
}

// TestJobsDataSource tests the jobs list data source
func TestJobsDataSource(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Fatal("KUBIYA_CONTROL_PLANE_API_KEY environment variable is not set")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../../testdata/jobs",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
			"TF_CLI_CONFIG_FILE":           os.Getenv("TF_CLI_CONFIG_FILE"),
			"HOME":                    os.Getenv("HOME"),
			"TF_SKIP_PROVIDER_VERIFY": "1",
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	// Verify list data source outputs
	jobsCount := terraform.Output(t, terraformOptions, "all_jobs_count")
	require.NotEmpty(t, jobsCount)
	assert.NotEqual(t, "0", jobsCount, "Should have at least one job")

	jobsList := terraform.Output(t, terraformOptions, "all_jobs_list")
	assert.NotEmpty(t, jobsList)

	t.Logf("✓ Jobs list data source test passed - found %s jobs", jobsCount)
}
