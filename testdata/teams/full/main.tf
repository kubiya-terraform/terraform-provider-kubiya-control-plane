
provider "controlplane" {}

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

# Full team with all optional fields
resource "controlplane_team" "full" {
  name        = "test-team-full"
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
    secrets = ["team-secret-1", "team-secret-2"]
    integration_ids = []
  })
}

# Data source test
data "controlplane_team" "full_lookup" {
  id = controlplane_team.full.id
}

# Outputs
output "team_id" {
  value = controlplane_team.full.id
}

output "team_name" {
  value = controlplane_team.full.name
}

output "team_description" {
  value = controlplane_team.full.description
}

output "team_status" {
  value = controlplane_team.full.status
}

output "team_runtime" {
  value = controlplane_team.full.runtime
}

output "team_skill_ids" {
  value = controlplane_team.full.skill_ids
}

output "data_team_description" {
  value = data.controlplane_team.full_lookup.description
}

output "data_team_skill_ids" {
  value = data.controlplane_team.full_lookup.skill_ids
}
