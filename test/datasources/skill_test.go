package test

import (
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// TestSkillDataSource tests the skill data source
func TestSkillDataSource(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Fatal("KUBIYA_CONTROL_PLANE_API_KEY environment variable is not set")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../../testdata/skills",
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
	dataName := terraform.Output(t, terraformOptions, "data_minimal_name")
	assert.Equal(t, "test-skill-minimal", dataName)

	dataType := terraform.Output(t, terraformOptions, "data_minimal_type")
	assert.Equal(t, "shell", dataType)

	dataDescription := terraform.Output(t, terraformOptions, "data_full_description")
	assert.Contains(t, dataDescription, "Comprehensive")

	dataEnabled := terraform.Output(t, terraformOptions, "data_full_enabled")
	assert.Equal(t, "true", dataEnabled)

	dataShellType := terraform.Output(t, terraformOptions, "data_shell_type")
	assert.Equal(t, "shell", dataShellType)

	dataDockerType := terraform.Output(t, terraformOptions, "data_docker_type")
	assert.Equal(t, "docker", dataDockerType)

	dataDisabledEnabled := terraform.Output(t, terraformOptions, "data_disabled_enabled")
	assert.Equal(t, "false", dataDisabledEnabled)

	t.Logf("✓ Skill data source test passed")
}
