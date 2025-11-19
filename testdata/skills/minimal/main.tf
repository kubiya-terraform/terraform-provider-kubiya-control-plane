terraform {
  required_providers {
    controlplane = {
      source = "kubiya/control-plane"
    }
  }
}

provider "controlplane" {}

# Minimal skill (required fields only)
resource "controlplane_skill" "minimal" {
  name = "test-skill-minimal"
  type = "shell"
}

# Data source test
data "controlplane_skill" "minimal_lookup" {
  id = controlplane_skill.minimal.id
}

# Outputs
output "skill_id" {
  value = controlplane_skill.minimal.id
}

output "skill_name" {
  value = controlplane_skill.minimal.name
}

output "skill_type" {
  value = controlplane_skill.minimal.type
}

output "data_skill_type" {
  value = data.controlplane_skill.minimal_lookup.type
}
