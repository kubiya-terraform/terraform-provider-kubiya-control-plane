
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

output "team_created_at" {
  value = controlplane_team.minimal.created_at
}

output "team_updated_at" {
  value = controlplane_team.minimal.updated_at
}

output "data_team_name" {
  value = data.controlplane_team.minimal_lookup.name
}
