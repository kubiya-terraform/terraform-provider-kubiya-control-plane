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

# ============================================================================
# Step 1: Create foundational resources (Environment, Project)
# ============================================================================

# Create a production environment
resource "controlplane_environment" "production" {
  name        = "production"
  description = "Production environment for agents"

  configuration = jsonencode({
    region         = "us-east-1"
    max_workers    = 10
    auto_scaling   = true
    retention_days = 90
  })

  execution_environment = jsonencode({
    env_vars = {
      LOG_LEVEL = "info"
      APP_ENV   = "production"
    }
  })
}

# Create a project
resource "controlplane_project" "platform" {
  name        = "platform-automation"
  description = "Platform automation and operations project"

  metadata = jsonencode({
    owner       = "platform-team"
    cost_center = "engineering"
  })
}

# ============================================================================
# Step 2: Create skills
# ============================================================================

# Create a shell skill
resource "controlplane_skill" "shell_ops" {
  name        = "shell-operations"
  description = "Shell command execution for operations"
  type        = "shell"
  enabled     = true

  configuration = jsonencode({
    allowed_commands = ["kubectl", "helm", "aws", "gcloud"]
    timeout          = 300
  })
}

# Create a file system skill
resource "controlplane_skill" "filesystem" {
  name        = "filesystem-access"
  description = "File system operations"
  type        = "file_system"
  enabled     = true

  configuration = jsonencode({
    allowed_paths = ["/app/configs", "/app/data"]
    max_file_size = 10485760
  })
}

# ============================================================================
# Step 3: Create policies
# ============================================================================

# Create a security policy
resource "controlplane_policy" "security" {
  name        = "production-security"
  description = "Security policy for production"
  enabled     = true

  policy_content = <<-EOT
    package kubiya.security

    # Deny destructive operations without approval
    deny[msg] {
      input.operation = "delete"
      count(input.approvals) < 2
      msg := "Delete operations require at least 2 approvals"
    }

    # Require MFA for sensitive operations
    deny[msg] {
      input.operation = "deploy"
      input.environment = "production"
      not input.mfa_verified
      msg := "Production deployments require MFA verification"
    }
  EOT

  tags = ["security", "production"]
}

# ============================================================================
# Step 4: Create a team
# ============================================================================

# Create a DevOps team
resource "controlplane_team" "devops" {
  name        = "devops-team"
  description = "DevOps and platform engineering team"

  configuration = jsonencode({
    max_agents      = 5
    default_runtime = "default"
    slack_webhook   = "https://hooks.slack.com/services/YOUR/WEBHOOK/URL"
  })

  capabilities = ["deployment", "monitoring", "incident_response"]
}

# ============================================================================
# Step 5: Create agents
# ============================================================================

# Create a deployment agent
resource "controlplane_agent" "deployer" {
  name        = "production-deployer"
  description = "Agent for production deployments"

  model_id = "gpt-4"
  llm_config = jsonencode({
    temperature = 0.3
    max_tokens  = 4000
  })

  runtime = "default"

  capabilities = ["kubernetes_deploy", "helm_deploy", "rollback"]

  configuration = jsonencode({
    max_retries     = 3
    timeout         = 600
    approval_needed = true
  })

  team_id = controlplane_team.devops.id
}

# Create a monitoring agent
resource "controlplane_agent" "monitor" {
  name        = "production-monitor"
  description = "Agent for monitoring and alerting"

  model_id = "gpt-4"
  llm_config = jsonencode({
    temperature = 0.5
    max_tokens  = 2000
  })

  runtime = "default"

  capabilities = ["metrics_collection", "alerting", "log_analysis"]

  configuration = jsonencode({
    check_interval = 60
    alert_channels = ["slack", "pagerduty"]
  })

  team_id = controlplane_team.devops.id
}

# ============================================================================
# Step 6: Register workers (optional - typically done at runtime)
# ============================================================================

resource "controlplane_worker" "worker_01" {
  environment_name = controlplane_environment.production.name
  hostname         = "worker-prod-01"

  metadata = jsonencode({
    region = "us-east-1"
    zone   = "us-east-1a"
  })
}

# ============================================================================
# Outputs
# ============================================================================

output "environment_id" {
  value       = controlplane_environment.production.id
  description = "Production environment ID"
}

output "project_id" {
  value       = controlplane_project.platform.id
  description = "Platform project ID"
}

output "team_id" {
  value       = controlplane_team.devops.id
  description = "DevOps team ID"
}

output "deployer_agent_id" {
  value       = controlplane_agent.deployer.id
  description = "Deployer agent ID"
}

output "monitor_agent_id" {
  value       = controlplane_agent.monitor.id
  description = "Monitor agent ID"
}

output "shell_skill_id" {
  value       = controlplane_skill.shell_ops.id
  description = "Shell operations skill ID"
}

output "security_policy_id" {
  value       = controlplane_policy.security.id
  description = "Security policy ID"
}

output "worker_id" {
  value       = controlplane_worker.worker_01.id
  description = "Worker ID"
}
