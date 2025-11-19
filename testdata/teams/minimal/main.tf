terraform {
  required_providers {
    controlplane = {
      source = "kubiya/control-plane"
    }
  }
}

provider "controlplane" {}

# Minimal team (required fields only)
resource "controlplane_team" "minimal" {
  name = "test-team-minimal"
}

# Data source test
data "controlplane_team" "minimal_lookup" {
  id = controlplane_team.minimal.id
}

# Outputs
output "team_id" {
  value = controlplane_team.minimal.id
}

output "team_name" {
  value = controlplane_team.minimal.name
}

output "data_team_name" {
  value = data.controlplane_team.minimal_lookup.name
}
