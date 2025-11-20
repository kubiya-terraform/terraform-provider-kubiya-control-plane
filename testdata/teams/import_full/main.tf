
provider "controlplane" {}

variable "team_id" {
  type        = string
  description = "Team ID to import"
}

resource "controlplane_team" "imported_full" {
  name        = "test-team-full"
  description = "Comprehensive test team with all fields configured"
  status      = "active"
  runtime     = "default"
}

output "imported_team_id" {
  value = controlplane_team.imported_full.id
}

output "imported_team_name" {
  value = controlplane_team.imported_full.name
}

output "imported_team_description" {
  value = controlplane_team.imported_full.description
}

output "imported_team_runtime" {
  value = controlplane_team.imported_full.runtime
}

output "imported_team_created_at" {
  value = controlplane_team.imported_full.created_at
}
