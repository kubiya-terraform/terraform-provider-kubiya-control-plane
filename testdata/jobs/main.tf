terraform {
  required_providers {
    controlplane = {
      source = "kubiya/control-plane"
    }
  }
}

provider "controlplane" {}

# Create agent for jobs
resource "controlplane_agent" "test" {
  name = "test-agent-for-job-ds"
}

# Minimal cron job (required fields only)
resource "controlplane_job" "minimal_cron" {
  name          = "test-job-minimal-cron"
  trigger_type  = "cron"
  cron_schedule = "0 9 * * *"
  prompt        = "Generate daily report"
  agent_id      = controlplane_agent.test.id
}

# Full cron job with all optional fields
resource "controlplane_job" "full_cron" {
  name          = "test-job-full-cron-ds"
  description   = "Comprehensive test cron job with all fields configured"
  trigger_type  = "cron"
  cron_schedule = "0 17 * * 1-5"
  cron_timezone = "America/New_York"
  enabled       = true

  prompt        = "Generate comprehensive weekly report with all metrics"
  agent_id      = controlplane_agent.test.id
  planning_mode = "predefined_agent"

  executor_type = "auto"

  input_parameters = jsonencode({
    report_type    = "weekly"
    include_charts = true
    recipients     = ["team@example.com"]
  })
}

# Minimal webhook job
resource "controlplane_job" "minimal_webhook" {
  name         = "test-job-minimal-webhook-ds"
  trigger_type = "webhook"
  prompt       = "Process webhook event"
  agent_id     = controlplane_agent.test.id
}

# Disabled job for enabled flag testing
resource "controlplane_job" "disabled" {
  name          = "test-job-disabled-ds"
  description   = "Disabled job for testing"
  trigger_type  = "cron"
  cron_schedule = "0 0 * * *"
  enabled       = false
  prompt        = "This job is disabled"
  agent_id      = controlplane_agent.test.id
}

# Data sources
data "controlplane_job" "minimal_cron_lookup" {
  id = controlplane_job.minimal_cron.id
}

data "controlplane_job" "full_cron_lookup" {
  id = controlplane_job.full_cron.id
}

data "controlplane_job" "minimal_webhook_lookup" {
  id = controlplane_job.minimal_webhook.id
}

data "controlplane_job" "disabled_lookup" {
  id = controlplane_job.disabled.id
}

# List data source
data "controlplane_jobs" "all" {}

# Outputs for tests
output "data_minimal_cron_name" {
  value = data.controlplane_job.minimal_cron_lookup.name
}

output "data_minimal_cron_trigger_type" {
  value = data.controlplane_job.minimal_cron_lookup.trigger_type
}

output "data_full_cron_description" {
  value = data.controlplane_job.full_cron_lookup.description
}

output "data_full_cron_cron_schedule" {
  value = data.controlplane_job.full_cron_lookup.cron_schedule
}

output "data_full_cron_planning_mode" {
  value = data.controlplane_job.full_cron_lookup.planning_mode
}

output "data_minimal_webhook_trigger_type" {
  value = data.controlplane_job.minimal_webhook_lookup.trigger_type
}

output "data_minimal_webhook_webhook_url" {
  value = data.controlplane_job.minimal_webhook_lookup.webhook_url
}

output "data_disabled_enabled" {
  value = tostring(data.controlplane_job.disabled_lookup.enabled)
}

output "all_jobs_count" {
  value = length(data.controlplane_jobs.all.jobs)
}

output "all_jobs_list" {
  value = jsonencode([for j in data.controlplane_jobs.all.jobs : j.name])
}
