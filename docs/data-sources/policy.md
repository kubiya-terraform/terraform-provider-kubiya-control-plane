---
page_title: "controlplane_policy Data Source"
subcategory: ""
description: |-
  Retrieves information about a Kubiya policy
---

# controlplane_policy (Data Source)

Retrieves information about an existing OPA policy in the Kubiya Control Plane.

## Example Usage

```terraform
data "controlplane_policy" "security" {
  id = "policy-uuid-here"
}

output "policy_name" {
  value = data.controlplane_policy.security.name
}

output "policy_enabled" {
  value = data.controlplane_policy.security.enabled
}

# Use policy content in another resource
output "policy_content" {
  value     = data.controlplane_policy.security.policy
  sensitive = true
}
```

## Schema

### Required

- `id` (String) The unique identifier of the policy to look up

### Read-Only

- `name` (String) The name of the policy
- `description` (String) Description of the policy
- `policy` (String) OPA Rego policy content
- `enabled` (Boolean) Whether the policy is enabled
- `created_at` (String) Creation timestamp
- `updated_at` (String) Last update timestamp
