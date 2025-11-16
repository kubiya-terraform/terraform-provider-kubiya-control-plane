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
  # KUBIYA_CONTROL_PLANE_BASE_URL (optional, defaults to https://control-plane.kubiya.ai)
}

# Test project resource
resource "controlplane_project" "test" {
  name        = "test-project"
  key         = "TEST"
  description = "Test project for automated testing"

  # Project settings
  settings = jsonencode({
    owner       = "devops-team"
    environment = "test"
    cost_center = "engineering"
  })
}

# Test data source lookup
data "controlplane_project" "test_lookup" {
  id = controlplane_project.test.id
}

output "project_id" {
  value = controlplane_project.test.id
}

output "project_name" {
  value = data.controlplane_project.test_lookup.name
}

output "project_description" {
  value = data.controlplane_project.test_lookup.description
}

output "project_settings" {
  value     = data.controlplane_project.test_lookup.settings
  sensitive = true
}
