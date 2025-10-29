---
page_title: "kubiya_control_plane_team Resource"
subcategory: ""
description: |-
  Manages a Kubiya team
---

# kubiya_control_plane_team (Resource)

Manages a team in the Kubiya Control Plane. Teams are used to organize agents and share configuration across multiple agents.

## Example Usage

```terraform
resource "kubiya_control_plane_team" "devops" {
  name        = "devops-team"
  description = "DevOps automation team"

  configuration = jsonencode({
    max_agents        = 10
    default_runtime   = "default"
    enable_monitoring = true
  })

  capabilities = ["deployment", "monitoring", "reporting"]
}
```

## Schema

### Required

- `name` (String) The name of the team

### Optional

- `description` (String) Description of the team's purpose
- `capabilities` (List of String) List of team capabilities
- `configuration` (String) Team configuration as JSON string

### Read-Only

- `id` (String) The unique identifier of the team
- `status` (String) Current status of the team
- `created_at` (String) Timestamp when the team was created
- `updated_at` (String) Timestamp when the team was last updated

## Import

Teams can be imported using their ID:

```shell
terraform import kubiya_control_plane_team.devops team-uuid-here
```
