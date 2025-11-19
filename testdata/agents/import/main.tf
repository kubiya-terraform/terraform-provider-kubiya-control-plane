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

variable "agent_name" {
  type        = string
  description = "Agent name for import configuration"
}

# Agent configuration for import
resource "controlplane_agent" "imported" {
  name = var.agent_name
}

# Outputs for test validation
output "imported_agent_id" {
  value = controlplane_agent.imported.id
}

output "imported_agent_name" {
  value = controlplane_agent.imported.name
}

output "imported_agent_created_at" {
  value = controlplane_agent.imported.created_at
}
