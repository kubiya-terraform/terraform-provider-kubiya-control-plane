
provider "controlplane" {}

# Full environment with all optional fields
resource "controlplane_environment" "full" {
  name         = "test-environment-full"
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

# Data source test
data "controlplane_environment" "full_lookup" {
  id = controlplane_environment.full.id
}

# Outputs
output "environment_id" {
  value = controlplane_environment.full.id
}

output "environment_name" {
  value = controlplane_environment.full.name
}

output "environment_display_name" {
  value = controlplane_environment.full.display_name
}

output "environment_description" {
  value = controlplane_environment.full.description
}

output "data_environment_description" {
  value = data.controlplane_environment.full_lookup.description
}

output "data_environment_display_name" {
  value = data.controlplane_environment.full_lookup.display_name
}
