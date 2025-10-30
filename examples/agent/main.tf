terraform {
  required_providers {
    controlplane = {
      source = "kubiya/control-plane"
    }
  }
}

provider "controlplane" {
  # Configuration is via environment variables:
  # KUBIYA_CONTROL_PLANE_API_KEY
  # KUBIYA_CONTROL_PLANE_ORG_ID
  # KUBIYA_CONTROL_PLANE_BASE_URL (optional, defaults to https://control-plane.kubiya.ai)
}

# Create an agent
resource "controlplane_agent" "example" {
  name        = "test-agent"
  description = "A test agent for demonstration"

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

  # Optional: Associate with a team
  # team_id = controlplane_team.example.id
}

# Look up an existing agent by ID
data "controlplane_agent" "existing" {
  id = "b21c26b6-b31d-4dd0-b85c-e45b50eeaad2"
}

# Output agent information
output "agent_id" {
  value       = controlplane_agent.example.id
  description = "The ID of the created agent"
}

output "agent_status" {
  value       = controlplane_agent.example.status
  description = "The current status of the agent"
}

output "existing_agent_name" {
  value       = data.controlplane_agent.existing.name
  description = "Name of the existing agent"
}
