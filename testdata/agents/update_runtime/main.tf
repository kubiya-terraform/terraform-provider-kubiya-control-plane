
provider "controlplane" {
  # Configuration via environment variables
}

variable "runtime" {
  type    = string
  default = "default"
}

# Agent with runtime that can be updated
resource "controlplane_agent" "runtime_test" {
  name    = "test-agent-update-runtime"
  runtime = var.runtime
}

# Outputs for test validation
output "agent_id" {
  value = controlplane_agent.runtime_test.id
}

output "agent_name" {
  value = controlplane_agent.runtime_test.name
}

output "agent_runtime" {
  value = controlplane_agent.runtime_test.runtime
}

output "agent_updated_at" {
  value = controlplane_agent.runtime_test.updated_at
}
