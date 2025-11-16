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
  runtime     = "default"  # or "claude_code"

  configuration = jsonencode({
    max_agents        = 10
    enable_monitoring = true
  })

  skill_ids = ["skill-id-1", "skill-id-2"]

  execution_environment = jsonencode({
    ENV_VAR = "value"
  })
}
```

## Schema

### Required

- `name` (String) The name of the team

### Optional

- `description` (String) Description of the team's purpose
- `runtime` (String) Runtime type for team leader: 'default' (Agno) or 'claude_code' (Claude Code SDK). Defaults to 'default'.
- `configuration` (String) Team configuration as JSON string
- `skill_ids` (List of String) List of skill IDs associated with the team
- `execution_environment` (String) Execution environment configuration as JSON string

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
