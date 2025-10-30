terraform {
  required_providers {
    controlplane = {
      source = "kubiya/control-plane"
    }
  }
}

provider "controlplane" {
  # Configuration via environment variables:
  # KUBIYA_CONTROL_PLANE_API_KEY
  # KUBIYA_CONTROL_PLANE_ORG_ID
  # KUBIYA_CONTROL_PLANE_BASE_URL (optional, defaults to https://control-plane.kubiya.ai)
}

# Test environment resource
resource "controlplane_environment" "test" {
  name        = "test-environment"
  description = "Test environment for automated testing"

  # Environment configuration
  settings = jsonencode({
    region         = "us-east-1"
    max_workers    = 5
    auto_scaling   = true
    retention_days = 30
  })

  # Execution environment variables (secrets, integrations)
  execution_environment = jsonencode({
    env_vars = {
      LOG_LEVEL = "info"
      APP_ENV   = "test"
    }
  })
}

# Test data source lookup
data "controlplane_environment" "test_lookup" {
  id = controlplane_environment.test.id
}

output "environment_id" {
  value = controlplane_environment.test.id
}

output "environment_name" {
  value = data.controlplane_environment.test_lookup.name
}

output "environment_settings" {
  value     = data.controlplane_environment.test_lookup.settings
  sensitive = true
}

output "environment_execution_environment" {
  value     = data.controlplane_environment.test_lookup.execution_environment
  sensitive = true
}
