---
page_title: "controlplane_team Resource"
subcategory: ""
description: |-
  Manages a Kubiya team
---

# controlplane_team (Resource)

Manages a team in the Kubiya Control Plane. Teams are used to organize agents and share configuration across multiple agents.

## Example Usage

```terraform
resource "controlplane_team" "devops" {
  name        = "devops-team"
  description = "DevOps automation team"

  configuration = jsonencode({
    max_agents        = 10
    default_runtime   = "default"
    enable_monitoring = true
  })
}
```

## Schema

### Required

- `name` (String) The name of the team

### Optional

- `description` (String) Description of the team's purpose
- `configuration` (String) Team configuration as JSON string

### Read-Only

- `id` (String) The unique identifier of the team
- `status` (String) Current status of the team
- `created_at` (String) Timestamp when the team was created
- `updated_at` (String) Timestamp when the team was last updated

## Import

Teams can be imported using their ID:

```shell
terraform import controlplane_team.devops team-uuid-here
```
