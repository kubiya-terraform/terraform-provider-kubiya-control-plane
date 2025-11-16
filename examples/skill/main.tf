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
  # KUBIYA_CONTROL_PLANE_BASE_URL (optional, defaults to https://control-plane.kubiya.ai)
}

# Create a file system skill
resource "controlplane_skill" "filesystem" {
  name        = "production-filesystem-skill"
  description = "Production file system operations skill"
  type        = "file_system"
  enabled     = true

  configuration = jsonencode({
    allowed_paths = ["/tmp", "/var/app"]
    max_file_size = 10485760 # 10MB
    operations    = ["read", "write", "list"]
  })
}

# Create a shell skill
resource "controlplane_skill" "shell" {
  name        = "production-shell-skill"
  description = "Production shell command execution skill"
  type        = "shell"
  enabled     = true

  configuration = jsonencode({
    allowed_commands = ["ls", "cat", "grep", "find"]
    timeout          = 30
    working_dir      = "/app"
  })
}

# Create a Docker skill
resource "controlplane_skill" "docker" {
  name        = "production-docker-skill"
  description = "Production docker operations skill"
  type        = "docker"
  enabled     = true

  configuration = jsonencode({
    allowed_registries = ["docker.io", "gcr.io"]
    max_containers     = 10
    network_mode       = "bridge"
  })
}

# Look up an existing skill by ID
data "controlplane_skill" "existing" {
  id = "8455d897-8e2c-4528-bcb6-26493df3d35d"
}

# Output skill information
output "filesystem_skill_id" {
  value       = controlplane_skill.filesystem.id
  description = "The ID of the filesystem skill"
}

output "shell_skill_id" {
  value       = controlplane_skill.shell.id
  description = "The ID of the shell skill"
}

output "docker_skill_id" {
  value       = controlplane_skill.docker.id
  description = "The ID of the docker skill"
}

output "existing_skill_name" {
  value       = data.controlplane_skill.existing.name
  description = "Name of the existing skill"
}

output "existing_skill_type" {
  value       = data.controlplane_skill.existing.type
  description = "Type of the existing skill"
}

output "existing_skill_enabled" {
  value       = data.controlplane_skill.existing.enabled
  description = "Whether the existing skill is enabled"
}
