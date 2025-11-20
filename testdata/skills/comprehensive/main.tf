
provider "controlplane" {
  # Configuration via environment variables:
  # KUBIYA_CONTROL_PLANE_API_KEY
  # KUBIYA_CONTROL_PLANE_BASE_URL (optional, defaults to https://control-plane.kubiya.ai)
}

# Test 1: Minimal skill (required fields only)
resource "controlplane_skill" "minimal" {
  name = "test-skill-minimal"
  type = "shell"
}

# Test 2: Full skill with all optional fields
resource "controlplane_skill" "full" {
  name        = "test-skill-full"
  type        = "python"
  description = "Comprehensive test skill with all fields configured"
  icon        = "python-icon"
  enabled     = true

  configuration = jsonencode({
    python_version = "3.11"
    packages       = ["requests", "pandas", "numpy"]
    timeout        = 300
    max_memory     = "512MB"
    env_vars = {
      PYTHONPATH = "/opt/python"
      DEBUG      = "true"
    }
  })
}

# Test 3: Shell skill type
resource "controlplane_skill" "shell" {
  name        = "test-skill-shell"
  type        = "shell"
  description = "Shell command execution skill"
  enabled     = true

  configuration = jsonencode({
    allowed_commands = ["ls", "pwd", "echo", "cat", "grep"]
    timeout          = 30
    working_dir      = "/tmp"
    shell            = "/bin/bash"
  })
}

# Test 4: Docker skill type
resource "controlplane_skill" "docker" {
  name        = "test-skill-docker"
  type        = "docker"
  description = "Docker container execution skill"
  icon        = "docker"
  enabled     = true

  configuration = jsonencode({
    image          = "alpine:latest"
    command        = ["sh", "-c"]
    network_mode   = "bridge"
    memory_limit   = "1GB"
    cpu_limit      = "1.0"
    volumes        = ["/data:/data"]
    env_vars = {
      DOCKER_ENV = "test"
    }
  })
}

# Test 5: Python skill type
resource "controlplane_skill" "python" {
  name        = "test-skill-python"
  type        = "python"
  description = "Python script execution skill"
  enabled     = true

  configuration = jsonencode({
    python_version = "3.10"
    packages       = ["boto3", "pymongo"]
    entry_point    = "main.py"
  })
}

# Test 6: File system skill type
resource "controlplane_skill" "file_system" {
  name        = "test-skill-file-system"
  type        = "file_system"
  description = "File system operations skill"
  enabled     = true

  configuration = jsonencode({
    allowed_paths      = ["/data", "/tmp"]
    allowed_operations = ["read", "write", "list"]
    max_file_size      = "10MB"
  })
}

# Test 7: File generation skill type
resource "controlplane_skill" "file_generation" {
  name        = "test-skill-file-generation"
  type        = "file_generation"
  description = "File generation skill"
  enabled     = true

  configuration = jsonencode({
    output_directory = "/output"
    templates = ["template1", "template2"]
    formats   = ["json", "yaml", "xml"]
  })
}

# Test 8: Custom skill type
resource "controlplane_skill" "custom" {
  name        = "test-skill-custom"
  type        = "custom"
  description = "Custom skill implementation"
  icon        = "custom-icon"
  enabled     = true

  configuration = jsonencode({
    handler = "custom_handler"
    config = {
      custom_param1 = "value1"
      custom_param2 = "value2"
    }
  })
}

# Test 9: Disabled skill
resource "controlplane_skill" "disabled" {
  name        = "test-skill-disabled"
  type        = "shell"
  description = "Disabled skill for testing"
  enabled     = false
}

# Test 10: Skill with icon
resource "controlplane_skill" "with_icon" {
  name        = "test-skill-with-icon"
  type        = "python"
  description = "Skill with icon"
  icon        = "python-logo"
  enabled     = true
}

# Test 11: Skill with empty description
resource "controlplane_skill" "empty_description" {
  name        = "test-skill-empty-description"
  type        = "shell"
  description = ""
  enabled     = true
}

# Test 12: Skill with complex nested configuration
resource "controlplane_skill" "complex_config" {
  name        = "test-skill-complex-config"
  type        = "docker"
  description = "Skill with complex nested configuration"
  enabled     = true

  configuration = jsonencode({
    image = "custom-image:v1.0"
    resources = {
      limits = {
        cpu    = "2"
        memory = "2Gi"
      }
      requests = {
        cpu    = "500m"
        memory = "512Mi"
      }
    }
    security = {
      run_as_user  = 1000
      capabilities = ["NET_ADMIN", "SYS_TIME"]
    }
    health_check = {
      enabled  = true
      interval = 30
      timeout  = 5
      retries  = 3
    }
  })
}

# Test 13: Skill for update testing
resource "controlplane_skill" "for_update" {
  name        = "test-skill-for-update"
  type        = "shell"
  description = "Initial description"
  enabled     = true

  configuration = jsonencode({
    version = 1
  })
}

# Data source tests
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

# Outputs for test validation
output "minimal_skill_id" {
  value = controlplane_skill.minimal.id
}

output "minimal_skill_name" {
  value = controlplane_skill.minimal.name
}

output "minimal_skill_type" {
  value = controlplane_skill.minimal.type
}

output "minimal_skill_created_at" {
  value = controlplane_skill.minimal.created_at
}

output "full_skill_id" {
  value = controlplane_skill.full.id
}

output "full_skill_name" {
  value = controlplane_skill.full.name
}

output "full_skill_type" {
  value = controlplane_skill.full.type
}

output "full_skill_description" {
  value = controlplane_skill.full.description
}

output "full_skill_icon" {
  value = controlplane_skill.full.icon
}

output "full_skill_enabled" {
  value = controlplane_skill.full.enabled
}

output "full_skill_configuration" {
  value     = controlplane_skill.full.configuration
  sensitive = true
}

output "full_skill_created_at" {
  value = controlplane_skill.full.created_at
}

output "full_skill_updated_at" {
  value = controlplane_skill.full.updated_at
}

output "shell_skill_id" {
  value = controlplane_skill.shell.id
}

output "shell_skill_type" {
  value = controlplane_skill.shell.type
}

output "docker_skill_id" {
  value = controlplane_skill.docker.id
}

output "docker_skill_type" {
  value = controlplane_skill.docker.type
}

output "docker_skill_icon" {
  value = controlplane_skill.docker.icon
}

output "python_skill_id" {
  value = controlplane_skill.python.id
}

output "python_skill_type" {
  value = controlplane_skill.python.type
}

output "file_system_skill_id" {
  value = controlplane_skill.file_system.id
}

output "file_system_skill_type" {
  value = controlplane_skill.file_system.type
}

output "file_generation_skill_id" {
  value = controlplane_skill.file_generation.id
}

output "file_generation_skill_type" {
  value = controlplane_skill.file_generation.type
}

output "custom_skill_id" {
  value = controlplane_skill.custom.id
}

output "custom_skill_type" {
  value = controlplane_skill.custom.type
}

output "disabled_skill_id" {
  value = controlplane_skill.disabled.id
}

output "disabled_skill_enabled" {
  value = controlplane_skill.disabled.enabled
}

output "with_icon_skill_id" {
  value = controlplane_skill.with_icon.id
}

output "with_icon_skill_icon" {
  value = controlplane_skill.with_icon.icon
}

output "complex_config_skill_id" {
  value = controlplane_skill.complex_config.id
}

output "complex_config_skill_configuration" {
  value     = controlplane_skill.complex_config.configuration
  sensitive = true
}

output "for_update_skill_id" {
  value = controlplane_skill.for_update.id
}

# Data source outputs for validation
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
  value = data.controlplane_skill.full_lookup.enabled
}

output "data_shell_type" {
  value = data.controlplane_skill.shell_lookup.type
}

output "data_docker_type" {
  value = data.controlplane_skill.docker_lookup.type
}

output "data_disabled_enabled" {
  value = data.controlplane_skill.disabled_lookup.enabled
}
