terraform {
  required_providers {
    controlplane = {
      source = "kubiya/control-plane"
    }
  }
}

provider "controlplane" {}

variable "team_name" {
  type    = string
  default = "test-team-update-name"
}

resource "controlplane_team" "test" {
  name = var.team_name
}

output "team_id" {
  value = controlplane_team.test.id
}

output "team_name" {
  value = controlplane_team.test.name
}

output "team_updated_at" {
  value = controlplane_team.test.updated_at
}
