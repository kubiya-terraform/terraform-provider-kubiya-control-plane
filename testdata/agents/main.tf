terraform {
  required_providers {
    kubiya_control_plane = {
      source = "kubiya/control-plane"
    }
  }
}

provider "kubiya_control_plane" {
  # Configuration via environment variables:
  # KUBIYA_CONTROL_PLANE_API_KEY
  # KUBIYA_CONTROL_PLANE_ORG_ID
}

# Test agent resource
resource "kubiya_control_plane_agent" "test" {
  name        = "test-agent"
  description = "Test agent for automated testing"
  model_id    = "gpt-4"
  runtime     = "default"

  llm_config = jsonencode({
    temperature = 0.7
    max_tokens  = 2000
  })
}

# Test data source lookup
data "kubiya_control_plane_agent" "test_lookup" {
  id = kubiya_control_plane_agent.test.id
}

output "agent_id" {
  value = kubiya_control_plane_agent.test.id
}

output "agent_name" {
  value = data.kubiya_control_plane_agent.test_lookup.name
}

output "agent_model_id" {
  value = data.kubiya_control_plane_agent.test_lookup.model_id
}
