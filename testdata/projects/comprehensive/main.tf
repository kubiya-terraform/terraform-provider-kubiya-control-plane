terraform {
  required_providers {
    controlplane = {
      source = "kubiya/control-plane"
    }
  }
}

provider "controlplane" {
  # Configuration via environment variables:
  # KUBIYA_CONTROL_PLANE_API_KEY
  # KUBIYA_CONTROL_PLANE_BASE_URL (optional, defaults to https://control-plane.kubiya.ai)
}

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

resource "controlplane_policy" "policy2" {
  name        = "test-policy-2-for-project"
  description = "Second test policy for project assignment"
  policy_type = "rego"
  enabled     = true
  policy_content = <<-EOT
    package project_test_2

    default allow = true
  EOT
}

# Test 1: Minimal project (required fields only)
resource "controlplane_project" "minimal" {
  name = "test-project-minimal"
  key  = "TMIN"
}

# Test 2: Full project with all optional fields
resource "controlplane_project" "full" {
  name        = "test-project-full"
  key         = "TFULL"
  description = "Comprehensive test project with all fields configured"

  # Project goals
  goals = "Complete testing of all project fields and configurations"

  # Status
  status = "active"

  # Visibility
  visibility = "org"

  # Restrict to environment
  restrict_to_environment = true

  # Default model
  default_model = "gpt-4"

  # Project settings - testing complex JSON
  settings = jsonencode({
    owner       = "devops-team"
    environment = "test"
    cost_center = "engineering"
    metadata = {
      created_by = "terraform"
      purpose    = "comprehensive_testing"
    }
    features = {
      auto_deploy       = true
      require_approval  = false
      enable_monitoring = true
    }
  })

  # Policy IDs - testing list field
  policy_ids = [
    controlplane_policy.policy1.id,
    controlplane_policy.policy2.id
  ]
}

# Test 3: Project with private visibility (default)
resource "controlplane_project" "private" {
  name        = "test-project-private"
  key         = "TPRIV"
  description = "Project with private visibility"
  visibility  = "private"
}

# Test 4: Project with archived status
resource "controlplane_project" "archived" {
  name        = "test-project-archived"
  key         = "TARCH"
  description = "Project with archived status"
  status      = "archived"
}

# Test 5: Project with paused status
resource "controlplane_project" "paused" {
  name        = "test-project-paused"
  key         = "TPAUS"
  description = "Project with paused status"
  status      = "paused"
}

# Test 6: Project with single policy
resource "controlplane_policy" "single_policy" {
  name        = "test-policy-single"
  description = "Single policy for project"
  policy_type = "rego"
  enabled     = true
  policy_content = <<-EOT
    package single_project_policy

    default allow = false
  EOT
}

resource "controlplane_project" "single_policy" {
  name       = "test-project-single-policy"
  key        = "TSPOL"
  policy_ids = [controlplane_policy.single_policy.id]
}

# Test 7: Project without environment restriction
resource "controlplane_project" "no_env_restriction" {
  name                    = "test-project-no-env-restriction"
  key                     = "TNOER"
  description             = "Project without environment restriction"
  restrict_to_environment = false
}

# Test 8: Project with custom default model
resource "controlplane_project" "custom_model" {
  name          = "test-project-custom-model"
  key           = "TCMOD"
  description   = "Project with custom default model"
  default_model = "claude-3-5-sonnet-20241022"
}

# Test 9: Project with complex settings
resource "controlplane_project" "complex_settings" {
  name        = "test-project-complex-settings"
  key         = "TCSET"
  description = "Project with complex settings configuration"

  settings = jsonencode({
    workflows = {
      ci_cd = {
        enabled = true
        stages  = ["build", "test", "deploy"]
      }
      notifications = {
        slack   = true
        email   = true
        webhook = "https://example.com/webhook"
      }
    }
    integrations = {
      github = {
        enabled = true
        org     = "kubiya"
      }
      jira = {
        enabled    = true
        project_id = "JIRA-123"
      }
    }
  })
}

# Test 10: Project with empty optional fields
resource "controlplane_project" "empty_optionals" {
  name        = "test-project-empty-optionals"
  key         = "TEMPT"
  description = ""
  goals       = ""
  policy_ids  = []
}

# Test 11: Project for update testing
resource "controlplane_project" "for_update" {
  name        = "test-project-for-update"
  key         = "TUPDT"
  description = "Initial description"
  status      = "active"
  visibility  = "private"

  settings = jsonencode({
    version = 1
  })
}

# Data source tests
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

# Outputs for test validation
output "minimal_project_id" {
  value = controlplane_project.minimal.id
}

output "minimal_project_name" {
  value = controlplane_project.minimal.name
}

output "minimal_project_key" {
  value = controlplane_project.minimal.key
}

output "minimal_project_created_at" {
  value = controlplane_project.minimal.created_at
}

output "full_project_id" {
  value = controlplane_project.full.id
}

output "full_project_name" {
  value = controlplane_project.full.name
}

output "full_project_key" {
  value = controlplane_project.full.key
}

output "full_project_description" {
  value = controlplane_project.full.description
}

output "full_project_goals" {
  value = controlplane_project.full.goals
}

output "full_project_status" {
  value = controlplane_project.full.status
}

output "full_project_visibility" {
  value = controlplane_project.full.visibility
}

output "full_project_restrict_to_environment" {
  value = controlplane_project.full.restrict_to_environment
}

output "full_project_default_model" {
  value = controlplane_project.full.default_model
}

output "full_project_settings" {
  value     = controlplane_project.full.settings
  sensitive = true
}

output "full_project_policy_ids" {
  value = controlplane_project.full.policy_ids
}

output "full_project_created_at" {
  value = controlplane_project.full.created_at
}

output "full_project_updated_at" {
  value = controlplane_project.full.updated_at
}

output "private_project_id" {
  value = controlplane_project.private.id
}

output "private_project_visibility" {
  value = controlplane_project.private.visibility
}

output "archived_project_id" {
  value = controlplane_project.archived.id
}

output "archived_project_status" {
  value = controlplane_project.archived.status
}

output "paused_project_id" {
  value = controlplane_project.paused.id
}

output "paused_project_status" {
  value = controlplane_project.paused.status
}

output "single_policy_project_id" {
  value = controlplane_project.single_policy.id
}

output "single_policy_project_policy_ids" {
  value = controlplane_project.single_policy.policy_ids
}

output "no_env_restriction_project_id" {
  value = controlplane_project.no_env_restriction.id
}

output "no_env_restriction_project_restrict_to_environment" {
  value = controlplane_project.no_env_restriction.restrict_to_environment
}

output "custom_model_project_id" {
  value = controlplane_project.custom_model.id
}

output "custom_model_project_default_model" {
  value = controlplane_project.custom_model.default_model
}

output "complex_settings_project_id" {
  value = controlplane_project.complex_settings.id
}

output "complex_settings_project_settings" {
  value     = controlplane_project.complex_settings.settings
  sensitive = true
}

output "for_update_project_id" {
  value = controlplane_project.for_update.id
}

# Data source outputs for validation
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
  value = data.controlplane_project.full_lookup.policy_ids
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
