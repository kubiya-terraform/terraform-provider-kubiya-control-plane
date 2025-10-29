terraform {
  required_providers {
    kubiya_control_plane = {
      source = "kubiya/control-plane"
    }
  }
}

provider "kubiya_control_plane" {
  # Configuration via environment variables:
  # KUBIYA_CONTROL_PLANE_API_KEY
  # KUBIYA_CONTROL_PLANE_ORG_ID
}

# Test policy resource
resource "kubiya_control_plane_policy" "test" {
  name        = "test-policy"
  description = "Test policy for automated testing"
  enabled     = true

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
}

# Test data source lookup
data "kubiya_control_plane_policy" "test_lookup" {
  id = kubiya_control_plane_policy.test.id
}

output "policy_id" {
  value = kubiya_control_plane_policy.test.id
}

output "policy_name" {
  value = data.kubiya_control_plane_policy.test_lookup.name
}

output "policy_enabled" {
  value = data.kubiya_control_plane_policy.test_lookup.enabled
}

output "policy_content" {
  value     = data.kubiya_control_plane_policy.test_lookup.policy
  sensitive = true
}
