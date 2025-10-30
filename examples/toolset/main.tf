terraform {
  required_providers {
    controlplane = {
      source = "kubiya/control-plane"
    }
  }
}

provider "controlplane" {
  # Configuration is via environment variables:
  # KUBIYA_CONTROL_PLANE_API_KEY
  # KUBIYA_CONTROL_PLANE_ORG_ID
  # KUBIYA_CONTROL_PLANE_BASE_URL (optional, defaults to https://control-plane.kubiya.ai)
}

# Create a file system toolset
resource "controlplane_toolset" "filesystem" {
  name        = "production-filesystem-toolset"
  description = "Production file system operations toolset"
  type        = "file_system"
  enabled     = true

  configuration = jsonencode({
    allowed_paths = ["/tmp", "/var/app"]
    max_file_size = 10485760 # 10MB
    operations    = ["read", "write", "list"]
  })
}

# Create a shell toolset
resource "controlplane_toolset" "shell" {
  name        = "production-shell-toolset"
  description = "Production shell command execution toolset"
  type        = "shell"
  enabled     = true

  configuration = jsonencode({
    allowed_commands = ["ls", "cat", "grep", "find"]
    timeout          = 30
    working_dir      = "/app"
  })
}

# Create a Docker toolset
resource "controlplane_toolset" "docker" {
  name        = "production-docker-toolset"
  description = "Production docker operations toolset"
  type        = "docker"
  enabled     = true

  configuration = jsonencode({
    allowed_registries = ["docker.io", "gcr.io"]
    max_containers     = 10
    network_mode       = "bridge"
  })
}

# Look up an existing toolset by ID
data "controlplane_toolset" "existing" {
  id = "8455d897-8e2c-4528-bcb6-26493df3d35d"
}

# Output toolset information
output "filesystem_toolset_id" {
  value       = controlplane_toolset.filesystem.id
  description = "The ID of the filesystem toolset"
}

output "shell_toolset_id" {
  value       = controlplane_toolset.shell.id
  description = "The ID of the shell toolset"
}

output "docker_toolset_id" {
  value       = controlplane_toolset.docker.id
  description = "The ID of the docker toolset"
}

output "existing_toolset_name" {
  value       = data.controlplane_toolset.existing.name
  description = "Name of the existing toolset"
}

output "existing_toolset_type" {
  value       = data.controlplane_toolset.existing.type
  description = "Type of the existing toolset"
}

output "existing_toolset_enabled" {
  value       = data.controlplane_toolset.existing.enabled
  description = "Whether the existing toolset is enabled"
}
