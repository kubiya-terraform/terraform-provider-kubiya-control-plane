
provider "controlplane" {}

# Minimal environment (required fields only)
resource "controlplane_environment" "minimal" {
  name = "test-environment-minimal"
}

# Data source test
data "controlplane_environment" "minimal_lookup" {
  id = controlplane_environment.minimal.id
}

# Outputs
output "environment_id" {
  value = controlplane_environment.minimal.id
}

output "environment_name" {
  value = controlplane_environment.minimal.name
}

output "environment_status" {
  value = controlplane_environment.minimal.status
}

output "data_environment_name" {
  value = data.controlplane_environment.minimal_lookup.name
}

output "data_environment_status" {
  value = data.controlplane_environment.minimal_lookup.status
}
