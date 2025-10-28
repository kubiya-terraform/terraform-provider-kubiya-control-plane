terraform {
  required_providers {
    kubiya_control_plane = {
      source = "hashicorp.com/kubiya/kubiya"
    }
  }
}

provider "kubiya_control_plane" {
  # Configuration will be read from environment variables:
  # - KUBIYA_CONTROL_PLANE_API_KEY (required)
  # - KUBIYA_CONTROL_PLANE_ORG_ID (required)
  # - KUBIYA_CONTROL_PLANE_ENV (optional, defaults to "development")
  # - KUBIYA_CONTROL_PLANE_BASE_URL (optional, override base URL)
}

# Create a project for organizing agents and teams
resource "kubiya_control_plane_project" "ml_platform" {
  name        = "ML Platform"
  key         = "ML"
  description = "Machine Learning Platform Project"
  goals       = "Build and manage ML infrastructure"
  visibility  = "private"
}

# Create an environment for running agents
resource "kubiya_control_plane_environment" "production" {
  name         = "production"
  display_name = "Production Environment"
  description  = "Production environment for ML agents"
  tags         = ["production", "ml"]
}

resource "kubiya_control_plane_environment" "staging" {
  name         = "staging"
  display_name = "Staging Environment"
  description  = "Staging environment for testing"
  tags         = ["staging"]
}

# Create a team
resource "kubiya_control_plane_team" "ml_team" {
  name        = "ml-team"
  description = "Machine Learning Team"
  configuration = jsonencode({
    slack_channel = "#ml-team"
    timezone      = "UTC"
  })
}

# Create agents with different configurations
resource "kubiya_control_plane_agent" "data_analyst" {
  name        = "data-analyst"
  description = "AI agent for data analysis tasks"
  model_id    = "kubiya/claude-sonnet-4"
  runtime     = "claude_code"
  team_id     = kubiya_control_plane_team.ml_team.id

  llm_config = jsonencode({
    temperature = 0.7
    max_tokens  = 4096
  })

  configuration = jsonencode({
    role        = "data-analyst"
    permissions = ["read-data", "write-reports"]
  })
}

resource "kubiya_control_plane_agent" "model_trainer" {
  name        = "model-trainer"
  description = "AI agent for training ML models"
  model_id    = "kubiya/claude-sonnet-4"
  runtime     = "default"
  team_id     = kubiya_control_plane_team.ml_team.id

  llm_config = jsonencode({
    temperature = 0.5
    max_tokens  = 8192
  })
}

# Create another team for operations
resource "kubiya_control_plane_team" "ops_team" {
  name        = "ops-team"
  description = "Operations Team"
}

resource "kubiya_control_plane_agent" "ops_assistant" {
  name        = "ops-assistant"
  description = "AI agent for operations tasks"
  model_id    = "kubiya/claude-sonnet-4"
  runtime     = "claude_code"
  team_id     = kubiya_control_plane_team.ops_team.id
}

# Outputs
output "ml_project_id" {
  value       = kubiya_control_plane_project.ml_platform.id
  description = "ML Platform Project ID"
}

output "production_environment_id" {
  value       = kubiya_control_plane_environment.production.id
  description = "Production Environment ID"
}

output "ml_team_id" {
  value       = kubiya_control_plane_team.ml_team.id
  description = "ML Team ID"
}

output "agents" {
  value = {
    data_analyst   = kubiya_control_plane_agent.data_analyst.id
    model_trainer  = kubiya_control_plane_agent.model_trainer.id
    ops_assistant  = kubiya_control_plane_agent.ops_assistant.id
  }
  description = "Agent IDs"
}
