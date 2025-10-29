terraform {
  required_providers {
    kubiya_control_plane = {
      source = "kubiya/control-plane"
    }
  }
}

provider "kubiya_control_plane" {
  # Configuration is via environment variables:
  # KUBIYA_CONTROL_PLANE_API_KEY
  # KUBIYA_CONTROL_PLANE_ORG_ID
  # KUBIYA_CONTROL_PLANE_BASE_URL (optional, defaults to https://control-plane.kubiya.ai)
}

# Create a team
resource "kubiya_control_plane_team" "example" {
  name        = "example-team"
  description = "An example team for demonstration"

  # Team configuration
  configuration = jsonencode({
    max_agents        = 10
    default_runtime   = "default"
    enable_monitoring = true
  })

  # Optional: Assign capabilities to the team
  capabilities = ["deployment", "monitoring", "reporting"]
}

# Look up an existing team by ID
data "kubiya_control_plane_team" "existing" {
  id = "team-uuid-here"
}

# Output team information
output "team_id" {
  value       = kubiya_control_plane_team.example.id
  description = "The ID of the created team"
}

output "team_status" {
  value       = kubiya_control_plane_team.example.status
  description = "The current status of the team"
}

output "existing_team_name" {
  value       = data.kubiya_control_plane_team.existing.name
  description = "Name of the existing team"
}

output "existing_team_config" {
  value       = data.kubiya_control_plane_team.existing.configuration
  description = "Configuration of the existing team"
  sensitive   = true
}
