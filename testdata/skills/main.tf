
provider "controlplane" {}

# Minimal skill (required fields only)
resource "controlplane_skill" "minimal" {
  name = "test-skill-minimal"
  type = "shell"
}

# Full skill with all optional fields
resource "controlplane_skill" "full" {
  name        = "test-skill-full-ds"
  type        = "python"
  description = "Comprehensive test skill with all fields configured"
  enabled     = true

  content = <<-EOT
    def main():
        print("Test skill execution")
        return {"status": "success"}
  EOT

  configuration = jsonencode({
    timeout     = 300
    memory_mb   = 512
    environment = "production"
  })
}

# Shell type skill
resource "controlplane_skill" "shell" {
  name        = "test-skill-shell-ds"
  type        = "shell"
  description = "Shell type skill for testing"
  enabled     = true

  content = <<-EOT
    #!/bin/bash
    echo "Shell skill execution"
  EOT
}

# Docker type skill
resource "controlplane_skill" "docker" {
  name        = "test-skill-docker-ds"
  type        = "docker"
  description = "Docker type skill for testing"
  enabled     = true

  content = <<-EOT
    FROM alpine:latest
    RUN echo "Docker skill"
  EOT
}

# Disabled skill for enabled flag testing
resource "controlplane_skill" "disabled" {
  name        = "test-skill-disabled-ds"
  type        = "shell"
  description = "Disabled skill for testing"
  enabled     = false
}

# Data sources
data "controlplane_skill" "minimal_lookup" {
  id = controlplane_skill.minimal.id
}

data "controlplane_skill" "full_lookup" {
  id = controlplane_skill.full.id
}

data "controlplane_skill" "shell_lookup" {
  id = controlplane_skill.shell.id
}

data "controlplane_skill" "docker_lookup" {
  id = controlplane_skill.docker.id
}

data "controlplane_skill" "disabled_lookup" {
  id = controlplane_skill.disabled.id
}

# Outputs for tests
output "data_minimal_name" {
  value = data.controlplane_skill.minimal_lookup.name
}

output "data_minimal_type" {
  value = data.controlplane_skill.minimal_lookup.type
}

output "data_full_description" {
  value = data.controlplane_skill.full_lookup.description
}

output "data_full_enabled" {
  value = tostring(data.controlplane_skill.full_lookup.enabled)
}

output "data_shell_type" {
  value = data.controlplane_skill.shell_lookup.type
}

output "data_docker_type" {
  value = data.controlplane_skill.docker_lookup.type
}

output "data_disabled_enabled" {
  value = tostring(data.controlplane_skill.disabled_lookup.enabled)
}
