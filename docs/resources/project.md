---
page_title: "kubiya_control_plane_project Resource"
subcategory: ""
description: |-
  Manages a Kubiya project
---

# kubiya_control_plane_project (Resource)

Manages a project in the Kubiya Control Plane. Projects are used to group and organize related resources.

## Example Usage

```terraform
resource "kubiya_control_plane_project" "platform" {
  name        = "platform-automation"
  description = "Platform automation and operations project"

  metadata = jsonencode({
    owner       = "platform-team"
    cost_center = "engineering"
    environment = "production"
  })
}
```

## Schema

### Required

- `name` (String) The name of the project

### Optional

- `description` (String) Description of the project
- `metadata` (String) Project metadata as JSON string

### Read-Only

- `id` (String) The unique identifier of the project
- `status` (String) Current status of the project
- `created_at` (String) Timestamp when the project was created
- `updated_at` (String) Timestamp when the project was last updated

## Import

Projects can be imported using their ID:

```shell
terraform import kubiya_control_plane_project.platform project-uuid-here
```
