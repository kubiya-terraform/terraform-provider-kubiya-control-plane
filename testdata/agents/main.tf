
provider "controlplane" {
  # Configuration via environment variables:
  # KUBIYA_CONTROL_PLANE_API_KEY
  # KUBIYA_CONTROL_PLANE_BASE_URL (optional, defaults to https://control-plane.kubiya.ai)
}

# Full agent with all optional fields
resource "controlplane_agent" "full" {
  name        = "test-agent-full-ds"
  description = "Comprehensive test agent for data source testing"

  status = "idle"

  model_id = "gpt-4"
  llm_config = jsonencode({
    temperature = 0.7
    max_tokens  = 2000
  })

  runtime = "default"

  capabilities = ["code_execution", "file_operations", "web_search"]

  configuration = jsonencode({
    max_retries = 3
    timeout     = 300
  })
}

# Claude Code agent
resource "controlplane_agent" "claude_code" {
  name        = "test-agent-claude-code-ds"
  description = "Agent using Claude Code SDK runtime"
  runtime     = "claude_code"

  model_id = "claude-3-5-sonnet-20241022"
  llm_config = jsonencode({
    temperature = 0.5
  })

  capabilities = ["advanced_reasoning", "code_generation"]
}

# Data source lookups
data "controlplane_agent" "full_lookup" {
  id = controlplane_agent.full.id
}

data "controlplane_agent" "claude_code_lookup" {
  id = controlplane_agent.claude_code.id
}

# Outputs for test validation
output "full_agent_name" {
  value = controlplane_agent.full.name
}

output "data_full_description" {
  value = data.controlplane_agent.full_lookup.description
}

output "data_full_capabilities" {
  value = data.controlplane_agent.full_lookup.capabilities
}

output "data_claude_code_runtime" {
  value = data.controlplane_agent.claude_code_lookup.runtime
}
