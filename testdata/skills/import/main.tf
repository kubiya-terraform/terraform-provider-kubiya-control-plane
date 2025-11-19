terraform {
  required_providers {
    controlplane = {
      source = "kubiya/control-plane"
    }
  }
}

provider "controlplane" {}

variable "skill_id" {
  type        = string
  description = "Skill ID to import"
}

variable "skill_name" {
  type        = string
  description = "Skill name"
}

resource "controlplane_skill" "imported" {
  name = var.skill_name
}

output "imported_skill_id" {
  value = controlplane_skill.imported.id
}

output "imported_skill_name" {
  value = controlplane_skill.imported.name
}
