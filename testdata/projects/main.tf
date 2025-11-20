terraform {
  required_providers {
    controlplane = {
      source = "kubiya/control-plane"
    }
  }
}

provider "controlplane" {}

# Create policy for project assignment testing
resource "controlplane_policy" "policy1" {
  name           = "test-policy-1-for-project-ds"
  description    = "Test policy for project assignment"
  policy_type    = "rego"
  enabled        = true
  policy_content = <<-EOT
    package project_test

    default allow = false

    allow {
      input.user.role == "admin"
    }
  EOT
}

# Minimal project (required fields only)
resource "controlplane_project" "minimal" {
  name = "test-project-minimal"
  key  = "TMIN"
}

# Full project with all optional fields
resource "controlplane_project" "full" {
  name        = "test-project-full-ds"
  key         = "TFDS"
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

# Archived project for status testing
resource "controlplane_project" "archived" {
  name        = "test-project-archived-ds"
  key         = "TARCH"
  description = "Archived project for status testing"
  status      = "archived"
}

# Project with custom model
resource "controlplane_project" "custom_model" {
  name          = "test-project-custom-model-ds"
  key           = "TCUST"
  description   = "Project with custom model"
  default_model = "claude-3-5-sonnet-20241022"
}

# Data sources
data "controlplane_project" "minimal_lookup" {
  id = controlplane_project.minimal.id
}

data "controlplane_project" "full_lookup" {
  id = controlplane_project.full.id
}

data "controlplane_project" "archived_lookup" {
  id = controlplane_project.archived.id
}

data "controlplane_project" "custom_model_lookup" {
  id = controlplane_project.custom_model.id
}

# Outputs for tests
output "data_minimal_name" {
  value = data.controlplane_project.minimal_lookup.name
}

output "data_minimal_key" {
  value = data.controlplane_project.minimal_lookup.key
}

output "data_full_description" {
  value = data.controlplane_project.full_lookup.description
}

output "data_full_policy_ids" {
  value = jsonencode(data.controlplane_project.full_lookup.policy_ids)
}

output "data_full_visibility" {
  value = data.controlplane_project.full_lookup.visibility
}

output "data_archived_status" {
  value = data.controlplane_project.archived_lookup.status
}

output "data_custom_model_default_model" {
  value = data.controlplane_project.custom_model_lookup.default_model
}
