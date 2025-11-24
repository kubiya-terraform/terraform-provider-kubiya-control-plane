# Testing Guide

This document describes the testing strategy and practices for the Kubiya Control Plane Terraform Provider.

## Testing Strategy

The provider uses a comprehensive testing approach with two types of tests:

### 1. Unit Tests

**Location**: Alongside source files in `internal/` packages

**Purpose**: Test individual functions and components in isolation

**Characteristics**:
- Fast execution
- No external dependencies
- Mock API responses
- Test edge cases and error handling

**Example**:
```go
func TestParseJSON(t *testing.T) {
    result, err := parseJSON(`{"key": "value"}`)
    assert.NoError(t, err)
    assert.Equal(t, "value", result["key"])
}
```

### 2. Integration Tests

**Location**: `test/` directory

**Purpose**: Test complete resource lifecycle against live API

**Characteristics**:
- Slower execution (API calls, Terraform operations)
- Requires live credentials
- Tests full CRUD lifecycle
- Validates real-world scenarios

**Example**:
```go
func TestKubiyaControlPlaneAgent(t *testing.T) {
    terraformOptions := &terraform.Options{
        TerraformDir: "../examples/agent",
        EnvVars: map[string]string{
            "KUBIYA_CONTROL_PLANE_API_KEY": apiKey,
        },
    }
    defer terraform.Destroy(t, terraformOptions)
    terraform.InitAndApply(t, terraformOptions)
    agentID := terraform.Output(t, terraformOptions, "agent_id")
    // Verify agentID is valid
}
```

## Running Tests

### Quick Start

```bash
# Run all tests (unit + integration)
make test

# Run only unit tests (fast, no credentials needed)
make test-unit

# Run only integration tests (requires API credentials)
make test-integration

# Run tests with coverage report
make test-coverage
```

### Prerequisites for Integration Tests

Integration tests require:

1. Valid Kubiya Control Plane credentials
2. Environment variables set:
   ```bash
   export KUBIYA_CONTROL_PLANE_API_KEY="your-api-key"
   ```

## Test Organization

```
terraform-provider-kubiya-control-plane/
├── test/                               # Integration tests
│   ├── integration_test.go             # All integration tests
│   └── README.md                       # Integration test docs
├── testdata/                           # Test fixtures and examples
│   ├── agents/
│   ├── teams/
│   ├── projects/
│   └── ...
├── internal/
│   ├── provider/
│   │   ├── *_resource.go               # Resource implementations
│   │   ├── *_resource_test.go          # Unit tests (when added)
│   │   └── helpers.go                  # Helper functions
│   └── clients/
│       ├── client.go                   # API client
│       └── client_test.go              # Unit tests (when added)
└── examples/                           # Used by integration tests
    ├── agent/
    ├── team/
    └── ...
```

## Writing Tests

### Unit Test Guidelines

1. **Location**: Place test files next to the code they test
   - `client.go` → `client_test.go`
   - `helpers.go` → `helpers_test.go`

2. **Naming**: Use descriptive test names
   ```go
   func TestParseJSON_ValidInput(t *testing.T) { }
   func TestParseJSON_InvalidInput(t *testing.T) { }
   func TestParseJSON_EmptyString(t *testing.T) { }
   ```

3. **Structure**: Follow the AAA pattern
   ```go
   func TestExample(t *testing.T) {
       // Arrange - setup test data
       input := "test data"

       // Act - call the function
       result, err := functionUnderTest(input)

       // Assert - verify results
       assert.NoError(t, err)
       assert.Equal(t, "expected", result)
   }
   ```

4. **Table-Driven Tests**: Use for multiple scenarios
   ```go
   func TestParseJSON(t *testing.T) {
       tests := []struct {
           name    string
           input   string
           want    map[string]interface{}
           wantErr bool
       }{
           {"valid json", `{"key":"value"}`, map[string]interface{}{"key":"value"}, false},
           {"empty string", "", map[string]interface{}{}, false},
           {"invalid json", `{invalid}`, nil, true},
       }

       for _, tt := range tests {
           t.Run(tt.name, func(t *testing.T) {
               got, err := parseJSON(tt.input)
               if tt.wantErr {
                   assert.Error(t, err)
               } else {
                   assert.NoError(t, err)
                   assert.Equal(t, tt.want, got)
               }
           })
       }
   }
   ```

### Integration Test Guidelines

1. **Credentials Check**: Always verify credentials first
   ```go
   apiKey := os.Getenv("KUBIYA_CONTROL_PLANE_API_KEY")
   if apiKey == "" {
       t.Fatal("KUBIYA_CONTROL_PLANE_API_KEY environment variable is not set")
   }
   ```

2. **Cleanup**: Always use defer for cleanup
   ```go
   defer terraform.Destroy(t, terraformOptions)
   ```

3. **Parallelization**: Use `t.Parallel()` for independent tests
   ```go
   func TestIndependentResource(t *testing.T) {
       t.Parallel()
       // test code
   }
   ```

4. **Logging**: Log important information
   ```go
   resourceID := terraform.Output(t, terraformOptions, "resource_id")
   t.Logf("Created resource with ID: %s", resourceID)
   ```

5. **Timeouts**: Set appropriate timeouts
   ```go
   // In command line
   go test ./test -v -timeout 30m
   ```

## Test Coverage

### Generate Coverage Report

```bash
make test-coverage
```

This generates:
- `coverage.out` - Coverage data
- `coverage.html` - Visual coverage report

### View Coverage Report

```bash
# Open in browser
open coverage.html

# Or view in terminal
go tool cover -func=coverage.out
```

### Coverage Goals

- **Overall**: Aim for 70%+ coverage
- **Critical paths**: 90%+ coverage for:
  - Client API methods
  - Resource CRUD operations
  - Error handling paths
- **Helper functions**: 80%+ coverage

## Continuous Integration

### CI Test Strategy

The CI pipeline runs tests in this order:

1. **Code Quality Checks**
   ```bash
   make fmt-check    # Verify code formatting
   make vet          # Run go vet
   ```

2. **Unit Tests**
   ```bash
   make test-unit    # Fast, no credentials needed
   ```

3. **Integration Tests** (optional, on main branch)
   ```bash
   make test-integration    # Requires credentials
   ```

### GitHub Actions Example

```yaml
name: Tests
on: [push, pull_request]

jobs:
  unit-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - name: Run unit tests
        run: make test-unit

  integration-tests:
    runs-on: ubuntu-latest
    # Only run on main branch
    if: github.ref == 'refs/heads/main'
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - name: Run integration tests
        env:
          KUBIYA_CONTROL_PLANE_API_KEY: ${{ secrets.KUBIYA_API_KEY }}
        run: make test-integration
```

## Debugging Tests

### Run Single Test

```bash
go test ./test -v -run TestKubiyaControlPlaneAgent
```

### Enable Terraform Debug Logging

```bash
TF_LOG=DEBUG go test ./test -v -run TestKubiyaControlPlaneAgent
```

### Enable Provider Debug Logging

```bash
TF_LOG_PROVIDER=DEBUG go test ./test -v -run TestKubiyaControlPlaneAgent
```

### Skip Cleanup for Inspection

Temporarily comment out defer cleanup:
```go
// defer terraform.Destroy(t, terraformOptions)
terraform.InitAndApply(t, terraformOptions)
// Inspect resources in Kubiya Control Plane UI
```

**Remember**: Manually clean up resources afterward!

## Test Maintenance

### Regular Maintenance Tasks

1. **Update Dependencies**
   ```bash
   go get -u github.com/gruntwork-io/terratest
   go get -u github.com/stretchr/testify
   go mod tidy
   ```

2. **Review Test Coverage**
   ```bash
   make test-coverage
   # Review coverage.html for gaps
   ```

3. **Clean Up Stale Test Data**
   ```bash
   # Remove any leftover test resources
   cd examples/agent && terraform destroy
   cd examples/team && terraform destroy
   # etc.
   ```

4. **Update Test Examples**
   - Keep example configurations up-to-date
   - Ensure examples match latest provider schema
   - Update documentation if examples change

### Adding Tests for New Resources

When adding a new resource:

1. **Create Example Configuration**
   ```bash
   mkdir examples/new-resource
   # Create main.tf with resource example
   ```

2. **Add Integration Test**
   ```go
   func TestKubiyaControlPlaneNewResource(t *testing.T) {
       t.Parallel()
       // Follow existing test pattern
   }
   ```

3. **Add Unit Tests** (if applicable)
   ```go
   // In internal/clients/new_resource_test.go or internal/provider/new_resource_test.go
   ```

4. **Update Documentation**
   - Add test to test/README.md
   - Document any special requirements

## Troubleshooting

### Common Issues

**"API key is not set"**
```bash
export KUBIYA_CONTROL_PLANE_API_KEY="your-key"
```

**"Test timeout"**
```bash
go test ./test -v -timeout 60m
```

**"Resource already exists"**
- Previous test run left resources
- Clean up manually or use unique names

**"Cannot find terraform"**
```bash
# Install Terraform
brew install terraform  # macOS
# or download from https://terraform.io
```

## Best Practices

1. ✅ **DO** use `t.Parallel()` for independent tests
2. ✅ **DO** use `defer` for cleanup
3. ✅ **DO** log important information with `t.Logf()`
4. ✅ **DO** check for required environment variables
5. ✅ **DO** use descriptive test names
6. ✅ **DO** test error cases

7. ❌ **DON'T** commit credentials
8. ❌ **DON'T** skip cleanup handlers
9. ❌ **DON'T** use hardcoded resource IDs
10. ❌ **DON'T** let tests depend on each other
11. ❌ **DON'T** leave test resources running

## Resources

- [Terratest Documentation](https://terratest.gruntwork.io/)
- [Testify Documentation](https://github.com/stretchr/testify)
- [Go Testing Package](https://pkg.go.dev/testing)
- [Terraform Plugin Testing](https://developer.hashicorp.com/terraform/plugin/testing)

## Support

For help with tests:
- Check existing test files for examples
- Review test/README.md for integration test details
- Ask in team chat or open an issue
