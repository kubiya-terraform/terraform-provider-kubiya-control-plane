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

# Test team resource
resource "kubiya_control_plane_team" "test" {
  name        = "test-team"
  description = "Test team for automated testing"
}

# Test data source lookup
data "kubiya_control_plane_team" "test_lookup" {
  id = kubiya_control_plane_team.test.id
}

output "team_id" {
  value = kubiya_control_plane_team.test.id
}

output "team_name" {
  value = data.kubiya_control_plane_team.test_lookup.name
}

output "team_description" {
  value = data.kubiya_control_plane_team.test_lookup.description
}
