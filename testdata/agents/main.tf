terraform {
  required_providers {
    controlplane = {
      source = "kubiya/control-plane"
    }
  }
}

provider "controlplane" {
  # Configuration via environment variables:
  # KUBIYA_CONTROL_PLANE_API_KEY
  # KUBIYA_CONTROL_PLANE_ORG_ID
  # KUBIYA_CONTROL_PLANE_BASE_URL (optional, defaults to https://control-plane.kubiya.ai)
}

# Test agent resource
resource "controlplane_agent" "test" {
  name        = "test-agent"
  description = "Test agent for automated testing"

  # LLM configuration
  model_id = "gpt-4"
  llm_config = jsonencode({
    temperature = 0.7
    max_tokens  = 2000
  })

  # Runtime configuration
  runtime = "default"

  # Agent capabilities
  capabilities = ["code_execution", "file_operations"]

  # Agent configuration
  configuration = jsonencode({
    max_retries = 3
    timeout     = 300
  })
}

# Test data source lookup
data "controlplane_agent" "test_lookup" {
  id = controlplane_agent.test.id
}

output "agent_id" {
  value = controlplane_agent.test.id
}

output "agent_name" {
  value = data.controlplane_agent.test_lookup.name
}

output "agent_model_id" {
  value = data.controlplane_agent.test_lookup.model_id
}

output "agent_status" {
  value = data.controlplane_agent.test_lookup.status
}
