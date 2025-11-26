---
page_title: "controlplane_agent Resource"
subcategory: ""
description: |-
  Manages a Kubiya AI agent
---

# controlplane_agent (Resource)

Manages an AI agent in the Kubiya Control Plane. Agents are autonomous AI assistants that can execute tasks, interact with tools, and follow policies.

## Example Usage

```terraform
# Basic agent
resource "controlplane_agent" "example" {
  name        = "my-agent"
  description = "An example AI agent"
  model_id    = "gpt-4"
  runtime     = "default"

  llm_config = jsonencode({
    temperature = 0.7
    max_tokens  = 2000
  })

  capabilities = ["code_execution", "file_operations"]

  configuration = jsonencode({
    max_retries = 3
    timeout     = 300
  })
}

# Agent with skills
resource "controlplane_agent" "with_skills" {
  name        = "platform-engineer"
  description = "Platform engineering agent with shell and Python skills"
  model_id    = "claude-sonnet-4"
  runtime     = "claude_code"

  # Assign skills to the agent
  # Note: Use 'skills' (not 'skill_ids') in Terraform
  skills = [
    controlplane_skill.shell.id,
    controlplane_skill.python.id
  ]

  system_prompt = <<-EOT
    You are a platform engineering assistant.
    Use your skills to help with infrastructure tasks.
  EOT

  configuration = jsonencode({
    mcpServers = {
      kubiya = {
        command = "kubiya"
        args    = ["mcp", "serve"]
      }
    }
  })
}
```

## Schema

### Required

- `name` (String) The name of the agent

### Optional

- `description` (String) Description of the agent's purpose
- `model_id` (String) LiteLLM model identifier (e.g., "gpt-4", "claude-sonnet-4")
- `runtime` (String) Runtime type. Valid values: `default`, `claude_code`
- `capabilities` (List of String) List of agent capabilities
- `configuration` (String) Agent configuration as JSON string
- `llm_config` (String) LLM configuration as JSON string (temperature, max_tokens, etc.)
- `team_id` (String) Team ID to assign this agent to
- `system_prompt` (String) System prompt for the agent
- `skills` (List of String) List of skill resource IDs to assign to the agent. **Important:** Use `skills` in Terraform (this maps to `skill_ids` in the API). Example: `skills = [controlplane_skill.shell.id]`
- `execution_environment` (String) Execution environment configuration as JSON string

### Read-Only

- `id` (String) The unique identifier of the agent
- `status` (String) Current status of the agent (idle, running, paused, completed, failed, stopped)
- `created_at` (String) Timestamp when the agent was created
- `updated_at` (String) Timestamp when the agent was last updated

## Import

Agents can be imported using their ID:

```shell
terraform import controlplane_agent.example agent-uuid-here
```
