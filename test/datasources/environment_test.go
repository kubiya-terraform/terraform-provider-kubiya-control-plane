package test

import (
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// TestEnvironmentDataSource tests the environment data source
func TestEnvironmentDataSource(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Fatal("KUBIYA_CONTROL_PLANE_API_KEY environment variable is not set")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../../testdata/environments",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	// Verify data source outputs
	dataName := terraform.Output(t, terraformOptions, "data_minimal_name")
	assert.Equal(t, "test-environment-minimal", dataName)

	dataStatus := terraform.Output(t, terraformOptions, "data_minimal_status")
	assert.NotEmpty(t, dataStatus)

	dataDescription := terraform.Output(t, terraformOptions, "data_full_description")
	assert.Contains(t, dataDescription, "Comprehensive")

	dataTags := terraform.Output(t, terraformOptions, "data_full_tags")
	assert.NotEmpty(t, dataTags)

	dataDisplayName := terraform.Output(t, terraformOptions, "data_full_display_name")
	assert.Equal(t, "Test Environment Full", dataDisplayName)

	t.Logf("✓ Environment data source test passed")
}
