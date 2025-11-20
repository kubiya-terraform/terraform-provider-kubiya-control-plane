
provider "controlplane" {}

# Minimal environment (required fields only)
resource "controlplane_environment" "minimal" {
  name = "test-environment-minimal"
}

# Full environment with all optional fields
resource "controlplane_environment" "full" {
  name         = "test-environment-full-ds"
  display_name = "Test Environment Full"
  description  = "Comprehensive test environment with all fields configured"

  settings = jsonencode({
    region           = "us-east-1"
    instance_type    = "medium"
    auto_scaling     = true
    max_instances    = 10
    monitoring_level = "detailed"
  })

  tags = ["test", "comprehensive", "full-config"]
}

# Data sources
data "controlplane_environment" "minimal_lookup" {
  id = controlplane_environment.minimal.id
}

data "controlplane_environment" "full_lookup" {
  id = controlplane_environment.full.id
}

# Outputs for tests
output "data_minimal_name" {
  value = data.controlplane_environment.minimal_lookup.name
}

output "data_minimal_status" {
  value = data.controlplane_environment.minimal_lookup.status
}

output "data_full_description" {
  value = data.controlplane_environment.full_lookup.description
}

output "data_full_tags" {
  value = jsonencode(data.controlplane_environment.full_lookup.tags)
}

output "data_full_display_name" {
  value = data.controlplane_environment.full_lookup.display_name
}
