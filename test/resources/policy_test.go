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
