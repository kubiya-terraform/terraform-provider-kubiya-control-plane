---
page_title: "controlplane_skill Resource"
subcategory: ""
description: |-
  Manages a Kubiya skill
---

# controlplane_skill (Resource)

Manages a skill in the Kubiya Control Plane. Skills provide specific capabilities to agents such as filesystem access, shell commands, or Docker operations.

## Example Usage

```terraform
# File system skill
resource "controlplane_skill" "filesystem" {
  name        = "filesystem-ops"
  description = "File system operations"
  type        = "file_system"
  enabled     = true

  configuration = jsonencode({
    allowed_paths = ["/app", "/tmp"]
    max_file_size = 10485760
    operations    = ["read", "write", "list"]
  })
}

# Shell skill
resource "controlplane_skill" "shell" {
  name        = "shell-commands"
  description = "Shell command execution"
  type        = "shell"
  enabled     = true

  configuration = jsonencode({
    allowed_commands = ["kubectl", "helm", "aws"]
    timeout          = 300
  })
}

# Docker skill
resource "controlplane_skill" "docker" {
  name        = "docker-ops"
  description = "Docker operations"
  type        = "docker"
  enabled     = true

  configuration = jsonencode({
    allowed_registries = ["docker.io", "gcr.io"]
    max_containers     = 10
  })
}
```

## Schema

### Required

- `name` (String) The name of the skill
- `type` (String) Type of skill. Valid values: `file_system`, `shell`, `docker`, `python`, `file_generation`, `custom`

### Optional

- `description` (String) Description of the skill
- `enabled` (Boolean) Whether the skill is enabled
- `configuration` (String) Skill configuration as JSON string (type-specific settings)

### Read-Only

- `id` (String) The unique identifier of the skill
- `created_at` (String) Timestamp when the skill was created
- `updated_at` (String) Timestamp when the skill was last updated

## Import

Skills can be imported using their ID:

```shell
terraform import controlplane_skill.filesystem skill-uuid-here
```
