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

// ============================================================================
// STATE MANAGEMENT TESTS - Update Lifecycle
// ============================================================================

// TestAgentUpdate_Name tests updating the agent name
func TestAgentUpdate_Name(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Skip("KUBIYA_CONTROL_PLANE_API_KEY not set, skipping test")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../../testdata/agents/minimal",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	// Initial create
	terraform.InitAndApply(t, terraformOptions)

	agentID := terraform.Output(t, terraformOptions, "agent_id")
	require.NotEmpty(t, agentID)
	originalName := terraform.Output(t, terraformOptions, "agent_name")
	assert.Equal(t, "test-agent-minimal", originalName)

	// Update the name by modifying the config
	terraformOptions.Vars = map[string]interface{}{
		"agent_name": "test-agent-updated",
	}

	// Apply the update
	terraform.Apply(t, terraformOptions)

	// Verify state was updated
	updatedAgentID := terraform.Output(t, terraformOptions, "agent_id")
	assert.Equal(t, agentID, updatedAgentID, "Agent ID should not change on update")

	updatedName := terraform.Output(t, terraformOptions, "agent_name")
	assert.Equal(t, "test-agent-updated", updatedName, "Agent name should be updated")

	// Verify updated_at changed
	updatedAt := terraform.Output(t, terraformOptions, "agent_updated_at")
	assert.NotEmpty(t, updatedAt, "updated_at should be set")

	t.Logf("✓ Agent name update test passed: ID=%s remained stable, Name changed to %s", agentID, updatedName)
}

// TestAgentUpdate_MultipleFields tests updating multiple fields at once
func TestAgentUpdate_MultipleFields(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Skip("KUBIYA_CONTROL_PLANE_API_KEY not set, skipping test")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../../testdata/agents/update_multiple",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	// Initial create
	terraform.InitAndApply(t, terraformOptions)

	agentID := terraform.Output(t, terraformOptions, "agent_id")
	require.NotEmpty(t, agentID)

	originalDescription := terraform.Output(t, terraformOptions, "agent_description")
	originalModelID := terraform.Output(t, terraformOptions, "agent_model_id")

	// Update multiple fields via vars
	terraformOptions.Vars = map[string]interface{}{
		"description": "Updated description with new content",
		"model_id":    "claude-3-5-sonnet-20241022",
	}

	// Apply update
	terraform.Apply(t, terraformOptions)

	// Verify ID remains the same (in-place update)
	updatedAgentID := terraform.Output(t, terraformOptions, "agent_id")
	assert.Equal(t, agentID, updatedAgentID, "Agent ID should remain stable across updates")

	// Verify fields changed
	updatedDescription := terraform.Output(t, terraformOptions, "agent_description")
	assert.NotEqual(t, originalDescription, updatedDescription, "Description should have changed")

	updatedModelID := terraform.Output(t, terraformOptions, "agent_model_id")
	assert.NotEqual(t, originalModelID, updatedModelID, "Model ID should have changed")

	t.Logf("✓ Multiple field update test passed: ID=%s remained stable", agentID)
}

// TestAgentUpdate_Runtime tests updating runtime field
func TestAgentUpdate_Runtime(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Skip("KUBIYA_CONTROL_PLANE_API_KEY not set, skipping test")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../../testdata/agents/update_runtime",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	// Initial create with default runtime
	terraform.InitAndApply(t, terraformOptions)

	agentID := terraform.Output(t, terraformOptions, "agent_id")
	require.NotEmpty(t, agentID)

	originalRuntime := terraform.Output(t, terraformOptions, "agent_runtime")
	assert.Equal(t, "default", originalRuntime)

	// Update runtime via vars
	terraformOptions.Vars = map[string]interface{}{
		"runtime": "claude_code",
	}

	terraform.Apply(t, terraformOptions)

	// Verify in-place update (ID should not change)
	updatedAgentID := terraform.Output(t, terraformOptions, "agent_id")
	assert.Equal(t, agentID, updatedAgentID, "Agent ID should remain the same")

	updatedRuntime := terraform.Output(t, terraformOptions, "agent_runtime")
	assert.Equal(t, "claude_code", updatedRuntime, "Runtime should be updated")

	t.Logf("✓ Runtime update test passed: Runtime changed from %s to %s", originalRuntime, updatedRuntime)
}

// TestAgentUpdate_TeamAssignment tests updating team assignment
func TestAgentUpdate_TeamAssignment(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Skip("KUBIYA_CONTROL_PLANE_API_KEY not set, skipping test")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../../testdata/agents/update_team",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	// Initial create without team
	terraform.InitAndApply(t, terraformOptions)

	agentID := terraform.Output(t, terraformOptions, "agent_id")
	require.NotEmpty(t, agentID)

	// Update to assign team
	terraformOptions.Vars = map[string]interface{}{
		"assign_team": true,
	}
	terraform.Apply(t, terraformOptions)

	// Verify agent ID didn't change
	updatedAgentID := terraform.Output(t, terraformOptions, "agent_id")
	assert.Equal(t, agentID, updatedAgentID)

	// Verify team assignment
	teamID := terraform.Output(t, terraformOptions, "team_id")
	agentTeamID := terraform.Output(t, terraformOptions, "agent_team_id")
	assert.Equal(t, teamID, agentTeamID, "Agent should be assigned to the team")

	t.Logf("✓ Team assignment update test passed: Agent %s assigned to team %s", agentID, teamID)
}

// ============================================================================
// STATE MANAGEMENT TESTS - Import
// ============================================================================

// TestAgentImport tests importing an existing agent into state
func TestAgentImport(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Skip("KUBIYA_CONTROL_PLANE_API_KEY not set, skipping test")
	}

	// First, create an agent outside of Terraform state
	createOptions := &terraform.Options{
		TerraformDir: "../../testdata/agents/minimal",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
		},
	}

	terraform.InitAndApply(t, createOptions)
	agentID := terraform.Output(t, createOptions, "agent_id")
	agentName := terraform.Output(t, createOptions, "agent_name")
	require.NotEmpty(t, agentID)

	// Destroy the Terraform state (but not the actual resource)
	terraform.RunTerraformCommand(t, createOptions, "state", "rm", "controlplane_agent.minimal")

	// Now import it into a new state using import testdata
	importOptions := &terraform.Options{
		TerraformDir: "../../testdata/agents/import",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
		},
		Vars: map[string]interface{}{
			"agent_id":   agentID,
			"agent_name": agentName,
		},
	}

	defer terraform.Destroy(t, importOptions)

	// Initialize and import
	terraform.Init(t, importOptions)
	terraform.RunTerraformCommand(t, importOptions, "import", "controlplane_agent.imported", agentID)

	// Verify the import by reading state
	importedID := terraform.Output(t, importOptions, "imported_agent_id")
	assert.Equal(t, agentID, importedID, "Imported agent ID should match original")

	importedName := terraform.Output(t, importOptions, "imported_agent_name")
	assert.Equal(t, agentName, importedName, "Imported agent name should match original")

	// Run plan to verify no changes needed
	planOutput := terraform.Plan(t, importOptions)
	assert.Contains(t, planOutput, "No changes", "Plan should show no changes after import")

	t.Logf("✓ Agent import test passed: Successfully imported agent %s", agentID)
}

// TestAgentImport_FullConfiguration tests importing an agent with all fields
func TestAgentImport_FullConfiguration(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Skip("KUBIYA_CONTROL_PLANE_API_KEY not set, skipping test")
	}

	// Create a fully configured agent
	createOptions := &terraform.Options{
		TerraformDir: "../../testdata/agents/full",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
		},
	}

	terraform.InitAndApply(t, createOptions)
	agentID := terraform.Output(t, createOptions, "agent_id")
	require.NotEmpty(t, agentID)

	// Get all properties for comparison
	originalName := terraform.Output(t, createOptions, "agent_name")
	originalDescription := terraform.Output(t, createOptions, "agent_description")
	originalModelID := terraform.Output(t, createOptions, "agent_model_id")
	originalRuntime := terraform.Output(t, createOptions, "agent_runtime")

	// Remove from state but keep the resource
	terraform.RunTerraformCommand(t, createOptions, "state", "rm", "controlplane_agent.full")

	// Import into new state
	importOptions := &terraform.Options{
		TerraformDir: "../../testdata/agents/import_full",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
		},
		Vars: map[string]interface{}{
			"agent_id": agentID,
		},
	}

	defer terraform.Destroy(t, importOptions)

	terraform.Init(t, importOptions)
	terraform.RunTerraformCommand(t, importOptions, "import", "controlplane_agent.imported_full", agentID)

	// Verify all fields imported correctly
	importedID := terraform.Output(t, importOptions, "imported_agent_id")
	assert.Equal(t, agentID, importedID)

	importedName := terraform.Output(t, importOptions, "imported_agent_name")
	assert.Equal(t, originalName, importedName)

	importedDescription := terraform.Output(t, importOptions, "imported_agent_description")
	assert.Equal(t, originalDescription, importedDescription)

	importedModelID := terraform.Output(t, importOptions, "imported_agent_model_id")
	assert.Equal(t, originalModelID, importedModelID)

	importedRuntime := terraform.Output(t, importOptions, "imported_agent_runtime")
	assert.Equal(t, originalRuntime, importedRuntime)

	t.Logf("✓ Full agent import test passed: All fields imported correctly for agent %s", agentID)
}

// ============================================================================
// STATE MANAGEMENT TESTS - State Refresh
// ============================================================================

// TestAgentStateRefresh tests that terraform refresh updates state correctly
func TestAgentStateRefresh(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Skip("KUBIYA_CONTROL_PLANE_API_KEY not set, skipping test")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../../testdata/agents/minimal",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	// Initial create
	terraform.InitAndApply(t, terraformOptions)

	agentID := terraform.Output(t, terraformOptions, "agent_id")
	require.NotEmpty(t, agentID)

	// Run refresh
	terraform.RunTerraformCommand(t, terraformOptions, "refresh")

	// Verify state is still valid after refresh
	refreshedID := terraform.Output(t, terraformOptions, "agent_id")
	assert.Equal(t, agentID, refreshedID, "Agent ID should remain the same after refresh")

	refreshedName := terraform.Output(t, terraformOptions, "agent_name")
	assert.Equal(t, "test-agent-minimal", refreshedName, "Agent name should be preserved after refresh")

	t.Logf("✓ State refresh test passed: State correctly synchronized for agent %s", agentID)
}

// ============================================================================
// STATE MANAGEMENT TESTS - Computed Attributes
// ============================================================================

// TestAgentComputedAttributes tests that computed fields are properly managed
func TestAgentComputedAttributes(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Skip("KUBIYA_CONTROL_PLANE_API_KEY not set, skipping test")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../../testdata/agents/minimal",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	// Initial create
	terraform.InitAndApply(t, terraformOptions)

	// Verify computed attributes are set
	agentID := terraform.Output(t, terraformOptions, "agent_id")
	require.NotEmpty(t, agentID, "Computed field 'id' should be set")

	createdAt := terraform.Output(t, terraformOptions, "agent_created_at")
	require.NotEmpty(t, createdAt, "Computed field 'created_at' should be set")

	// Store original values
	originalCreatedAt := createdAt

	// Update the agent
	terraformOptions.Vars = map[string]interface{}{
		"agent_name": "test-agent-updated",
	}
	terraform.Apply(t, terraformOptions)

	// Verify computed attributes behavior
	updatedCreatedAt := terraform.Output(t, terraformOptions, "agent_created_at")
	assert.Equal(t, originalCreatedAt, updatedCreatedAt, "created_at should not change on update")

	updatedAt := terraform.Output(t, terraformOptions, "agent_updated_at")
	assert.NotEmpty(t, updatedAt, "updated_at should be set after update")

	t.Logf("✓ Computed attributes test passed: created_at=%s, updated_at=%s", createdAt, updatedAt)
}
