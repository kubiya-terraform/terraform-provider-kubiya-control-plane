
provider "controlplane" {
  # Configuration via environment variables:
  # KUBIYA_CONTROL_PLANE_API_KEY
  # KUBIYA_CONTROL_PLANE_BASE_URL (optional, defaults to https://control-plane.kubiya.ai)
}

# Create dependencies for job testing
resource "controlplane_agent" "test_agent" {
  name        = "test-agent-for-jobs"
  description = "Agent for job testing"
}

resource "controlplane_team" "test_team" {
  name        = "test-team-for-jobs"
  description = "Team for job testing"
}

resource "controlplane_environment" "test_env" {
  name        = "test-env-for-jobs"
  description = "Environment for job testing"
}

resource "controlplane_worker_queue" "test_queue" {
  environment_id = controlplane_environment.test_env.id
  name           = "test-queue-for-jobs"
  description    = "Worker queue for job testing"
}

# Test 1: Minimal cron job (required fields only)
resource "controlplane_job" "minimal_cron" {
  name            = "test-job-minimal-cron"
  trigger_type    = "cron"
  cron_schedule   = "0 9 * * *"
  prompt_template = "Execute daily report"
}

# Test 2: Full cron job with all optional fields
resource "controlplane_job" "full_cron" {
  name            = "test-job-full-cron"
  description     = "Comprehensive cron job with all fields"
  enabled         = true
  trigger_type    = "cron"
  cron_schedule   = "0 17 * * 1-5"  # 5pm weekdays
  cron_timezone   = "America/New_York"
  planning_mode   = "predefined_agent"
  entity_type     = "agent"
  entity_id       = controlplane_agent.test_agent.id
  prompt_template = "Daily report for {{date}} in {{environment}}"
  system_prompt   = "You are a helpful assistant that generates daily reports"
  executor_type   = "specific_queue"
  worker_queue_name = controlplane_worker_queue.test_queue.name

  config = jsonencode({
    timeout     = 3600
    max_retries = 3
    priority    = "high"
  })

  execution_env_vars = {
    REPORT_TYPE = "daily"
    OUTPUT_FORMAT = "pdf"
    NOTIFICATION_ENABLED = "true"
  }

  execution_secrets = ["api-key", "db-password"]
  execution_integrations = ["slack-integration", "email-integration"]
}

# Test 3: Minimal webhook job
resource "controlplane_job" "minimal_webhook" {
  name            = "test-job-minimal-webhook"
  trigger_type    = "webhook"
  prompt_template = "Process webhook event: {{event_type}}"
}

# Test 4: Full webhook job
resource "controlplane_job" "full_webhook" {
  name            = "test-job-full-webhook"
  description     = "Comprehensive webhook job"
  enabled         = true
  trigger_type    = "webhook"
  planning_mode   = "predefined_team"
  entity_type     = "team"
  entity_id       = controlplane_team.test_team.id
  prompt_template = "Handle {{event_type}} for {{resource_id}}"
  system_prompt   = "Process incoming webhook events efficiently"
  executor_type   = "environment"
  environment_name = controlplane_environment.test_env.name

  config = jsonencode({
    timeout = 1800
    webhook_validation = {
      enabled = true
      signature_header = "X-Webhook-Signature"
    }
  })

  execution_env_vars = {
    WEBHOOK_HANDLER = "default"
    VALIDATE_PAYLOAD = "true"
  }
}

# Test 5: Minimal manual job
resource "controlplane_job" "minimal_manual" {
  name            = "test-job-minimal-manual"
  trigger_type    = "manual"
  prompt_template = "Execute manual task"
}

# Test 6: Full manual job
resource "controlplane_job" "full_manual" {
  name            = "test-job-full-manual"
  description     = "Comprehensive manual job"
  enabled         = true
  trigger_type    = "manual"
  planning_mode   = "on_the_fly"
  prompt_template = "Execute task: {{task_description}}"
  system_prompt   = "Follow the task instructions carefully"
  executor_type   = "auto"

  config = jsonencode({
    require_confirmation = true
    allow_parameters = true
  })
}

# Test 7: Cron job with predefined_workflow planning mode
resource "controlplane_job" "cron_workflow" {
  name            = "test-job-cron-workflow"
  description     = "Cron job with workflow planning"
  trigger_type    = "cron"
  cron_schedule   = "0 0 * * 0"  # Weekly on Sunday midnight
  cron_timezone   = "UTC"
  planning_mode   = "predefined_workflow"
  entity_type     = "workflow"
  entity_id       = "workflow-123"
  prompt_template = "Run weekly workflow"
}

# Test 8: Job with specific_queue executor
resource "controlplane_job" "specific_queue" {
  name              = "test-job-specific-queue"
  description       = "Job routed to specific queue"
  trigger_type      = "cron"
  cron_schedule     = "0 */4 * * *"  # Every 4 hours
  executor_type     = "specific_queue"
  worker_queue_name = controlplane_worker_queue.test_queue.name
  prompt_template   = "Execute on specific queue"
}

# Test 9: Job with environment executor
resource "controlplane_job" "environment_executor" {
  name             = "test-job-environment-executor"
  description      = "Job routed to environment"
  trigger_type     = "cron"
  cron_schedule    = "30 12 * * *"  # Daily at 12:30pm
  executor_type    = "environment"
  environment_name = controlplane_environment.test_env.name
  prompt_template  = "Execute in environment"
}

# Test 10: Job with auto executor (default)
resource "controlplane_job" "auto_executor" {
  name            = "test-job-auto-executor"
  description     = "Job with auto executor routing"
  trigger_type    = "manual"
  executor_type   = "auto"
  prompt_template = "Execute with auto routing"
}

# Test 11: Cron job with different timezones
resource "controlplane_job" "timezone_pst" {
  name            = "test-job-timezone-pst"
  trigger_type    = "cron"
  cron_schedule   = "0 9 * * *"
  cron_timezone   = "America/Los_Angeles"
  prompt_template = "PST morning job"
}

resource "controlplane_job" "timezone_tokyo" {
  name            = "test-job-timezone-tokyo"
  trigger_type    = "cron"
  cron_schedule   = "0 9 * * *"
  cron_timezone   = "Asia/Tokyo"
  prompt_template = "Tokyo morning job"
}

# Test 12: Job with execution_env_vars only
resource "controlplane_job" "env_vars_only" {
  name            = "test-job-env-vars-only"
  description     = "Job with execution environment variables"
  trigger_type    = "manual"
  prompt_template = "Execute with env vars"

  execution_env_vars = {
    VAR1 = "value1"
    VAR2 = "value2"
    VAR3 = "value3"
  }
}

# Test 13: Job with execution_secrets only
resource "controlplane_job" "secrets_only" {
  name            = "test-job-secrets-only"
  description     = "Job with execution secrets"
  trigger_type    = "manual"
  prompt_template = "Execute with secrets"

  execution_secrets = ["secret-alpha", "secret-beta", "secret-gamma"]
}

# Test 14: Job with execution_integrations only
resource "controlplane_job" "integrations_only" {
  name            = "test-job-integrations-only"
  description     = "Job with execution integrations"
  trigger_type    = "manual"
  prompt_template = "Execute with integrations"

  execution_integrations = ["integration-1", "integration-2"]
}

# Test 15: Job with complex config
resource "controlplane_job" "complex_config" {
  name            = "test-job-complex-config"
  description     = "Job with complex configuration"
  trigger_type    = "cron"
  cron_schedule   = "0 2 * * *"
  prompt_template = "Execute with complex config"

  config = jsonencode({
    timeout     = 7200
    max_retries = 5
    retry_delay = 60
    priority    = "critical"
    notifications = {
      on_success = true
      on_failure = true
      channels   = ["email", "slack"]
    }
    resource_limits = {
      cpu    = "2"
      memory = "4Gi"
    }
    metadata = {
      owner = "platform-team"
      cost_center = "engineering"
    }
  })
}

# Test 16: Disabled job
resource "controlplane_job" "disabled" {
  name            = "test-job-disabled"
  description     = "Disabled job for testing"
  enabled         = false
  trigger_type    = "cron"
  cron_schedule   = "0 0 * * *"
  prompt_template = "This job is disabled"
}

# Test 17: Job with template variables
resource "controlplane_job" "template_variables" {
  name            = "test-job-template-variables"
  description     = "Job with multiple template variables"
  trigger_type    = "manual"
  prompt_template = "Process {{entity}} with {{action}} in {{environment}} for {{user}}"
  system_prompt   = "Handle the request according to the parameters provided"
}

# Test 18: Frequent cron job (every minute)
resource "controlplane_job" "frequent_cron" {
  name            = "test-job-frequent"
  description     = "Frequently executing cron job"
  trigger_type    = "cron"
  cron_schedule   = "* * * * *"  # Every minute
  prompt_template = "Execute frequently"
}

# Test 19: Monthly cron job
resource "controlplane_job" "monthly_cron" {
  name            = "test-job-monthly"
  description     = "Monthly cron job"
  trigger_type    = "cron"
  cron_schedule   = "0 0 1 * *"  # First day of month at midnight
  cron_timezone   = "UTC"
  prompt_template = "Execute monthly report"
}

# Test 20: Job for update testing
resource "controlplane_job" "for_update" {
  name            = "test-job-for-update"
  description     = "Initial description"
  enabled         = true
  trigger_type    = "cron"
  cron_schedule   = "0 12 * * *"
  prompt_template = "Initial prompt"

  config = jsonencode({
    version = 1
  })
}

# Data source tests
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

# Test jobs data source (list all jobs)
data "controlplane_jobs" "all_jobs" {
}

# Outputs for test validation
output "minimal_cron_job_id" {
  value = controlplane_job.minimal_cron.id
}

output "minimal_cron_job_name" {
  value = controlplane_job.minimal_cron.name
}

output "minimal_cron_job_trigger_type" {
  value = controlplane_job.minimal_cron.trigger_type
}

output "minimal_cron_job_cron_schedule" {
  value = controlplane_job.minimal_cron.cron_schedule
}

output "minimal_cron_job_status" {
  value = controlplane_job.minimal_cron.status
}

output "minimal_cron_job_created_at" {
  value = controlplane_job.minimal_cron.created_at
}

output "full_cron_job_id" {
  value = controlplane_job.full_cron.id
}

output "full_cron_job_name" {
  value = controlplane_job.full_cron.name
}

output "full_cron_job_description" {
  value = controlplane_job.full_cron.description
}

output "full_cron_job_enabled" {
  value = controlplane_job.full_cron.enabled
}

output "full_cron_job_trigger_type" {
  value = controlplane_job.full_cron.trigger_type
}

output "full_cron_job_cron_schedule" {
  value = controlplane_job.full_cron.cron_schedule
}

output "full_cron_job_cron_timezone" {
  value = controlplane_job.full_cron.cron_timezone
}

output "full_cron_job_planning_mode" {
  value = controlplane_job.full_cron.planning_mode
}

output "full_cron_job_entity_type" {
  value = controlplane_job.full_cron.entity_type
}

output "full_cron_job_entity_id" {
  value = controlplane_job.full_cron.entity_id
}

output "full_cron_job_prompt_template" {
  value = controlplane_job.full_cron.prompt_template
}

output "full_cron_job_system_prompt" {
  value = controlplane_job.full_cron.system_prompt
}

output "full_cron_job_executor_type" {
  value = controlplane_job.full_cron.executor_type
}

output "full_cron_job_worker_queue_name" {
  value = controlplane_job.full_cron.worker_queue_name
}

output "full_cron_job_config" {
  value     = controlplane_job.full_cron.config
  sensitive = true
}

output "full_cron_job_execution_env_vars" {
  value     = controlplane_job.full_cron.execution_env_vars
  sensitive = true
}

output "full_cron_job_execution_secrets" {
  value     = controlplane_job.full_cron.execution_secrets
  sensitive = true
}

output "full_cron_job_execution_integrations" {
  value = controlplane_job.full_cron.execution_integrations
}

output "full_cron_job_status" {
  value = controlplane_job.full_cron.status
}

output "full_cron_job_created_at" {
  value = controlplane_job.full_cron.created_at
}

output "full_cron_job_updated_at" {
  value = controlplane_job.full_cron.updated_at
}

output "minimal_webhook_job_id" {
  value = controlplane_job.minimal_webhook.id
}

output "minimal_webhook_job_trigger_type" {
  value = controlplane_job.minimal_webhook.trigger_type
}

output "minimal_webhook_job_webhook_url" {
  value     = controlplane_job.minimal_webhook.webhook_url
  sensitive = true
}

output "minimal_webhook_job_webhook_secret" {
  value     = controlplane_job.minimal_webhook.webhook_secret
  sensitive = true
}

output "full_webhook_job_id" {
  value = controlplane_job.full_webhook.id
}

output "full_webhook_job_webhook_url" {
  value     = controlplane_job.full_webhook.webhook_url
  sensitive = true
}

output "full_webhook_job_planning_mode" {
  value = controlplane_job.full_webhook.planning_mode
}

output "full_webhook_job_entity_type" {
  value = controlplane_job.full_webhook.entity_type
}

output "minimal_manual_job_id" {
  value = controlplane_job.minimal_manual.id
}

output "minimal_manual_job_trigger_type" {
  value = controlplane_job.minimal_manual.trigger_type
}

output "full_manual_job_id" {
  value = controlplane_job.full_manual.id
}

output "full_manual_job_planning_mode" {
  value = controlplane_job.full_manual.planning_mode
}

output "cron_workflow_job_id" {
  value = controlplane_job.cron_workflow.id
}

output "cron_workflow_job_planning_mode" {
  value = controlplane_job.cron_workflow.planning_mode
}

output "cron_workflow_job_entity_type" {
  value = controlplane_job.cron_workflow.entity_type
}

output "specific_queue_job_id" {
  value = controlplane_job.specific_queue.id
}

output "specific_queue_job_executor_type" {
  value = controlplane_job.specific_queue.executor_type
}

output "specific_queue_job_worker_queue_name" {
  value = controlplane_job.specific_queue.worker_queue_name
}

output "environment_executor_job_id" {
  value = controlplane_job.environment_executor.id
}

output "environment_executor_job_executor_type" {
  value = controlplane_job.environment_executor.executor_type
}

output "environment_executor_job_environment_name" {
  value = controlplane_job.environment_executor.environment_name
}

output "auto_executor_job_id" {
  value = controlplane_job.auto_executor.id
}

output "auto_executor_job_executor_type" {
  value = controlplane_job.auto_executor.executor_type
}

output "timezone_pst_job_id" {
  value = controlplane_job.timezone_pst.id
}

output "timezone_pst_job_cron_timezone" {
  value = controlplane_job.timezone_pst.cron_timezone
}

output "timezone_tokyo_job_id" {
  value = controlplane_job.timezone_tokyo.id
}

output "timezone_tokyo_job_cron_timezone" {
  value = controlplane_job.timezone_tokyo.cron_timezone
}

output "disabled_job_id" {
  value = controlplane_job.disabled.id
}

output "disabled_job_enabled" {
  value = controlplane_job.disabled.enabled
}

output "complex_config_job_id" {
  value = controlplane_job.complex_config.id
}

output "complex_config_job_config" {
  value     = controlplane_job.complex_config.config
  sensitive = true
}

output "for_update_job_id" {
  value = controlplane_job.for_update.id
}

# Data source outputs for validation
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
  value     = data.controlplane_job.minimal_webhook_lookup.webhook_url
  sensitive = true
}

output "data_disabled_enabled" {
  value = data.controlplane_job.disabled_lookup.enabled
}

# Jobs data source outputs
output "all_jobs_count" {
  value = length(data.controlplane_jobs.all_jobs.jobs)
}

output "all_jobs_list" {
  value = [for j in data.controlplane_jobs.all_jobs.jobs : j.name]
}
