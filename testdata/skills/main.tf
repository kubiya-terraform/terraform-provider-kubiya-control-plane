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

# Test skill resource
resource "controlplane_skill" "test" {
  name        = "test-skill"
  description = "Test skill for automated testing"
  type        = "shell"
  enabled     = true

  configuration = jsonencode({
    allowed_commands = ["ls", "pwd", "echo"]
    timeout          = 30
    working_dir      = "/tmp"
  })
}

# Test data source lookup
data "controlplane_skill" "test_lookup" {
  id = controlplane_skill.test.id
}

output "skill_id" {
  value = controlplane_skill.test.id
}

output "skill_name" {
  value = data.controlplane_skill.test_lookup.name
}

output "skill_type" {
  value = data.controlplane_skill.test_lookup.type
}

output "skill_enabled" {
  value = data.controlplane_skill.test_lookup.enabled
}

output "skill_configuration" {
  value     = data.controlplane_skill.test_lookup.configuration
  sensitive = true
}
