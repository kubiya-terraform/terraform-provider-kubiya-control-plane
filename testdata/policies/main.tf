terraform {
  required_providers {
    controlplane = {
      source = "kubiya/control-plane"
    }
  }
}

provider "controlplane" {
  # Configuration via environment variables:
  # KUBIYA_CONTROL_PLANE_API_KEY
  # KUBIYA_CONTROL_PLANE_ORG_ID
  # KUBIYA_CONTROL_PLANE_BASE_URL (optional, defaults to https://control-plane.kubiya.ai)
}

# Test policy resource
resource "controlplane_policy" "test" {
  name        = "test-policy"
  description = "Test policy for automated testing"
  enabled     = true

  # OPA Rego policy content
  policy_content = <<-EOT
    package kubiya.test

    # Simple test policy
    default allow = true

    deny[msg] {
      input.operation = "delete"
      count(input.approvals) < 1
      msg := "Delete operations require at least one approval"
    }

    deny[msg] {
      input.resource.sensitive = true
      not input.user.role = "admin"
      msg := "Only admins can access sensitive resources"
    }
  EOT

  tags = ["test", "automated-testing"]
}

# Test data source lookup
data "controlplane_policy" "test_lookup" {
  id = controlplane_policy.test.id
}

output "policy_id" {
  value = controlplane_policy.test.id
}

output "policy_name" {
  value = data.controlplane_policy.test_lookup.name
}

output "policy_enabled" {
  value = data.controlplane_policy.test_lookup.enabled
}

output "policy_content" {
  value     = data.controlplane_policy.test_lookup.policy_content
  sensitive = true
}
