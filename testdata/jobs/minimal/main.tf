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
  name = "test-agent-for-job"
}

# Minimal cron job (required fields only)
resource "controlplane_job" "minimal_cron" {
  name         = "test-job-minimal-cron"
  trigger_type = "cron"
  cron_schedule = "0 9 * * *"
  prompt       = "Generate daily report"
  agent_id     = controlplane_agent.test.id
}

# Data source test
data "controlplane_job" "minimal_cron_lookup" {
  id = controlplane_job.minimal_cron.id
}

# Outputs
output "job_id" {
  value = controlplane_job.minimal_cron.id
}

output "job_name" {
  value = controlplane_job.minimal_cron.name
}

output "job_trigger_type" {
  value = controlplane_job.minimal_cron.trigger_type
}

output "job_cron_schedule" {
  value = controlplane_job.minimal_cron.cron_schedule
}

output "data_job_trigger_type" {
  value = data.controlplane_job.minimal_cron_lookup.trigger_type
}
