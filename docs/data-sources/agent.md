---
page_title: "controlplane_agent Data Source"
subcategory: ""
description: |-
  Retrieves information about a Kubiya agent
---

# controlplane_agent (Data Source)

Retrieves information about an existing AI agent in the Kubiya Control Plane.

## Example Usage

```terraform
# Look up an agent by ID
data "controlplane_agent" "existing" {
  id = "agent-uuid-here"
}

# Use the agent information
output "agent_name" {
  value = data.controlplane_agent.existing.name
}

output "agent_status" {
  value = data.controlplane_agent.existing.status
}

# Reference in another resource
resource "controlplane_agent" "new" {
  name    = "new-agent"
  team_id = data.controlplane_agent.existing.team_id
  # ...
}
```

## Schema

### Required

- `id` (String) The unique identifier of the agent to look up

### Read-Only

- `name` (String) The name of the agent
- `description` (String) Description of the agent
- `status` (String) Current status of the agent
- `capabilities` (List of String) List of agent capabilities
- `configuration` (String) Agent configuration as JSON string
- `model_id` (String) LiteLLM model identifier
- `llm_config` (String) LLM configuration as JSON string
- `runtime` (String) Runtime type
- `team_id` (String) ID of the team this agent belongs to
- `created_at` (String) Creation timestamp
- `updated_at` (String) Last update timestamp
