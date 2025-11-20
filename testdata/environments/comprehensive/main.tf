
provider "controlplane" {
  # Configuration via environment variables:
  # KUBIYA_CONTROL_PLANE_API_KEY
  # KUBIYA_CONTROL_PLANE_BASE_URL (optional, defaults to https://control-plane.kubiya.ai)
}

# Test 1: Minimal environment (required fields only)
resource "controlplane_environment" "minimal" {
  name = "test-environment-minimal"
}

# Test 2: Full environment with all optional fields
resource "controlplane_environment" "full" {
  name         = "test-environment-full"
  display_name = "Test Environment Full"
  description  = "Comprehensive test environment with all fields configured"

  # Tags - testing list field
  tags = ["production", "critical", "automated", "testing"]

  # Environment settings - testing complex JSON
  settings = jsonencode({
    region         = "us-east-1"
    max_workers    = 5
    auto_scaling   = true
    retention_days = 30
    features = {
      logging    = true
      monitoring = true
      alerting   = true
    }
    thresholds = {
      cpu_percent    = 80
      memory_percent = 85
      disk_percent   = 90
    }
  })

  # Execution environment - testing complex JSON with all sub-fields
  execution_environment = jsonencode({
    env_vars = {
      LOG_LEVEL     = "debug"
      APP_ENV       = "test"
      FEATURE_FLAGS = "all"
      DEBUG_MODE    = "true"
      API_TIMEOUT   = "30"
    }
    secrets = ["api-key-secret", "db-password", "jwt-secret"]
    integration_ids = ["integration-1", "integration-2", "integration-3"]
  })
}

# Test 3: Environment with display_name only
resource "controlplane_environment" "with_display_name" {
  name         = "test-env-display"
  display_name = "Test Environment Display Name"
  description  = "Environment with display name"
}

# Test 4: Environment with tags only
resource "controlplane_environment" "with_tags" {
  name        = "test-env-tags"
  description = "Environment with tags"
  tags        = ["dev", "low-priority", "testing"]
}

# Test 5: Environment with settings only
resource "controlplane_environment" "with_settings" {
  name        = "test-env-settings"
  description = "Environment with settings"

  settings = jsonencode({
    region      = "us-west-2"
    max_workers = 10
    config = {
      timeout = 600
      retries = 3
    }
  })
}

# Test 6: Environment with execution_environment only
resource "controlplane_environment" "with_exec_env" {
  name        = "test-env-exec-env"
  description = "Environment with execution environment"

  execution_environment = jsonencode({
    env_vars = {
      RUNTIME_ENV = "testing"
      VERBOSE     = "true"
    }
    secrets = ["secret-1"]
    integration_ids = []
  })
}

# Test 7: Environment with env_vars only in execution_environment
resource "controlplane_environment" "env_vars_only" {
  name        = "test-env-vars-only"
  description = "Environment with env vars only"

  execution_environment = jsonencode({
    env_vars = {
      VAR1 = "value1"
      VAR2 = "value2"
    }
  })
}

# Test 8: Environment with secrets only in execution_environment
resource "controlplane_environment" "secrets_only" {
  name        = "test-env-secrets-only"
  description = "Environment with secrets only"

  execution_environment = jsonencode({
    secrets = ["secret-alpha", "secret-beta", "secret-gamma"]
  })
}

# Test 9: Environment with integration_ids only in execution_environment
resource "controlplane_environment" "integrations_only" {
  name        = "test-env-integrations-only"
  description = "Environment with integration IDs only"

  execution_environment = jsonencode({
    integration_ids = ["int-1", "int-2"]
  })
}

# Test 10: Environment with empty optional fields
resource "controlplane_environment" "empty_optionals" {
  name         = "test-env-empty-optionals"
  display_name = ""
  description  = ""
  tags         = []
}

# Test 11: Environment with complex nested settings
resource "controlplane_environment" "complex_settings" {
  name        = "test-env-complex-settings"
  description = "Environment with complex nested settings"

  settings = jsonencode({
    infrastructure = {
      provider = "aws"
      region   = "us-east-1"
      vpc = {
        cidr_block = "10.0.0.0/16"
        subnets = [
          { az = "us-east-1a", cidr = "10.0.1.0/24" },
          { az = "us-east-1b", cidr = "10.0.2.0/24" }
        ]
      }
    }
    compute = {
      instance_type = "t3.medium"
      min_count     = 2
      max_count     = 10
    }
    monitoring = {
      enabled = true
      tools   = ["cloudwatch", "datadog"]
    }
  })
}

# Test 12: Environment for update testing
resource "controlplane_environment" "for_update" {
  name        = "test-env-for-update"
  description = "Initial description"

  settings = jsonencode({
    version = 1
  })
}

# Data source tests
data "controlplane_environment" "minimal_lookup" {
  id = controlplane_environment.minimal.id
}

data "controlplane_environment" "full_lookup" {
  id = controlplane_environment.full.id
}

data "controlplane_environment" "with_tags_lookup" {
  id = controlplane_environment.with_tags.id
}

data "controlplane_environment" "with_exec_env_lookup" {
  id = controlplane_environment.with_exec_env.id
}

# Outputs for test validation
output "minimal_environment_id" {
  value = controlplane_environment.minimal.id
}

output "minimal_environment_name" {
  value = controlplane_environment.minimal.name
}

output "minimal_environment_status" {
  value = controlplane_environment.minimal.status
}

output "minimal_environment_created_at" {
  value = controlplane_environment.minimal.created_at
}

output "full_environment_id" {
  value = controlplane_environment.full.id
}

output "full_environment_name" {
  value = controlplane_environment.full.name
}

output "full_environment_display_name" {
  value = controlplane_environment.full.display_name
}

output "full_environment_description" {
  value = controlplane_environment.full.description
}

output "full_environment_tags" {
  value = controlplane_environment.full.tags
}

output "full_environment_status" {
  value = controlplane_environment.full.status
}

output "full_environment_settings" {
  value     = controlplane_environment.full.settings
  sensitive = true
}

output "full_environment_execution_environment" {
  value     = controlplane_environment.full.execution_environment
  sensitive = true
}

output "full_environment_created_at" {
  value = controlplane_environment.full.created_at
}

output "full_environment_updated_at" {
  value = controlplane_environment.full.updated_at
}

output "with_display_name_environment_id" {
  value = controlplane_environment.with_display_name.id
}

output "with_display_name_environment_display_name" {
  value = controlplane_environment.with_display_name.display_name
}

output "with_tags_environment_id" {
  value = controlplane_environment.with_tags.id
}

output "with_tags_environment_tags" {
  value = controlplane_environment.with_tags.tags
}

output "with_settings_environment_id" {
  value = controlplane_environment.with_settings.id
}

output "with_settings_environment_settings" {
  value     = controlplane_environment.with_settings.settings
  sensitive = true
}

output "with_exec_env_environment_id" {
  value = controlplane_environment.with_exec_env.id
}

output "with_exec_env_environment_execution_environment" {
  value     = controlplane_environment.with_exec_env.execution_environment
  sensitive = true
}

output "env_vars_only_environment_id" {
  value = controlplane_environment.env_vars_only.id
}

output "secrets_only_environment_id" {
  value = controlplane_environment.secrets_only.id
}

output "integrations_only_environment_id" {
  value = controlplane_environment.integrations_only.id
}

output "complex_settings_environment_id" {
  value = controlplane_environment.complex_settings.id
}

output "for_update_environment_id" {
  value = controlplane_environment.for_update.id
}

# Data source outputs for validation
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
  value = data.controlplane_environment.full_lookup.tags
}

output "data_full_display_name" {
  value = data.controlplane_environment.full_lookup.display_name
}

output "data_with_tags_tags" {
  value = data.controlplane_environment.with_tags_lookup.tags
}

output "data_with_exec_env_execution_environment" {
  value     = data.controlplane_environment.with_exec_env_lookup.execution_environment
  sensitive = true
}
