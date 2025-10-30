---
page_title: "controlplane_worker_queues Data Source"
subcategory: ""
description: |-
  Fetches all Kubiya worker queues in an environment
---

# controlplane_worker_queues (Data Source)

Fetches all worker queues in a specific environment from the Kubiya Control Plane.

## Example Usage

```terraform
data "controlplane_environment" "production" {
  name = "production"
}

data "controlplane_worker_queues" "all" {
  environment_id = data.controlplane_environment.production.id
}

output "queue_count" {
  value = length(data.controlplane_worker_queues.all.queues)
}

output "queue_names" {
  value = [for q in data.controlplane_worker_queues.all.queues : q.name]
}

output "total_active_workers" {
  value = sum([for q in data.controlplane_worker_queues.all.queues : q.active_workers])
}
```

## Schema

### Required

- `environment_id` (String) Environment ID to list worker queues from

### Read-Only

- `queues` (List of Object) List of worker queues with the following attributes:
  - `id` (String) Worker Queue ID
  - `environment_id` (String) Environment ID
  - `name` (String) Worker queue name
  - `display_name` (String) User-friendly display name
  - `description` (String) Queue description
  - `status` (String) Worker queue status
  - `max_workers` (Number) Maximum number of workers allowed
  - `heartbeat_interval` (Number) Seconds between heartbeats
  - `tags` (List of String) Tags for the worker queue
  - `settings` (Map of String) Additional settings as key-value pairs
  - `created_at` (String) Creation timestamp
  - `updated_at` (String) Last update timestamp
  - `active_workers` (Number) Number of currently active workers
  - `task_queue_name` (String) Temporal task queue name
