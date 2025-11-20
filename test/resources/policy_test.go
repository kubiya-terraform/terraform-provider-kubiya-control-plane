package test

import (
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPolicyBasic tests the basic policy resource lifecycle using the example
func TestPolicyBasic(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Fatal("KUBIYA_CONTROL_PLANE_API_KEY environment variable is not set")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../../examples/policy",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
			"TF_CLI_CONFIG_FILE":           os.Getenv("TF_CLI_CONFIG_FILE"),
			"HOME":                    os.Getenv("HOME"),
			"TF_SKIP_PROVIDER_VERIFY": "1",
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	policyID := terraform.Output(t, terraformOptions, "policy_id")
	t.Logf("Created policy with ID: %s", policyID)
}

// TestPolicyMinimal tests minimal policy configuration
func TestPolicyMinimal(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Skip("KUBIYA_CONTROL_PLANE_API_KEY not set, skipping testdata test")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../../testdata/policies/minimal",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
			"TF_CLI_CONFIG_FILE":           os.Getenv("TF_CLI_CONFIG_FILE"),
			"HOME":                    os.Getenv("HOME"),
			"TF_SKIP_PROVIDER_VERIFY": "1",
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	policyID := terraform.Output(t, terraformOptions, "policy_id")
	require.NotEmpty(t, policyID)

	policyName := terraform.Output(t, terraformOptions, "policy_name")
	assert.Equal(t, "test-policy-minimal", policyName)

	policyType := terraform.Output(t, terraformOptions, "policy_type")
	assert.Equal(t, "rego", policyType)

	// Verify data source
	dataPolicyName := terraform.Output(t, terraformOptions, "data_policy_name")
	assert.Equal(t, policyName, dataPolicyName)

	t.Logf("✓ Minimal policy test passed: ID=%s, Name=%s", policyID, policyName)
}

// TestPolicyFull tests policy with all optional fields
func TestPolicyFull(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Skip("KUBIYA_CONTROL_PLANE_API_KEY not set, skipping testdata test")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../../testdata/policies/full",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
			"TF_CLI_CONFIG_FILE":           os.Getenv("TF_CLI_CONFIG_FILE"),
			"HOME":                    os.Getenv("HOME"),
			"TF_SKIP_PROVIDER_VERIFY": "1",
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	policyID := terraform.Output(t, terraformOptions, "policy_id")
	require.NotEmpty(t, policyID)

	policyName := terraform.Output(t, terraformOptions, "policy_name")
	assert.Equal(t, "test-policy-full", policyName)

	policyDescription := terraform.Output(t, terraformOptions, "policy_description")
	assert.Contains(t, policyDescription, "Comprehensive")

	policyEnabled := terraform.Output(t, terraformOptions, "policy_enabled")
	assert.Equal(t, "true", policyEnabled)

	// Verify data source
	dataDescription := terraform.Output(t, terraformOptions, "data_policy_description")
	assert.Equal(t, policyDescription, dataDescription)

	t.Logf("✓ Full policy test passed: ID=%s", policyID)
}

// TestPolicyComprehensive tests all policy resource scenarios together
func TestPolicyComprehensive(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Fatal("KUBIYA_CONTROL_PLANE_API_KEY environment variable is not set")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../../testdata/policies/comprehensive",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
			"TF_CLI_CONFIG_FILE":           os.Getenv("TF_CLI_CONFIG_FILE"),
			"HOME":                    os.Getenv("HOME"),
			"TF_SKIP_PROVIDER_VERIFY": "1",
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	// Test minimal policy
	minimalID := terraform.Output(t, terraformOptions, "minimal_policy_id")
	require.NotEmpty(t, minimalID)
	minimalName := terraform.Output(t, terraformOptions, "minimal_policy_name")
	assert.Equal(t, "test-policy-minimal", minimalName)

	// Test full policy
	fullID := terraform.Output(t, terraformOptions, "full_policy_id")
	require.NotEmpty(t, fullID)
	fullEnabled := terraform.Output(t, terraformOptions, "full_policy_enabled")
	assert.Equal(t, "true", fullEnabled)

	t.Logf("✓ Comprehensive policy tests passed")
}

// ============================================================================
// STATE MANAGEMENT TESTS - Update Lifecycle & Import
// ============================================================================

// TestPolicyUpdate_Fields tests updating policy fields
func TestPolicyUpdate_Fields(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Skip("KUBIYA_CONTROL_PLANE_API_KEY not set, skipping test")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../../testdata/policies/update_fields",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
			"TF_CLI_CONFIG_FILE":           os.Getenv("TF_CLI_CONFIG_FILE"),
			"HOME":                    os.Getenv("HOME"),
			"TF_SKIP_PROVIDER_VERIFY": "1",
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	policyID := terraform.Output(t, terraformOptions, "policy_id")
	require.NotEmpty(t, policyID)

	// Update description
	terraformOptions.Vars = map[string]interface{}{
		"description": "Updated policy description",
	}
	terraform.Apply(t, terraformOptions)

	// Verify in-place update
	updatedPolicyID := terraform.Output(t, terraformOptions, "policy_id")
	assert.Equal(t, policyID, updatedPolicyID)

	updatedDescription := terraform.Output(t, terraformOptions, "policy_description")
	assert.Equal(t, "Updated policy description", updatedDescription)

	t.Logf("✓ Policy update test passed: ID=%s remained stable", policyID)
}

// TestPolicyImport tests importing an existing policy
func TestPolicyImport(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Skip("KUBIYA_CONTROL_PLANE_API_KEY not set, skipping test")
	}

	createOptions := &terraform.Options{
		TerraformDir: "../../testdata/policies/minimal",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
			"TF_CLI_CONFIG_FILE":           os.Getenv("TF_CLI_CONFIG_FILE"),
			"HOME":                    os.Getenv("HOME"),
			"TF_SKIP_PROVIDER_VERIFY": "1",
		},
	}

	terraform.InitAndApply(t, createOptions)
	policyID := terraform.Output(t, createOptions, "policy_id")
	policyName := terraform.Output(t, createOptions, "policy_name")
	require.NotEmpty(t, policyID)

	terraform.RunTerraformCommand(t, createOptions, "state", "rm", "controlplane_policy.minimal")

	importOptions := &terraform.Options{
		TerraformDir: "../../testdata/policies/import",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
			"TF_CLI_CONFIG_FILE":           os.Getenv("TF_CLI_CONFIG_FILE"),
			"HOME":                    os.Getenv("HOME"),
			"TF_SKIP_PROVIDER_VERIFY": "1",
		},
		Vars: map[string]interface{}{
			"policy_id":   policyID,
			"policy_name": policyName,
		},
	}

	defer terraform.Destroy(t, importOptions)

	terraform.Init(t, importOptions)
	terraform.RunTerraformCommand(t, importOptions, "import", "controlplane_policy.imported", policyID)

	importedID := terraform.Output(t, importOptions, "imported_policy_id")
	assert.Equal(t, policyID, importedID)

	importedName := terraform.Output(t, importOptions, "imported_policy_name")
	assert.Equal(t, policyName, importedName)

	t.Logf("✓ Policy import test passed: Successfully imported policy %s", policyID)
}

// TestPolicyStateRefresh tests terraform refresh
func TestPolicyStateRefresh(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Skip("KUBIYA_CONTROL_PLANE_API_KEY not set, skipping test")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../../testdata/policies/minimal",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
			"TF_CLI_CONFIG_FILE":           os.Getenv("TF_CLI_CONFIG_FILE"),
			"HOME":                    os.Getenv("HOME"),
			"TF_SKIP_PROVIDER_VERIFY": "1",
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	policyID := terraform.Output(t, terraformOptions, "policy_id")
	require.NotEmpty(t, policyID)

	terraform.RunTerraformCommand(t, terraformOptions, "refresh")

	refreshedID := terraform.Output(t, terraformOptions, "policy_id")
	assert.Equal(t, policyID, refreshedID)

	t.Logf("✓ State refresh test passed for policy %s", policyID)
}
