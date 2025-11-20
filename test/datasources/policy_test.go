package test

import (
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// TestPolicyDataSource tests the policy data source
func TestPolicyDataSource(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Fatal("KUBIYA_CONTROL_PLANE_API_KEY environment variable is not set")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../../testdata/policies",
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
	assert.Equal(t, "test-policy-minimal", dataName)

	dataDescription := terraform.Output(t, terraformOptions, "data_full_description")
	assert.Contains(t, dataDescription, "Comprehensive")

	dataEnabled := terraform.Output(t, terraformOptions, "data_full_enabled")
	assert.Equal(t, "true", dataEnabled)

	dataDisabledEnabled := terraform.Output(t, terraformOptions, "data_disabled_enabled")
	assert.Equal(t, "false", dataDisabledEnabled)

	dataTags := terraform.Output(t, terraformOptions, "data_rbac_tags")
	assert.NotEmpty(t, dataTags)

	t.Logf("✓ Policy data source test passed")
}
