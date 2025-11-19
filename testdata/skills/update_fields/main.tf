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
  default = "Original skill description"
}

resource "controlplane_skill" "test" {
  name        = "test-skill-update"
  description = var.description
}

output "skill_id" {
  value = controlplane_skill.test.id
}

output "skill_name" {
  value = controlplane_skill.test.name
}

output "skill_description" {
  value = controlplane_skill.test.description
}

output "skill_updated_at" {
  value = controlplane_skill.test.updated_at
}
