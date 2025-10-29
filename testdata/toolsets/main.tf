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

# Test toolset resource
resource "kubiya_control_plane_toolset" "test" {
  name    = "test-toolset"
  type    = "shell"
  enabled = true

  configuration = jsonencode({
    allowed_commands = ["ls", "pwd", "echo"]
    timeout          = 30
  })
}

# Test data source lookup
data "kubiya_control_plane_toolset" "test_lookup" {
  id = kubiya_control_plane_toolset.test.id
}

output "toolset_id" {
  value = kubiya_control_plane_toolset.test.id
}

output "toolset_name" {
  value = data.kubiya_control_plane_toolset.test_lookup.name
}

output "toolset_type" {
  value = data.kubiya_control_plane_toolset.test_lookup.type
}

output "toolset_enabled" {
  value = data.kubiya_control_plane_toolset.test_lookup.enabled
}
