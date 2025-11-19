terraform {
  required_providers {
    controlplane = {
      source = "kubiya/control-plane"
    }
  }
}

provider "controlplane" {}

# Minimal project (required fields only)
resource "controlplane_project" "minimal" {
  name = "test-project-minimal"
  key  = "TMIN"
}

# Data source test
data "controlplane_project" "minimal_lookup" {
  id = controlplane_project.minimal.id
}

# Outputs
output "project_id" {
  value = controlplane_project.minimal.id
}

output "project_name" {
  value = controlplane_project.minimal.name
}

output "project_key" {
  value = controlplane_project.minimal.key
}

output "data_project_name" {
  value = data.controlplane_project.minimal_lookup.name
}

output "data_project_key" {
  value = data.controlplane_project.minimal_lookup.key
}
