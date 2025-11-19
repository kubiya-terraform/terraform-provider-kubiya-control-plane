# Terraform Provider Tests

Comprehensive test suite for the Kubiya Control Plane Terraform Provider.

## Quick Start

```bash
# Set API key
export KUBIYA_CONTROL_PLANE_API_KEY="your-api-key"

# Run all tests
make test

# Run integration tests only
make test-integration

# Run specific resource tests
go test ./test/resources -run TestAgent -v
```

## Test Organization

### Directory Structure

```
test/
├── README.md                    # This file
├── STATE_MANAGEMENT_TESTS.md   # Detailed state management test documentation
├── helpers/                     # Reusable test helpers
│   └── lifecycle_test_helper.go
├── resources/                   # Resource lifecycle tests
│   ├── agent_test.go
│   ├── team_test.go
│   ├── project_test.go
│   ├── environment_test.go
│   ├── skill_test.go
│   ├── policy_test.go
│   ├── worker_queue_test.go
│   ├── job_test.go
│   └── complete_setup_test.go
└── datasources/                 # Data source tests
    ├── agent_test.go
    ├── team_test.go
    └── ...
```

### Test Types

#### 1. Basic CRUD Tests

Test fundamental create, read, and delete operations.

**Files**: `*_test.go` - Basic test functions
**Examples**:
- `TestAgentBasic`
- `TestTeamMinimal`
- `TestProjectFull`

#### 2. State Management Tests

Comprehensive tests for Terraform state lifecycle.

**Documentation**: See [STATE_MANAGEMENT_TESTS.md](./STATE_MANAGEMENT_TESTS.md)

**Test Categories**:
- **Update Tests**: Verify in-place updates without recreation
- **Import Tests**: Verify resource import functionality
- **State Refresh Tests**: Verify state synchronization
- **Computed Attributes Tests**: Verify computed field behavior

**Examples**:
- `TestAgentUpdate_Name`
- `TestAgentUpdate_Runtime`
- `TestAgentImport`
- `TestAgentStateRefresh`

#### 3. Data Source Tests

Test data source read operations.

**Files**: `datasources/*_test.go`
**Examples**:
- `TestAgentDataSource`
- `TestTeamDataSource`

#### 4. Integration Tests

Test complete multi-resource scenarios.

**Files**: `complete_setup_test.go`
**Example**: `TestCompleteSetup`

## Test Data

### Location

Test configurations are in `testdata/` directory (one level up from `test/`):

```
testdata/
├── agents/
│   ├── minimal/          # Minimal valid configuration
│   ├── full/             # All optional fields
│   ├── comprehensive/    # Multiple scenarios
│   ├── update_*/         # Update test configurations
│   ├── import/           # Import test configurations
│   └── import_full/      # Full import configurations
├── teams/
│   └── [similar structure]
└── [other resources]/
    └── [similar structure]
```

## Running Tests

### All Tests

```bash
# Format, vet, and run all tests
make check

# Run unit tests only (fast)
make test-unit

# Run integration tests only
make test-integration
```

### Filtered Tests

```bash
# Run tests for specific resource
go test ./test/resources -run TestAgent -v
go test ./test/resources -run TestTeam -v

# Run specific test category
go test ./test/resources -run Update -v        # All update tests
go test ./test/resources -run Import -v        # All import tests
go test ./test/resources -run StateRefresh -v  # All refresh tests

# Run single test
go test ./test/resources -run TestAgentUpdate_Name -v
```

### Test Options

```bash
# Run with longer timeout
go test ./test/resources -timeout 30m -v

# Run with race detector
go test -race ./test/resources -v

# Run with coverage
go test -cover ./test/resources -v
```

## Test Statistics

### Coverage Summary

| Resource Type | Basic CRUD | Update Tests | Import Tests | State Refresh |
|--------------|------------|--------------|--------------|---------------|
| Agent        | ✅ 6 tests | ✅ 5 tests   | ✅ 2 tests   | ✅ 2 tests    |
| Team         | ✅ 4 tests | ✅ 4 tests   | ✅ 2 tests   | ✅ 2 tests    |
| Project      | ✅ 3 tests | ✅ 1 test    | ✅ 1 test    | ✅ 1 test     |
| Environment  | ✅ 3 tests | ✅ 1 test    | ✅ 1 test    | ✅ 1 test     |
| Skill        | ✅ 3 tests | ✅ 1 test    | ✅ 1 test    | ✅ 1 test     |
| Policy       | ✅ 3 tests | ✅ 1 test    | ✅ 1 test    | ✅ 1 test     |
| Worker Queue | ✅ 3 tests | ✅ 1 test    | ✅ 1 test    | ✅ 1 test     |
| Job          | ✅ 3 tests | ✅ 1 test    | ✅ 1 test    | ✅ 1 test     |

**Total: 80+ comprehensive state management tests**

## Prerequisites

### Environment Variables

```bash
# Required for integration tests
export KUBIYA_CONTROL_PLANE_API_KEY="your-api-key-here"

# Optional: Override default API URL
export KUBIYA_CONTROL_PLANE_BASE_URL="https://custom-api.example.com"
```

## Resources

- [Terratest Documentation](https://terratest.gruntwork.io/)
- [Terraform Testing Best Practices](https://www.terraform.io/plugin/sdkv2/testing)
- [State Management Tests Documentation](./STATE_MANAGEMENT_TESTS.md)
