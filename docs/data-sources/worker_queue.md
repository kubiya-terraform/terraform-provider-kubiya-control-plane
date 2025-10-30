---
page_title: "controlplane_worker_queue Data Source"
subcategory: ""
description: |-
  Fetches a Kubiya worker queue by ID
---

# controlplane_worker_queue (Data Source)

Fetches a worker queue from the Kubiya Control Plane by its ID.

## Example Usage

```terraform
data "controlplane_worker_queue" "example" {
  id = "queue-uuid-here"
}

output "queue_name" {
  value = data.controlplane_worker_queue.example.name
}

output "active_workers" {
  value = data.controlplane_worker_queue.example.active_workers
}
```

## Schema

### Required

- `id` (String) Worker Queue ID

### Read-Only

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
