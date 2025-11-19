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

# Minimal agent (required fields only)
resource "controlplane_agent" "minimal" {
  name = "test-agent-minimal"
}

# Data source test
data "controlplane_agent" "minimal_lookup" {
  id = controlplane_agent.minimal.id
}

# Outputs for test validation
output "agent_id" {
  value = controlplane_agent.minimal.id
}

output "agent_name" {
  value = controlplane_agent.minimal.name
}

output "agent_created_at" {
  value = controlplane_agent.minimal.created_at
}

output "data_agent_name" {
  value = data.controlplane_agent.minimal_lookup.name
}
