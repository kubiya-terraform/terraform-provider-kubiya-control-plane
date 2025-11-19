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

variable "description" {
  type    = string
  default = "Original description"
}

variable "model_id" {
  type    = string
  default = "gpt-4"
}

# Agent with multiple fields that will be updated
resource "controlplane_agent" "update_test" {
  name        = "test-agent-update-multiple"
  description = var.description
  model_id    = var.model_id
  runtime     = "default"
}

# Outputs for test validation
output "agent_id" {
  value = controlplane_agent.update_test.id
}

output "agent_name" {
  value = controlplane_agent.update_test.name
}

output "agent_description" {
  value = controlplane_agent.update_test.description
}

output "agent_model_id" {
  value = controlplane_agent.update_test.model_id
}

output "agent_runtime" {
  value = controlplane_agent.update_test.runtime
}

output "agent_updated_at" {
  value = controlplane_agent.update_test.updated_at
}
