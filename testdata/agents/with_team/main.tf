
provider "controlplane" {
  # Configuration via environment variables:
  # KUBIYA_CONTROL_PLANE_API_KEY
  # KUBIYA_CONTROL_PLANE_BASE_URL (optional, defaults to https://control-plane.kubiya.ai)
}

# First create a team
resource "controlplane_team" "test_team" {
  name        = "test-team-for-agent"
  description = "Team for agent assignment testing"
}

# Agent with team assignment
resource "controlplane_agent" "with_team" {
  name        = "test-agent-with-team"
  description = "Agent assigned to a team"
  team_id     = controlplane_team.test_team.id

  capabilities = ["team_collaboration"]
}

# Data source test
data "controlplane_agent" "with_team_lookup" {
  id = controlplane_agent.with_team.id
}

# Outputs for test validation
output "agent_id" {
  value = controlplane_agent.with_team.id
}

output "agent_team_id" {
  value = controlplane_agent.with_team.team_id
}

output "team_id" {
  value = controlplane_team.test_team.id
}

output "data_agent_team_id" {
  value = data.controlplane_agent.with_team_lookup.team_id
}
