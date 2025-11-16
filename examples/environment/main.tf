terraform {
  required_providers {
    controlplane = {
      source = "kubiya/control-plane"
    }
  }
}

provider "controlplane" {
  # Configuration is via environment variables:
  # KUBIYA_CONTROL_PLANE_API_KEY
  # KUBIYA_CONTROL_PLANE_BASE_URL (optional, defaults to https://control-plane.kubiya.ai)
}

# Create an environment
resource "controlplane_environment" "example" {
  name        = "test-env"
  description = "Production environment for agents"

  # Environment configuration
  settings = jsonencode({
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
data "controlplane_environment" "existing" {
  id = "1fe58e91-5a7f-460d-bafd-97dc8f6cb125"
}

# Output environment information
output "environment_id" {
  value       = controlplane_environment.example.id
  description = "The ID of the created environment"
}

output "environment_status" {
  value       = controlplane_environment.example.status
  description = "The current status of the environment"
}

output "existing_environment_name" {
  value       = data.controlplane_environment.existing.name
  description = "Name of the existing environment"
}

output "existing_environment_config" {
  value       = data.controlplane_environment.existing.settings
  description = "Configuration of the existing environment"
  sensitive   = true
}
