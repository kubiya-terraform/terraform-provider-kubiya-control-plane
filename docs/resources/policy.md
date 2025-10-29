---
page_title: "kubiya_control_plane_policy Resource"
subcategory: ""
description: |-
  Manages a Kubiya OPA policy
---

# kubiya_control_plane_policy (Resource)

Manages an OPA (Open Policy Agent) Rego policy in the Kubiya Control Plane. Policies define governance rules for agent operations, security, compliance, and cost control.

## Example Usage

```terraform
# Security policy
resource "kubiya_control_plane_policy" "security" {
  name        = "production-security-policy"
  description = "Security policy for production environments"
  enabled     = true

  policy_content = <<-EOT
    package kubiya.security

    # Deny operations during business hours
    deny[msg] {
      input.operation = "deploy"
      input.environment = "production"
      is_business_hours
      msg := "Production deployments not allowed during business hours"
    }

    # Require approvals for delete operations
    deny[msg] {
      input.operation = "delete"
      count(input.approvals) < 2
      msg := "Delete operations require at least 2 approvals"
    }

    is_business_hours {
      hour := time.clock(time.now_ns())[0]
      hour >= 9
      hour < 17
    }
  EOT

  tags = ["security", "production", "compliance"]
}

# Cost control policy
resource "kubiya_control_plane_policy" "cost_control" {
  name        = "cost-control-policy"
  description = "Policy to control infrastructure costs"
  enabled     = true

  policy_content = <<-EOT
    package kubiya.cost

    # Limit instance sizes
    deny[msg] {
      input.action = "create_instance"
      input.instance_type = "x2.32xlarge"
      msg := "Instance type too large"
    }

    # Require cost tags
    deny[msg] {
      input.action = "create_resource"
      not input.tags.cost_center
      msg := "All resources must have a cost_center tag"
    }
  EOT

  tags = ["cost", "governance"]
}
```

## Schema

### Required

- `name` (String) The name of the policy
- `policy_content` (String) OPA Rego policy content
- `enabled` (Boolean) Whether the policy is enabled

### Optional

- `description` (String) Description of the policy
- `tags` (List of String) Tags for categorizing the policy

### Read-Only

- `id` (String) The unique identifier of the policy
- `created_at` (String) Timestamp when the policy was created
- `updated_at` (String) Timestamp when the policy was last updated

## Import

Policies can be imported using their ID:

```shell
terraform import kubiya_control_plane_policy.security policy-uuid-here
```
