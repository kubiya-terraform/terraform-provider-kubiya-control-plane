---
page_title: "controlplane_toolset Data Source"
subcategory: ""
description: |-
  Retrieves information about a Kubiya toolset
---

# controlplane_toolset (Data Source)

Retrieves information about an existing toolset in the Kubiya Control Plane.

## Example Usage

```terraform
data "controlplane_toolset" "shell" {
  id = "toolset-uuid-here"
}

output "toolset_name" {
  value = data.controlplane_toolset.shell.name
}

output "toolset_type" {
  value = data.controlplane_toolset.shell.type
}

output "toolset_enabled" {
  value = data.controlplane_toolset.shell.enabled
}
```

## Schema

### Required

- `id` (String) The unique identifier of the toolset to look up

### Read-Only

- `name` (String) The name of the toolset
- `description` (String) Description of the toolset
- `type` (String) Type of toolset (file_system, shell, docker, etc.)
- `configuration` (String) Toolset configuration as JSON string
- `enabled` (Boolean) Whether the toolset is enabled
- `created_at` (String) Creation timestamp
- `updated_at` (String) Last update timestamp
