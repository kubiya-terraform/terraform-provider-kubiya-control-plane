terraform {
  required_providers {
    controlplane = {
      source = "kubiya/control-plane"
    }
  }
}

provider "controlplane" {}

# Full skill with all optional fields
resource "controlplane_skill" "full" {
  name        = "test-skill-full"
  type        = "python"
  description = "Comprehensive test skill with all fields configured"
  enabled     = true

  content = <<-EOT
    def main():
        print("Test skill execution")
        return {"status": "success"}
  EOT

  configuration = jsonencode({
    timeout     = 300
    memory_mb   = 512
    environment = "production"
  })
}

# Data source test
data "controlplane_skill" "full_lookup" {
  id = controlplane_skill.full.id
}

# Outputs
output "skill_id" {
  value = controlplane_skill.full.id
}

output "skill_name" {
  value = controlplane_skill.full.name
}

output "skill_type" {
  value = controlplane_skill.full.type
}

output "skill_description" {
  value = controlplane_skill.full.description
}

output "skill_enabled" {
  value = controlplane_skill.full.enabled
}

output "data_skill_description" {
  value = data.controlplane_skill.full_lookup.description
}
