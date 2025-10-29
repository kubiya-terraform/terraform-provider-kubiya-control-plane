---
page_title: "kubiya_control_plane_worker Resource"
subcategory: ""
description: |-
  Manages a Kubiya worker registration
---

# kubiya_control_plane_worker (Resource)

Manages a worker registration in the Kubiya Control Plane. Workers are execution nodes that run agent tasks. Note that workers typically self-register at runtime, so this resource is primarily for pre-registration and discovery.

## Example Usage

```terraform
resource "kubiya_control_plane_worker" "example" {
  environment_name = "production"
  hostname         = "worker-node-01"

  metadata = jsonencode({
    region     = "us-east-1"
    datacenter = "dc1"
    capacity   = "high"
    tags = {
      environment = "production"
      team        = "platform"
    }
  })
}
```

## Schema

### Required

- `environment_name` (String) Name of the environment this worker belongs to

### Optional

- `hostname` (String) Hostname of the worker
- `metadata` (String) Worker metadata as JSON string

### Read-Only

- `id` (String) The unique identifier of the worker
- `status` (String) Current status of the worker
- `registered_at` (String) Timestamp when the worker was registered
- `last_heartbeat` (String) Timestamp of the last heartbeat
- `updated_at` (String) Timestamp when the worker was last updated

## Important Notes

- Workers are runtime entities that manage their own lifecycle
- Workers typically self-register when they start
- This resource is primarily for pre-registration and discovery
- Workers connect to the control plane using the environment's worker token

## Import

Workers can be imported using their ID:

```shell
terraform import kubiya_control_plane_worker.example worker-uuid-here
```
