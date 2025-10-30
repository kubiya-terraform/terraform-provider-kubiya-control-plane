---
page_title: "controlplane_team Data Source"
subcategory: ""
description: |-
  Retrieves information about a Kubiya team
---

# controlplane_team (Data Source)

Retrieves information about an existing team in the Kubiya Control Plane.

## Example Usage

```terraform
data "controlplane_team" "devops" {
  id = "team-uuid-here"
}

output "team_name" {
  value = data.controlplane_team.devops.name
}

output "team_config" {
  value     = data.controlplane_team.devops.configuration
  sensitive = true
}
```

## Schema

### Required

- `id` (String) The unique identifier of the team to look up

### Read-Only

- `name` (String) The name of the team
- `description` (String) Description of the team
- `status` (String) Current status of the team
- `capabilities` (List of String) List of team capabilities
- `configuration` (String) Team configuration as JSON string
- `created_at` (String) Creation timestamp
- `updated_at` (String) Last update timestamp
