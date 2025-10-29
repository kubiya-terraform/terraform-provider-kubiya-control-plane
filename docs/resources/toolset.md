---
page_title: "kubiya_control_plane_toolset Resource"
subcategory: ""
description: |-
  Manages a Kubiya toolset
---

# kubiya_control_plane_toolset (Resource)

Manages a toolset in the Kubiya Control Plane. Toolsets provide specific capabilities to agents such as filesystem access, shell commands, or Docker operations.

## Example Usage

```terraform
# File system toolset
resource "kubiya_control_plane_toolset" "filesystem" {
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

# Shell toolset
resource "kubiya_control_plane_toolset" "shell" {
  name        = "shell-commands"
  description = "Shell command execution"
  type        = "shell"
  enabled     = true

  configuration = jsonencode({
    allowed_commands = ["kubectl", "helm", "aws"]
    timeout          = 300
  })
}

# Docker toolset
resource "kubiya_control_plane_toolset" "docker" {
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

- `name` (String) The name of the toolset
- `type` (String) Type of toolset. Valid values: `file_system`, `shell`, `docker`, `python`, `file_generation`, `custom`

### Optional

- `description` (String) Description of the toolset
- `enabled` (Boolean) Whether the toolset is enabled
- `configuration` (String) Toolset configuration as JSON string (type-specific settings)

### Read-Only

- `id` (String) The unique identifier of the toolset
- `created_at` (String) Timestamp when the toolset was created
- `updated_at` (String) Timestamp when the toolset was last updated

## Import

Toolsets can be imported using their ID:

```shell
terraform import kubiya_control_plane_toolset.filesystem toolset-uuid-here
```
