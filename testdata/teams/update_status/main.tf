
provider "controlplane" {}

variable "status" {
  type    = string
  default = "active"
}

resource "controlplane_team" "test" {
  name   = "test-team-update-status"
  status = var.status
}

output "team_id" {
  value = controlplane_team.test.id
}

output "team_status" {
  value = controlplane_team.test.status
}

output "team_updated_at" {
  value = controlplane_team.test.updated_at
}
