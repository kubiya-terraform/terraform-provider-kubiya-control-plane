package test

import (
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestTeamBasic tests the basic team resource lifecycle using the example
func TestTeamBasic(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Fatal("KUBIYA_CONTROL_PLANE_API_KEY environment variable is not set")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../../examples/team",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	teamID := terraform.Output(t, terraformOptions, "team_id")
	t.Logf("Created team with ID: %s", teamID)
}

// TestTeamMinimal tests minimal team configuration
func TestTeamMinimal(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Skip("KUBIYA_CONTROL_PLANE_API_KEY not set, skipping testdata test")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../../testdata/teams/minimal",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	teamID := terraform.Output(t, terraformOptions, "team_id")
	require.NotEmpty(t, teamID)

	teamName := terraform.Output(t, terraformOptions, "team_name")
	assert.Equal(t, "test-team-minimal", teamName)

	// Verify data source
	dataTeamName := terraform.Output(t, terraformOptions, "data_team_name")
	assert.Equal(t, teamName, dataTeamName)

	t.Logf("✓ Minimal team test passed: ID=%s, Name=%s", teamID, teamName)
}

// TestTeamFull tests team with all optional fields
func TestTeamFull(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Skip("KUBIYA_CONTROL_PLANE_API_KEY not set, skipping testdata test")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../../testdata/teams/full",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	teamID := terraform.Output(t, terraformOptions, "team_id")
	require.NotEmpty(t, teamID)

	teamName := terraform.Output(t, terraformOptions, "team_name")
	assert.Equal(t, "test-team-full", teamName)

	teamDescription := terraform.Output(t, terraformOptions, "team_description")
	assert.Contains(t, teamDescription, "Comprehensive")

	runtime := terraform.Output(t, terraformOptions, "team_runtime")
	assert.Equal(t, "default", runtime)

	status := terraform.Output(t, terraformOptions, "team_status")
	assert.Equal(t, "active", status)

	skillIDs := terraform.Output(t, terraformOptions, "team_skill_ids")
	assert.NotEmpty(t, skillIDs, "Team should have skills assigned")

	// Verify data source
	dataDescription := terraform.Output(t, terraformOptions, "data_team_description")
	assert.Equal(t, teamDescription, dataDescription)

	t.Logf("✓ Full team test passed: ID=%s", teamID)
}

// TestTeamClaudeCode tests team with claude_code runtime
func TestTeamClaudeCode(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Skip("KUBIYA_CONTROL_PLANE_API_KEY not set, skipping testdata test")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../../testdata/teams/claude_code",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	teamID := terraform.Output(t, terraformOptions, "team_id")
	require.NotEmpty(t, teamID)

	runtime := terraform.Output(t, terraformOptions, "team_runtime")
	assert.Equal(t, "claude_code", runtime)

	// Verify data source
	dataRuntime := terraform.Output(t, terraformOptions, "data_team_runtime")
	assert.Equal(t, "claude_code", dataRuntime)

	t.Logf("✓ Claude Code team test passed: ID=%s", teamID)
}

// TestTeamComprehensive tests all team resource scenarios together
func TestTeamComprehensive(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Fatal("KUBIYA_CONTROL_PLANE_API_KEY environment variable is not set")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../../testdata/teams/comprehensive",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	// Test minimal team
	minimalID := terraform.Output(t, terraformOptions, "minimal_team_id")
	require.NotEmpty(t, minimalID)
	minimalName := terraform.Output(t, terraformOptions, "minimal_team_name")
	assert.Equal(t, "test-team-minimal", minimalName)

	// Test full team
	fullID := terraform.Output(t, terraformOptions, "full_team_id")
	require.NotEmpty(t, fullID)
	fullName := terraform.Output(t, terraformOptions, "full_team_name")
	assert.Equal(t, "test-team-full", fullName)
	fullRuntime := terraform.Output(t, terraformOptions, "full_team_runtime")
	assert.Equal(t, "default", fullRuntime)

	// Test team statuses
	inactiveID := terraform.Output(t, terraformOptions, "inactive_team_id")
	require.NotEmpty(t, inactiveID)
	inactiveStatus := terraform.Output(t, terraformOptions, "inactive_team_status")
	assert.Equal(t, "inactive", inactiveStatus)

	archivedID := terraform.Output(t, terraformOptions, "archived_team_id")
	require.NotEmpty(t, archivedID)
	archivedStatus := terraform.Output(t, terraformOptions, "archived_team_status")
	assert.Equal(t, "archived", archivedStatus)

	// Test claude_code runtime
	claudeCodeRuntime := terraform.Output(t, terraformOptions, "claude_code_team_runtime")
	assert.Equal(t, "claude_code", claudeCodeRuntime)

	// Test data sources
	dataInactiveStatus := terraform.Output(t, terraformOptions, "data_inactive_status")
	assert.Equal(t, "inactive", dataInactiveStatus)

	t.Logf("✓ Comprehensive team tests passed")
}
