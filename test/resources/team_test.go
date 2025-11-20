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
			"TF_CLI_CONFIG_FILE":           os.Getenv("TF_CLI_CONFIG_FILE"),
			"HOME":                    os.Getenv("HOME"),
			"TF_SKIP_PROVIDER_VERIFY": "1",
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
			"TF_CLI_CONFIG_FILE":           os.Getenv("TF_CLI_CONFIG_FILE"),
			"HOME":                    os.Getenv("HOME"),
			"TF_SKIP_PROVIDER_VERIFY": "1",
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
			"TF_CLI_CONFIG_FILE":           os.Getenv("TF_CLI_CONFIG_FILE"),
			"HOME":                    os.Getenv("HOME"),
			"TF_SKIP_PROVIDER_VERIFY": "1",
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
			"TF_CLI_CONFIG_FILE":           os.Getenv("TF_CLI_CONFIG_FILE"),
			"HOME":                    os.Getenv("HOME"),
			"TF_SKIP_PROVIDER_VERIFY": "1",
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
			"TF_CLI_CONFIG_FILE":           os.Getenv("TF_CLI_CONFIG_FILE"),
			"HOME":                    os.Getenv("HOME"),
			"TF_SKIP_PROVIDER_VERIFY": "1",
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

// ============================================================================
// STATE MANAGEMENT TESTS - Update Lifecycle
// ============================================================================

// TestTeamUpdate_Name tests updating the team name
func TestTeamUpdate_Name(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Skip("KUBIYA_CONTROL_PLANE_API_KEY not set, skipping test")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../../testdata/teams/update_name",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
			"TF_CLI_CONFIG_FILE":           os.Getenv("TF_CLI_CONFIG_FILE"),
			"HOME":                    os.Getenv("HOME"),
			"TF_SKIP_PROVIDER_VERIFY": "1",
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	// Initial create
	terraform.InitAndApply(t, terraformOptions)

	teamID := terraform.Output(t, terraformOptions, "team_id")
	require.NotEmpty(t, teamID)
	originalName := terraform.Output(t, terraformOptions, "team_name")

	// Update the name
	terraformOptions.Vars = map[string]interface{}{
		"team_name": "test-team-updated",
	}
	terraform.Apply(t, terraformOptions)

	// Verify state was updated
	updatedTeamID := terraform.Output(t, terraformOptions, "team_id")
	assert.Equal(t, teamID, updatedTeamID, "Team ID should not change on update")

	updatedName := terraform.Output(t, terraformOptions, "team_name")
	assert.NotEqual(t, originalName, updatedName, "Team name should be updated")

	t.Logf("✓ Team name update test passed: ID=%s remained stable", teamID)
}

// TestTeamUpdate_Status tests updating the team status
func TestTeamUpdate_Status(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Skip("KUBIYA_CONTROL_PLANE_API_KEY not set, skipping test")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../../testdata/teams/update_status",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
			"TF_CLI_CONFIG_FILE":           os.Getenv("TF_CLI_CONFIG_FILE"),
			"HOME":                    os.Getenv("HOME"),
			"TF_SKIP_PROVIDER_VERIFY": "1",
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	// Initial create with active status
	terraform.InitAndApply(t, terraformOptions)

	teamID := terraform.Output(t, terraformOptions, "team_id")
	require.NotEmpty(t, teamID)

	originalStatus := terraform.Output(t, terraformOptions, "team_status")
	assert.Equal(t, "active", originalStatus)

	// Update to inactive
	terraformOptions.Vars = map[string]interface{}{
		"status": "inactive",
	}
	terraform.Apply(t, terraformOptions)

	// Verify in-place update
	updatedTeamID := terraform.Output(t, terraformOptions, "team_id")
	assert.Equal(t, teamID, updatedTeamID)

	updatedStatus := terraform.Output(t, terraformOptions, "team_status")
	assert.Equal(t, "inactive", updatedStatus)

	t.Logf("✓ Team status update test passed: Status changed from %s to %s", originalStatus, updatedStatus)
}

// TestTeamUpdate_Runtime tests updating the team runtime
func TestTeamUpdate_Runtime(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Skip("KUBIYA_CONTROL_PLANE_API_KEY not set, skipping test")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../../testdata/teams/update_runtime",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
			"TF_CLI_CONFIG_FILE":           os.Getenv("TF_CLI_CONFIG_FILE"),
			"HOME":                    os.Getenv("HOME"),
			"TF_SKIP_PROVIDER_VERIFY": "1",
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	// Initial create with default runtime
	terraform.InitAndApply(t, terraformOptions)

	teamID := terraform.Output(t, terraformOptions, "team_id")
	require.NotEmpty(t, teamID)

	originalRuntime := terraform.Output(t, terraformOptions, "team_runtime")

	// Update runtime
	terraformOptions.Vars = map[string]interface{}{
		"runtime": "claude_code",
	}
	terraform.Apply(t, terraformOptions)

	// Verify in-place update
	updatedTeamID := terraform.Output(t, terraformOptions, "team_id")
	assert.Equal(t, teamID, updatedTeamID)

	updatedRuntime := terraform.Output(t, terraformOptions, "team_runtime")
	assert.Equal(t, "claude_code", updatedRuntime)

	t.Logf("✓ Team runtime update test passed: Runtime changed from %s to %s", originalRuntime, updatedRuntime)
}

// TestTeamUpdate_MultipleFields tests updating multiple fields at once
func TestTeamUpdate_MultipleFields(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Skip("KUBIYA_CONTROL_PLANE_API_KEY not set, skipping test")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../../testdata/teams/update_multiple",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
			"TF_CLI_CONFIG_FILE":           os.Getenv("TF_CLI_CONFIG_FILE"),
			"HOME":                    os.Getenv("HOME"),
			"TF_SKIP_PROVIDER_VERIFY": "1",
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	// Initial create
	terraform.InitAndApply(t, terraformOptions)

	teamID := terraform.Output(t, terraformOptions, "team_id")
	require.NotEmpty(t, teamID)

	// Update multiple fields
	terraformOptions.Vars = map[string]interface{}{
		"description": "Updated team description",
		"status":      "inactive",
		"runtime":     "claude_code",
	}
	terraform.Apply(t, terraformOptions)

	// Verify ID remains the same
	updatedTeamID := terraform.Output(t, terraformOptions, "team_id")
	assert.Equal(t, teamID, updatedTeamID)

	// Verify fields changed
	updatedDescription := terraform.Output(t, terraformOptions, "team_description")
	assert.Equal(t, "Updated team description", updatedDescription)

	updatedStatus := terraform.Output(t, terraformOptions, "team_status")
	assert.Equal(t, "inactive", updatedStatus)

	updatedRuntime := terraform.Output(t, terraformOptions, "team_runtime")
	assert.Equal(t, "claude_code", updatedRuntime)

	t.Logf("✓ Multiple field update test passed: ID=%s remained stable", teamID)
}

// ============================================================================
// STATE MANAGEMENT TESTS - Import
// ============================================================================

// TestTeamImport tests importing an existing team into state
func TestTeamImport(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Skip("KUBIYA_CONTROL_PLANE_API_KEY not set, skipping test")
	}

	// First, create a team
	createOptions := &terraform.Options{
		TerraformDir: "../../testdata/teams/minimal",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
			"TF_CLI_CONFIG_FILE":           os.Getenv("TF_CLI_CONFIG_FILE"),
			"HOME":                    os.Getenv("HOME"),
			"TF_SKIP_PROVIDER_VERIFY": "1",
		},
	}

	terraform.InitAndApply(t, createOptions)
	teamID := terraform.Output(t, createOptions, "team_id")
	teamName := terraform.Output(t, createOptions, "team_name")
	require.NotEmpty(t, teamID)

	// Remove from state
	terraform.RunTerraformCommand(t, createOptions, "state", "rm", "controlplane_team.minimal")

	// Import into new state
	importOptions := &terraform.Options{
		TerraformDir: "../../testdata/teams/import",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
			"TF_CLI_CONFIG_FILE":           os.Getenv("TF_CLI_CONFIG_FILE"),
			"HOME":                    os.Getenv("HOME"),
			"TF_SKIP_PROVIDER_VERIFY": "1",
		},
		Vars: map[string]interface{}{
			"team_id":   teamID,
			"team_name": teamName,
		},
	}

	defer terraform.Destroy(t, importOptions)

	terraform.Init(t, importOptions)
	terraform.RunTerraformCommand(t, importOptions, "import", "controlplane_team.imported", teamID)

	// Verify the import
	importedID := terraform.Output(t, importOptions, "imported_team_id")
	assert.Equal(t, teamID, importedID)

	importedName := terraform.Output(t, importOptions, "imported_team_name")
	assert.Equal(t, teamName, importedName)

	t.Logf("✓ Team import test passed: Successfully imported team %s", teamID)
}

// TestTeamImport_FullConfiguration tests importing a team with all fields
func TestTeamImport_FullConfiguration(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Skip("KUBIYA_CONTROL_PLANE_API_KEY not set, skipping test")
	}

	// Create a fully configured team
	createOptions := &terraform.Options{
		TerraformDir: "../../testdata/teams/full",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
			"TF_CLI_CONFIG_FILE":           os.Getenv("TF_CLI_CONFIG_FILE"),
			"HOME":                    os.Getenv("HOME"),
			"TF_SKIP_PROVIDER_VERIFY": "1",
		},
	}

	terraform.InitAndApply(t, createOptions)
	teamID := terraform.Output(t, createOptions, "team_id")
	require.NotEmpty(t, teamID)

	originalName := terraform.Output(t, createOptions, "team_name")
	originalDescription := terraform.Output(t, createOptions, "team_description")
	originalRuntime := terraform.Output(t, createOptions, "team_runtime")

	// Remove from state
	terraform.RunTerraformCommand(t, createOptions, "state", "rm", "controlplane_team.full")

	// Import into new state
	importOptions := &terraform.Options{
		TerraformDir: "../../testdata/teams/import_full",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
			"TF_CLI_CONFIG_FILE":           os.Getenv("TF_CLI_CONFIG_FILE"),
			"HOME":                    os.Getenv("HOME"),
			"TF_SKIP_PROVIDER_VERIFY": "1",
		},
		Vars: map[string]interface{}{
			"team_id": teamID,
		},
	}

	defer terraform.Destroy(t, importOptions)

	terraform.Init(t, importOptions)
	terraform.RunTerraformCommand(t, importOptions, "import", "controlplane_team.imported_full", teamID)

	// Verify all fields
	importedID := terraform.Output(t, importOptions, "imported_team_id")
	assert.Equal(t, teamID, importedID)

	importedName := terraform.Output(t, importOptions, "imported_team_name")
	assert.Equal(t, originalName, importedName)

	importedDescription := terraform.Output(t, importOptions, "imported_team_description")
	assert.Equal(t, originalDescription, importedDescription)

	importedRuntime := terraform.Output(t, importOptions, "imported_team_runtime")
	assert.Equal(t, originalRuntime, importedRuntime)

	t.Logf("✓ Full team import test passed: All fields imported correctly for team %s", teamID)
}

// ============================================================================
// STATE MANAGEMENT TESTS - State Refresh
// ============================================================================

// TestTeamStateRefresh tests that terraform refresh updates state correctly
func TestTeamStateRefresh(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Skip("KUBIYA_CONTROL_PLANE_API_KEY not set, skipping test")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../../testdata/teams/minimal",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
			"TF_CLI_CONFIG_FILE":           os.Getenv("TF_CLI_CONFIG_FILE"),
			"HOME":                    os.Getenv("HOME"),
			"TF_SKIP_PROVIDER_VERIFY": "1",
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	// Initial create
	terraform.InitAndApply(t, terraformOptions)

	teamID := terraform.Output(t, terraformOptions, "team_id")
	require.NotEmpty(t, teamID)

	// Run refresh
	terraform.RunTerraformCommand(t, terraformOptions, "refresh")

	// Verify state is still valid
	refreshedID := terraform.Output(t, terraformOptions, "team_id")
	assert.Equal(t, teamID, refreshedID)

	refreshedName := terraform.Output(t, terraformOptions, "team_name")
	assert.Equal(t, "test-team-minimal", refreshedName)

	t.Logf("✓ State refresh test passed: State correctly synchronized for team %s", teamID)
}

// TestTeamComputedAttributes tests that computed fields are properly managed
func TestTeamComputedAttributes(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Skip("KUBIYA_CONTROL_PLANE_API_KEY not set, skipping test")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../../testdata/teams/minimal",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
			"TF_CLI_CONFIG_FILE":           os.Getenv("TF_CLI_CONFIG_FILE"),
			"HOME":                    os.Getenv("HOME"),
			"TF_SKIP_PROVIDER_VERIFY": "1",
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	// Initial create
	terraform.InitAndApply(t, terraformOptions)

	// Verify computed attributes
	teamID := terraform.Output(t, terraformOptions, "team_id")
	require.NotEmpty(t, teamID)

	createdAt := terraform.Output(t, terraformOptions, "team_created_at")
	require.NotEmpty(t, createdAt)

	t.Logf("✓ Computed attributes test passed: created_at=%s", createdAt)
}
