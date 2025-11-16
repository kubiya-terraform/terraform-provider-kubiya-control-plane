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
  # KUBIYA_CONTROL_PLANE_BASE_URL (optional, defaults to https://control-plane.kubiya.ai)
}

# Create a security policy for production environments
resource "controlplane_policy" "production_security" {
  name        = "production-security-policy"
  description = "Security policy for production environments"
  enabled     = true

  # OPA Rego policy content
  policy_content = <<-EOT
    package kubiya.security

    # Deny operations during business hours
    deny[msg] {
      input.operation = "deploy"
      input.environment = "production"
      is_business_hours
      msg := "Production deployments are not allowed during business hours"
    }

    # Check for required approvals
    deny[msg] {
      input.operation = "delete"
      count(input.approvals) < 2
      msg := "Delete operations require at least 2 approvals"
    }

    # Helper function to check business hours
    is_business_hours {
      hour := time.clock(time.now_ns())[0]
      hour >= 9
      hour < 17
    }
  EOT

  tags = ["security", "production", "compliance"]
}

# Create a cost control policy
resource "controlplane_policy" "cost_control" {
  name        = "cost-control-policy"
  description = "Policy to control infrastructure costs"
  enabled     = true

  policy_content = <<-EOT
    package kubiya.cost

    # Limit instance sizes
    deny[msg] {
      input.action = "create_instance"
      input.instance_type = "x2.32xlarge"
      msg := "Instance type too large, maximum allowed is m5.2xlarge"
    }

    # Require cost tags
    deny[msg] {
      input.action = "create_resource"
      not input.tags.cost_center
      msg := "All resources must have a cost_center tag"
    }
  EOT

  tags = ["cost", "governance"]
}

# Create a compliance policy
resource "controlplane_policy" "compliance" {
  name        = "compliance-policy"
  description = "Policy for regulatory compliance"
  enabled     = true

  policy_content = <<-EOT
    package kubiya.compliance

    # Require encryption
    deny[msg] {
      input.resource_type = "database"
      not input.encrypted
      msg := "All databases must have encryption enabled"
    }

    # Require data residency compliance
    deny[msg] {
      input.region = "us-west-2"
      input.data_classification = "sensitive"
      msg := "Sensitive data must be stored in approved regions only"
    }

    # Audit log requirement
    deny[msg] {
      input.resource_type = "api_gateway"
      not input.logging_enabled
      msg := "API gateways must have logging enabled for compliance"
    }
  EOT

  tags = ["compliance", "audit", "security"]
}

# Look up an existing policy by ID
data "controlplane_policy" "existing" {
  id = "2db93293-2c30-4557-a08b-db6f80f9ef57"
}

# Output policy information
output "production_security_policy_id" {
  value       = controlplane_policy.production_security.id
  description = "The ID of the production security policy"
}

output "cost_control_policy_id" {
  value       = controlplane_policy.cost_control.id
  description = "The ID of the cost control policy"
}

output "compliance_policy_id" {
  value       = controlplane_policy.compliance.id
  description = "The ID of the compliance policy"
}

output "existing_policy_name" {
  value       = data.controlplane_policy.existing.name
  description = "Name of the existing policy"
}

output "existing_policy_enabled" {
  value       = data.controlplane_policy.existing.enabled
  description = "Whether the existing policy is enabled"
}
