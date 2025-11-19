terraform {
  required_providers {
    controlplane = {
      source = "kubiya/control-plane"
    }
  }
}

provider "controlplane" {}

variable "environment_id" {
  type        = string
  description = "Environment ID to import"
}

variable "environment_name" {
  type        = string
  description = "Environment name"
}

resource "controlplane_environment" "imported" {
  name = var.environment_name
}

output "imported_environment_id" {
  value = controlplane_environment.imported.id
}

output "imported_environment_name" {
  value = controlplane_environment.imported.name
}
