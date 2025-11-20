package test

import (
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// TestTeamDataSource tests the team data source
func TestTeamDataSource(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Fatal("KUBIYA_CONTROL_PLANE_API_KEY environment variable is not set")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../../testdata/teams",
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
	assert.Equal(t, "test-team-minimal", dataName)

	dataDescription := terraform.Output(t, terraformOptions, "data_full_description")
	assert.Contains(t, dataDescription, "Comprehensive")

	dataSkillIDs := terraform.Output(t, terraformOptions, "data_full_skill_ids")
	assert.NotEmpty(t, dataSkillIDs)

	dataStatus := terraform.Output(t, terraformOptions, "data_inactive_status")
	assert.Equal(t, "inactive", dataStatus)

	t.Logf("✓ Team data source test passed")
}
