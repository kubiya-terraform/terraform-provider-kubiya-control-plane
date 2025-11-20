
provider "controlplane" {}

# Create policies for project assignment testing
resource "controlplane_policy" "policy1" {
  name        = "test-policy-1-for-project"
  description = "First test policy for project assignment"
  policy_type = "rego"
  enabled     = true
  policy_content = <<-EOT
    package project_test

    default allow = false

    allow {
      input.user.role == "admin"
    }
  EOT
}

# Full project with all optional fields
resource "controlplane_project" "full" {
  name        = "test-project-full"
  key         = "TFULL"
  description = "Comprehensive test project with all fields configured"

  goals                   = "Complete testing of all project fields and configurations"
  status                  = "active"
  visibility              = "org"
  restrict_to_environment = true
  default_model           = "gpt-4"

  settings = jsonencode({
    owner       = "devops-team"
    environment = "test"
    cost_center = "engineering"
    metadata = {
      created_by = "terraform"
      purpose    = "comprehensive_testing"
    }
  })

  policy_ids = [controlplane_policy.policy1.id]
}

# Data source test
data "controlplane_project" "full_lookup" {
  id = controlplane_project.full.id
}

# Outputs
output "project_id" {
  value = controlplane_project.full.id
}

output "project_name" {
  value = controlplane_project.full.name
}

output "project_key" {
  value = controlplane_project.full.key
}

output "project_description" {
  value = controlplane_project.full.description
}

output "project_visibility" {
  value = controlplane_project.full.visibility
}

output "project_default_model" {
  value = controlplane_project.full.default_model
}

output "project_status" {
  value = controlplane_project.full.status
}

output "data_project_description" {
  value = data.controlplane_project.full_lookup.description
}

output "data_project_visibility" {
  value = data.controlplane_project.full_lookup.visibility
}
