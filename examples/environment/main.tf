terraform {
  required_providers {
    kubiya_control_plane = {
      source = "kubiya/control-plane"
    }
  }
}

provider "kubiya_control_plane" {
  # Configuration is via environment variables:
  # KUBIYA_CONTROL_PLANE_API_KEY
  # KUBIYA_CONTROL_PLANE_ORG_ID
  # KUBIYA_CONTROL_PLANE_BASE_URL (optional, defaults to https://control-plane.kubiya.ai)
}

# Create an environment
resource "kubiya_control_plane_environment" "example" {
  name        = "production"
  description = "Production environment for agents"

  # Environment configuration
  configuration = jsonencode({
    region           = "us-east-1"
    max_workers      = 5
    auto_scaling     = true
    retention_days   = 30
    notification_url = "https://hooks.slack.com/example"
  })

  # Execution environment variables (secrets, integrations)
  execution_environment = jsonencode({
    env_vars = {
      LOG_LEVEL = "info"
      APP_ENV   = "production"
    }
  })
}

# Look up an existing environment by ID
data "kubiya_control_plane_environment" "existing" {
  id = "environment-uuid-here"
}

# Output environment information
output "environment_id" {
  value       = kubiya_control_plane_environment.example.id
  description = "The ID of the created environment"
}

output "environment_status" {
  value       = kubiya_control_plane_environment.example.status
  description = "The current status of the environment"
}

output "existing_environment_name" {
  value       = data.kubiya_control_plane_environment.existing.name
  description = "Name of the existing environment"
}

output "existing_environment_config" {
  value       = data.kubiya_control_plane_environment.existing.configuration
  description = "Configuration of the existing environment"
  sensitive   = true
}
