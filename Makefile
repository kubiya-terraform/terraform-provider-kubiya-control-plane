.PHONY: build install test clean fmt vet lint validate ci check tidy version help

# Default target
default: help

# Build the provider
build:
	@echo "Building provider..."
	go build -o kubiya-control-plane

# Install the provider locally
install:
	@echo "Installing provider..."
	go install

# Run tests
test:
	@echo "Running tests..."
	go test ./... -v

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f kubiya-control-plane coverage.out coverage.html

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Check if code is formatted
fmt-check:
	@echo "Checking code formatting..."
	@if [ -n "$$(go fmt ./...)" ]; then \
		echo "Code is not formatted. Run 'make fmt' to fix."; \
		exit 1; \
	fi
	@echo "Code is properly formatted!"

# Run go vet
vet:
	@echo "Running go vet..."
	go vet ./...

# Run linters
lint: fmt-check vet
	@echo "Running linters..."
	@which golangci-lint > /dev/null || (echo "golangci-lint not installed. Install from https://golangci-lint.run/usage/install/" && exit 1)
	golangci-lint run ./...

# Validate example Terraform configurations
validate:
	@echo "Validating Terraform examples..."
	@cd examples && terraform init -backend=false || true
	@cd examples && terraform validate || true

# Run CI checks locally
ci: fmt-check vet test
	@echo "All CI checks passed!"

# Run all checks (fmt, vet, test)
check: fmt vet test
	@echo "All checks passed!"

# Tidy go modules
tidy:
	@echo "Tidying go modules..."
	go mod tidy

# Verify go mod tidy doesn't change anything
tidy-check:
	@echo "Checking go mod tidy..."
	@go mod tidy
	@if [ -n "$$(git status --porcelain go.mod go.sum)" ]; then \
		echo "go.mod or go.sum is not tidy. Run 'make tidy' and commit changes."; \
		git diff go.mod go.sum; \
		exit 1; \
	fi
	@echo "go.mod and go.sum are tidy!"

# Show version
version:
	@cat VERSION

# Bump version (usage: make version-bump VERSION=x.y.z)
version-bump:
	@if [ -z "$(VERSION)" ]; then \
		echo "Usage: make version-bump VERSION=x.y.z"; \
		exit 1; \
	fi
	@echo $(VERSION) > VERSION
	@echo "Version bumped to $(VERSION)"
	@echo "Don't forget to update CHANGELOG.md and create a git tag!"

# Create a release (usage: make release VERSION=x.y.z)
release:
	@if [ -z "$(VERSION)" ]; then \
		echo "Usage: make release VERSION=x.y.z"; \
		exit 1; \
	fi
	@echo "Creating release $(VERSION)..."
	@make version-bump VERSION=$(VERSION)
	@git add VERSION CHANGELOG.md
	@git commit -m "Release $(VERSION)"
	@git tag -a "v$(VERSION)" -m "Release $(VERSION)"
	@echo "Release $(VERSION) created. Push with: git push && git push --tags"

# Help
help:
	@echo "Available targets:"
	@echo "  build          - Build the provider binary"
	@echo "  install        - Install the provider locally"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  clean          - Clean build artifacts"
	@echo "  fmt            - Format code"
	@echo "  fmt-check      - Check if code is formatted"
	@echo "  vet            - Run go vet"
	@echo "  lint           - Run all linters"
	@echo "  validate       - Validate Terraform examples"
	@echo "  ci             - Run CI checks locally (fmt-check, vet, test)"
	@echo "  check          - Run all checks (fmt, vet, test)"
	@echo "  tidy           - Tidy go modules"
	@echo "  tidy-check     - Check if go.mod is tidy"
	@echo "  version        - Show current version"
	@echo "  version-bump   - Bump version (usage: make version-bump VERSION=x.y.z)"
	@echo "  release        - Create a release (usage: make release VERSION=x.y.z)"
	@echo "  help           - Show this help message"
