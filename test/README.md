# Integration Tests

This directory contains integration tests for the Kubiya Control Plane Terraform Provider. These tests verify the full lifecycle of resources against a live Control Plane API.

## Prerequisites

Before running integration tests, you need:

1. **Terraform installed** (version 1.0 or later)
2. **Go installed** (version 1.21 or later)
3. **Kubiya Control Plane credentials**:
   - API Key
   - Organization ID

## Configuration

Set the required environment variables:

```bash
export KUBIYA_CONTROL_PLANE_API_KEY="your-api-key-here"
export KUBIYA_CONTROL_PLANE_ORG_ID="your-org-id-here"

# Optional: Override the default API URL
# export KUBIYA_CONTROL_PLANE_BASE_URL="https://custom-url.example.com"
```

## Running Tests

### Run All Integration Tests

```bash
make test-integration
```

This will run all integration tests in parallel (where applicable). The tests will:
1. Initialize Terraform in each example directory
2. Apply the configuration
3. Verify outputs
4. Destroy all created resources

### Run Specific Test

You can run a specific test using Go's test filtering:

```bash
# Run only agent tests
go test ./test -v -run TestKubiyaControlPlaneAgent

# Run only team tests
go test ./test -v -run TestKubiyaControlPlaneTeam

# Run only the complete setup test
go test ./test -v -run TestKubiyaControlPlaneCompleteSetup
```

### Run Tests with Timeout

Integration tests can take time. Use a custom timeout:

```bash
go test ./test -v -timeout 45m
```

## Test Structure

Each test follows this pattern:

```go
func TestKubiyaControlPlaneResource(t *testing.T) {
    // 1. Check for required environment variables
    apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
    if apiKey == "" {
        t.Fatal("KUBIYA_CONTROL_PLANE_API_KEY environment variable is not set")
    }

    // 2. Configure Terraform options
    terraformOptions := &terraform.Options{
        TerraformDir: "../examples/resource",
        EnvVars: map[string]string{
            "KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
            "KUBIYA_CONTROL_PLANE_ORG_ID":  orgID,
        },
    }

    // 3. Cleanup (always runs, even on test failure)
    defer terraform.Destroy(t, terraformOptions)

    // 4. Initialize and apply configuration
    terraform.InitAndApply(t, terraformOptions)

    // 5. Verify outputs
    resourceID := terraform.Output(t, terraformOptions, "resource_id")
    t.Logf("Created resource with ID: %s", resourceID)
}
```

## Available Tests

### Resource Tests

- **TestKubiyaControlPlaneAgent** - Tests agent resource lifecycle
  - Creates agent with LLM configuration
  - Verifies agent creation
  - Tests agent destruction

- **TestKubiyaControlPlaneTeam** - Tests team resource lifecycle
  - Creates team with configuration
  - Verifies team creation
  - Tests team destruction

- **TestKubiyaControlPlaneProject** - Tests project resource lifecycle
  - Creates project
  - Verifies project creation
  - Tests project destruction

- **TestKubiyaControlPlaneEnvironment** - Tests environment resource lifecycle
  - Creates environment with configuration
  - Verifies environment creation
  - Tests environment destruction

- **TestKubiyaControlPlaneToolSet** - Tests toolset resource lifecycle
  - Creates toolset (shell, docker, etc.)
  - Verifies toolset configuration
  - Tests toolset destruction

- **TestKubiyaControlPlanePolicy** - Tests policy resource lifecycle
  - Creates OPA Rego policy
  - Verifies policy content
  - Tests policy destruction

- **TestKubiyaControlPlaneWorker** - Tests worker resource lifecycle
  - Registers worker
  - Verifies worker registration
  - Tests worker deregistration

### End-to-End Tests

- **TestKubiyaControlPlaneCompleteSetup** - Tests complete setup
  - Creates multiple resources (project, environment, team, agent)
  - Tests resource dependencies
  - Verifies all resources work together
  - Tests cleanup of all resources

## Test Parallelization

Most individual resource tests run in parallel using `t.Parallel()`. This speeds up test execution significantly.

The complete setup test does NOT run in parallel as it tests resource interactions and dependencies.

## Debugging Tests

### Enable Detailed Logging

```bash
# Run with verbose output
go test ./test -v -run TestKubiyaControlPlaneAgent

# Run with even more detail
TF_LOG=DEBUG go test ./test -v -run TestKubiyaControlPlaneAgent
```

### Skip Cleanup for Debugging

If you want to inspect resources after a test failure, you can modify the test to comment out the defer:

```go
// Temporarily comment out to prevent cleanup
// defer terraform.Destroy(t, terraformOptions)
```

**WARNING**: Remember to manually clean up resources if you do this!

### Manual Cleanup

If tests fail and leave resources behind:

```bash
# Navigate to the example directory
cd examples/agent

# Destroy manually
terraform destroy
```

## Best Practices

1. **Always set cleanup with defer**: Ensures resources are destroyed even if tests fail
2. **Use t.Parallel() when possible**: Speeds up test execution
3. **Check environment variables**: Fail fast if credentials are missing
4. **Log important information**: Use `t.Logf()` to log resource IDs and important info
5. **Use meaningful test names**: Make it clear what each test does
6. **Set appropriate timeouts**: Integration tests can take time

## CI/CD Integration

To run integration tests in CI/CD:

```yaml
# Example GitHub Actions workflow
- name: Run Integration Tests
  env:
    KUBIYA_CONTROL_PLANE_API_KEY: ${{ secrets.KUBIYA_API_KEY }}
    KUBIYA_CONTROL_PLANE_ORG_ID: ${{ secrets.KUBIYA_ORG_ID }}
  run: make test-integration
```

## Troubleshooting

### "API key is not set" Error

Make sure you've exported the required environment variables:
```bash
export KUBIYA_CONTROL_PLANE_API_KEY="your-key"
export KUBIYA_CONTROL_PLANE_ORG_ID="your-org-id"
```

### "Organization ID is required" Error

Both API key and Organization ID must be set:
```bash
export KUBIYA_CONTROL_PLANE_ORG_ID="your-org-id"
```

### Test Timeout

If tests are timing out, increase the timeout:
```bash
go test ./test -v -timeout 60m
```

### Authentication Errors

- Verify your API key is valid
- Check that your API key has appropriate permissions
- Ensure you're using the correct organization ID

### Resource Conflicts

If you get conflicts (resource already exists):
- Check if previous test runs left resources behind
- Manually clean up using `terraform destroy` in the example directories
- Consider using unique resource names in tests

## Contributing

When adding new tests:

1. Follow the existing test structure
2. Use descriptive test names
3. Add defer cleanup handlers
4. Use t.Parallel() if the test is independent
5. Document any special requirements
6. Update this README with the new test

## Test Coverage

To see what's covered by tests:

```bash
make test-coverage
```

This generates a coverage report in `coverage.html`.

## Dependencies

The tests use:

- **Terratest** (github.com/gruntwork-io/terratest) - Terraform testing framework
- **Testify** (github.com/stretchr/testify) - Assertion library
- **Standard Go testing** - Native Go test framework

See `go.mod` for exact versions.
