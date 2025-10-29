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

# Test project resource
resource "kubiya_control_plane_project" "test" {
  name        = "test-project"
  description = "Test project for automated testing"
}

# Test data source lookup
data "kubiya_control_plane_project" "test_lookup" {
  id = kubiya_control_plane_project.test.id
}

output "project_id" {
  value = kubiya_control_plane_project.test.id
}

output "project_name" {
  value = data.kubiya_control_plane_project.test_lookup.name
}

output "project_description" {
  value = data.kubiya_control_plane_project.test_lookup.description
}
