terraform {
  required_providers {
    controlplane = {
      source = "kubiya/control-plane"
    }
  }
}

provider "controlplane" {
  # Configuration via environment variables
}

variable "assign_team" {
  type    = bool
  default = false
}

# Team to assign to agent
resource "controlplane_team" "test" {
  name = "test-team-for-agent"
}

# Agent without team initially, team added via update
resource "controlplane_agent" "team_test" {
  name    = "test-agent-update-team"
  team_id = var.assign_team ? controlplane_team.test.id : null
}

# Outputs for test validation
output "agent_id" {
  value = controlplane_agent.team_test.id
}

output "team_id" {
  value = controlplane_team.test.id
}

output "agent_team_id" {
  value = controlplane_agent.team_test.team_id
}

output "agent_name" {
  value = controlplane_agent.team_test.name
}
