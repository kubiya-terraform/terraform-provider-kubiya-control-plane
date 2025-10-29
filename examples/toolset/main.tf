terraform {
  required_providers {
    kubiya_control_plane = {
      source = "kubiya/control-plane"
    }
  }
}

provider "kubiya_control_plane" {
  # Configuration is via environment variables:
  # KUBIYA_CONTROL_PLANE_API_KEY
  # KUBIYA_CONTROL_PLANE_ORG_ID
  # KUBIYA_CONTROL_PLANE_BASE_URL (optional, defaults to https://control-plane.kubiya.ai)
}

# Create a file system toolset
resource "kubiya_control_plane_toolset" "filesystem" {
  name        = "example-filesystem-toolset"
  description = "File system operations toolset"
  type        = "file_system"
  enabled     = true

  configuration = jsonencode({
    allowed_paths = ["/tmp", "/var/app"]
    max_file_size = 10485760 # 10MB
    operations    = ["read", "write", "list"]
  })
}

# Create a shell toolset
resource "kubiya_control_plane_toolset" "shell" {
  name        = "example-shell-toolset"
  description = "Shell command execution toolset"
  type        = "shell"
  enabled     = true

  configuration = jsonencode({
    allowed_commands = ["ls", "cat", "grep", "find"]
    timeout          = 30
    working_dir      = "/app"
  })
}

# Create a Docker toolset
resource "kubiya_control_plane_toolset" "docker" {
  name        = "example-docker-toolset"
  description = "Docker operations toolset"
  type        = "docker"
  enabled     = true

  configuration = jsonencode({
    allowed_registries = ["docker.io", "gcr.io"]
    max_containers     = 10
    network_mode       = "bridge"
  })
}

# Look up an existing toolset by ID
data "kubiya_control_plane_toolset" "existing" {
  id = "toolset-uuid-here"
}

# Output toolset information
output "filesystem_toolset_id" {
  value       = kubiya_control_plane_toolset.filesystem.id
  description = "The ID of the filesystem toolset"
}

output "shell_toolset_id" {
  value       = kubiya_control_plane_toolset.shell.id
  description = "The ID of the shell toolset"
}

output "docker_toolset_id" {
  value       = kubiya_control_plane_toolset.docker.id
  description = "The ID of the docker toolset"
}

output "existing_toolset_name" {
  value       = data.kubiya_control_plane_toolset.existing.name
  description = "Name of the existing toolset"
}

output "existing_toolset_type" {
  value       = data.kubiya_control_plane_toolset.existing.type
  description = "Type of the existing toolset"
}

output "existing_toolset_enabled" {
  value       = data.kubiya_control_plane_toolset.existing.enabled
  description = "Whether the existing toolset is enabled"
}
