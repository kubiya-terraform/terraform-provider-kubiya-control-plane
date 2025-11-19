terraform {
  required_providers {
    controlplane = {
      source = "kubiya/control-plane"
    }
  }
}

provider "controlplane" {}

# Full policy with all optional fields
resource "controlplane_policy" "full" {
  name        = "test-policy-full"
  description = "Comprehensive test policy with all fields configured"
  policy_type = "rego"
  enabled     = true

  policy_content = <<-EOT
    package test_full

    default allow = false

    allow {
      input.user.role == "admin"
    }

    allow {
      input.user.role == "developer"
      input.action == "read"
    }
  EOT

  tags = ["test", "comprehensive", "rbac"]
}

# Data source test
data "controlplane_policy" "full_lookup" {
  id = controlplane_policy.full.id
}

# Outputs
output "policy_id" {
  value = controlplane_policy.full.id
}

output "policy_name" {
  value = controlplane_policy.full.name
}

output "policy_description" {
  value = controlplane_policy.full.description
}

output "policy_type" {
  value = controlplane_policy.full.policy_type
}

output "policy_enabled" {
  value = controlplane_policy.full.enabled
}

output "data_policy_description" {
  value = data.controlplane_policy.full_lookup.description
}

output "data_policy_enabled" {
  value = data.controlplane_policy.full_lookup.enabled
}
