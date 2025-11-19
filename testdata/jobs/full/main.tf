terraform {
  required_providers {
    controlplane = {
      source = "kubiya/control-plane"
    }
  }
}

provider "controlplane" {}

# Create agent for job
resource "controlplane_agent" "test" {
  name = "test-agent-for-full-job"
}

# Full cron job with all optional fields
resource "controlplane_job" "full_cron" {
  name          = "test-job-full-cron"
  description   = "Comprehensive test cron job with all fields configured"
  trigger_type  = "cron"
  cron_schedule = "0 17 * * 1-5"
  cron_timezone = "America/New_York"
  enabled       = true

  prompt       = "Generate comprehensive weekly report with all metrics"
  agent_id     = controlplane_agent.test.id
  planning_mode = "predefined_agent"

  executor_type = "auto"

  input_parameters = jsonencode({
    report_type = "weekly"
    include_charts = true
    recipients = ["team@example.com"]
  })
}

# Data source test
data "controlplane_job" "full_cron_lookup" {
  id = controlplane_job.full_cron.id
}

# Outputs
output "job_id" {
  value = controlplane_job.full_cron.id
}

output "job_name" {
  value = controlplane_job.full_cron.name
}

output "job_description" {
  value = controlplane_job.full_cron.description
}

output "job_trigger_type" {
  value = controlplane_job.full_cron.trigger_type
}

output "job_cron_schedule" {
  value = controlplane_job.full_cron.cron_schedule
}

output "job_cron_timezone" {
  value = controlplane_job.full_cron.cron_timezone
}

output "job_enabled" {
  value = controlplane_job.full_cron.enabled
}

output "job_planning_mode" {
  value = controlplane_job.full_cron.planning_mode
}

output "data_job_description" {
  value = data.controlplane_job.full_cron_lookup.description
}

output "data_job_planning_mode" {
  value = data.controlplane_job.full_cron_lookup.planning_mode
}
