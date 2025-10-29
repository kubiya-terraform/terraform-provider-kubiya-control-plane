---
page_title: "kubiya_control_plane_environment Resource"
subcategory: ""
description: |-
  Manages a Kubiya execution environment
---

# kubiya_control_plane_environment (Resource)

Manages an execution environment in the Kubiya Control Plane. Environments define the execution context for agents including variables, secrets, and integrations.

## Example Usage

```terraform
resource "kubiya_control_plane_environment" "production" {
  name        = "production"
  description = "Production environment for agents"

  configuration = jsonencode({
    region           = "us-east-1"
    max_workers      = 10
    auto_scaling     = true
    retention_days   = 90
    notification_url = "https://hooks.slack.com/example"
  })

  execution_environment = jsonencode({
    env_vars = {
      LOG_LEVEL = "info"
      APP_ENV   = "production"
    }
  })
}
```

## Schema

### Required

- `name` (String) The name of the environment

### Optional

- `description` (String) Description of the environment
- `configuration` (String) Environment configuration as JSON string
- `execution_environment` (String) Execution environment variables and settings as JSON string

### Read-Only

- `id` (String) The unique identifier of the environment
- `status` (String) Current status of the environment
- `created_at` (String) Timestamp when the environment was created
- `updated_at` (String) Timestamp when the environment was last updated

## Import

Environments can be imported using their ID:

```shell
terraform import kubiya_control_plane_environment.production environment-uuid-here
```
