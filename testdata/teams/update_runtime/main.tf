terraform {
  required_providers {
    controlplane = {
      source = "kubiya/control-plane"
    }
  }
}

provider "controlplane" {}

variable "runtime" {
  type    = string
  default = "default"
}

resource "controlplane_team" "test" {
  name    = "test-team-update-runtime"
  runtime = var.runtime
}

output "team_id" {
  value = controlplane_team.test.id
}

output "team_runtime" {
  value = controlplane_team.test.runtime
}

output "team_updated_at" {
  value = controlplane_team.test.updated_at
}
