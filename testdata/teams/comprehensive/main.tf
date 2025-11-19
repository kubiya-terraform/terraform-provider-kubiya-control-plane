terraform {
  required_providers {
    controlplane = {
      source = "kubiya/control-plane"
    }
  }
}

provider "controlplane" {
  # Configuration via environment variables:
  # KUBIYA_CONTROL_PLANE_API_KEY
  # KUBIYA_CONTROL_PLANE_BASE_URL (optional, defaults to https://control-plane.kubiya.ai)
}

# Create skills for team assignment testing
resource "controlplane_skill" "skill1" {
  name        = "test-skill-1-for-team"
  type        = "shell"
  description = "First test skill for team assignment"
  enabled     = true
}

resource "controlplane_skill" "skill2" {
  name        = "test-skill-2-for-team"
  type        = "python"
  description = "Second test skill for team assignment"
  enabled     = true
}

# Test 1: Minimal team (required fields only)
resource "controlplane_team" "minimal" {
  name = "test-team-minimal"
}

# Test 2: Full team with all optional fields
resource "controlplane_team" "full" {
  name        = "test-team-full"
  description = "Comprehensive test team with all fields configured"
  status      = "active"
  runtime     = "default"

  # Team configuration - testing complex JSON
  configuration = jsonencode({
    max_agents        = 10
    default_runtime   = "default"
    enable_monitoring = true
    policies = {
      max_execution_time = 3600
      require_approval   = false
    }
    settings = {
      verbose_logging = true
      auto_retry      = true
    }
  })

  # Skills assignment - testing list field
  skill_ids = [
    controlplane_skill.skill1.id,
    controlplane_skill.skill2.id
  ]

  # Execution environment - testing complex JSON
  execution_environment = jsonencode({
    env_vars = {
      TEAM_ENV_VAR_1 = "value1"
      TEAM_ENV_VAR_2 = "value2"
      DEBUG_MODE     = "true"
    }
    secrets = ["team-secret-1", "team-secret-2"]
    integration_ids = []
  })
}

# Test 3: Team with claude_code runtime
resource "controlplane_team" "claude_code" {
  name        = "test-team-claude-code"
  description = "Team using Claude Code SDK runtime"
  runtime     = "claude_code"
  status      = "active"

  configuration = jsonencode({
    sdk_version = "latest"
  })
}

# Test 4: Team with inactive status
resource "controlplane_team" "inactive" {
  name        = "test-team-inactive"
  description = "Team with inactive status"
  status      = "inactive"
}

# Test 5: Team with archived status
resource "controlplane_team" "archived" {
  name        = "test-team-archived"
  description = "Team with archived status"
  status      = "archived"
}

# Test 6: Team with single skill
resource "controlplane_skill" "single_skill" {
  name    = "test-skill-single"
  type    = "docker"
  enabled = true
}

resource "controlplane_team" "single_skill" {
  name        = "test-team-single-skill"
  description = "Team with single skill"
  skill_ids   = [controlplane_skill.single_skill.id]
}

# Test 7: Team with execution environment only
resource "controlplane_team" "exec_env_only" {
  name        = "test-team-exec-env-only"
  description = "Team with execution environment configuration"

  execution_environment = jsonencode({
    env_vars = {
      APP_ENV  = "test"
      LOG_LEVEL = "debug"
    }
    secrets = ["api-key"]
    integration_ids = ["integration-1", "integration-2"]
  })
}

# Test 8: Team with empty optional fields
resource "controlplane_team" "empty_optionals" {
  name        = "test-team-empty-optionals"
  description = ""
  skill_ids   = []
}

# Test 9: Team for update testing
resource "controlplane_team" "for_update" {
  name        = "test-team-for-update"
  description = "Initial description"
  runtime     = "default"
  status      = "active"

  configuration = jsonencode({
    version = 1
  })
}

# Data source tests
data "controlplane_team" "minimal_lookup" {
  id = controlplane_team.minimal.id
}

data "controlplane_team" "full_lookup" {
  id = controlplane_team.full.id
}

data "controlplane_team" "claude_code_lookup" {
  id = controlplane_team.claude_code.id
}

data "controlplane_team" "inactive_lookup" {
  id = controlplane_team.inactive.id
}

# Outputs for test validation
output "minimal_team_id" {
  value = controlplane_team.minimal.id
}

output "minimal_team_name" {
  value = controlplane_team.minimal.name
}

output "minimal_team_created_at" {
  value = controlplane_team.minimal.created_at
}

output "full_team_id" {
  value = controlplane_team.full.id
}

output "full_team_name" {
  value = controlplane_team.full.name
}

output "full_team_description" {
  value = controlplane_team.full.description
}

output "full_team_status" {
  value = controlplane_team.full.status
}

output "full_team_runtime" {
  value = controlplane_team.full.runtime
}

output "full_team_configuration" {
  value     = controlplane_team.full.configuration
  sensitive = true
}

output "full_team_skill_ids" {
  value = controlplane_team.full.skill_ids
}

output "full_team_execution_environment" {
  value     = controlplane_team.full.execution_environment
  sensitive = true
}

output "full_team_created_at" {
  value = controlplane_team.full.created_at
}

output "full_team_updated_at" {
  value = controlplane_team.full.updated_at
}

output "claude_code_team_id" {
  value = controlplane_team.claude_code.id
}

output "claude_code_team_runtime" {
  value = controlplane_team.claude_code.runtime
}

output "inactive_team_id" {
  value = controlplane_team.inactive.id
}

output "inactive_team_status" {
  value = controlplane_team.inactive.status
}

output "archived_team_id" {
  value = controlplane_team.archived.id
}

output "archived_team_status" {
  value = controlplane_team.archived.status
}

output "single_skill_team_id" {
  value = controlplane_team.single_skill.id
}

output "single_skill_team_skill_ids" {
  value = controlplane_team.single_skill.skill_ids
}

output "exec_env_only_team_id" {
  value = controlplane_team.exec_env_only.id
}

output "exec_env_only_team_execution_environment" {
  value     = controlplane_team.exec_env_only.execution_environment
  sensitive = true
}

output "for_update_team_id" {
  value = controlplane_team.for_update.id
}

# Data source outputs for validation
output "data_minimal_name" {
  value = data.controlplane_team.minimal_lookup.name
}

output "data_full_description" {
  value = data.controlplane_team.full_lookup.description
}

output "data_full_skill_ids" {
  value = data.controlplane_team.full_lookup.skill_ids
}

output "data_claude_code_runtime" {
  value = data.controlplane_team.claude_code_lookup.runtime
}

output "data_inactive_status" {
  value = data.controlplane_team.inactive_lookup.status
}
