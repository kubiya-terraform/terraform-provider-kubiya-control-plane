# Contributing to Kubiya Control Plane Terraform Provider

Thank you for your interest in contributing to the Kubiya Control Plane Terraform Provider! This document provides guidelines and instructions for contributing to this project.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Making Changes](#making-changes)
- [Testing](#testing)
- [Documentation](#documentation)
- [Submitting Changes](#submitting-changes)
- [Provider Development Guidelines](#provider-development-guidelines)

## Code of Conduct

This project adheres to a Code of Conduct that all contributors are expected to follow. Please read [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md) before contributing.

## Getting Started

### Prerequisites

- [Go](https://golang.org/doc/install) >= 1.22
- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Git](https://git-scm.com/downloads)
- Access to a Kubiya Control Plane API instance (for integration testing)

### Fork and Clone

1. Fork the repository on GitHub
2. Clone your fork locally:
   ```bash
   git clone https://github.com/YOUR_USERNAME/terraform-provider-kubiya-control-plane.git
   cd terraform-provider-kubiya-control-plane
   ```
3. Add the upstream repository as a remote:
   ```bash
   git remote add upstream https://github.com/kubiya-terraform/terraform-provider-kubiya-control-plane.git
   ```

## Development Setup

### Install Dependencies

```bash
go mod download
```

### Build the Provider

```bash
go build -o terraform-provider-kubiya-control-plane
```

Or install it to your `$GOPATH/bin`:

```bash
go install
```

### Local Development Configuration

For local testing, configure Terraform to use your local build by creating or modifying `~/.terraformrc`:

```hcl
provider_installation {
  dev_overrides {
    "kubiya/kubiya-control-plane" = "/path/to/your/GOPATH/bin"
  }

  direct {}
}
```

### Environment Variables

Set up your development environment:

```bash
export KUBIYA_CONTROL_PLANE_API_KEY=your_api_key_here
export KUBIYA_CONTROL_PLANE_BASE_URL=http://localhost:7777  # Optional
```

## Making Changes

### Create a Feature Branch

Always create a new branch for your changes:

```bash
git checkout -b feature/your-feature-name
```

Use descriptive branch names:
- `feature/add-webhook-resource` for new features
- `fix/agent-update-bug` for bug fixes
- `docs/improve-readme` for documentation updates
- `refactor/simplify-client-code` for refactoring

### Commit Messages

Write clear and descriptive commit messages:

- Use the present tense ("Add feature" not "Added feature")
- Use the imperative mood ("Move cursor to..." not "Moves cursor to...")
- Limit the first line to 72 characters or less
- Reference issues and pull requests liberally after the first line

Example:
```
Add support for webhook resource

- Implement webhook CRUD operations
- Add webhook data source
- Include comprehensive tests
- Update documentation

Fixes #123
```

## Testing

### Unit Tests

Run unit tests before submitting your changes:

```bash
go test ./...
```

Run tests with coverage:

```bash
go test -v -cover ./...
```

### Integration Tests

Integration tests require a running Kubiya Control Plane API instance:

```bash
export KUBIYA_CONTROL_PLANE_API_KEY=your_api_key
export KUBIYA_CONTROL_PLANE_BASE_URL=http://localhost:7777
go test -v -tags=integration ./test/...
```

### Acceptance Tests

Run Terraform acceptance tests:

```bash
TF_ACC=1 go test ./... -v -timeout 120m
```

### Writing Tests

- Write unit tests for all new functions
- Add integration tests for API client operations
- Include acceptance tests for new resources and data sources
- Ensure tests are idempotent and clean up after themselves
- Use table-driven tests where appropriate

Example test structure:
```go
func TestResourceCreate(t *testing.T) {
    tests := []struct {
        name    string
        input   ResourceInput
        want    ResourceOutput
        wantErr bool
    }{
        // test cases
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // test implementation
        })
    }
}
```

## Documentation

### Code Documentation

- Add godoc comments for all exported functions, types, and constants
- Keep comments clear, concise, and up-to-date
- Use complete sentences with proper punctuation

### Resource Documentation

When adding or modifying resources:

1. Update or create the resource documentation in `docs/resources/`
2. Include all schema attributes with descriptions
3. Provide complete usage examples
4. Document any limitations or special behaviors

### Example Configurations

Add example Terraform configurations in the `examples/` directory for:
- New resources
- Complex use cases
- Integration scenarios

## Submitting Changes

### Before Submitting

1. Ensure all tests pass locally
2. Run code formatting: `go fmt ./...`
3. Run linting: `go vet ./...`
4. Update documentation as needed
5. Add or update tests for your changes
6. Update CHANGELOG.md following [Keep a Changelog](https://keepachangelog.com/) format

### Pull Request Process

1. Push your changes to your fork:
   ```bash
   git push origin feature/your-feature-name
   ```

2. Create a Pull Request from your fork to the main repository

3. Fill out the Pull Request template completely:
   - Provide a clear description of the changes
   - Reference any related issues
   - Include testing instructions
   - Note any breaking changes

4. Ensure all CI checks pass

5. Address review feedback promptly

6. Once approved, a maintainer will merge your PR

### Pull Request Guidelines

- Keep PRs focused on a single feature or fix
- Break large changes into smaller, logical PRs
- Include relevant tests and documentation updates
- Ensure backward compatibility unless explicitly discussed
- Update CHANGELOG.md with your changes

## Provider Development Guidelines

### Adding New Resources

When adding a new resource:

1. **Define the Entity Model** (`internal/entities/`)
   - Create a struct representing the API entity
   - Include JSON tags for API serialization
   - Add any helper methods

2. **Implement the API Client** (`internal/clients/`)
   - Implement CRUD operations
   - Handle API errors appropriately
   - Add comprehensive error messages

3. **Create the Resource** (`internal/provider/`)
   - Implement the Terraform resource schema
   - Implement Create, Read, Update, Delete functions
   - Add import support where applicable
   - Handle state management properly

4. **Register the Resource** (`internal/provider/provider.go`)
   - Add the resource to the provider's resource map

5. **Add Examples** (`examples/`)
   - Create example usage in `examples/resources/`

6. **Document the Resource** (`docs/resources/`)
   - Create comprehensive documentation
   - Include all attributes and their types
   - Provide usage examples

7. **Write Tests**
   - Unit tests for the client
   - Acceptance tests for the resource

### Schema Design

- Use appropriate attribute types (String, Int64, Bool, etc.)
- Mark required vs. optional attributes correctly
- Add clear descriptions for all attributes
- Use validators where appropriate
- Consider whether attributes should be Computed
- For complex nested structures, use appropriate types (List, Set, Map, Object)

### State Management

- Always handle nil values safely
- Properly convert between API models and Terraform state
- Use appropriate null/unknown handling for optional attributes
- Ensure Read function updates state correctly
- Handle deleted resources gracefully

### Error Handling

- Provide clear, actionable error messages
- Include relevant context in errors
- Use appropriate diagnostic severity (Error vs. Warning)
- Handle API-specific errors appropriately

### API Client Guidelines

- Use consistent patterns across all clients
- Implement proper timeout handling
- Log important operations for debugging
- Handle rate limiting appropriately
- Parse and return meaningful API errors

## Getting Help

- Open an issue for bugs or feature requests
- Join our community discussions
- Check existing issues and pull requests before creating new ones
- Ask questions in your pull request if you're unsure about something

## License

By contributing to this project, you agree that your contributions will be licensed under the GNU Affero General Public License v3.0 (AGPL-3.0).

## Recognition

Contributors will be acknowledged in the project's release notes and documentation.

Thank you for contributing to the Kubiya Control Plane Terraform Provider!
