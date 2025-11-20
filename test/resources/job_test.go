package test

import (
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestJobBasic tests the basic job resource lifecycle using the example
func TestJobBasic(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Fatal("KUBIYA_CONTROL_PLANE_API_KEY environment variable is not set")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../../examples/job",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
			"TF_CLI_CONFIG_FILE":           os.Getenv("TF_CLI_CONFIG_FILE"),
			"HOME":                    os.Getenv("HOME"),
			"TF_SKIP_PROVIDER_VERIFY": "1",
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	jobID := terraform.Output(t, terraformOptions, "daily_report_job_id")
	t.Logf("Created job with ID: %s", jobID)

	webhookURL := terraform.Output(t, terraformOptions, "webhook_url")
	t.Logf("Webhook URL: %s", webhookURL)
}

// TestJobMinimal tests minimal job configuration
func TestJobMinimal(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Skip("KUBIYA_CONTROL_PLANE_API_KEY not set, skipping testdata test")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../../testdata/jobs/minimal",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
			"TF_CLI_CONFIG_FILE":           os.Getenv("TF_CLI_CONFIG_FILE"),
			"HOME":                    os.Getenv("HOME"),
			"TF_SKIP_PROVIDER_VERIFY": "1",
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	jobID := terraform.Output(t, terraformOptions, "job_id")
	require.NotEmpty(t, jobID)

	jobName := terraform.Output(t, terraformOptions, "job_name")
	assert.Equal(t, "test-job-minimal-cron", jobName)

	triggerType := terraform.Output(t, terraformOptions, "job_trigger_type")
	assert.Equal(t, "cron", triggerType)

	cronSchedule := terraform.Output(t, terraformOptions, "job_cron_schedule")
	assert.Equal(t, "0 9 * * *", cronSchedule)

	// Verify data source
	dataTriggerType := terraform.Output(t, terraformOptions, "data_job_trigger_type")
	assert.Equal(t, triggerType, dataTriggerType)

	t.Logf("✓ Minimal job test passed: ID=%s, Name=%s", jobID, jobName)
}

// TestJobFull tests job with all optional fields
func TestJobFull(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Skip("KUBIYA_CONTROL_PLANE_API_KEY not set, skipping testdata test")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../../testdata/jobs/full",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
			"TF_CLI_CONFIG_FILE":           os.Getenv("TF_CLI_CONFIG_FILE"),
			"HOME":                    os.Getenv("HOME"),
			"TF_SKIP_PROVIDER_VERIFY": "1",
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	jobID := terraform.Output(t, terraformOptions, "job_id")
	require.NotEmpty(t, jobID)

	jobName := terraform.Output(t, terraformOptions, "job_name")
	assert.Equal(t, "test-job-full-cron", jobName)

	jobDescription := terraform.Output(t, terraformOptions, "job_description")
	assert.Contains(t, jobDescription, "Comprehensive")

	cronTimezone := terraform.Output(t, terraformOptions, "job_cron_timezone")
	assert.Equal(t, "America/New_York", cronTimezone)

	planningMode := terraform.Output(t, terraformOptions, "job_planning_mode")
	assert.Equal(t, "predefined_agent", planningMode)

	jobEnabled := terraform.Output(t, terraformOptions, "job_enabled")
	assert.Equal(t, "true", jobEnabled)

	// Verify data source
	dataDescription := terraform.Output(t, terraformOptions, "data_job_description")
	assert.Equal(t, jobDescription, dataDescription)

	t.Logf("✓ Full job test passed: ID=%s", jobID)
}

// TestJobComprehensive tests all job resource scenarios together
func TestJobComprehensive(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Fatal("KUBIYA_CONTROL_PLANE_API_KEY environment variable is not set")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../../testdata/jobs/comprehensive",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
			"TF_CLI_CONFIG_FILE":           os.Getenv("TF_CLI_CONFIG_FILE"),
			"HOME":                    os.Getenv("HOME"),
			"TF_SKIP_PROVIDER_VERIFY": "1",
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	// Test minimal cron job
	minimalCronID := terraform.Output(t, terraformOptions, "minimal_cron_job_id")
	require.NotEmpty(t, minimalCronID)
	minimalCronTriggerType := terraform.Output(t, terraformOptions, "minimal_cron_job_trigger_type")
	assert.Equal(t, "cron", minimalCronTriggerType)

	// Test full cron job
	fullCronID := terraform.Output(t, terraformOptions, "full_cron_job_id")
	require.NotEmpty(t, fullCronID)
	fullCronEnabled := terraform.Output(t, terraformOptions, "full_cron_job_enabled")
	assert.Equal(t, "true", fullCronEnabled)

	t.Logf("✓ Comprehensive job tests passed")
}

// ============================================================================
// STATE MANAGEMENT TESTS - Update Lifecycle & Import
// ============================================================================

// TestJobUpdate_Fields tests updating job fields
func TestJobUpdate_Fields(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Skip("KUBIYA_CONTROL_PLANE_API_KEY not set, skipping test")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../../testdata/jobs/update_fields",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
			"TF_CLI_CONFIG_FILE":           os.Getenv("TF_CLI_CONFIG_FILE"),
			"HOME":                    os.Getenv("HOME"),
			"TF_SKIP_PROVIDER_VERIFY": "1",
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	jobID := terraform.Output(t, terraformOptions, "job_id")
	require.NotEmpty(t, jobID)

	// Update description
	terraformOptions.Vars = map[string]interface{}{
		"description": "Updated job description",
	}
	terraform.Apply(t, terraformOptions)

	// Verify in-place update
	updatedJobID := terraform.Output(t, terraformOptions, "job_id")
	assert.Equal(t, jobID, updatedJobID)

	updatedDescription := terraform.Output(t, terraformOptions, "job_description")
	assert.Equal(t, "Updated job description", updatedDescription)

	t.Logf("✓ Job update test passed: ID=%s remained stable", jobID)
}

// TestJobImport tests importing an existing job
func TestJobImport(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Skip("KUBIYA_CONTROL_PLANE_API_KEY not set, skipping test")
	}

	createOptions := &terraform.Options{
		TerraformDir: "../../testdata/jobs/minimal",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
			"TF_CLI_CONFIG_FILE":           os.Getenv("TF_CLI_CONFIG_FILE"),
			"HOME":                    os.Getenv("HOME"),
			"TF_SKIP_PROVIDER_VERIFY": "1",
		},
	}

	terraform.InitAndApply(t, createOptions)
	jobID := terraform.Output(t, createOptions, "job_id")
	jobName := terraform.Output(t, createOptions, "job_name")
	require.NotEmpty(t, jobID)

	terraform.RunTerraformCommand(t, createOptions, "state", "rm", "controlplane_job.minimal")

	importOptions := &terraform.Options{
		TerraformDir: "../../testdata/jobs/import",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
			"TF_CLI_CONFIG_FILE":           os.Getenv("TF_CLI_CONFIG_FILE"),
			"HOME":                    os.Getenv("HOME"),
			"TF_SKIP_PROVIDER_VERIFY": "1",
		},
		Vars: map[string]interface{}{
			"job_id":   jobID,
			"job_name": jobName,
		},
	}

	defer terraform.Destroy(t, importOptions)

	terraform.Init(t, importOptions)
	terraform.RunTerraformCommand(t, importOptions, "import", "controlplane_job.imported", jobID)

	importedID := terraform.Output(t, importOptions, "imported_job_id")
	assert.Equal(t, jobID, importedID)

	importedName := terraform.Output(t, importOptions, "imported_job_name")
	assert.Equal(t, jobName, importedName)

	t.Logf("✓ Job import test passed: Successfully imported job %s", jobID)
}

// TestJobStateRefresh tests terraform refresh
func TestJobStateRefresh(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Skip("KUBIYA_CONTROL_PLANE_API_KEY not set, skipping test")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../../testdata/jobs/minimal",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
			"TF_CLI_CONFIG_FILE":           os.Getenv("TF_CLI_CONFIG_FILE"),
			"HOME":                    os.Getenv("HOME"),
			"TF_SKIP_PROVIDER_VERIFY": "1",
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	jobID := terraform.Output(t, terraformOptions, "job_id")
	require.NotEmpty(t, jobID)

	terraform.RunTerraformCommand(t, terraformOptions, "refresh")

	refreshedID := terraform.Output(t, terraformOptions, "job_id")
	assert.Equal(t, jobID, refreshedID)

	t.Logf("✓ State refresh test passed for job %s", jobID)
}
