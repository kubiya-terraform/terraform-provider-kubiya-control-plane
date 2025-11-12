package test

import (
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// TestAgentConfiguration tests the agent configuration from testdata
func TestAgentConfiguration(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Skip("KUBIYA_CONTROL_PLANE_API_KEY not set, skipping testdata test")
	}

	orgID := os.Getenv("KUBIYA_CONTROL_PLANE_ORG_ID")
	if orgID == "" {
		t.Skip("KUBIYA_CONTROL_PLANE_ORG_ID not set, skipping testdata test")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../testdata/agents",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
			"KUBIYA_CONTROL_PLANE_ORG_ID":  orgID,
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	// Test outputs
	agentID := terraform.Output(t, terraformOptions, "agent_id")
	assert.NotEmpty(t, agentID, "Agent ID should not be empty")

	agentName := terraform.Output(t, terraformOptions, "agent_name")
	assert.Equal(t, "test-agent", agentName, "Agent name should match")

	modelID := terraform.Output(t, terraformOptions, "agent_model_id")
	assert.Equal(t, "gpt-4", modelID, "Model ID should be gpt-4")

	t.Logf("Created agent: ID=%s, Name=%s, Model=%s", agentID, agentName, modelID)
}

// TestTeamConfiguration tests the team configuration from testdata
func TestTeamConfiguration(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Skip("KUBIYA_CONTROL_PLANE_API_KEY not set, skipping testdata test")
	}

	orgID := os.Getenv("KUBIYA_CONTROL_PLANE_ORG_ID")
	if orgID == "" {
		t.Skip("KUBIYA_CONTROL_PLANE_ORG_ID not set, skipping testdata test")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../testdata/teams",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
			"KUBIYA_CONTROL_PLANE_ORG_ID":  orgID,
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	// Test outputs
	teamID := terraform.Output(t, terraformOptions, "team_id")
	assert.NotEmpty(t, teamID, "Team ID should not be empty")

	teamName := terraform.Output(t, terraformOptions, "team_name")
	assert.Equal(t, "test-team", teamName, "Team name should match")

	t.Logf("Created team: ID=%s, Name=%s", teamID, teamName)
}

// TestProjectConfiguration tests the project configuration from testdata
func TestProjectConfiguration(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Skip("KUBIYA_CONTROL_PLANE_API_KEY not set, skipping testdata test")
	}

	orgID := os.Getenv("KUBIYA_CONTROL_PLANE_ORG_ID")
	if orgID == "" {
		t.Skip("KUBIYA_CONTROL_PLANE_ORG_ID not set, skipping testdata test")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../testdata/projects",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
			"KUBIYA_CONTROL_PLANE_ORG_ID":  orgID,
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	// Test outputs
	projectID := terraform.Output(t, terraformOptions, "project_id")
	assert.NotEmpty(t, projectID, "Project ID should not be empty")

	projectName := terraform.Output(t, terraformOptions, "project_name")
	assert.Equal(t, "test-project", projectName, "Project name should match")

	t.Logf("Created project: ID=%s, Name=%s", projectID, projectName)
}

// TestEnvironmentConfiguration tests the environment configuration from testdata
func TestEnvironmentConfiguration(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Skip("KUBIYA_CONTROL_PLANE_API_KEY not set, skipping testdata test")
	}

	orgID := os.Getenv("KUBIYA_CONTROL_PLANE_ORG_ID")
	if orgID == "" {
		t.Skip("KUBIYA_CONTROL_PLANE_ORG_ID not set, skipping testdata test")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../testdata/environments",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
			"KUBIYA_CONTROL_PLANE_ORG_ID":  orgID,
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	// Test outputs
	environmentID := terraform.Output(t, terraformOptions, "environment_id")
	assert.NotEmpty(t, environmentID, "Environment ID should not be empty")

	environmentName := terraform.Output(t, terraformOptions, "environment_name")
	assert.Equal(t, "test-environment", environmentName, "Environment name should match")

	configuration := terraform.Output(t, terraformOptions, "environment_configuration")
	assert.NotEmpty(t, configuration, "Environment configuration should not be empty")

	t.Logf("Created environment: ID=%s, Name=%s", environmentID, environmentName)
}

// TestSkillConfiguration tests the skill configuration from testdata
func TestSkillConfiguration(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Skip("KUBIYA_CONTROL_PLANE_API_KEY not set, skipping testdata test")
	}

	orgID := os.Getenv("KUBIYA_CONTROL_PLANE_ORG_ID")
	if orgID == "" {
		t.Skip("KUBIYA_CONTROL_PLANE_ORG_ID not set, skipping testdata test")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../testdata/skills",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
			"KUBIYA_CONTROL_PLANE_ORG_ID":  orgID,
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	// Test outputs
	skillID := terraform.Output(t, terraformOptions, "skill_id")
	assert.NotEmpty(t, skillID, "Skill ID should not be empty")

	skillName := terraform.Output(t, terraformOptions, "skill_name")
	assert.Equal(t, "test-skill", skillName, "Skill name should match")

	skillType := terraform.Output(t, terraformOptions, "skill_type")
	assert.Equal(t, "shell", skillType, "Skill type should be shell")

	skillEnabled := terraform.Output(t, terraformOptions, "skill_enabled")
	assert.Equal(t, "true", skillEnabled, "Skill should be enabled")

	t.Logf("Created skill: ID=%s, Name=%s, Type=%s", skillID, skillName, skillType)
}

// TestPolicyConfiguration tests the policy configuration from testdata
func TestPolicyConfiguration(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Skip("KUBIYA_CONTROL_PLANE_API_KEY not set, skipping testdata test")
	}

	orgID := os.Getenv("KUBIYA_CONTROL_PLANE_ORG_ID")
	if orgID == "" {
		t.Skip("KUBIYA_CONTROL_PLANE_ORG_ID not set, skipping testdata test")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../testdata/policies",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
			"KUBIYA_CONTROL_PLANE_ORG_ID":  orgID,
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	// Test outputs
	policyID := terraform.Output(t, terraformOptions, "policy_id")
	assert.NotEmpty(t, policyID, "Policy ID should not be empty")

	policyName := terraform.Output(t, terraformOptions, "policy_name")
	assert.Equal(t, "test-policy", policyName, "Policy name should match")

	policyEnabled := terraform.Output(t, terraformOptions, "policy_enabled")
	assert.Equal(t, "true", policyEnabled, "Policy should be enabled")

	t.Logf("Created policy: ID=%s, Name=%s, Enabled=%s", policyID, policyName, policyEnabled)
}

// TestWorkerConfiguration tests the worker configuration from testdata
func TestWorkerConfiguration(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Skip("KUBIYA_CONTROL_PLANE_API_KEY not set, skipping testdata test")
	}

	orgID := os.Getenv("KUBIYA_CONTROL_PLANE_ORG_ID")
	if orgID == "" {
		t.Skip("KUBIYA_CONTROL_PLANE_ORG_ID not set, skipping testdata test")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../testdata/workers",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
			"KUBIYA_CONTROL_PLANE_ORG_ID":  orgID,
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	// Test outputs
	workerID := terraform.Output(t, terraformOptions, "worker_id")
	assert.NotEmpty(t, workerID, "Worker ID should not be empty")

	workerName := terraform.Output(t, terraformOptions, "worker_name")
	assert.Equal(t, "test-worker", workerName, "Worker name should match")

	t.Logf("Created worker: ID=%s, Name=%s", workerID, workerName)
}
