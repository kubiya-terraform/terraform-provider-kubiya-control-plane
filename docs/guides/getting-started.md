---
page_title: "Getting Started with Kubiya Control Plane"
subcategory: "Guides"
description: |-
  A guide to getting started with the Kubiya Control Plane Terraform Provider
---

# Getting Started with Kubiya Control Plane

This guide will walk you through setting up and using the Kubiya Control Plane Terraform Provider.

## Prerequisites

- Terraform 1.0 or later installed
- A Kubiya Control Plane account
- API credentials (API Key and Organization ID)

## Step 1: Obtain Credentials

1. Log in to your Kubiya Control Plane account
2. Navigate to **Settings** → **API Keys**
3. Generate a new API key
4. Copy your Organization ID from the account settings

## Step 2: Configure Environment

Set the required environment variables:

```bash
export KUBIYA_CONTROL_PLANE_API_KEY="your-api-key-here"
export KUBIYA_CONTROL_PLANE_ORG_ID="your-org-id-here"
# Optional: Override the default API URL
# export KUBIYA_CONTROL_PLANE_BASE_URL="https://custom-url.example.com"
```

## Step 3: Create Your First Configuration

Create a file named `main.tf`:

```terraform
terraform {
  required_providers {
    kubiya_control_plane = {
      source  = "kubiya/control-plane"
      version = "~> 1.0"
    }
  }
}

provider "kubiya_control_plane" {
  # Configuration via environment variables
}

# Create an environment
resource "kubiya_control_plane_environment" "dev" {
  name        = "development"
  description = "Development environment"

  configuration = jsonencode({
    region = "us-east-1"
  })
}

# Create a team
resource "kubiya_control_plane_team" "platform" {
  name        = "platform-team"
  description = "Platform engineering team"
}

# Create an agent
resource "kubiya_control_plane_agent" "assistant" {
  name        = "dev-assistant"
  description = "Development assistant agent"
  model_id    = "gpt-4"
  runtime     = "default"
  team_id     = kubiya_control_plane_team.platform.id

  llm_config = jsonencode({
    temperature = 0.7
    max_tokens  = 2000
  })
}

# Output important information
output "environment_id" {
  value       = kubiya_control_plane_environment.dev.id
  description = "Development environment ID"
}

output "agent_id" {
  value       = kubiya_control_plane_agent.assistant.id
  description = "Assistant agent ID"
}
```

## Step 4: Initialize and Apply

Initialize Terraform:

```bash
terraform init
```

Review the planned changes:

```bash
terraform plan
```

Apply the configuration:

```bash
terraform apply
```

## Step 5: Verify Resources

After applying, you can verify the created resources:

```bash
# View outputs
terraform output

# Look up resources using data sources
```

Create a file named `data.tf`:

```terraform
# Look up the created agent
data "kubiya_control_plane_agent" "assistant" {
  id = kubiya_control_plane_agent.assistant.id
}

output "agent_status" {
  value = data.kubiya_control_plane_agent.assistant.status
}
```

## Step 6: Add More Resources

Expand your configuration by adding toolsets and policies:

```terraform
# Add a shell toolset
resource "kubiya_control_plane_toolset" "shell" {
  name    = "shell-commands"
  type    = "shell"
  enabled = true

  configuration = jsonencode({
    allowed_commands = ["ls", "cat", "grep"]
    timeout          = 30
  })
}

# Add a security policy
resource "kubiya_control_plane_policy" "security" {
  name        = "basic-security"
  description = "Basic security policy"
  enabled     = true

  policy_content = <<-EOT
    package kubiya.security

    deny[msg] {
      input.operation = "delete"
      count(input.approvals) < 1
      msg := "Delete operations require approval"
    }
  EOT

  tags = ["security"]
}
```

## Best Practices

1. **Use Variables**: Define reusable values in variables
2. **Remote State**: Store state in S3 or Terraform Cloud
3. **Modules**: Create reusable modules for common patterns
4. **Version Pinning**: Pin provider versions for stability
5. **Sensitive Data**: Mark sensitive outputs appropriately
6. **Documentation**: Document your configurations

## Next Steps

- Explore the [example configurations](https://github.com/kubiya/terraform-provider-kubiya-control-plane/tree/main/examples)
- Review the [resource documentation](../resources)
- Check the [data source documentation](../data-sources)
- Join the community for support and discussions

## Cleanup

When you're done experimenting, clean up resources:

```bash
terraform destroy
```

## Troubleshooting

### Authentication Issues

If you encounter authentication errors:
- Verify environment variables are set correctly
- Check that your API key hasn't expired
- Ensure you have proper permissions

### Resource Not Found

If resources cannot be found:
- Verify the resource ID is correct
- Check you're in the correct environment
- Ensure the resource exists in your organization

### State Issues

If you have state conflicts:
- Check for concurrent Terraform runs
- Verify state file location
- Consider using state locking with remote backends

## Support

For help and support:
- GitHub Issues: [Report a bug](https://github.com/kubiya/terraform-provider-kubiya-control-plane/issues)
- Documentation: [Full documentation](https://docs.kubiya.ai)
- Email: support@kubiya.ai
