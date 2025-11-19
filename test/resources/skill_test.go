package test

import (
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSkillBasic tests the basic skill resource lifecycle using the example
func TestSkillBasic(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Fatal("KUBIYA_CONTROL_PLANE_API_KEY environment variable is not set")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../../examples/skill",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	skillID := terraform.Output(t, terraformOptions, "skill_id")
	t.Logf("Created skill with ID: %s", skillID)
}

// TestSkillMinimal tests minimal skill configuration
func TestSkillMinimal(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Skip("KUBIYA_CONTROL_PLANE_API_KEY not set, skipping testdata test")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../../testdata/skills/minimal",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	skillID := terraform.Output(t, terraformOptions, "skill_id")
	require.NotEmpty(t, skillID)

	skillName := terraform.Output(t, terraformOptions, "skill_name")
	assert.Equal(t, "test-skill-minimal", skillName)

	skillType := terraform.Output(t, terraformOptions, "skill_type")
	assert.Equal(t, "shell", skillType)

	// Verify data source
	dataSkillType := terraform.Output(t, terraformOptions, "data_skill_type")
	assert.Equal(t, skillType, dataSkillType)

	t.Logf("✓ Minimal skill test passed: ID=%s, Name=%s, Type=%s", skillID, skillName, skillType)
}

// TestSkillFull tests skill with all optional fields
func TestSkillFull(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Skip("KUBIYA_CONTROL_PLANE_API_KEY not set, skipping testdata test")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../../testdata/skills/full",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	skillID := terraform.Output(t, terraformOptions, "skill_id")
	require.NotEmpty(t, skillID)

	skillName := terraform.Output(t, terraformOptions, "skill_name")
	assert.Equal(t, "test-skill-full", skillName)

	skillType := terraform.Output(t, terraformOptions, "skill_type")
	assert.Equal(t, "python", skillType)

	skillDescription := terraform.Output(t, terraformOptions, "skill_description")
	assert.Contains(t, skillDescription, "Comprehensive")

	skillEnabled := terraform.Output(t, terraformOptions, "skill_enabled")
	assert.Equal(t, "true", skillEnabled)

	// Verify data source
	dataDescription := terraform.Output(t, terraformOptions, "data_skill_description")
	assert.Equal(t, skillDescription, dataDescription)

	t.Logf("✓ Full skill test passed: ID=%s", skillID)
}

// TestSkillComprehensive tests all skill resource scenarios together
func TestSkillComprehensive(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Fatal("KUBIYA_CONTROL_PLANE_API_KEY environment variable is not set")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../../testdata/skills/comprehensive",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	// Test minimal skill
	minimalID := terraform.Output(t, terraformOptions, "minimal_skill_id")
	require.NotEmpty(t, minimalID)
	minimalType := terraform.Output(t, terraformOptions, "minimal_skill_type")
	assert.Equal(t, "shell", minimalType)

	// Test full skill
	fullID := terraform.Output(t, terraformOptions, "full_skill_id")
	require.NotEmpty(t, fullID)
	fullType := terraform.Output(t, terraformOptions, "full_skill_type")
	assert.Equal(t, "python", fullType)

	t.Logf("✓ Comprehensive skill tests passed")
}

// ============================================================================
// STATE MANAGEMENT TESTS - Update Lifecycle & Import
// ============================================================================

// TestSkillUpdate_Fields tests updating skill fields
func TestSkillUpdate_Fields(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Skip("KUBIYA_CONTROL_PLANE_API_KEY not set, skipping test")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../../testdata/skills/update_fields",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	skillID := terraform.Output(t, terraformOptions, "skill_id")
	require.NotEmpty(t, skillID)

	// Update description
	terraformOptions.Vars = map[string]interface{}{
		"description": "Updated skill description",
	}
	terraform.Apply(t, terraformOptions)

	// Verify in-place update
	updatedSkillID := terraform.Output(t, terraformOptions, "skill_id")
	assert.Equal(t, skillID, updatedSkillID)

	updatedDescription := terraform.Output(t, terraformOptions, "skill_description")
	assert.Equal(t, "Updated skill description", updatedDescription)

	t.Logf("✓ Skill update test passed: ID=%s remained stable", skillID)
}

// TestSkillImport tests importing an existing skill
func TestSkillImport(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Skip("KUBIYA_CONTROL_PLANE_API_KEY not set, skipping test")
	}

	createOptions := &terraform.Options{
		TerraformDir: "../../testdata/skills/minimal",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
		},
	}

	terraform.InitAndApply(t, createOptions)
	skillID := terraform.Output(t, createOptions, "skill_id")
	skillName := terraform.Output(t, createOptions, "skill_name")
	require.NotEmpty(t, skillID)

	terraform.RunTerraformCommand(t, createOptions, "state", "rm", "controlplane_skill.minimal")

	importOptions := &terraform.Options{
		TerraformDir: "../../testdata/skills/import",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
		},
		Vars: map[string]interface{}{
			"skill_id":   skillID,
			"skill_name": skillName,
		},
	}

	defer terraform.Destroy(t, importOptions)

	terraform.Init(t, importOptions)
	terraform.RunTerraformCommand(t, importOptions, "import", "controlplane_skill.imported", skillID)

	importedID := terraform.Output(t, importOptions, "imported_skill_id")
	assert.Equal(t, skillID, importedID)

	importedName := terraform.Output(t, importOptions, "imported_skill_name")
	assert.Equal(t, skillName, importedName)

	t.Logf("✓ Skill import test passed: Successfully imported skill %s", skillID)
}

// TestSkillStateRefresh tests terraform refresh
func TestSkillStateRefresh(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Skip("KUBIYA_CONTROL_PLANE_API_KEY not set, skipping test")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../../testdata/skills/minimal",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	skillID := terraform.Output(t, terraformOptions, "skill_id")
	require.NotEmpty(t, skillID)

	terraform.RunTerraformCommand(t, terraformOptions, "refresh")

	refreshedID := terraform.Output(t, terraformOptions, "skill_id")
	assert.Equal(t, skillID, refreshedID)

	t.Logf("✓ State refresh test passed for skill %s", skillID)
}
