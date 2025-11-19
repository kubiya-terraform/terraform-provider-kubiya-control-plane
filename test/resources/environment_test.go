package test

import (
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestEnvironmentBasic tests the basic environment resource lifecycle using the example
func TestEnvironmentBasic(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Fatal("KUBIYA_CONTROL_PLANE_API_KEY environment variable is not set")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../../examples/environment",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	environmentID := terraform.Output(t, terraformOptions, "environment_id")
	t.Logf("Created environment with ID: %s", environmentID)
}

// TestEnvironmentMinimal tests minimal environment configuration
func TestEnvironmentMinimal(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Skip("KUBIYA_CONTROL_PLANE_API_KEY not set, skipping testdata test")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../../testdata/environments/minimal",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	environmentID := terraform.Output(t, terraformOptions, "environment_id")
	require.NotEmpty(t, environmentID)

	environmentName := terraform.Output(t, terraformOptions, "environment_name")
	assert.Equal(t, "test-environment-minimal", environmentName)

	environmentStatus := terraform.Output(t, terraformOptions, "environment_status")
	assert.NotEmpty(t, environmentStatus)

	// Verify data source
	dataEnvironmentStatus := terraform.Output(t, terraformOptions, "data_environment_status")
	assert.Equal(t, environmentStatus, dataEnvironmentStatus)

	t.Logf("✓ Minimal environment test passed: ID=%s, Name=%s", environmentID, environmentName)
}

// TestEnvironmentFull tests environment with all optional fields
func TestEnvironmentFull(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Skip("KUBIYA_CONTROL_PLANE_API_KEY not set, skipping testdata test")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../../testdata/environments/full",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	environmentID := terraform.Output(t, terraformOptions, "environment_id")
	require.NotEmpty(t, environmentID)

	environmentName := terraform.Output(t, terraformOptions, "environment_name")
	assert.Equal(t, "test-environment-full", environmentName)

	displayName := terraform.Output(t, terraformOptions, "environment_display_name")
	assert.Equal(t, "Test Environment Full", displayName)

	description := terraform.Output(t, terraformOptions, "environment_description")
	assert.Contains(t, description, "Comprehensive")

	// Verify data source
	dataDescription := terraform.Output(t, terraformOptions, "data_environment_description")
	assert.Equal(t, description, dataDescription)

	t.Logf("✓ Full environment test passed: ID=%s", environmentID)
}

// TestEnvironmentComprehensive tests all environment resource scenarios together
func TestEnvironmentComprehensive(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
	if apiKey == "" {
		t.Fatal("KUBIYA_CONTROL_PLANE_API_KEY environment variable is not set")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../../testdata/environments/comprehensive",
		EnvVars: map[string]string{
			"KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	// Test minimal environment
	minimalID := terraform.Output(t, terraformOptions, "minimal_environment_id")
	require.NotEmpty(t, minimalID)
	minimalName := terraform.Output(t, terraformOptions, "minimal_environment_name")
	assert.Equal(t, "test-environment-minimal", minimalName)
	minimalStatus := terraform.Output(t, terraformOptions, "minimal_environment_status")
	assert.NotEmpty(t, minimalStatus)

	// Test full environment
	fullID := terraform.Output(t, terraformOptions, "full_environment_id")
	require.NotEmpty(t, fullID)
	fullDisplayName := terraform.Output(t, terraformOptions, "full_environment_display_name")
	assert.Equal(t, "Test Environment Full", fullDisplayName)

	// Test data sources
	dataMinimalStatus := terraform.Output(t, terraformOptions, "data_minimal_status")
	assert.NotEmpty(t, dataMinimalStatus)

	t.Logf("✓ Comprehensive environment tests passed")
}
