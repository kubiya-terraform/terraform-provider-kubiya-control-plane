terraform {
  required_providers {
    controlplane = {
      source = "kubiya/control-plane"
    }
  }
}

provider "controlplane" {}

# Team with claude_code runtime
resource "controlplane_team" "claude_code" {
  name        = "test-team-claude-code"
  description = "Team using Claude Code SDK runtime"
  runtime     = "claude_code"
  status      = "active"

  configuration = jsonencode({
    sdk_version = "latest"
  })
}

# Data source test
data "controlplane_team" "claude_code_lookup" {
  id = controlplane_team.claude_code.id
}

# Outputs
output "team_id" {
  value = controlplane_team.claude_code.id
}

output "team_runtime" {
  value = controlplane_team.claude_code.runtime
}

output "data_team_runtime" {
  value = data.controlplane_team.claude_code_lookup.runtime
}
