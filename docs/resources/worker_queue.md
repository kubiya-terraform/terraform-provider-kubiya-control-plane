---
page_title: "controlplane_worker_queue Resource"
subcategory: ""
description: |-
  Manages a Kubiya worker queue
---

# controlplane_worker_queue (Resource)

Manages a worker queue in the Kubiya Control Plane. Worker queues are used to organize and manage workers within an environment, providing fine-grained control over worker resources and task distribution.

## Example Usage

```terraform
resource "controlplane_environment" "production" {
  name         = "production"
  display_name = "Production Environment"
  description  = "Main production environment"
}

resource "controlplane_worker_queue" "default" {
  environment_id     = controlplane_environment.production.id
  name               = "default-queue"
  display_name       = "Default Worker Queue"
  description        = "Main worker queue for production"
  heartbeat_interval = 30
  max_workers        = 10

  tags = ["production", "primary"]

  settings = {
    region = "us-east-1"
    tier   = "production"
  }
}
```

## Schema

### Required

- `environment_id` (String) ID of the environment this worker queue belongs to
- `name` (String) Worker queue name (lowercase, no spaces, 2-50 characters)

### Optional

- `display_name` (String) User-friendly display name
- `description` (String) Queue description
- `status` (String) Worker queue status (active, inactive, paused). Default: "active"
- `max_workers` (Number) Maximum number of workers allowed (null = unlimited)
- `heartbeat_interval` (Number) Seconds between heartbeats (10-300). Default: 30
- `tags` (List of String) Tags for the worker queue
- `settings` (Map of String) Additional settings as key-value pairs

### Read-Only

- `id` (String) The unique identifier of the worker queue
- `created_at` (String) Timestamp when the worker queue was created
- `updated_at` (String) Timestamp when the worker queue was last updated
- `active_workers` (Number) Number of currently active workers in the queue
- `task_queue_name` (String) Temporal task queue name for this worker queue

## Important Notes

- Worker queues provide fine-grained control over worker resources within an environment
- The `task_queue_name` should be used when starting workers to connect to this specific queue
- Workers automatically register with the queue when they start using the queue ID
- Active workers are tracked in real-time using heartbeats
- Deleting a queue with active workers will fail - workers must be stopped first

## Import

Worker queues can be imported using their ID:

```shell
terraform import controlplane_worker_queue.example queue-uuid-here
```
