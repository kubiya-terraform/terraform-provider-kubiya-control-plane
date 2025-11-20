package test

import (
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// TestProjectDataSource tests the project data source
func TestProjectDataSource(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Fatal("KUBIYA_CONTROL_PLANE_API_KEY environment variable is not set")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../../testdata/projects",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
			"TF_CLI_CONFIG_FILE":           os.Getenv("TF_CLI_CONFIG_FILE"),
			"HOME":                         os.Getenv("HOME"),
			"TF_SKIP_PROVIDER_VERIFY":      "1",
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	// Verify data source outputs
	dataName := terraform.Output(t, terraformOptions, "data_minimal_name")
	assert.Equal(t, "test-project-minimal", dataName)

	dataKey := terraform.Output(t, terraformOptions, "data_minimal_key")
	assert.Equal(t, "TMIN", dataKey)

	dataDescription := terraform.Output(t, terraformOptions, "data_full_description")
	assert.Contains(t, dataDescription, "Comprehensive")

	dataPolicyIDs := terraform.Output(t, terraformOptions, "data_full_policy_ids")
	assert.NotEmpty(t, dataPolicyIDs)

	dataVisibility := terraform.Output(t, terraformOptions, "data_full_visibility")
	assert.Equal(t, "org", dataVisibility)

	dataStatus := terraform.Output(t, terraformOptions, "data_archived_status")
	assert.Equal(t, "archived", dataStatus)

	dataModel := terraform.Output(t, terraformOptions, "data_custom_model_default_model")
	assert.Equal(t, "claude-3-5-sonnet-20241022", dataModel)

	t.Logf("✓ Project data source test passed")
}
