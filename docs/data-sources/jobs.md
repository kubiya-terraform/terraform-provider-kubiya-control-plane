---
page_title: "controlplane_jobs Data Source"
subcategory: ""
description: |-
  Fetches all Kubiya jobs
---

# controlplane_jobs (Data Source)

Fetches all jobs from the Kubiya Control Plane.

## Example Usage

```terraform
data "controlplane_jobs" "all" {}

output "job_count" {
  value = length(data.controlplane_jobs.all.jobs)
}

output "enabled_jobs" {
  value = [for j in data.controlplane_jobs.all.jobs : j.name if j.enabled]
}

output "cron_jobs" {
  value = [for j in data.controlplane_jobs.all.jobs : j.name if j.trigger_type == "cron"]
}

# Find jobs by status
output "active_jobs" {
  value = [for j in data.controlplane_jobs.all.jobs : j.name if j.status == "active"]
}
```

## Schema

### Read-Only

- `jobs` (List of Object) List of jobs with the following attributes:
  - `id` (String) Job ID
  - `name` (String) Job name
  - `description` (String) Job description
  - `enabled` (Boolean) Whether the job is enabled
  - `status` (String) Job status
  - `trigger_type` (String) Trigger type: `cron`, `webhook`, or `manual`
  - `cron_schedule` (String) Cron expression
  - `cron_timezone` (String) Timezone for cron schedule
  - `webhook_url` (String) Full webhook URL (for webhook triggers)
  - `planning_mode` (String) Planning mode
  - `entity_type` (String) Entity type: `agent`, `team`, or `workflow`
  - `entity_id` (String) Entity ID
  - `prompt_template` (String) Prompt template
  - `system_prompt` (String) System prompt
  - `executor_type` (String) Executor routing type
  - `worker_queue_name` (String) Worker queue name
  - `environment_name` (String) Environment name
  - `created_at` (String) Timestamp when the job was created
  - `updated_at` (String) Timestamp when the job was last updated
