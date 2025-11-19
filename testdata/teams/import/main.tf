terraform {
  required_providers {
    controlplane = {
      source = "kubiya/control-plane"
    }
  }
}

provider "controlplane" {}

variable "team_id" {
  type        = string
  description = "Team ID to import"
}

variable "team_name" {
  type        = string
  description = "Team name for import configuration"
}

resource "controlplane_team" "imported" {
  name = var.team_name
}

output "imported_team_id" {
  value = controlplane_team.imported.id
}

output "imported_team_name" {
  value = controlplane_team.imported.name
}

output "imported_team_created_at" {
  value = controlplane_team.imported.created_at
}
