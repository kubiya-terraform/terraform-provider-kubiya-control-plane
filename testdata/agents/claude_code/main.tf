
provider "controlplane" {
  # Configuration via environment variables:
  # KUBIYA_CONTROL_PLANE_API_KEY
  # KUBIYA_CONTROL_PLANE_BASE_URL (optional, defaults to https://control-plane.kubiya.ai)
}

# Agent with claude_code runtime
resource "controlplane_agent" "claude_code" {
  name        = "test-agent-claude-code"
  description = "Agent using Claude Code SDK runtime"
  runtime     = "claude_code"

  model_id = "claude-3-5-sonnet-20241022"
  llm_config = jsonencode({
    temperature = 0.5
  })

  capabilities = ["advanced_reasoning", "code_generation"]
}

# Data source test
data "controlplane_agent" "claude_code_lookup" {
  id = controlplane_agent.claude_code.id
}

# Outputs for test validation
output "agent_id" {
  value = controlplane_agent.claude_code.id
}

output "agent_runtime" {
  value = controlplane_agent.claude_code.runtime
}

output "agent_model_id" {
  value = controlplane_agent.claude_code.model_id
}

output "data_agent_runtime" {
  value = data.controlplane_agent.claude_code_lookup.runtime
}
