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
  # KUBIYA_CONTROL_PLANE_BASE_URL (optional, defaults to https://control-plane.kubiya.ai)
}

# Full agent with all optional fields
resource "controlplane_agent" "full" {
  name        = "test-agent-full"
  description = "Comprehensive test agent with all fields configured"

  # Status - testing optional status field
  status = "idle"

  # LLM configuration
  model_id = "gpt-4"
  llm_config = jsonencode({
    temperature = 0.7
    max_tokens  = 2000
    top_p       = 0.9
  })

  # Runtime configuration
  runtime = "default"

  # Agent capabilities - testing list field
  capabilities = ["code_execution", "file_operations", "web_search", "data_analysis"]

  # Agent configuration - testing complex JSON
  configuration = jsonencode({
    max_retries    = 3
    timeout        = 300
    retry_delay    = 5
    enable_logging = true
    settings = {
      verbose = true
      debug   = false
    }
  })
}

# Data source test
data "controlplane_agent" "full_lookup" {
  id = controlplane_agent.full.id
}

# Outputs for test validation
output "agent_id" {
  value = controlplane_agent.full.id
}

output "agent_name" {
  value = controlplane_agent.full.name
}

output "agent_description" {
  value = controlplane_agent.full.description
}

output "agent_status" {
  value = controlplane_agent.full.status
}

output "agent_model_id" {
  value = controlplane_agent.full.model_id
}

output "agent_llm_config" {
  value = controlplane_agent.full.llm_config
}

output "agent_runtime" {
  value = controlplane_agent.full.runtime
}

output "agent_capabilities" {
  value = controlplane_agent.full.capabilities
}

output "agent_configuration" {
  value = controlplane_agent.full.configuration
}

output "agent_created_at" {
  value = controlplane_agent.full.created_at
}

output "agent_updated_at" {
  value = controlplane_agent.full.updated_at
}

output "data_agent_description" {
  value = data.controlplane_agent.full_lookup.description
}

output "data_agent_capabilities" {
  value = data.controlplane_agent.full_lookup.capabilities
}
