---
page_title: "controlplane_project Data Source"
subcategory: ""
description: |-
  Retrieves information about a Kubiya project
---

# controlplane_project (Data Source)

Retrieves information about an existing project in the Kubiya Control Plane.

## Example Usage

```terraform
data "controlplane_project" "platform" {
  id = "project-uuid-here"
}

output "project_name" {
  value = data.controlplane_project.platform.name
}

output "project_metadata" {
  value = data.controlplane_project.platform.metadata
}
```

## Schema

### Required

- `id` (String) The unique identifier of the project to look up

### Read-Only

- `name` (String) The name of the project
- `description` (String) Description of the project
- `status` (String) Current status of the project
- `metadata` (String) Project metadata as JSON string
- `created_at` (String) Creation timestamp
- `updated_at` (String) Last update timestamp
