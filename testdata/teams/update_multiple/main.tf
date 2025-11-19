terraform {
  required_providers {
    controlplane = {
      source = "kubiya/control-plane"
    }
  }
}

provider "controlplane" {}

variable "description" {
  type    = string
  default = "Original team description"
}

variable "status" {
  type    = string
  default = "active"
}

variable "runtime" {
  type    = string
  default = "default"
}

resource "controlplane_team" "test" {
  name        = "test-team-update-multiple"
  description = var.description
  status      = var.status
  runtime     = var.runtime
}

output "team_id" {
  value = controlplane_team.test.id
}

output "team_description" {
  value = controlplane_team.test.description
}

output "team_status" {
  value = controlplane_team.test.status
}

output "team_runtime" {
  value = controlplane_team.test.runtime
}

output "team_updated_at" {
  value = controlplane_team.test.updated_at
}
