terraform {
  required_providers {
    controlplane = {
      source = "kubiya/control-plane"
    }
  }
}

provider "controlplane" {}

# Create skills for team assignment testing
resource "controlplane_skill" "skill1" {
  name        = "test-skill-1-for-team-ds"
  type        = "shell"
  description = "First test skill for team assignment"
  enabled     = true
}

resource "controlplane_skill" "skill2" {
  name        = "test-skill-2-for-team-ds"
  type        = "python"
  description = "Second test skill for team assignment"
  enabled     = true
}

# Minimal team (required fields only)
resource "controlplane_team" "minimal" {
  name = "test-team-minimal"
}

# Full team with all optional fields
resource "controlplane_team" "full" {
  name        = "test-team-full-ds"
  description = "Comprehensive test team with all fields configured"
  status      = "active"
  runtime     = "default"

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

  skill_ids = [
    controlplane_skill.skill1.id,
    controlplane_skill.skill2.id
  ]

  execution_environment = jsonencode({
    env_vars = {
      TEAM_ENV_VAR_1 = "value1"
      TEAM_ENV_VAR_2 = "value2"
      DEBUG_MODE     = "true"
    }
    secrets         = ["team-secret-1", "team-secret-2"]
    integration_ids = []
  })
}

# Inactive team for status testing
resource "controlplane_team" "inactive" {
  name        = "test-team-inactive-ds"
  description = "Inactive team for status testing"
  status      = "inactive"
}

# Data sources
data "controlplane_team" "minimal_lookup" {
  id = controlplane_team.minimal.id
}

data "controlplane_team" "full_lookup" {
  id = controlplane_team.full.id
}

data "controlplane_team" "inactive_lookup" {
  id = controlplane_team.inactive.id
}

# Outputs for tests
output "data_minimal_name" {
  value = data.controlplane_team.minimal_lookup.name
}

output "data_full_description" {
  value = data.controlplane_team.full_lookup.description
}

output "data_full_skill_ids" {
  value = jsonencode(data.controlplane_team.full_lookup.skill_ids)
}

output "data_inactive_status" {
  value = data.controlplane_team.inactive_lookup.status
}
