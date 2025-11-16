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

  settings = jsonencode(merge(
    {
      region         = "us-east-1"
      max_workers    = 10
      auto_scaling   = true
      retention_days = 30
    },
    each.value.settings
  ))

  execution_environment = jsonencode(merge(
    {
      env_vars = {
        LOG_LEVEL = "info"
      }
    },
    each.value.execution_environment
  ))
}

# ============================================================================
# Projects
# ============================================================================

resource "controlplane_project" "this" {
  for_each = var.projects

  name        = each.key
  key         = each.value.key
  description = each.value.description

  settings = jsonencode(merge(
    {
      owner       = "engineering"
      environment = "production"
    },
    each.value.settings
  ))
}

# ============================================================================
# Teams
# ============================================================================

resource "controlplane_team" "this" {
  for_each = var.teams

  name        = each.key
  description = each.value.description
  runtime     = each.value.runtime

  configuration = jsonencode(merge(
    {
      max_agents = 10
    },
    each.value.configuration
  ))

  capabilities = each.value.capabilities
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

  configuration = jsonencode(each.value.configuration)
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

  llm_config = jsonencode(merge(
    {
      temperature = 0.7
      max_tokens  = 2000
    },
    each.value.llm_config
  ))

  capabilities = each.value.capabilities

  configuration = jsonencode(merge(
    {
      max_retries = 3
      timeout     = 300
    },
    each.value.configuration
  ))

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

  config = each.value.config != null ? jsonencode(each.value.config) : null

  depends_on = [
    controlplane_agent.this,
    controlplane_team.this,
    controlplane_environment.this
  ]
}
