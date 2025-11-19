terraform {
  required_providers {
    controlplane = {
      source = "kubiya/control-plane"
    }
  }
}

provider "controlplane" {}

# Minimal policy (required fields only)
resource "controlplane_policy" "minimal" {
  name         = "test-policy-minimal"
  policy_type  = "rego"
  policy_content = <<-EOT
    package test_minimal
    default allow = true
  EOT
}

# Data source test
data "controlplane_policy" "minimal_lookup" {
  id = controlplane_policy.minimal.id
}

# Outputs
output "policy_id" {
  value = controlplane_policy.minimal.id
}

output "policy_name" {
  value = controlplane_policy.minimal.name
}

output "policy_type" {
  value = controlplane_policy.minimal.policy_type
}

output "data_policy_name" {
  value = data.controlplane_policy.minimal_lookup.name
}
