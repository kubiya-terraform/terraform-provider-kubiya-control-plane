---
page_title: "controlplane_job Data Source"
subcategory: ""
description: |-
  Fetches a Kubiya job
---

# controlplane_job (Data Source)

Fetches a job from the Kubiya Control Plane by ID.

## Example Usage

```terraform
data "controlplane_job" "daily_report" {
  id = "job_123abc"
}

output "job_trigger_type" {
  value = data.controlplane_job.daily_report.trigger_type
}

output "job_status" {
  value = data.controlplane_job.daily_report.status
}

# Use the webhook URL if it's a webhook-triggered job
output "webhook_url" {
  value = data.controlplane_job.daily_report.webhook_url
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Required) Job ID to fetch.

## Attribute Reference

In addition to the argument above, the following attributes are exported:

* `name` - Job name.
* `description` - Job description.
* `enabled` - Whether the job is enabled.
* `status` - Job status.
* `trigger_type` - Trigger type: `cron`, `webhook`, or `manual`.
* `cron_schedule` - Cron expression.
* `cron_timezone` - Timezone for cron schedule.
* `webhook_url` - Full webhook URL (for webhook triggers).
* `planning_mode` - Planning mode.
* `entity_type` - Entity type: `agent`, `team`, or `workflow`.
* `entity_id` - Entity ID.
* `prompt_template` - Prompt template.
* `system_prompt` - System prompt.
* `executor_type` - Executor routing type.
* `worker_queue_name` - Worker queue name.
* `environment_name` - Environment name.
* `created_at` - Timestamp when the job was created.
* `updated_at` - Timestamp when the job was last updated.
