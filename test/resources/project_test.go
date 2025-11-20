package test

import (
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestProjectBasic tests the basic project resource lifecycle using the example
func TestProjectBasic(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Fatal("KUBIYA_CONTROL_PLANE_API_KEY environment variable is not set")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../../examples/project",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
			"TF_CLI_CONFIG_FILE":           os.Getenv("TF_CLI_CONFIG_FILE"),
			"HOME":                         os.Getenv("HOME"),
			"TF_SKIP_PROVIDER_VERIFY":      "1",
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	projectID := terraform.Output(t, terraformOptions, "project_id")
	t.Logf("Created project with ID: %s", projectID)
}

// TestProjectMinimal tests minimal project configuration
func TestProjectMinimal(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Skip("KUBIYA_CONTROL_PLANE_API_KEY not set, skipping testdata test")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../../testdata/projects/minimal",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
			"TF_CLI_CONFIG_FILE":           os.Getenv("TF_CLI_CONFIG_FILE"),
			"HOME":                         os.Getenv("HOME"),
			"TF_SKIP_PROVIDER_VERIFY":      "1",
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	projectID := terraform.Output(t, terraformOptions, "project_id")
	require.NotEmpty(t, projectID)

	projectName := terraform.Output(t, terraformOptions, "project_name")
	assert.Equal(t, "test-project-minimal", projectName)

	projectKey := terraform.Output(t, terraformOptions, "project_key")
	assert.Equal(t, "TMIN", projectKey)

	// Verify data source
	dataProjectKey := terraform.Output(t, terraformOptions, "data_project_key")
	assert.Equal(t, projectKey, dataProjectKey)

	t.Logf("✓ Minimal project test passed: ID=%s, Name=%s, Key=%s", projectID, projectName, projectKey)
}

// TestProjectFull tests project with all optional fields
func TestProjectFull(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Skip("KUBIYA_CONTROL_PLANE_API_KEY not set, skipping testdata test")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../../testdata/projects/full",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
			"TF_CLI_CONFIG_FILE":           os.Getenv("TF_CLI_CONFIG_FILE"),
			"HOME":                         os.Getenv("HOME"),
			"TF_SKIP_PROVIDER_VERIFY":      "1",
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	projectID := terraform.Output(t, terraformOptions, "project_id")
	require.NotEmpty(t, projectID)

	projectName := terraform.Output(t, terraformOptions, "project_name")
	assert.Equal(t, "test-project-full", projectName)

	projectDescription := terraform.Output(t, terraformOptions, "project_description")
	assert.Contains(t, projectDescription, "Comprehensive")

	visibility := terraform.Output(t, terraformOptions, "project_visibility")
	assert.Equal(t, "org", visibility)

	defaultModel := terraform.Output(t, terraformOptions, "project_default_model")
	assert.Equal(t, "gpt-4", defaultModel)

	status := terraform.Output(t, terraformOptions, "project_status")
	assert.Equal(t, "active", status)

	// Verify data source
	dataDescription := terraform.Output(t, terraformOptions, "data_project_description")
	assert.Equal(t, projectDescription, dataDescription)

	t.Logf("✓ Full project test passed: ID=%s", projectID)
}

// TestProjectComprehensive tests all project resource scenarios together
func TestProjectComprehensive(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Fatal("KUBIYA_CONTROL_PLANE_API_KEY environment variable is not set")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../../testdata/projects/comprehensive",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
			"TF_CLI_CONFIG_FILE":           os.Getenv("TF_CLI_CONFIG_FILE"),
			"HOME":                         os.Getenv("HOME"),
			"TF_SKIP_PROVIDER_VERIFY":      "1",
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	// Test minimal project
	minimalID := terraform.Output(t, terraformOptions, "minimal_project_id")
	require.NotEmpty(t, minimalID)
	minimalName := terraform.Output(t, terraformOptions, "minimal_project_name")
	assert.Equal(t, "test-project-minimal", minimalName)
	minimalKey := terraform.Output(t, terraformOptions, "minimal_project_key")
	assert.Equal(t, "TMIN", minimalKey)

	// Test full project
	fullID := terraform.Output(t, terraformOptions, "full_project_id")
	require.NotEmpty(t, fullID)
	fullVisibility := terraform.Output(t, terraformOptions, "full_project_visibility")
	assert.Equal(t, "org", fullVisibility)
	fullDefaultModel := terraform.Output(t, terraformOptions, "full_project_default_model")
	assert.Equal(t, "gpt-4", fullDefaultModel)

	// Test project statuses
	archivedStatus := terraform.Output(t, terraformOptions, "archived_project_status")
	assert.Equal(t, "archived", archivedStatus)
	pausedStatus := terraform.Output(t, terraformOptions, "paused_project_status")
	assert.Equal(t, "paused", pausedStatus)

	// Test custom model
	customModel := terraform.Output(t, terraformOptions, "custom_model_project_default_model")
	assert.Equal(t, "claude-3-5-sonnet-20241022", customModel)

	// Test data sources
	dataFullVisibility := terraform.Output(t, terraformOptions, "data_full_visibility")
	assert.Equal(t, "org", dataFullVisibility)

	t.Logf("✓ Comprehensive project tests passed")
}

// ============================================================================
// STATE MANAGEMENT TESTS - Update Lifecycle & Import
// ============================================================================

// TestProjectUpdate_Fields tests updating project fields
func TestProjectUpdate_Fields(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Skip("KUBIYA_CONTROL_PLANE_API_KEY not set, skipping test")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../../testdata/projects/update_fields",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
			"TF_CLI_CONFIG_FILE":           os.Getenv("TF_CLI_CONFIG_FILE"),
			"HOME":                         os.Getenv("HOME"),
			"TF_SKIP_PROVIDER_VERIFY":      "1",
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	// Initial create
	terraform.InitAndApply(t, terraformOptions)

	projectID := terraform.Output(t, terraformOptions, "project_id")
	require.NotEmpty(t, projectID)

	// Update fields
	terraformOptions.Vars = map[string]interface{}{
		"description": "Updated project description",
		"visibility":  "private",
	}
	terraform.Apply(t, terraformOptions)

	// Verify in-place update
	updatedProjectID := terraform.Output(t, terraformOptions, "project_id")
	assert.Equal(t, projectID, updatedProjectID)

	updatedDescription := terraform.Output(t, terraformOptions, "project_description")
	assert.Equal(t, "Updated project description", updatedDescription)

	updatedVisibility := terraform.Output(t, terraformOptions, "project_visibility")
	assert.Equal(t, "private", updatedVisibility)

	t.Logf("✓ Project update test passed: ID=%s remained stable", projectID)
}

// TestProjectImport tests importing an existing project
func TestProjectImport(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Skip("KUBIYA_CONTROL_PLANE_API_KEY not set, skipping test")
	}

	// Create project
	createOptions := &terraform.Options{
		TerraformDir: "../../testdata/projects/minimal",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
			"TF_CLI_CONFIG_FILE":           os.Getenv("TF_CLI_CONFIG_FILE"),
			"HOME":                         os.Getenv("HOME"),
			"TF_SKIP_PROVIDER_VERIFY":      "1",
		},
	}

	terraform.InitAndApply(t, createOptions)
	projectID := terraform.Output(t, createOptions, "project_id")
	projectName := terraform.Output(t, createOptions, "project_name")
	projectKey := terraform.Output(t, createOptions, "project_key")
	require.NotEmpty(t, projectID)

	// Remove from state
	terraform.RunTerraformCommand(t, createOptions, "state", "rm", "controlplane_project.minimal")

	// Import
	importOptions := &terraform.Options{
		TerraformDir: "../../testdata/projects/import",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
			"TF_CLI_CONFIG_FILE":           os.Getenv("TF_CLI_CONFIG_FILE"),
			"HOME":                         os.Getenv("HOME"),
			"TF_SKIP_PROVIDER_VERIFY":      "1",
		},
		Vars: map[string]interface{}{
			"project_id":   projectID,
			"project_name": projectName,
			"project_key":  projectKey,
		},
	}

	defer terraform.Destroy(t, importOptions)

	terraform.Init(t, importOptions)
	terraform.RunTerraformCommand(t, importOptions, "import", "controlplane_project.imported", projectID)

	// Verify import
	importedID := terraform.Output(t, importOptions, "imported_project_id")
	assert.Equal(t, projectID, importedID)

	importedName := terraform.Output(t, importOptions, "imported_project_name")
	assert.Equal(t, projectName, importedName)

	t.Logf("✓ Project import test passed: Successfully imported project %s", projectID)
}

// TestProjectStateRefresh tests terraform refresh
func TestProjectStateRefresh(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Skip("KUBIYA_CONTROL_PLANE_API_KEY not set, skipping test")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../../testdata/projects/minimal",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
			"TF_CLI_CONFIG_FILE":           os.Getenv("TF_CLI_CONFIG_FILE"),
			"HOME":                         os.Getenv("HOME"),
			"TF_SKIP_PROVIDER_VERIFY":      "1",
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	projectID := terraform.Output(t, terraformOptions, "project_id")
	require.NotEmpty(t, projectID)

	// Run refresh
	terraform.RunTerraformCommand(t, terraformOptions, "refresh")

	// Verify state is still valid
	refreshedID := terraform.Output(t, terraformOptions, "project_id")
	assert.Equal(t, projectID, refreshedID)

	t.Logf("✓ State refresh test passed for project %s", projectID)
}
