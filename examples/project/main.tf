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

# Create a project
resource "controlplane_project" "example" {
  name        = "updated-example-project"
  key         = "PROJ"
  description = "An example project for demonstration"

  # Project settings
  settings = jsonencode({
    owner       = "devops-team"
    environment = "production"
    cost_center = "engineering"
  })
}

# Look up an existing project by ID
data "controlplane_project" "existing" {
  id = "65a206de-f17d-4e65-93bb-7fed5c72d567"
}

# Output project information
output "project_id" {
  value       = controlplane_project.example.id
  description = "The ID of the created project"
}

output "project_status" {
  value       = controlplane_project.example.status
  description = "The current status of the project"
}

output "existing_project_name" {
  value       = data.controlplane_project.existing.name
  description = "Name of the existing project"
}

output "existing_project_settings" {
  value       = data.controlplane_project.existing.settings
  description = "Settings of the existing project"
}
