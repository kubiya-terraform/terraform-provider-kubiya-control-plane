# State Management Tests

This document describes the comprehensive state management tests for the Kubiya Terraform Provider.

## Overview

These tests verify that the Terraform provider correctly manages state through all lifecycle operations:

- **Create**: Resources are created and state is properly initialized
- **Read**: State can be read and refreshed from the API
- **Update**: Resources can be updated in-place without recreation
- **Delete**: Resources are properly destroyed and removed from state
- **Import**: Existing resources can be imported into Terraform state

## Test Coverage

### Resources Tested

All resources have comprehensive state management tests:

1. **Agent** (`controlplane_agent`)
2. **Team** (`controlplane_team`)
3. **Project** (`controlplane_project`)
4. **Environment** (`controlplane_environment`)
5. **Skill** (`controlplane_skill`)
6. **Policy** (`controlplane_policy`)
7. **Worker Queue** (`controlplane_worker_queue`)
8. **Job** (`controlplane_job`)

### Test Categories

#### 1. Update Lifecycle Tests

These tests verify that resources support in-place updates without being destroyed and recreated.

**What they test:**
- Resource ID remains stable across updates
- Individual field updates work correctly
- Multiple field updates work correctly
- Computed fields (created_at, updated_at) behave correctly
- No unnecessary resource recreation (ForceNew behavior)

**Example tests:**
- `TestAgentUpdate_Name`
- `TestAgentUpdate_Runtime`
- `TestAgentUpdate_MultipleFields`
- `TestTeamUpdate_Status`

#### 2. Import Tests

These tests verify that existing resources can be imported into Terraform state.

**What they test:**
- Resources can be imported using their ID
- Imported state matches actual resource state
- All fields are correctly populated after import
- Terraform plan shows no changes after successful import

**Example tests:**
- `TestAgentImport`
- `TestAgentImport_FullConfiguration`
- `TestTeamImport`

#### 3. State Refresh Tests

These tests verify that `terraform refresh` correctly synchronizes state with the API.

**What they test:**
- State remains valid after refresh
- Resource attributes are correctly updated from API
- No unnecessary diffs are generated

**Example tests:**
- `TestAgentStateRefresh`
- `TestTeamStateRefresh`

#### 4. Computed Attributes Tests

These tests verify that computed fields (ID, created_at, updated_at) are properly managed.

**What they test:**
- ID is generated and remains stable
- created_at is set on creation and never changes
- updated_at is updated when resource changes
- Computed fields don't cause unnecessary diffs

**Example tests:**
- `TestAgentComputedAttributes`
- `TestTeamComputedAttributes`

## Running the Tests

### Prerequisites

1. Set the API key environment variable:
   ```bash
   export KUBIYA_CONTROL_PLANE_API_KEY="your-api-key"
   ```

2. Ensure the provider is built:
   ```bash
   make build
   ```

### Run All Tests

```bash
# Run all tests (unit + integration)
make test

# Run integration tests only
make test-integration
```

### Run Specific Test Categories

```bash
# Run all agent tests
go test ./test/resources -run TestAgent -v

# Run only update tests for all resources
go test ./test/resources -run Update -v

# Run only import tests for all resources
go test ./test/resources -run Import -v

# Run state refresh tests
go test ./test/resources -run StateRefresh -v
```

### Run Individual Tests

```bash
# Run a specific test
go test ./test/resources -run TestAgentUpdate_Name -v

# Run a specific test with timeout
go test ./test/resources -run TestAgentImport -v -timeout 10m
```

## Test Data Structure

Test configurations are organized under `testdata/` directory:

```
testdata/
├── agents/
│   ├── minimal/          # Minimal configuration tests
│   ├── full/             # Full configuration tests
│   ├── comprehensive/    # Multiple scenarios
│   ├── update_multiple/  # Update test configs
│   ├── update_runtime/   # Runtime update configs
│   ├── update_team/      # Team assignment configs
│   ├── import/           # Import test configs
│   └── import_full/      # Full import configs
├── teams/
│   ├── minimal/
│   ├── full/
│   ├── update_name/
│   ├── update_status/
│   ├── update_runtime/
│   ├── update_multiple/
│   ├── import/
│   └── import_full/
└── [similar structure for other resources]
```

## Test Design Principles

### 1. Isolation

- Each test runs in parallel (`t.Parallel()`)
- Tests clean up after themselves (`defer terraform.Destroy()`)
- No shared state between tests

### 2. Comprehensive Coverage

- Test both minimal and full configurations
- Test single-field and multi-field updates
- Test edge cases and error conditions

### 3. Clear Assertions

- Use descriptive assertion messages
- Verify both positive and negative cases
- Log test progress for debugging

### 4. Real API Testing

- All tests run against the real API
- No mocking (ensures real-world compatibility)
- Tests verify actual state management behavior

## Common Test Patterns

### Update Test Pattern

```go
func TestResourceUpdate_Field(t *testing.T) {
    // 1. Create resource with initial values
    terraform.InitAndApply(t, options)
    originalID := terraform.Output(t, options, "resource_id")

    // 2. Update via variables
    options.Vars = map[string]interface{}{
        "field": "new_value",
    }
    terraform.Apply(t, options)

    // 3. Verify ID didn't change (in-place update)
    updatedID := terraform.Output(t, options, "resource_id")
    assert.Equal(t, originalID, updatedID)

    // 4. Verify field was updated
    newValue := terraform.Output(t, options, "resource_field")
    assert.Equal(t, "new_value", newValue)
}
```

### Import Test Pattern

```go
func TestResourceImport(t *testing.T) {
    // 1. Create resource
    terraform.InitAndApply(t, createOptions)
    resourceID := terraform.Output(t, createOptions, "resource_id")

    // 2. Remove from state
    terraform.RunTerraformCommand(t, createOptions, "state", "rm", "...")

    // 3. Import into new state
    terraform.RunTerraformCommand(t, importOptions, "import", "...", resourceID)

    // 4. Verify imported data matches
    importedID := terraform.Output(t, importOptions, "imported_resource_id")
    assert.Equal(t, resourceID, importedID)
}
```

## Troubleshooting

### Test Failures

1. **API Key Issues**
   - Ensure `KUBIYA_CONTROL_PLANE_API_KEY` is set
   - Verify API key has correct permissions

2. **Timeout Errors**
   - Increase test timeout: `-timeout 30m`
   - Check API availability

3. **Resource Already Exists**
   - Clean up orphaned resources from previous failed tests
   - Use unique names with timestamps

4. **State Lock Errors**
   - Ensure no parallel runs using same state
   - Clean up `.terraform/` directories

### Debugging Tests

```bash
# Run with verbose output
go test ./test/resources -run TestAgentUpdate -v

# Run with race detector
go test -race ./test/resources -run TestAgentUpdate -v

# Run with coverage
go test -cover ./test/resources -run TestAgentUpdate -v
```

## CI/CD Integration

### GitHub Actions

```yaml
- name: Run State Management Tests
  env:
    KUBIYA_CONTROL_PLANE_API_KEY: ${{ secrets.API_KEY }}
  run: |
    make test-integration
```

### Test Selection for CI

```bash
# Run fast tests only in PR checks
go test ./test/resources -run 'Update|Import' -short

# Run full suite in main branch
make test-integration
```

## Future Improvements

### Planned Enhancements

1. **Drift Detection Tests**
   - Modify resources outside Terraform
   - Verify drift is detected in plan

2. **ForceNew Validation Tests**
   - Verify which attributes trigger replacement
   - Test replacement behavior

3. **Dependency Tests**
   - Test resources with dependencies
   - Verify correct ordering

4. **Error Handling Tests**
   - Test partial apply failures
   - Verify state consistency after errors

5. **Performance Tests**
   - Measure update operation timing
   - Track API call efficiency

## Contributing

When adding new resources or fields:

1. **Add Update Tests**
   - Test all updatable fields
   - Verify in-place updates

2. **Add Import Tests**
   - Test minimal import
   - Test full configuration import

3. **Add State Refresh Tests**
   - Verify refresh behavior

4. **Update Documentation**
   - Add examples to this document
   - Update test data structure

## Resources

- [Terraform Plugin Testing](https://www.terraform.io/plugin/sdkv2/testing)
- [Terratest Documentation](https://terratest.gruntwork.io/)
- [Provider Development Best Practices](https://www.terraform.io/docs/extend/best-practices/index.html)

## Support

For questions or issues:

1. Check existing test examples
2. Review test output carefully
3. Open an issue with:
   - Test name
   - Error message
   - Steps to reproduce
