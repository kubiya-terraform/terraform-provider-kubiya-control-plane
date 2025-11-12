---
page_title: "controlplane_skill Data Source"
subcategory: ""
description: |-
  Retrieves information about a Kubiya skill
---

# controlplane_skill (Data Source)

Retrieves information about an existing skill in the Kubiya Control Plane.

## Example Usage

```terraform
data "controlplane_skill" "shell" {
  id = "skill-uuid-here"
}

output "skill_name" {
  value = data.controlplane_skill.shell.name
}

output "skill_type" {
  value = data.controlplane_skill.shell.type
}

output "skill_enabled" {
  value = data.controlplane_skill.shell.enabled
}
```

## Schema

### Required

- `id` (String) The unique identifier of the skill to look up

### Read-Only

- `name` (String) The name of the skill
- `description` (String) Description of the skill
- `type` (String) Type of skill (file_system, shell, docker, etc.)
- `configuration` (String) Skill configuration as JSON string
- `enabled` (Boolean) Whether the skill is enabled
- `created_at` (String) Creation timestamp
- `updated_at` (String) Last update timestamp
