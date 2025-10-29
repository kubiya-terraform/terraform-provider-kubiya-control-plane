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

# Create a project
resource "kubiya_control_plane_project" "example" {
  name        = "example-project"
  description = "An example project for demonstration"

  # Project metadata
  metadata = jsonencode({
    owner       = "devops-team"
    environment = "production"
    cost_center = "engineering"
  })
}

# Look up an existing project by ID
data "kubiya_control_plane_project" "existing" {
  id = "project-uuid-here"
}

# Output project information
output "project_id" {
  value       = kubiya_control_plane_project.example.id
  description = "The ID of the created project"
}

output "project_status" {
  value       = kubiya_control_plane_project.example.status
  description = "The current status of the project"
}

output "existing_project_name" {
  value       = data.kubiya_control_plane_project.existing.name
  description = "Name of the existing project"
}

output "existing_project_metadata" {
  value       = data.kubiya_control_plane_project.existing.metadata
  description = "Metadata of the existing project"
}
