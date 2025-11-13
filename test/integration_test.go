package test

import (
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
)

// TestKubiyaControlPlaneAgent tests the agent resource lifecycle
func TestKubiyaControlPlaneAgent(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Fatal("KUBIYA_CONTROL_PLANE_API_KEY environment variable is not set")
	}

	orgID := os.Getenv("KUBIYA_CONTROL_PLANE_ORG_ID")
	if orgID == "" {
		t.Fatal("KUBIYA_CONTROL_PLANE_ORG_ID environment variable is not set")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../examples/agent",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
			"KUBIYA_CONTROL_PLANE_ORG_ID":  orgID,
		},
	}

	// Cleanup after test
	defer terraform.Destroy(t, terraformOptions)

	// Run init and apply
	terraform.InitAndApply(t, terraformOptions)

	// Verify outputs
	agentID := terraform.Output(t, terraformOptions, "agent_id")
	t.Logf("Created agent with ID: %s", agentID)
}

// TestKubiyaControlPlaneTeam tests the team resource lifecycle
func TestKubiyaControlPlaneTeam(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Fatal("KUBIYA_CONTROL_PLANE_API_KEY environment variable is not set")
	}

	orgID := os.Getenv("KUBIYA_CONTROL_PLANE_ORG_ID")
	if orgID == "" {
		t.Fatal("KUBIYA_CONTROL_PLANE_ORG_ID environment variable is not set")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../examples/team",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
			"KUBIYA_CONTROL_PLANE_ORG_ID":  orgID,
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	teamID := terraform.Output(t, terraformOptions, "team_id")
	t.Logf("Created team with ID: %s", teamID)
}

// TestKubiyaControlPlaneProject tests the project resource lifecycle
func TestKubiyaControlPlaneProject(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Fatal("KUBIYA_CONTROL_PLANE_API_KEY environment variable is not set")
	}

	orgID := os.Getenv("KUBIYA_CONTROL_PLANE_ORG_ID")
	if orgID == "" {
		t.Fatal("KUBIYA_CONTROL_PLANE_ORG_ID environment variable is not set")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../examples/project",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
			"KUBIYA_CONTROL_PLANE_ORG_ID":  orgID,
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	projectID := terraform.Output(t, terraformOptions, "project_id")
	t.Logf("Created project with ID: %s", projectID)
}

// TestKubiyaControlPlaneEnvironment tests the environment resource lifecycle
func TestKubiyaControlPlaneEnvironment(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Fatal("KUBIYA_CONTROL_PLANE_API_KEY environment variable is not set")
	}

	orgID := os.Getenv("KUBIYA_CONTROL_PLANE_ORG_ID")
	if orgID == "" {
		t.Fatal("KUBIYA_CONTROL_PLANE_ORG_ID environment variable is not set")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../examples/environment",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
			"KUBIYA_CONTROL_PLANE_ORG_ID":  orgID,
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	environmentID := terraform.Output(t, terraformOptions, "environment_id")
	t.Logf("Created environment with ID: %s", environmentID)
}

// TestKubiyaControlPlaneSkill tests the skill resource lifecycle
func TestKubiyaControlPlaneSkill(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Fatal("KUBIYA_CONTROL_PLANE_API_KEY environment variable is not set")
	}

	orgID := os.Getenv("KUBIYA_CONTROL_PLANE_ORG_ID")
	if orgID == "" {
		t.Fatal("KUBIYA_CONTROL_PLANE_ORG_ID environment variable is not set")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../examples/skill",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
			"KUBIYA_CONTROL_PLANE_ORG_ID":  orgID,
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	skillID := terraform.Output(t, terraformOptions, "skill_id")
	t.Logf("Created skill with ID: %s", skillID)
}

// TestKubiyaControlPlanePolicy tests the policy resource lifecycle
func TestKubiyaControlPlanePolicy(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Fatal("KUBIYA_CONTROL_PLANE_API_KEY environment variable is not set")
	}

	orgID := os.Getenv("KUBIYA_CONTROL_PLANE_ORG_ID")
	if orgID == "" {
		t.Fatal("KUBIYA_CONTROL_PLANE_ORG_ID environment variable is not set")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../examples/policy",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
			"KUBIYA_CONTROL_PLANE_ORG_ID":  orgID,
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	policyID := terraform.Output(t, terraformOptions, "policy_id")
	t.Logf("Created policy with ID: %s", policyID)
}

// TestKubiyaControlPlaneWorkerQueue tests the worker queue resource lifecycle
func TestKubiyaControlPlaneWorkerQueue(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Fatal("KUBIYA_CONTROL_PLANE_API_KEY environment variable is not set")
	}

	orgID := os.Getenv("KUBIYA_CONTROL_PLANE_ORG_ID")
	if orgID == "" {
		t.Fatal("KUBIYA_CONTROL_PLANE_ORG_ID environment variable is not set")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../examples/worker_queue",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
			"KUBIYA_CONTROL_PLANE_ORG_ID":  orgID,
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	queueID := terraform.Output(t, terraformOptions, "queue_id")
	t.Logf("Created worker queue with ID: %s", queueID)
}

// TestKubiyaControlPlaneJob tests the job resource lifecycle
func TestKubiyaControlPlaneJob(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Fatal("KUBIYA_CONTROL_PLANE_API_KEY environment variable is not set")
	}

	orgID := os.Getenv("KUBIYA_CONTROL_PLANE_ORG_ID")
	if orgID == "" {
		t.Fatal("KUBIYA_CONTROL_PLANE_ORG_ID environment variable is not set")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../examples/job",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
			"KUBIYA_CONTROL_PLANE_ORG_ID":  orgID,
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	jobID := terraform.Output(t, terraformOptions, "daily_report_job_id")
	t.Logf("Created job with ID: %s", jobID)

	webhookURL := terraform.Output(t, terraformOptions, "webhook_url")
	t.Logf("Webhook URL: %s", webhookURL)
}

// TestKubiyaControlPlaneCompleteSetup tests the complete setup example
func TestKubiyaControlPlaneCompleteSetup(t *testing.T) {
	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Fatal("KUBIYA_CONTROL_PLANE_API_KEY environment variable is not set")
	}

	orgID := os.Getenv("KUBIYA_CONTROL_PLANE_ORG_ID")
	if orgID == "" {
		t.Fatal("KUBIYA_CONTROL_PLANE_ORG_ID environment variable is not set")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../examples/complete-setup",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
			"KUBIYA_CONTROL_PLANE_ORG_ID":  orgID,
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
