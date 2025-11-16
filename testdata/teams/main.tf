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

# Test team resource
resource "controlplane_team" "test" {
  name        = "test-team"
  description = "Test team for automated testing"

  # Team configuration
  configuration = jsonencode({
    max_agents        = 10
    default_runtime   = "default"
    enable_monitoring = true
  })
}

# Test data source lookup
data "controlplane_team" "test_lookup" {
  id = controlplane_team.test.id
}

output "team_id" {
  value = controlplane_team.test.id
}

output "team_name" {
  value = data.controlplane_team.test_lookup.name
}

output "team_description" {
  value = data.controlplane_team.test_lookup.description
}

output "team_configuration" {
  value     = data.controlplane_team.test_lookup.configuration
  sensitive = true
}
