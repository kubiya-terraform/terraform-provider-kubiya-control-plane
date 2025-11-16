# ============================================================================
# Environment Outputs
# ============================================================================

output "environments" {
  description = "Map of created environments"
  value = {
    for k, v in controlplane_environment.this : k => {
      id     = v.id
      name   = v.name
      status = v.status
    }
  }
}

output "environment_ids" {
  description = "Map of environment names to IDs"
  value       = { for k, v in controlplane_environment.this : k => v.id }
}

# ============================================================================
# Project Outputs
# ============================================================================

output "projects" {
  description = "Map of created projects"
  value = {
    for k, v in controlplane_project.this : k => {
      id     = v.id
      name   = v.name
      key    = v.key
      status = v.status
    }
  }
}

output "project_ids" {
  description = "Map of project names to IDs"
  value       = { for k, v in controlplane_project.this : k => v.id }
}

# ============================================================================
# Team Outputs
# ============================================================================

output "teams" {
  description = "Map of created teams"
  value = {
    for k, v in controlplane_team.this : k => {
      id      = v.id
      name    = v.name
      runtime = v.runtime
      status  = v.status
    }
  }
}

output "team_ids" {
  description = "Map of team names to IDs"
  value       = { for k, v in controlplane_team.this : k => v.id }
}

# ============================================================================
# Skill Outputs
# ============================================================================

output "skills" {
  description = "Map of created skills"
  value = {
    for k, v in controlplane_skill.this : k => {
      id      = v.id
      name    = v.name
      type    = v.type
      enabled = v.enabled
    }
  }
}

output "skill_ids" {
  description = "Map of skill names to IDs"
  value       = { for k, v in controlplane_skill.this : k => v.id }
}

# ============================================================================
# Policy Outputs
# ============================================================================

output "policies" {
  description = "Map of created policies"
  value = {
    for k, v in controlplane_policy.this : k => {
      id      = v.id
      name    = v.name
      enabled = v.enabled
      tags    = v.tags
    }
  }
}

output "policy_ids" {
  description = "Map of policy names to IDs"
  value       = { for k, v in controlplane_policy.this : k => v.id }
}

# ============================================================================
# Agent Outputs
# ============================================================================

output "agents" {
  description = "Map of created agents"
  value = {
    for k, v in controlplane_agent.this : k => {
      id       = v.id
      name     = v.name
      model_id = v.model_id
      runtime  = v.runtime
      status   = v.status
      team_id  = v.team_id
    }
  }
}

output "agent_ids" {
  description = "Map of agent names to IDs"
  value       = { for k, v in controlplane_agent.this : k => v.id }
}

# ============================================================================
# Worker Queue Outputs
# ============================================================================

output "worker_queues" {
  description = "Map of created worker queues"
  value = {
    for k, v in controlplane_worker_queue.this : k => {
      id                = v.id
      name              = v.name
      display_name      = v.display_name
      task_queue_name   = v.task_queue_name
      environment_id    = v.environment_id
      active_workers    = v.active_workers
      max_workers       = v.max_workers
    }
  }
}

output "worker_queue_ids" {
  description = "Map of worker queue names to IDs"
  value       = { for k, v in controlplane_worker_queue.this : k => v.id }
}

output "worker_queue_task_names" {
  description = "Map of worker queue names to task queue names (for worker registration)"
  value       = { for k, v in controlplane_worker_queue.this : k => v.task_queue_name }
}

# ============================================================================
# Job Outputs
# ============================================================================

output "jobs" {
  description = "Map of created jobs"
  value = {
    for k, v in controlplane_job.this : k => {
      id           = v.id
      name         = v.name
      trigger_type = v.trigger_type
      enabled      = v.enabled
      status       = v.status
      webhook_url  = v.webhook_url
    }
  }
}

output "job_ids" {
  description = "Map of job names to IDs"
  value       = { for k, v in controlplane_job.this : k => v.id }
}

output "job_webhook_urls" {
  description = "Map of webhook-triggered job names to webhook URLs"
  value = {
    for k, v in controlplane_job.this : k => v.webhook_url
    if v.trigger_type == "webhook"
  }
}

# ============================================================================
# Summary Outputs
# ============================================================================

output "summary" {
  description = "Summary of created resources"
  value = {
    environments   = length(controlplane_environment.this)
    projects       = length(controlplane_project.this)
    teams          = length(controlplane_team.this)
    agents         = length(controlplane_agent.this)
    skills         = length(controlplane_skill.this)
    policies       = length(controlplane_policy.this)
    worker_queues  = length(controlplane_worker_queue.this)
    jobs           = length(controlplane_job.this)
  }
}
