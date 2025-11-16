# ============================================================================
# Engineering Organization Module
# ============================================================================
# This module creates a complete engineering organization setup including:
# - Environments
# - Projects
# - Teams
# - Agents
# - Skills
# - Policies
# - Worker Queues
# - Jobs
# ============================================================================

# ============================================================================
# Environments
# ============================================================================

resource "controlplane_environment" "this" {
  for_each = var.environments

  name        = each.key
  description = each.value.description

  settings = coalesce(
    each.value.settings,
    jsonencode({
      region         = "us-east-1"
      max_workers    = 10
      auto_scaling   = true
      retention_days = 30
    })
  )

  execution_environment = coalesce(
    each.value.execution_environment,
    jsonencode({
      env_vars = {
        LOG_LEVEL = "info"
      }
    })
  )
}

# ============================================================================
# Projects
# ============================================================================

resource "controlplane_project" "this" {
  for_each = var.projects

  name        = each.key
  key         = each.value.key
  description = each.value.description

  settings = coalesce(
    each.value.settings,
    jsonencode({
      owner       = "engineering"
      environment = "production"
    })
  )
}

# ============================================================================
# Teams
# ============================================================================

resource "controlplane_team" "this" {
  for_each = var.teams

  name        = each.key
  description = each.value.description
  runtime     = each.value.runtime

  configuration = coalesce(
    each.value.configuration,
    jsonencode({
      max_agents = 10
    })
  )
}

# ============================================================================
# Skills
# ============================================================================

resource "controlplane_skill" "this" {
  for_each = var.skills

  name        = each.key
  description = each.value.description
  type        = each.value.type
  enabled     = each.value.enabled

  configuration = coalesce(each.value.configuration, jsonencode({}))
}

# ============================================================================
# Policies
# ============================================================================

resource "controlplane_policy" "this" {
  for_each = var.policies

  name           = each.key
  description    = each.value.description
  enabled        = each.value.enabled
  policy_content = each.value.policy_content
  tags           = each.value.tags
}

# ============================================================================
# Agents
# ============================================================================

resource "controlplane_agent" "this" {
  for_each = var.agents

  name        = each.key
  description = each.value.description
  model_id    = each.value.model_id
  runtime     = each.value.runtime

  llm_config = coalesce(
    each.value.llm_config,
    jsonencode({
      temperature = 0.7
      max_tokens  = 2000
    })
  )

  capabilities = each.value.capabilities

  configuration = coalesce(
    each.value.configuration,
    jsonencode({
      max_retries = 3
      timeout     = 300
    })
  )

  team_id = each.value.team_name != null ? controlplane_team.this[each.value.team_name].id : null

  depends_on = [controlplane_team.this]
}

# ============================================================================
# Worker Queues
# ============================================================================

resource "controlplane_worker_queue" "this" {
  for_each = var.worker_queues

  environment_id     = controlplane_environment.this[each.value.environment_name].id
  name               = each.key
  display_name       = each.value.display_name
  description        = each.value.description
  heartbeat_interval = each.value.heartbeat_interval
  max_workers        = each.value.max_workers
  tags               = each.value.tags

  settings = each.value.settings

  depends_on = [controlplane_environment.this]
}

# ============================================================================
# Jobs
# ============================================================================

resource "controlplane_job" "this" {
  for_each = var.jobs

  name         = each.key
  description  = each.value.description
  enabled      = each.value.enabled
  trigger_type = each.value.trigger_type

  # Cron-specific fields
  cron_schedule = each.value.trigger_type == "cron" ? each.value.cron_schedule : null
  cron_timezone = each.value.trigger_type == "cron" ? each.value.cron_timezone : null

  # Planning configuration
  planning_mode   = each.value.planning_mode
  entity_type     = each.value.entity_type
  entity_id       = each.value.entity_type == "agent" && each.value.entity_name != null ? controlplane_agent.this[each.value.entity_name].id : (each.value.entity_type == "team" && each.value.entity_name != null ? controlplane_team.this[each.value.entity_name].id : null)
  prompt_template = each.value.prompt_template
  system_prompt   = each.value.system_prompt

  # Executor configuration
  executor_type    = each.value.executor_type
  environment_name = each.value.environment_name != null ? controlplane_environment.this[each.value.environment_name].name : null

  # Execution configuration
  execution_env_vars = each.value.execution_env_vars
  execution_secrets  = each.value.execution_secrets

  config = each.value.config

  depends_on = [
    controlplane_agent.this,
    controlplane_team.this,
    controlplane_environment.this
  ]
}
