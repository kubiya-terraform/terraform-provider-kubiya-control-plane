---
page_title: "controlplane_job Resource"
subcategory: ""
description: |-
  Manages a Kubiya job
---

# controlplane_job (Resource)

Manages a job in the Kubiya Control Plane. Jobs can be triggered by cron schedules, webhooks, or manually, and can execute agents, teams, or workflows with configurable prompts.

## Example Usage

```terraform
# Cron-triggered job
resource "controlplane_job" "daily_report" {
  name         = "daily-report"
  description  = "Generate daily report at 5pm"
  enabled      = true
  trigger_type = "cron"
  cron_schedule = "0 17 * * *"
  cron_timezone = "America/New_York"

  planning_mode   = "predefined_agent"
  entity_type     = "agent"
  entity_id       = controlplane_agent.reporter.id
  prompt_template = "Generate a daily report for {{date}}"

  executor_type = "auto"

  execution_env_vars = {
    REPORT_FORMAT = "pdf"
  }

  execution_secrets = ["slack_token"]
}

# Webhook-triggered job
resource "controlplane_job" "github_webhook" {
  name         = "github-pr-handler"
  description  = "Handle GitHub PR events"
  enabled      = true
  trigger_type = "webhook"

  planning_mode   = "predefined_team"
  entity_type     = "team"
  entity_id       = controlplane_team.devops.id
  prompt_template = "Process GitHub PR: {{pr_number}} - {{pr_title}}"
  system_prompt   = "You are a DevOps assistant handling pull requests."

  executor_type = "environment"
  environment_name = "production"

  execution_integrations = ["github_integration_id"]
}

# Manual job with specific worker queue
resource "controlplane_job" "manual_deployment" {
  name         = "manual-deployment"
  description  = "Manual deployment trigger"
  enabled      = true
  trigger_type = "manual"

  planning_mode   = "predefined_workflow"
  entity_type     = "workflow"
  entity_id       = "workflow_123"
  prompt_template = "Deploy version {{version}} to {{environment}}"

  executor_type = "specific_queue"
  worker_queue_name = "deployment-queue"

  config = jsonencode({
    timeout = 3600
    retry_policy = {
      max_attempts = 3
    }
  })
}

# On-the-fly planning job
resource "controlplane_job" "flexible_task" {
  name         = "flexible-task"
  description  = "Job with on-the-fly planning"
  enabled      = true
  trigger_type = "manual"

  planning_mode   = "on_the_fly"
  prompt_template = "Analyze and handle: {{task_description}}"
  system_prompt   = "You are a general-purpose AI assistant."

  executor_type = "auto"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Job name.
* `trigger_type` - (Required) Trigger type. Must be one of: `cron`, `webhook`, or `manual`.
* `prompt_template` - (Required) Prompt template. Can include `{{variables}}` for dynamic parameters.
* `description` - (Optional) Job description.
* `enabled` - (Optional) Whether the job is enabled. Defaults to `true`.
* `cron_schedule` - (Optional) Cron expression (e.g., `0 17 * * *` for daily at 5pm). Required when `trigger_type` is `cron`.
* `cron_timezone` - (Optional) Timezone for cron schedule (e.g., `America/New_York`). Defaults to `UTC`.
* `planning_mode` - (Optional) Planning mode. Must be one of: `on_the_fly`, `predefined_agent`, `predefined_team`, or `predefined_workflow`. Defaults to `predefined_agent`.
* `entity_type` - (Optional) Entity type: `agent`, `team`, or `workflow`. Required when `planning_mode` is not `on_the_fly`.
* `entity_id` - (Optional) Entity ID (agent_id, team_id, or workflow_id). Required when `planning_mode` is not `on_the_fly`.
* `system_prompt` - (Optional) System prompt for the job execution.
* `executor_type` - (Optional) Executor routing. Must be one of: `auto`, `specific_queue`, or `environment`. Defaults to `auto`.
* `worker_queue_name` - (Optional) Worker queue name. Required when `executor_type` is `specific_queue`.
* `environment_name` - (Optional) Environment name. Required when `executor_type` is `environment`.
* `config` - (Optional) Additional execution config as JSON string (timeout, retry, etc.).
* `execution_env_vars` - (Optional) Map of environment variables to inject into execution.
* `execution_secrets` - (Optional) List of secret names to inject into execution environment.
* `execution_integrations` - (Optional) List of integration IDs to inject into execution environment.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Job ID.
* `status` - Job status.
* `webhook_url` - Full webhook URL (generated for webhook triggers).
* `webhook_secret` - Webhook HMAC secret for signature verification (sensitive).
* `created_at` - Timestamp when the job was created.
* `updated_at` - Timestamp when the job was last updated.

## Import

Jobs can be imported using their ID:

```shell
terraform import controlplane_job.example job_123abc
```
