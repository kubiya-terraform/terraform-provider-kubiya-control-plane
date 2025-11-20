
provider "controlplane" {}

# Minimal policy (required fields only)
resource "controlplane_policy" "minimal" {
  name           = "test-policy-minimal"
  policy_type    = "rego"
  policy_content = <<-EOT
    package test_minimal
    default allow = true
  EOT
}

# Full policy with all optional fields
resource "controlplane_policy" "full" {
  name        = "test-policy-full-ds"
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

# Disabled policy for enabled flag testing
resource "controlplane_policy" "disabled" {
  name        = "test-policy-disabled-ds"
  description = "Disabled policy for testing"
  policy_type = "rego"
  enabled     = false

  policy_content = <<-EOT
    package test_disabled
    default allow = false
  EOT
}

# RBAC policy with tags
resource "controlplane_policy" "rbac" {
  name        = "test-policy-rbac-ds"
  description = "RBAC policy with tags"
  policy_type = "rego"
  enabled     = true

  policy_content = <<-EOT
    package rbac
    default allow = false
  EOT

  tags = ["rbac", "security", "access-control"]
}

# Data sources
data "controlplane_policy" "minimal_lookup" {
  id = controlplane_policy.minimal.id
}

data "controlplane_policy" "full_lookup" {
  id = controlplane_policy.full.id
}

data "controlplane_policy" "disabled_lookup" {
  id = controlplane_policy.disabled.id
}

data "controlplane_policy" "rbac_lookup" {
  id = controlplane_policy.rbac.id
}

# Outputs for tests
output "data_minimal_name" {
  value = data.controlplane_policy.minimal_lookup.name
}

output "data_full_description" {
  value = data.controlplane_policy.full_lookup.description
}

output "data_full_enabled" {
  value = tostring(data.controlplane_policy.full_lookup.enabled)
}

output "data_disabled_enabled" {
  value = tostring(data.controlplane_policy.disabled_lookup.enabled)
}

output "data_rbac_tags" {
  value = jsonencode(data.controlplane_policy.rbac_lookup.tags)
}
