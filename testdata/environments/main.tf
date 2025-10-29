terraform {
  required_providers {
    kubiya_control_plane = {
      source = "kubiya/control-plane"
    }
  }
}

provider "kubiya_control_plane" {
  # Configuration via environment variables:
  # KUBIYA_CONTROL_PLANE_API_KEY
  # KUBIYA_CONTROL_PLANE_ORG_ID
}

# Test environment resource
resource "kubiya_control_plane_environment" "test" {
  name        = "test-environment"
  description = "Test environment for automated testing"

  configuration = jsonencode({
    region = "us-east-1"
    tier   = "test"
  })
}

# Test data source lookup
data "kubiya_control_plane_environment" "test_lookup" {
  id = kubiya_control_plane_environment.test.id
}

output "environment_id" {
  value = kubiya_control_plane_environment.test.id
}

output "environment_name" {
  value = data.kubiya_control_plane_environment.test_lookup.name
}

output "environment_configuration" {
  value = data.kubiya_control_plane_environment.test_lookup.configuration
}
