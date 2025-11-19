package test

import (
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAgentBasic tests the basic agent resource lifecycle using the example
func TestAgentBasic(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Fatal("KUBIYA_CONTROL_PLANE_API_KEY environment variable is not set")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../../examples/agent",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	agentID := terraform.Output(t, terraformOptions, "agent_id")
	t.Logf("Created agent with ID: %s", agentID)
}

// TestAgentMinimal tests minimal agent configuration
func TestAgentMinimal(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Skip("KUBIYA_CONTROL_PLANE_API_KEY not set, skipping testdata test")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../../testdata/agents/minimal",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	agentID := terraform.Output(t, terraformOptions, "agent_id")
	require.NotEmpty(t, agentID, "Agent ID should not be empty")

	agentName := terraform.Output(t, terraformOptions, "agent_name")
	assert.Equal(t, "test-agent-minimal", agentName)

	// Verify data source
	dataAgentName := terraform.Output(t, terraformOptions, "data_agent_name")
	assert.Equal(t, agentName, dataAgentName)

	t.Logf("✓ Minimal agent test passed: ID=%s, Name=%s", agentID, agentName)
}

// TestAgentFull tests agent with all optional fields
func TestAgentFull(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Skip("KUBIYA_CONTROL_PLANE_API_KEY not set, skipping testdata test")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../../testdata/agents/full",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	agentID := terraform.Output(t, terraformOptions, "agent_id")
	require.NotEmpty(t, agentID)

	agentName := terraform.Output(t, terraformOptions, "agent_name")
	assert.Equal(t, "test-agent-full", agentName)

	agentDescription := terraform.Output(t, terraformOptions, "agent_description")
	assert.Contains(t, agentDescription, "Comprehensive")

	runtime := terraform.Output(t, terraformOptions, "agent_runtime")
	assert.Equal(t, "default", runtime)

	modelID := terraform.Output(t, terraformOptions, "agent_model_id")
	assert.Equal(t, "gpt-4", modelID)

	// Verify data source
	dataDescription := terraform.Output(t, terraformOptions, "data_agent_description")
	assert.Equal(t, agentDescription, dataDescription)

	t.Logf("✓ Full agent test passed: ID=%s", agentID)
}

// TestAgentClaudeCode tests agent with claude_code runtime
func TestAgentClaudeCode(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Skip("KUBIYA_CONTROL_PLANE_API_KEY not set, skipping testdata test")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../../testdata/agents/claude_code",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	agentID := terraform.Output(t, terraformOptions, "agent_id")
	require.NotEmpty(t, agentID)

	runtime := terraform.Output(t, terraformOptions, "agent_runtime")
	assert.Equal(t, "claude_code", runtime)

	modelID := terraform.Output(t, terraformOptions, "agent_model_id")
	assert.Equal(t, "claude-3-5-sonnet-20241022", modelID)

	// Verify data source
	dataRuntime := terraform.Output(t, terraformOptions, "data_agent_runtime")
	assert.Equal(t, "claude_code", dataRuntime)

	t.Logf("✓ Claude Code agent test passed: ID=%s", agentID)
}

// TestAgentWithTeam tests agent with team assignment
func TestAgentWithTeam(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Skip("KUBIYA_CONTROL_PLANE_API_KEY not set, skipping testdata test")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../../testdata/agents/with_team",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	agentID := terraform.Output(t, terraformOptions, "agent_id")
	require.NotEmpty(t, agentID)

	teamID := terraform.Output(t, terraformOptions, "team_id")
	require.NotEmpty(t, teamID)

	agentTeamID := terraform.Output(t, terraformOptions, "agent_team_id")
	assert.Equal(t, teamID, agentTeamID, "Agent should be assigned to the created team")

	// Verify data source
	dataAgentTeamID := terraform.Output(t, terraformOptions, "data_agent_team_id")
	assert.Equal(t, teamID, dataAgentTeamID)

	t.Logf("✓ Agent with team test passed: AgentID=%s, TeamID=%s", agentID, teamID)
}

// TestAgentComprehensive tests all agent resource scenarios together
// This test is useful for comprehensive validation but individual tests above provide better granularity
func TestAgentComprehensive(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Fatal("KUBIYA_CONTROL_PLANE_API_KEY environment variable is not set")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../../testdata/agents/comprehensive",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	// Test minimal agent
	minimalID := terraform.Output(t, terraformOptions, "minimal_agent_id")
	require.NotEmpty(t, minimalID)
	minimalName := terraform.Output(t, terraformOptions, "minimal_agent_name")
	assert.Equal(t, "test-agent-minimal", minimalName)

	// Test full agent
	fullID := terraform.Output(t, terraformOptions, "full_agent_id")
	require.NotEmpty(t, fullID)
	fullName := terraform.Output(t, terraformOptions, "full_agent_name")
	assert.Equal(t, "test-agent-full", fullName)

	// Test claude_code agent
	claudeCodeID := terraform.Output(t, terraformOptions, "claude_code_agent_id")
	require.NotEmpty(t, claudeCodeID)
	claudeCodeRuntime := terraform.Output(t, terraformOptions, "claude_code_agent_runtime")
	assert.Equal(t, "claude_code", claudeCodeRuntime)

	// Test agent with team
	withTeamID := terraform.Output(t, terraformOptions, "with_team_agent_id")
	require.NotEmpty(t, withTeamID)
	withTeamTeamID := terraform.Output(t, terraformOptions, "with_team_agent_team_id")
	assert.NotEmpty(t, withTeamTeamID)

	// Test data sources
	dataMinimalName := terraform.Output(t, terraformOptions, "data_minimal_name")
	assert.Equal(t, "test-agent-minimal", dataMinimalName)

	t.Logf("✓ Comprehensive agent tests passed")
}
