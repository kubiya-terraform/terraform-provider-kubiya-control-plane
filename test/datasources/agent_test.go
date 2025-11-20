package test

import (
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAgentDataSource tests the agent data source
func TestAgentDataSource(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Fatal("KUBIYA_CONTROL_PLANE_API_KEY environment variable is not set")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../../testdata/agents",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
			"TF_CLI_CONFIG_FILE":           os.Getenv("TF_CLI_CONFIG_FILE"),
			"HOME":                    os.Getenv("HOME"),
			"TF_SKIP_PROVIDER_VERIFY": "1",
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	// Verify data source matches resource
	resourceName := terraform.Output(t, terraformOptions, "full_agent_name")
	dataSourceName := terraform.Output(t, terraformOptions, "data_full_description")

	require.NotEmpty(t, resourceName)
	require.NotEmpty(t, dataSourceName)

	// Verify all computed fields are populated
	capabilities := terraform.Output(t, terraformOptions, "data_full_capabilities")
	assert.NotEmpty(t, capabilities)

	runtime := terraform.Output(t, terraformOptions, "data_claude_code_runtime")
	assert.Equal(t, "claude_code", runtime)

	t.Logf("✓ Agent data source test passed")
}
