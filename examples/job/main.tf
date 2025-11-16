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
  # KUBIYA_CONTROL_PLANE_ORG_ID
  # KUBIYA_CONTROL_PLANE_BASE_URL (optional, defaults to https://control-plane.kubiya.ai)
}

# Create an agent to use in jobs
resource "controlplane_agent" "example" {
  name        = "example-job-agent"
  description = "Agent for job examples"
  model_id    = "kubiya/claude-sonnet-4"
  runtime     = "default"

  llm_config = jsonencode({
    temperature = 0.7
    max_tokens  = 4096
  })

  configuration = jsonencode({
    environment_variables = {
      LOG_LEVEL = "info"
    }
  })
}

# Create a cron-triggered job
resource "controlplane_job" "daily_report" {
  name          = "daily-report-job"
  description   = "Generate daily report at 5pm Eastern Time"
  enabled       = true
  trigger_type  = "cron"
  cron_schedule = "0 17 * * *" # 5 PM daily
  cron_timezone = "America/New_York"

  planning_mode   = "predefined_agent"
  entity_type     = "agent"
  entity_id       = controlplane_agent.example.id
  prompt_template = "Generate a daily report for {{date}}"

  executor_type = "auto"

  execution_env_vars = {
    REPORT_FORMAT = "pdf"
    OUTPUT_PATH   = "/reports"
  }

  execution_secrets = ["slack_webhook_token"]
}

# Create a webhook-triggered job
resource "controlplane_job" "webhook_handler" {
  name         = "webhook-handler-job"
  description  = "Handle incoming webhook events"
  enabled      = true
  trigger_type = "webhook"

  planning_mode   = "predefined_agent"
  entity_type     = "agent"
  entity_id       = controlplane_agent.example.id
  prompt_template = "Process webhook event: {{event_type}} - {{payload}}"
  system_prompt   = "You are an event processing assistant. Parse webhook payloads and take appropriate actions."

  executor_type = "auto"

  config = jsonencode({
    timeout = 300
    retry_policy = {
      max_attempts = 3
      backoff      = "exponential"
    }
  })
}

# Create a manual job
resource "controlplane_job" "manual_task" {
  name         = "manual-task-job"
  description  = "Manually triggered task"
  enabled      = true
  trigger_type = "manual"

  planning_mode   = "on_the_fly"
  prompt_template = "Execute task: {{task_description}}"
  system_prompt   = "You are a general-purpose automation assistant."

  executor_type = "auto"
}

# Fetch a specific job
data "controlplane_job" "daily_report" {
  id = controlplane_job.daily_report.id
}

# Fetch all jobs
data "controlplane_jobs" "all" {}

# Outputs
output "daily_report_job_id" {
  value       = controlplane_job.daily_report.id
  description = "The ID of the daily report job"
}

output "daily_report_status" {
  value       = controlplane_job.daily_report.status
  description = "The status of the daily report job"
}

output "webhook_job_id" {
  value       = controlplane_job.webhook_handler.id
  description = "The ID of the webhook handler job"
}

output "webhook_url" {
  value       = controlplane_job.webhook_handler.webhook_url
  description = "The webhook URL for triggering the job"
  sensitive   = false
}

output "manual_job_id" {
  value       = controlplane_job.manual_task.id
  description = "The ID of the manual task job"
}

output "total_jobs_count" {
  value       = length(data.controlplane_jobs.all.jobs)
  description = "Total number of jobs in the organization"
}

output "enabled_jobs" {
  value       = [for j in data.controlplane_jobs.all.jobs : j.name if j.enabled]
  description = "List of enabled job names"
}

output "cron_jobs" {
  value       = [for j in data.controlplane_jobs.all.jobs : j.name if j.trigger_type == "cron"]
  description = "List of cron-triggered job names"
}
