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
  # KUBIYA_CONTROL_PLANE_BASE_URL (optional, defaults to https://control-plane.kubiya.ai)
}

# Test 1: Minimal policy (required fields only)
resource "controlplane_policy" "minimal" {
  name = "test-policy-minimal"
  policy_content = <<-EOT
    package minimal_test
    default allow = true
  EOT
}

# Test 2: Full policy with all optional fields
resource "controlplane_policy" "full" {
  name        = "test-policy-full"
  description = "Comprehensive test policy with all fields configured"
  policy_type = "rego"
  enabled     = true

  policy_content = <<-EOT
    package kubiya.full_test

    # Comprehensive policy with multiple rules
    import future.keywords.if
    import future.keywords.in

    default allow = false

    # Allow admins to do anything
    allow if {
      input.user.role == "admin"
    }

    # Allow users to read their own resources
    allow if {
      input.operation == "read"
      input.resource.owner == input.user.id
    }

    # Deny sensitive operations without approval
    deny[msg] {
      input.operation in ["delete", "update"]
      input.resource.sensitive == true
      count(input.approvals) < 2
      msg := "Sensitive operations require at least 2 approvals"
    }

    # Deny production changes during business hours
    deny[msg] {
      input.environment == "production"
      is_business_hours
      not input.user.role == "admin"
      msg := "Production changes during business hours require admin role"
    }

    is_business_hours if {
      hour := time.clock([time.now_ns(), "UTC"])[0]
      hour >= 9
      hour < 17
    }
  EOT

  tags = ["comprehensive", "test", "production"]
}

# Test 3: Policy with simple allow rule
resource "controlplane_policy" "allow_all" {
  name        = "test-policy-allow-all"
  description = "Policy that allows all operations"
  enabled     = true

  policy_content = <<-EOT
    package allow_all_test

    default allow = true
  EOT

  tags = ["permissive"]
}

# Test 4: Policy with deny rule
resource "controlplane_policy" "deny_deletes" {
  name        = "test-policy-deny-deletes"
  description = "Policy that denies all delete operations"
  enabled     = true

  policy_content = <<-EOT
    package deny_deletes_test

    deny[msg] {
      input.operation == "delete"
      msg := "Delete operations are not allowed"
    }
  EOT

  tags = ["restrictive", "safety"]
}

# Test 5: Policy with RBAC rules
resource "controlplane_policy" "rbac" {
  name        = "test-policy-rbac"
  description = "Role-based access control policy"
  enabled     = true

  policy_content = <<-EOT
    package rbac_test

    import future.keywords.if

    default allow = false

    # Admins can do everything
    allow if {
      input.user.role == "admin"
    }

    # Developers can read and create
    allow if {
      input.user.role == "developer"
      input.operation in ["read", "create"]
    }

    # Viewers can only read
    allow if {
      input.user.role == "viewer"
      input.operation == "read"
    }
  EOT

  tags = ["rbac", "security"]
}

# Test 6: Policy with resource-based rules
resource "controlplane_policy" "resource_policy" {
  name        = "test-policy-resource-based"
  description = "Policy based on resource attributes"
  enabled     = true

  policy_content = <<-EOT
    package resource_test

    import future.keywords.if

    deny[msg] {
      input.resource.type == "secret"
      not input.user.clearance_level >= 3
      msg := "Insufficient clearance level for accessing secrets"
    }

    deny[msg] {
      input.resource.classification == "confidential"
      not input.user.department == input.resource.department
      msg := "Cross-department access to confidential resources is not allowed"
    }
  EOT

  tags = ["resource-based"]
}

# Test 7: Policy with time-based rules
resource "controlplane_policy" "time_based" {
  name        = "test-policy-time-based"
  description = "Policy with time-based restrictions"
  enabled     = true

  policy_content = <<-EOT
    package time_test

    import future.keywords.if

    deny[msg] {
      input.operation == "deploy"
      is_weekend
      msg := "Deployments are not allowed on weekends"
    }

    is_weekend if {
      day := time.weekday([time.now_ns(), "UTC"])
      day in [0, 6]  # Sunday = 0, Saturday = 6
    }
  EOT

  tags = ["time-based", "deployment"]
}

# Test 8: Disabled policy
resource "controlplane_policy" "disabled" {
  name        = "test-policy-disabled"
  description = "Disabled policy for testing"
  enabled     = false

  policy_content = <<-EOT
    package disabled_test

    default allow = true
  EOT
}

# Test 9: Policy with JSON type
resource "controlplane_policy" "json_type" {
  name        = "test-policy-json"
  description = "Policy with JSON type"
  policy_type = "json"
  enabled     = true

  policy_content = <<-EOT
    package json_test

    default allow = true
  EOT
}

# Test 10: Policy with no tags
resource "controlplane_policy" "no_tags" {
  name        = "test-policy-no-tags"
  description = "Policy without tags"
  enabled     = true

  policy_content = <<-EOT
    package no_tags_test

    default allow = true
  EOT

  tags = []
}

# Test 11: Policy with single tag
resource "controlplane_policy" "single_tag" {
  name        = "test-policy-single-tag"
  description = "Policy with single tag"
  enabled     = true

  policy_content = <<-EOT
    package single_tag_test

    default allow = true
  EOT

  tags = ["single"]
}

# Test 12: Policy with empty description
resource "controlplane_policy" "empty_description" {
  name        = "test-policy-empty-description"
  description = ""
  enabled     = true

  policy_content = <<-EOT
    package empty_desc_test

    default allow = true
  EOT
}

# Test 13: Policy with complex approval workflow
resource "controlplane_policy" "approval_workflow" {
  name        = "test-policy-approval-workflow"
  description = "Complex approval workflow policy"
  enabled     = true

  policy_content = <<-EOT
    package approval_workflow_test

    import future.keywords.if
    import future.keywords.in

    # Different approval requirements based on risk level
    deny[msg] {
      input.risk_level == "critical"
      count(input.approvals) < 3
      msg := "Critical changes require 3 approvals"
    }

    deny[msg] {
      input.risk_level == "high"
      count(input.approvals) < 2
      msg := "High risk changes require 2 approvals"
    }

    deny[msg] {
      input.risk_level == "medium"
      count(input.approvals) < 1
      msg := "Medium risk changes require 1 approval"
    }

    # Require security team approval for security-related changes
    deny[msg] {
      input.category == "security"
      not has_security_approval
      msg := "Security changes require approval from security team"
    }

    has_security_approval if {
      some approval in input.approvals
      approval.team == "security"
    }
  EOT

  tags = ["approval", "workflow", "compliance"]
}

# Test 14: Policy for update testing
resource "controlplane_policy" "for_update" {
  name        = "test-policy-for-update"
  description = "Initial description"
  enabled     = true

  policy_content = <<-EOT
    package update_test

    # Version 1
    default allow = true
  EOT

  tags = ["v1"]
}

# Data source tests
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

# Outputs for test validation
output "minimal_policy_id" {
  value = controlplane_policy.minimal.id
}

output "minimal_policy_name" {
  value = controlplane_policy.minimal.name
}

output "minimal_policy_version" {
  value = controlplane_policy.minimal.version
}

output "minimal_policy_created_at" {
  value = controlplane_policy.minimal.created_at
}

output "full_policy_id" {
  value = controlplane_policy.full.id
}

output "full_policy_name" {
  value = controlplane_policy.full.name
}

output "full_policy_description" {
  value = controlplane_policy.full.description
}

output "full_policy_policy_type" {
  value = controlplane_policy.full.policy_type
}

output "full_policy_enabled" {
  value = controlplane_policy.full.enabled
}

output "full_policy_policy_content" {
  value     = controlplane_policy.full.policy_content
  sensitive = true
}

output "full_policy_tags" {
  value = controlplane_policy.full.tags
}

output "full_policy_version" {
  value = controlplane_policy.full.version
}

output "full_policy_created_at" {
  value = controlplane_policy.full.created_at
}

output "full_policy_updated_at" {
  value = controlplane_policy.full.updated_at
}

output "allow_all_policy_id" {
  value = controlplane_policy.allow_all.id
}

output "deny_deletes_policy_id" {
  value = controlplane_policy.deny_deletes.id
}

output "rbac_policy_id" {
  value = controlplane_policy.rbac.id
}

output "rbac_policy_tags" {
  value = controlplane_policy.rbac.tags
}

output "resource_policy_id" {
  value = controlplane_policy.resource_policy.id
}

output "time_based_policy_id" {
  value = controlplane_policy.time_based.id
}

output "disabled_policy_id" {
  value = controlplane_policy.disabled.id
}

output "disabled_policy_enabled" {
  value = controlplane_policy.disabled.enabled
}

output "json_type_policy_id" {
  value = controlplane_policy.json_type.id
}

output "json_type_policy_policy_type" {
  value = controlplane_policy.json_type.policy_type
}

output "no_tags_policy_id" {
  value = controlplane_policy.no_tags.id
}

output "no_tags_policy_tags" {
  value = controlplane_policy.no_tags.tags
}

output "single_tag_policy_id" {
  value = controlplane_policy.single_tag.id
}

output "single_tag_policy_tags" {
  value = controlplane_policy.single_tag.tags
}

output "approval_workflow_policy_id" {
  value = controlplane_policy.approval_workflow.id
}

output "for_update_policy_id" {
  value = controlplane_policy.for_update.id
}

# Data source outputs for validation
output "data_minimal_name" {
  value = data.controlplane_policy.minimal_lookup.name
}

output "data_full_description" {
  value = data.controlplane_policy.full_lookup.description
}

output "data_full_enabled" {
  value = data.controlplane_policy.full_lookup.enabled
}

output "data_disabled_enabled" {
  value = data.controlplane_policy.disabled_lookup.enabled
}

output "data_rbac_tags" {
  value = data.controlplane_policy.rbac_lookup.tags
}
