package helpers

import (
	"fmt"
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ResourceLifecycleTestConfig configures a lifecycle test
type ResourceLifecycleTestConfig struct {
	ResourceType     string            // e.g., "agent", "team", "project"
	ResourceName     string            // e.g., "minimal", "test"
	TestDataDir      string            // e.g., "../../testdata/agents/minimal"
	UpdateVars       map[string]interface{} // Variables for update test
	ExpectedOutputs  map[string]string // Expected output values after update
}

// TestResourceUpdate tests updating a resource with in-place modifications
func TestResourceUpdate(t *testing.T, config ResourceLifecycleTestConfig) {
	t.Helper()

	apiKey := getAPIKey(t)

	terraformOptions := &terraform.Options{
		TerraformDir: config.TestDataDir,
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	// Initial create
	terraform.InitAndApply(t, terraformOptions)

	resourceID := terraform.Output(t, terraformOptions, fmt.Sprintf("%s_id", config.ResourceType))
	require.NotEmpty(t, resourceID, "Resource ID should not be empty")

	// Apply update with vars
	if len(config.UpdateVars) > 0 {
		terraformOptions.Vars = config.UpdateVars
		terraform.Apply(t, terraformOptions)

		// Verify ID didn't change (in-place update)
		updatedResourceID := terraform.Output(t, terraformOptions, fmt.Sprintf("%s_id", config.ResourceType))
		assert.Equal(t, resourceID, updatedResourceID, "Resource ID should remain stable across updates")

		// Verify expected outputs
		for outputKey, expectedValue := range config.ExpectedOutputs {
			actualValue := terraform.Output(t, terraformOptions, outputKey)
			assert.Equal(t, expectedValue, actualValue, fmt.Sprintf("Output %s should match expected value", outputKey))
		}

		t.Logf("✓ Update test passed: %s ID=%s remained stable", config.ResourceType, resourceID)
	}
}

// TestResourceImport tests importing an existing resource into Terraform state
func TestResourceImport(t *testing.T, createDir, importDir, resourceType, resourceName string) {
	t.Helper()

	apiKey := getAPIKey(t)

	// Create resource
	createOptions := &terraform.Options{
		TerraformDir: createDir,
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
		},
	}

	terraform.InitAndApply(t, createOptions)
	resourceID := terraform.Output(t, createOptions, fmt.Sprintf("%s_id", resourceType))
	resourceNameValue := terraform.Output(t, createOptions, fmt.Sprintf("%s_name", resourceType))
	require.NotEmpty(t, resourceID)

	// Remove from state
	terraform.RunTerraformCommand(t, createOptions, "state", "rm", fmt.Sprintf("controlplane_%s.%s", resourceType, resourceName))

	// Import
	importOptions := &terraform.Options{
		TerraformDir: importDir,
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
		},
		Vars: map[string]interface{}{
			fmt.Sprintf("%s_id", resourceType):   resourceID,
			fmt.Sprintf("%s_name", resourceType): resourceNameValue,
		},
	}

	defer terraform.Destroy(t, importOptions)

	terraform.Init(t, importOptions)
	terraform.RunTerraformCommand(t, importOptions, "import", fmt.Sprintf("controlplane_%s.imported", resourceType), resourceID)

	// Verify import
	importedID := terraform.Output(t, importOptions, fmt.Sprintf("imported_%s_id", resourceType))
	assert.Equal(t, resourceID, importedID)

	importedName := terraform.Output(t, importOptions, fmt.Sprintf("imported_%s_name", resourceType))
	assert.Equal(t, resourceNameValue, importedName)

	t.Logf("✓ Import test passed: Successfully imported %s %s", resourceType, resourceID)
}

// TestResourceStateRefresh tests that terraform refresh correctly synchronizes state
func TestResourceStateRefresh(t *testing.T, testDataDir, resourceType string) {
	t.Helper()

	apiKey := getAPIKey(t)

	terraformOptions := &terraform.Options{
		TerraformDir: testDataDir,
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	// Initial create
	terraform.InitAndApply(t, terraformOptions)

	resourceID := terraform.Output(t, terraformOptions, fmt.Sprintf("%s_id", resourceType))
	require.NotEmpty(t, resourceID)

	// Run refresh
	terraform.RunTerraformCommand(t, terraformOptions, "refresh")

	// Verify state is still valid
	refreshedID := terraform.Output(t, terraformOptions, fmt.Sprintf("%s_id", resourceType))
	assert.Equal(t, resourceID, refreshedID, "Resource ID should remain the same after refresh")

	t.Logf("✓ State refresh test passed: State correctly synchronized for %s %s", resourceType, resourceID)
}

// TestResourceComputedAttributes tests that computed fields are set correctly
func TestResourceComputedAttributes(t *testing.T, testDataDir, resourceType string, updateVars map[string]interface{}) {
	t.Helper()

	apiKey := getAPIKey(t)

	terraformOptions := &terraform.Options{
		TerraformDir: testDataDir,
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	// Initial create
	terraform.InitAndApply(t, terraformOptions)

	// Verify computed attributes are set
	resourceID := terraform.Output(t, terraformOptions, fmt.Sprintf("%s_id", resourceType))
	require.NotEmpty(t, resourceID, "Computed field 'id' should be set")

	createdAt := terraform.Output(t, terraformOptions, fmt.Sprintf("%s_created_at", resourceType))
	require.NotEmpty(t, createdAt, "Computed field 'created_at' should be set")

	// If update vars provided, test that created_at stays stable
	if len(updateVars) > 0 {
		originalCreatedAt := createdAt

		terraformOptions.Vars = updateVars
		terraform.Apply(t, terraformOptions)

		updatedCreatedAt := terraform.Output(t, terraformOptions, fmt.Sprintf("%s_created_at", resourceType))
		assert.Equal(t, originalCreatedAt, updatedCreatedAt, "created_at should not change on update")

		updatedAt := terraform.Output(t, terraformOptions, fmt.Sprintf("%s_updated_at", resourceType))
		assert.NotEmpty(t, updatedAt, "updated_at should be set after update")

		t.Logf("✓ Computed attributes test passed: created_at=%s remains stable, updated_at=%s", createdAt, updatedAt)
	} else {
		t.Logf("✓ Computed attributes test passed: created_at=%s", createdAt)
	}
}

// getAPIKey retrieves the API key from environment
func getAPIKey(t *testing.T) string {
	t.Helper()
	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Skip("KUBIYA_CONTROL_PLANE_API_KEY not set, skipping test")
	}
	return apiKey
}
