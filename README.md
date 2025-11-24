# Control Plane Terraform Provider

[![CI](https://github.com/kubiya-terraform/terraform-provider-kubiya-control-plane/actions/workflows/ci.yml/badge.svg)](https://github.com/kubiya-terraform/terraform-provider-kubiya-control-plane/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/kubiya-terraform/terraform-provider-kubiya-control-plane/branch/main/graph/badge.svg)](https://codecov.io/gh/kubiya-terraform/terraform-provider-kubiya-control-plane)
[![Go Version](https://img.shields.io/github/go-mod/go-version/kubiya-terraform/terraform-provider-kubiya-control-plane)](https://go.dev/)
[![License](https://img.shields.io/badge/license-AGPL--3.0-blue.svg)](LICENSE)

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

The provider requires an API key to authenticate with the Kubiya Control Plane API. Set the following environment variables:

```shell
export KUBIYA_CONTROL_PLANE_API_KEY=your_api_key_here
export KUBIYA_CONTROL_PLANE_BASE_URL=http://localhost:7777  # Optional: override base URL (defaults to https://control-plane.kubiya.ai)
```

### Example Usage

```terraform
terraform {
  required_providers {
    kubiya = {
      source = "kubiya/kubiya-control-plane"
    }
  }
}

provider "kubiya" {
  # Configuration is read from environment variables:
  # - KUBIYA_CONTROL_PLANE_API_KEY (required)
  # - KUBIYA_CONTROL_PLANE_BASE_URL (optional, defaults to https://control-plane.kubiya.ai)
}

# Create a project
resource "kubiya_project" "example" {
  name        = "example-project"
  key         = "EX"
  description = "Example project"

  settings = jsonencode({
    owner       = "devops-team"
    environment = "production"
  })
}

# Create an environment
resource "kubiya_environment" "production" {
  name        = "production"
  description = "Production environment for agents"

  settings = jsonencode({
    region     = "us-east-1"
    max_workers = 5
  })
}

# Create a team
resource "kubiya_team" "example" {
  name        = "example-team"
  description = "Example team"

  configuration = jsonencode({
    max_agents = 10
  })
}

# Create an agent
resource "kubiya_agent" "example" {
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

- `kubiya_agent` - Manage AI agents
- `kubiya_team` - Manage teams
- `kubiya_project` - Manage projects
- `kubiya_environment` - Manage execution environments
- `kubiya_skill` - Manage skills (filesystem, shell, docker)
- `kubiya_policy` - Manage OPA Rego governance policies
- `kubiya_worker` - Register and manage workers

## Data Sources

The provider supports the following data sources for read-only lookups:

- `kubiya_agent` - Lookup existing agents by ID
- `kubiya_team` - Lookup existing teams by ID
- `kubiya_project` - Lookup existing projects by ID
- `kubiya_environment` - Lookup existing environments by ID
- `kubiya_skill` - Lookup existing skills by ID
- `kubiya_policy` - Lookup existing policies by ID

### Example Data Source Usage

```terraform
# Lookup an existing agent
data "kubiya_agent" "existing" {
  id = "agent-uuid-here"
}

# Use the data source in other resources
resource "kubiya_agent" "new_agent" {
  name        = "new-agent"
  description = "New agent based on ${data.kubiya_agent.existing.name}"
  model_id    = data.kubiya_agent.existing.model_id
  runtime     = data.kubiya_agent.existing.runtime

  llm_config = data.kubiya_agent.existing.llm_config
}

# Lookup a project and use its information
data "kubiya_project" "ml_project" {
  id = "project-uuid-here"
}

output "project_info" {
  value = {
    name     = data.kubiya_project.ml_project.name
    key      = data.kubiya_project.ml_project.key
    settings = data.kubiya_project.ml_project.settings
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
    "kubiya/kubiya-control-plane" = "/path/to/your/GOPATH/bin"
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

This project is licensed under the GNU Affero General Public License v3.0 (AGPL-3.0). See the [LICENSE](LICENSE) file for details.

Copyright (C) 2025 Kubiya

This program is free software: you can redistribute it and/or modify it under the terms of the GNU Affero General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

For more information, visit https://www.gnu.org/licenses/agpl-3.0.html
