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
  default = "Original environment description"
}

resource "controlplane_environment" "test" {
  name        = "test-environment-update"
  description = var.description
}

output "environment_id" {
  value = controlplane_environment.test.id
}

output "environment_name" {
  value = controlplane_environment.test.name
}

output "environment_description" {
  value = controlplane_environment.test.description
}

output "environment_updated_at" {
  value = controlplane_environment.test.updated_at
}
