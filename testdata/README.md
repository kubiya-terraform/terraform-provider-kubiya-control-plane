# Test Data Configurations

This directory contains Terraform configurations used for testing the Kubiya Control Plane Provider. Each subdirectory contains a focused test configuration for a specific resource type.

## Purpose

These test configurations are designed to:

1. **Test resource creation and lifecycle** - Verify that resources can be created, read, and destroyed
2. **Test data source lookups** - Validate that data sources can correctly retrieve resource information
3. **Provide isolated test scenarios** - Each configuration is self-contained and independent
4. **Serve as validation examples** - Demonstrate correct usage patterns for each resource type

## Directory Structure

```
testdata/
├── agents/         # Agent resource test configuration
├── teams/          # Team resource test configuration
├── projects/       # Project resource test configuration
├── environments/   # Environment resource test configuration
├── toolsets/       # ToolSet resource test configuration
├── policies/       # Policy resource test configuration
└── workers/        # Worker resource test configuration
```

## Test Configurations

### agents/main.tf

Tests the agent resource with:
- Basic agent creation (name, description, model_id, runtime)
- LLM configuration (temperature, max_tokens)
- Data source lookup of the created agent
- Output validation

### teams/main.tf

Tests the team resource with:
- Basic team creation (name, description)
- Data source lookup of the created team
- Output validation

### projects/main.tf

Tests the project resource with:
- Basic project creation (name, description)
- Data source lookup of the created project
- Output validation

### environments/main.tf

Tests the environment resource with:
- Environment creation with JSON configuration
- Region and tier settings
- Data source lookup of the created environment
- Configuration output validation

### toolsets/main.tf

Tests the toolset resource with:
- ToolSet creation (name, type, enabled)
- Shell toolset with allowed commands
- Timeout configuration
- Data source lookup of the created toolset
- Output validation including type and enabled status

### policies/main.tf

Tests the policy resource with:
- Policy creation with OPA Rego content
- Simple test policy with approval rules
- Data source lookup of the created policy
- Sensitive output handling for policy content

### workers/main.tf

Tests the worker resource with:
- Worker registration (name, description)
- Configuration with max_concurrent_tasks and timeout
- Output validation including status

## Running Tests

### Prerequisites

Set environment variables:
```bash
export KUBIYA_CONTROL_PLANE_API_KEY="your-api-key"
export KUBIYA_CONTROL_PLANE_ORG_ID="your-org-id"
```

### Run All Testdata Tests

```bash
# Run all testdata-based tests
go test ./test -v -run "Test.*Configuration"
```

### Run Specific Testdata Test

```bash
# Test agent configuration
go test ./test -v -run TestAgentConfiguration

# Test team configuration
go test ./test -v -run TestTeamConfiguration

# Test policy configuration
go test ./test -v -run TestPolicyConfiguration
```

### Manual Testing

You can also test configurations manually:

```bash
# Navigate to a testdata directory
cd testdata/agents

# Initialize and apply
terraform init
terraform apply

# View outputs
terraform output

# Clean up
terraform destroy
```

## Test Features

### Resource Creation
Each configuration creates a test resource with appropriate attributes:
- Descriptive names (e.g., "test-agent", "test-team")
- Test descriptions indicating automated testing
- Realistic configuration values

### Data Source Validation
Each configuration includes:
- A data source that looks up the created resource by ID
- Output blocks that expose data source attributes
- Validation that resource and data source values match

### Output Validation
Tests verify:
- Resource IDs are generated correctly
- Names match expected values
- Types and configurations are correct
- Data source lookups return consistent data

### Assertions
The `testdata_test.go` file includes assertions using testify:
- `assert.NotEmpty()` - Validates required fields are populated
- `assert.Equal()` - Validates expected values match actual values
- Descriptive test logs for debugging

## Test Characteristics

### Isolation
- Each test runs in its own directory
- No dependencies between tests
- Tests run in parallel for speed

### Cleanup
- All tests use `defer terraform.Destroy()`
- Resources are cleaned up even on test failure
- No manual cleanup required

### Skipping
- Tests are skipped if environment variables are not set
- Graceful degradation without credentials
- No false failures in CI without credentials

## Adding New Test Configurations

When adding a new resource type:

1. **Create directory**
   ```bash
   mkdir testdata/new-resource
   ```

2. **Create main.tf**
   - Include provider configuration
   - Create test resource
   - Add data source lookup
   - Define outputs

3. **Add test in testdata_test.go**
   ```go
   func TestNewResourceConfiguration(t *testing.T) {
       // Follow existing pattern
   }
   ```

4. **Update this README**
   - Add directory to structure
   - Document test configuration
   - Add usage examples

## Best Practices

1. ✅ **Use "test-" prefix** for resource names
2. ✅ **Include "Test ... for automated testing"** in descriptions
3. ✅ **Add data source lookups** to validate read operations
4. ✅ **Define useful outputs** for validation
5. ✅ **Keep configurations simple** and focused
6. ✅ **Use realistic values** but avoid production data
7. ✅ **Document special requirements** in comments

## Differences from Examples

### testdata/ vs examples/

| Aspect | testdata/ | examples/ |
|--------|-----------|-----------|
| Purpose | Automated testing | Documentation & manual testing |
| Complexity | Simple, focused | Realistic, comprehensive |
| Names | test-* | production-like |
| Data sources | Included | Optional |
| Assertions | Yes (in tests) | No |
| Usage | CI/CD pipelines | User learning |

### When to Use Which

- **Use testdata/** for:
  - Automated test suites
  - CI/CD validation
  - Unit-style resource testing
  - Data source validation

- **Use examples/** for:
  - User documentation
  - Integration testing
  - Real-world scenarios
  - Manual validation

## Troubleshooting

### Tests are skipped

Make sure environment variables are set:
```bash
export KUBIYA_CONTROL_PLANE_API_KEY="your-key"
export KUBIYA_CONTROL_PLANE_ORG_ID="your-org-id"
```

### Resources not cleaned up

If tests fail, clean up manually:
```bash
cd testdata/agents
terraform destroy
```

### Configuration errors

Validate syntax:
```bash
cd testdata/agents
terraform init
terraform validate
```

## Related Documentation

- [../test/README.md](../test/README.md) - Integration test documentation
- [../TESTING.md](../TESTING.md) - Overall testing guide
- [../examples/README.md](../examples/README.md) - Example configurations

## Contributing

When contributing test configurations:

1. Follow existing patterns
2. Test your configuration manually first
3. Ensure cleanup works properly
4. Add appropriate assertions in tests
5. Update this README
6. Run all tests before submitting
