---
page_title: "kubiya_control_plane_agent Resource"
subcategory: ""
description: |-
  Manages a Kubiya AI agent
---

# kubiya_control_plane_agent (Resource)

Manages an AI agent in the Kubiya Control Plane. Agents are autonomous AI assistants that can execute tasks, interact with tools, and follow policies.

## Example Usage

```terraform
# Basic agent
resource "kubiya_control_plane_agent" "example" {
  name        = "my-agent"
  description = "An example AI agent"
  model_id    = "gpt-4"
  runtime     = "default"

  llm_config = jsonencode({
    temperature = 0.7
    max_tokens  = 2000
  })

  capabilities = ["code_execution", "file_operations"]
}

# Agent with team association
resource "kubiya_control_plane_agent" "team_agent" {
  name     = "team-agent"
  model_id = "gpt-4"
  runtime  = "default"
  team_id  = kubiya_control_plane_team.devops.id
}
```

## Schema

### Required

- `name` (String) The name of the agent
- `model_id` (String) LiteLLM model identifier (e.g., "gpt-4", "claude-3-opus")
- `runtime` (String) Runtime type. Valid values: `default`, `claude_code`

### Optional

- `description` (String) Description of the agent's purpose
- `capabilities` (List of String) List of agent capabilities
- `configuration` (String) Agent configuration as JSON string
- `llm_config` (String) LLM configuration as JSON string (temperature, max_tokens, etc.)
- `team_id` (String) ID of the team this agent belongs to

### Read-Only

- `id` (String) The unique identifier of the agent
- `status` (String) Current status of the agent (idle, running, paused, completed, failed, stopped)
- `created_at` (String) Timestamp when the agent was created
- `updated_at` (String) Timestamp when the agent was last updated

## Import

Agents can be imported using their ID:

```shell
terraform import kubiya_control_plane_agent.example agent-uuid-here
```
