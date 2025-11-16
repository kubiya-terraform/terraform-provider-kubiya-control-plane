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
  # KUBIYA_CONTROL_PLANE_BASE_URL (optional, defaults to https://control-plane.kubiya.ai)
}

# Create a team
resource "controlplane_team" "example" {
  name        = "test-team"
  description = "A test team for demonstration"

  # Runtime type: "default" (Agno) or "claude_code" (Claude Code SDK)
  runtime = "default"

  # Team configuration
  configuration = jsonencode({
    max_agents        = 10
    enable_monitoring = true
  })
}

# Look up an existing team by ID
data "controlplane_team" "existing" {
  id = "71b17bbb-4143-4599-933a-d55b20994732"
}

# Output team information
output "team_id" {
  value       = controlplane_team.example.id
  description = "The ID of the created team"
}

output "team_status" {
  value       = controlplane_team.example.status
  description = "The current status of the team"
}

output "existing_team_name" {
  value       = data.controlplane_team.existing.name
  description = "Name of the existing team"
}

output "existing_team_config" {
  value       = data.controlplane_team.existing.configuration
  description = "Configuration of the existing team"
  sensitive   = true
}
