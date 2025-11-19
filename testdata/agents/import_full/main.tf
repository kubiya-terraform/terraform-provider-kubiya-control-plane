terraform {
  required_providers {
    controlplane = {
      source = "kubiya/control-plane"
    }
  }
}

provider "controlplane" {
  # Configuration via environment variables
}

variable "agent_id" {
  type        = string
  description = "Agent ID to import"
}

# Fully configured agent for import
resource "controlplane_agent" "imported_full" {
  name        = "test-agent-full"
  description = "Comprehensive test agent with all fields configured"
  status      = "idle"
  model_id    = "gpt-4"
  llm_config = jsonencode({
    temperature = 0.7
    max_tokens  = 2000
    top_p       = 0.9
  })
  runtime = "default"
  capabilities = ["code_execution", "file_operations", "web_search", "data_analysis"]
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

# Outputs for test validation
output "imported_agent_id" {
  value = controlplane_agent.imported_full.id
}

output "imported_agent_name" {
  value = controlplane_agent.imported_full.name
}

output "imported_agent_description" {
  value = controlplane_agent.imported_full.description
}

output "imported_agent_model_id" {
  value = controlplane_agent.imported_full.model_id
}

output "imported_agent_runtime" {
  value = controlplane_agent.imported_full.runtime
}

output "imported_agent_created_at" {
  value = controlplane_agent.imported_full.created_at
}
