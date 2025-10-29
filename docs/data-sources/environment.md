---
page_title: "kubiya_control_plane_environment Data Source"
subcategory: ""
description: |-
  Retrieves information about a Kubiya environment
---

# kubiya_control_plane_environment (Data Source)

Retrieves information about an existing execution environment in the Kubiya Control Plane.

## Example Usage

```terraform
data "kubiya_control_plane_environment" "production" {
  id = "environment-uuid-here"
}

output "environment_name" {
  value = data.kubiya_control_plane_environment.production.name
}

output "environment_config" {
  value     = data.kubiya_control_plane_environment.production.configuration
  sensitive = true
}
```

## Schema

### Required

- `id` (String) The unique identifier of the environment to look up

### Read-Only

- `name` (String) The name of the environment
- `description` (String) Description of the environment
- `status` (String) Current status of the environment
- `configuration` (String) Environment configuration as JSON string
- `execution_environment` (String) Execution environment settings as JSON string
- `created_at` (String) Creation timestamp
- `updated_at` (String) Last update timestamp
