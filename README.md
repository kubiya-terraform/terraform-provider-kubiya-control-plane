# Control Plane Terraform Provider

[![CI](https://github.com/kubiya-terraform/kubiya-control-plane/actions/workflows/ci.yml/badge.svg)](https://github.com/kubiya-terraform/kubiya-control-plane/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/kubiya-terraform/kubiya-control-plane/branch/main/graph/badge.svg)](https://codecov.io/gh/kubiya-terraform/kubiya-control-plane)
[![Go Version](https://img.shields.io/github/go-mod/go-version/kubiya-terraform/kubiya-control-plane)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

Terraform provider for managing Kubiya Control Plane resources.

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.22

## Building The Provider

1. Clone the repository
2. Enter the repository directory
3. Build the provider using the Go `install` command:

```shell
go install
```

## Using the Provider

### Authentication

The provider requires an API key and organization ID to authenticate with the Kubiya Control Plane API. Set the following environment variables:

```shell
export KUBIYA_CONTROL_PLANE_API_KEY=your_api_key_here
export KUBIYA_CONTROL_PLANE_ORG_ID=your_organization_id_here
export KUBIYA_CONTROL_PLANE_BASE_URL=http://localhost:7777  # Optional: override base URL (defaults to https://control-plane.kubiya.ai)
```

### Example Usage

```terraform
terraform {
  required_providers {
    controlplane = {
      source = "kubiya/control-plane"
    }
  }
}

provider "controlplane" {
  # Configuration is read from environment variables:
  # - KUBIYA_CONTROL_PLANE_API_KEY (required)
  # - KUBIYA_CONTROL_PLANE_ORG_ID (required)
  # - KUBIYA_CONTROL_PLANE_BASE_URL (optional, defaults to https://control-plane.kubiya.ai)
}

# Create a project
resource "controlplane_project" "example" {
  name        = "example-project"
  key         = "EX"
  description = "Example project"

  settings = jsonencode({
    owner       = "devops-team"
    environment = "production"
  })
}

# Create an environment
resource "controlplane_environment" "production" {
  name        = "production"
  description = "Production environment for agents"

  settings = jsonencode({
    region     = "us-east-1"
    max_workers = 5
  })
}

# Create a team
resource "controlplane_team" "example" {
  name        = "example-team"
  description = "Example team"

  configuration = jsonencode({
    max_agents = 10
  })
}

# Create an agent
resource "controlplane_agent" "example" {
  name        = "example-agent"
  description = "Example AI agent"
  model_id    = "gpt-4"
  runtime     = "default"

  llm_config = jsonencode({
    temperature = 0.7
    max_tokens  = 2000
  })
}
```

## Resources

The provider currently supports the following resources:

- `controlplane_agent` - Manage AI agents
- `controlplane_team` - Manage teams
- `controlplane_project` - Manage projects
- `controlplane_environment` - Manage execution environments
- `controlplane_skill` - Manage skills (filesystem, shell, docker)
- `controlplane_policy` - Manage OPA Rego governance policies
- `controlplane_worker` - Register and manage workers

## Data Sources

The provider supports the following data sources for read-only lookups:

- `controlplane_agent` - Lookup existing agents by ID
- `controlplane_team` - Lookup existing teams by ID
- `controlplane_project` - Lookup existing projects by ID
- `controlplane_environment` - Lookup existing environments by ID
- `controlplane_skill` - Lookup existing skills by ID
- `controlplane_policy` - Lookup existing policies by ID

### Example Data Source Usage

```terraform
# Lookup an existing agent
data "controlplane_agent" "existing" {
  id = "agent-uuid-here"
}

# Use the data source in other resources
resource "controlplane_agent" "new_agent" {
  name        = "new-agent"
  description = "New agent based on ${data.controlplane_agent.existing.name}"
  model_id    = data.controlplane_agent.existing.model_id
  runtime     = data.controlplane_agent.existing.runtime

  llm_config = data.controlplane_agent.existing.llm_config
}

# Lookup a project and use its information
data "controlplane_project" "ml_project" {
  id = "project-uuid-here"
}

output "project_info" {
  value = {
    name     = data.controlplane_project.ml_project.name
    key      = data.controlplane_project.ml_project.key
    settings = data.controlplane_project.ml_project.settings
  }
}
```

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `go generate`.

### Testing

```shell
go test ./...
```

### Local Development

For local development and testing, you can use the following configuration in your `~/.terraformrc` file:

```hcl
provider_installation {
  dev_overrides {
    "kubiya/control-plane" = "/path/to/your/GOPATH/bin"
  }

  direct {}
}
```

Then build and install the provider:

```shell
go install
```

Now you can use the provider in your Terraform configurations without having to publish it to a registry.

## Directory Structure

- `internal/provider/` - Provider implementation
- `internal/clients/` - API client implementations
- `internal/entities/` - Data models and entities
- `examples/` - Example Terraform configurations
- `docs/` - Provider documentation
- `test/` - Integration tests

## Adding New Resources

1. Create a new file in `internal/entities/` for your data model
2. Create a new file in `internal/clients/` for API operations
3. Create a new file in `internal/provider/` for the resource implementation
4. Register the resource in `internal/provider/provider.go`
5. Add examples in `examples/`
6. Add documentation in `docs/resources/`

## License

[Your License Here]
