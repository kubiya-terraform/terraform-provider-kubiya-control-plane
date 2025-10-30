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

# Test toolset resource
resource "controlplane_toolset" "test" {
  name        = "test-toolset"
  description = "Test toolset for automated testing"
  type        = "shell"
  enabled     = true

  configuration = jsonencode({
    allowed_commands = ["ls", "pwd", "echo"]
    timeout          = 30
    working_dir      = "/tmp"
  })
}

# Test data source lookup
data "controlplane_toolset" "test_lookup" {
  id = controlplane_toolset.test.id
}

output "toolset_id" {
  value = controlplane_toolset.test.id
}

output "toolset_name" {
  value = data.controlplane_toolset.test_lookup.name
}

output "toolset_type" {
  value = data.controlplane_toolset.test_lookup.type
}

output "toolset_enabled" {
  value = data.controlplane_toolset.test_lookup.enabled
}

output "toolset_configuration" {
  value     = data.controlplane_toolset.test_lookup.configuration
  sensitive = true
}
