---
page_title: "Kubiya Control Plane Provider"
subcategory: ""
description: |-
  Terraform provider for managing Kubiya Control Plane resources including agents, teams, projects, environments, skills, and policies.
---

# Kubiya Control Plane Provider

The Kubiya Control Plane provider allows you to manage Kubiya platform resources using Terraform. This includes AI agents, teams, projects, execution environments, skills, and governance policies.

## Features

- **Agent Management**: Create and configure AI agents with custom LLM settings
- **Team Organization**: Organize agents into teams with shared configuration
- **Project Management**: Group resources into projects for better organization
- **Environment Configuration**: Manage execution environments with custom settings
- **Skill Integration**: Configure various skills (filesystem, shell, docker, etc.)
- **Policy Governance**: Implement OPA Rego policies for security and compliance
- **Worker Queue Management**: Configure and manage worker queues for task execution
- **Job Scheduling**: Create and manage scheduled, webhook-triggered, and manual jobs

## Authentication

The provider requires authentication via environment variables.

## Example Usage

```terraform
terraform {
  required_providers {
    controlplane = {
      source  = "kubiya/control-plane"
      version = "~> 1.0"
    }
  }
}

provider "controlplane" {
  # Configuration is via environment variables:
  # KUBIYA_CONTROL_PLANE_API_KEY
  # KUBIYA_CONTROL_PLANE_BASE_URL (optional)
}

resource "controlplane_agent" "example" {
  name     = "my-agent"
  model_id = "gpt-4"
  runtime  = "default"

  llm_config = jsonencode({
    temperature = 0.7
    max_tokens  = 2000
  })
}
```

## Environment Variables

- `KUBIYA_CONTROL_PLANE_API_KEY` (required) - Your Kubiya API key
- `KUBIYA_CONTROL_PLANE_BASE_URL` (optional) - Custom API base URL (defaults to https://control-plane.kubiya.ai)
