package test

import (
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
)

// TestCompleteSetup tests the complete setup example
func TestCompleteSetup(t *testing.T) {
	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Fatal("KUBIYA_CONTROL_PLANE_API_KEY environment variable is not set")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../../examples/complete-setup",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	// Verify multiple outputs from the complete setup
	projectID := terraform.Output(t, terraformOptions, "project_id")
	t.Logf("Created project with ID: %s", projectID)

	environmentID := terraform.Output(t, terraformOptions, "environment_id")
	t.Logf("Created environment with ID: %s", environmentID)

	teamID := terraform.Output(t, terraformOptions, "team_id")
	t.Logf("Created team with ID: %s", teamID)
}
