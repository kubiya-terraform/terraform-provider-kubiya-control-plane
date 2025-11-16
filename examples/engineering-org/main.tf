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

# ============================================================================
# Example 1: Using all defaults
# ============================================================================

# Commented out to avoid conflicts with existing resources
# Uncomment if you want to test the defaults module separately
# module "engineering_org_defaults" {
#   source = "../../modules/engineering-org"
#
#   # This will create:
#   # - 1 production environment
#   # - 1 platform project
#   # - 1 devops team
#   # - 2 agents (deployer, monitor)
#   # - 2 skills (shell, filesystem)
#   # - 1 security policy
#   # - 1 worker queue
#   # - 1 health check job
# }

# ============================================================================
# Example 2: Custom configuration with multiple items (with unique names)
# ============================================================================

module "engineering_org_custom" {
  source = "../../modules/engineering-org"

  # Multiple environments
  environments = {
    prod_example = {
      description = "Production environment (example)"
      settings = jsonencode({
        region         = "us-east-1"
        max_workers    = 20
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
    staging_example = {
      description = "Staging environment (example)"
      settings = jsonencode({
        region         = "us-west-2"
        max_workers    = 10
        auto_scaling   = true
        retention_days = 30
      })
      execution_environment = jsonencode({
        env_vars = {
          LOG_LEVEL = "debug"
          APP_ENV   = "staging"
        }
      })
    }
    dev_example = {
      description = "Development environment (example)"
      settings = jsonencode({
        region         = "us-west-2"
        max_workers    = 5
        auto_scaling   = false
        retention_days = 7
      })
      execution_environment = jsonencode({
        env_vars = {
          LOG_LEVEL = "debug"
          APP_ENV   = "development"
        }
      })
    }
  }

  # Multiple projects
  projects = {
    platform_example = {
      key         = "PLATEX"
      description = "Platform engineering project (example)"
      settings = jsonencode({
        owner       = "platform-team"
        cost_center = "engineering"
      })
    }
    data_example = {
      key         = "DATAEX"
      description = "Data engineering project (example)"
      settings = jsonencode({
        owner       = "data-team"
        cost_center = "data"
      })
    }
    security_example = {
      key         = "SECEX"
      description = "Security and compliance project (example)"
      settings = jsonencode({
        owner       = "security-team"
        cost_center = "security"
      })
    }
  }

  # Multiple teams
  teams = {
    devops_example = {
      description = "DevOps and platform engineering team (example)"
      runtime     = "default"
      configuration = jsonencode({
        max_agents        = 15
        enable_monitoring = true
      })
    }
    sre_example = {
      description = "Site reliability engineering team (example)"
      runtime     = "default"
      configuration = jsonencode({
        max_agents        = 10
        enable_monitoring = true
      })
    }
    data_eng_example = {
      description = "Data engineering team (example)"
      runtime     = "default"
      configuration = jsonencode({
        max_agents = 8
      })
    }
  }

  # Multiple skills
  skills = {
    shell_example = {
      description = "Shell command execution (example)"
      type        = "shell"
      enabled     = true
      configuration = jsonencode({
        allowed_commands = ["kubectl", "helm", "aws", "terraform", "ansible"]
        timeout          = 600
        working_dir      = "/app"
      })
    }
    filesystem_example = {
      description = "File system operations (example)"
      type        = "file_system"
      enabled     = true
      configuration = jsonencode({
        allowed_paths = ["/app/configs", "/app/data", "/tmp"]
        max_file_size = 52428800 # 50MB
        operations    = ["read", "write", "list", "delete"]
      })
    }
    docker_example = {
      description = "Docker operations (example)"
      type        = "docker"
      enabled     = true
      configuration = jsonencode({
        allowed_registries = ["docker.io", "gcr.io", "ghcr.io"]
        max_containers     = 20
        network_mode       = "bridge"
      })
    }
    api_example = {
      description = "API integrations (example)"
      type        = "custom"
      enabled     = true
      configuration = jsonencode({
        allowed_domains = ["api.github.com", "api.slack.com", "hooks.slack.com"]
        timeout         = 30
      })
    }
  }

  # Multiple policies
  policies = {
    security = {
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
    cost_control = {
      description = "Cost control policy"
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
    compliance = {
      description = "Compliance policy"
      enabled     = true
      policy_content = <<-EOT
        package kubiya.compliance

        # Require encryption
        deny[msg] {
          input.resource_type = "database"
          not input.encrypted
          msg := "All databases must have encryption enabled"
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
  }

  # Multiple agents
  agents = {
    deployer = {
      description = "Production deployment agent"
      model_id    = "gpt-4"
      runtime     = "default"
      llm_config = jsonencode({
        temperature = 0.3
        max_tokens  = 4000
      })
      capabilities = ["kubernetes_deploy", "helm_deploy", "rollback"]
      configuration = jsonencode({
        max_retries     = 3
        timeout         = 900
        approval_needed = true
      })
      team_name = "devops_example"
    }
    monitor = {
      description = "Monitoring and alerting agent"
      model_id    = "gpt-4"
      runtime     = "default"
      llm_config = jsonencode({
        temperature = 0.5
        max_tokens  = 2000
      })
      capabilities = ["metrics_collection", "alerting", "log_analysis"]
      configuration = jsonencode({
        check_interval = 60
        alert_channels = ["slack", "pagerduty"]
      })
      team_name = "sre_example"
    }
    incident_responder = {
      description = "Incident response agent"
      model_id    = "gpt-4"
      runtime     = "default"
      llm_config = jsonencode({
        temperature = 0.4
        max_tokens  = 3000
      })
      capabilities = ["incident_management", "root_cause_analysis", "remediation"]
      configuration = jsonencode({
        escalation_timeout = 600
      })
      team_name = "sre_example"
    }
    data_pipeline = {
      description = "Data pipeline management agent"
      model_id    = "gpt-4"
      runtime     = "default"
      llm_config = jsonencode({
        temperature = 0.6
        max_tokens  = 2500
      })
      capabilities = ["etl", "data_quality", "pipeline_monitoring"]
      configuration = jsonencode({
        max_retries = 5
      })
      team_name = "data_eng_example"
    }
  }

  # Multiple worker queues
  worker_queues = {
    production-primary = {
      environment_name   = "prod_example"
      display_name       = "Production Primary Queue"
      description        = "Primary worker queue for production workloads"
      heartbeat_interval = 60
      max_workers        = 20
      tags               = ["production", "primary", "high-priority"]
      settings = {
        region   = "us-east-1"
        tier     = "production"
        priority = "high"
      }
    }
    production-secondary = {
      environment_name   = "prod_example"
      display_name       = "Production Secondary Queue"
      description        = "Secondary worker queue for batch jobs"
      heartbeat_interval = 120
      max_workers        = 10
      tags               = ["production", "secondary", "batch"]
      settings = {
        region   = "us-east-1"
        tier     = "production"
        priority = "normal"
      }
    }
    staging-default = {
      environment_name   = "staging_example"
      display_name       = "Staging Default Queue"
      description        = "Default worker queue for staging"
      heartbeat_interval = 60
      max_workers        = 10
      tags               = ["staging", "default"]
      settings = {
        region = "us-west-2"
        tier   = "staging"
      }
    }
    development-default = {
      environment_name   = "dev_example"
      display_name       = "Development Default Queue"
      description        = "Default worker queue for development"
      heartbeat_interval = 120
      max_workers        = 5
      tags               = ["development", "default"]
      settings = {
        region = "us-west-2"
        tier   = "development"
      }
    }
  }

  # Multiple jobs
  jobs = {
    health_check = {
      description   = "Daily health check"
      enabled       = true
      trigger_type  = "cron"
      cron_schedule = "0 9 * * *" # 9 AM UTC daily
      cron_timezone = "UTC"
      planning_mode = "predefined_agent"
      entity_type   = "agent"
      entity_name   = "monitor"
      prompt_template = "Run daily health check for all production services"
      system_prompt = "Check the health of all production services and report any issues"
      executor_type = "auto"
      execution_env_vars = {
        CHECK_TYPE       = "comprehensive"
        ALERT_ON_FAILURE = "true"
      }
    }
    nightly_backup = {
      description   = "Nightly database backup"
      enabled       = true
      trigger_type  = "cron"
      cron_schedule = "0 2 * * *" # 2 AM UTC daily
      cron_timezone = "UTC"
      planning_mode = "predefined_team"
      entity_type   = "team"
      entity_name   = "devops_example"
      prompt_template = "Perform nightly backup of all production databases"
      system_prompt = "Execute database backup procedures and verify completion"
      executor_type = "environment"
      environment_name = "prod_example"
      execution_env_vars = {
        BACKUP_TYPE       = "full"
        RETENTION_DAYS    = "30"
        COMPRESSION_LEVEL = "9"
      }
      execution_secrets = ["db_credentials", "s3_backup_bucket"]
    }
    deployment_webhook = {
      description  = "Handle deployment webhook events"
      enabled      = true
      trigger_type = "webhook"
      planning_mode = "predefined_agent"
      entity_type   = "agent"
      entity_name   = "deployer"
      prompt_template = "Process deployment: {{service_name}} version {{version}} to {{environment}}"
      system_prompt = "Process deployment requests and verify prerequisites"
      executor_type = "environment"
      environment_name = "prod_example"
      config = jsonencode({
        timeout = 1800 # 30 minutes
        retry_policy = {
          max_attempts = 3
          backoff      = "exponential"
        }
      })
    }
    incident_response = {
      description  = "Manual incident response"
      enabled      = true
      trigger_type = "manual"
      planning_mode = "predefined_agent"
      entity_type   = "agent"
      entity_name   = "incident_responder"
      prompt_template = "Handle incident: {{incident_id}} - {{description}}"
      system_prompt = "Coordinate incident response and resolution"
      executor_type = "auto"
      execution_secrets = ["pagerduty_token", "slack_webhook"]
    }
    data_quality_check = {
      description   = "Hourly data quality check"
      enabled       = true
      trigger_type  = "cron"
      cron_schedule = "0 * * * *" # Every hour
      cron_timezone = "UTC"
      planning_mode = "predefined_agent"
      entity_type   = "agent"
      entity_name   = "data_pipeline"
      prompt_template = "Run data quality checks for {{pipeline_name}}"
      system_prompt = "Validate data quality and report anomalies"
      executor_type = "auto"
      execution_env_vars = {
        CHECK_LEVEL = "standard"
        ALERT_THRESHOLD = "0.95"
      }
    }
  }
}

# ============================================================================
# Outputs
# ============================================================================

# Commented out since defaults module is commented out
# output "defaults_summary" {
#   description = "Summary of resources created with defaults"
#   value       = module.engineering_org_defaults.summary
# }

output "custom_summary" {
  description = "Summary of resources created with custom config"
  value       = module.engineering_org_custom.summary
}

output "custom_environment_ids" {
  description = "Environment IDs from custom module"
  value       = module.engineering_org_custom.environment_ids
}

output "custom_team_ids" {
  description = "Team IDs from custom module"
  value       = module.engineering_org_custom.team_ids
}

output "custom_agent_ids" {
  description = "Agent IDs from custom module"
  value       = module.engineering_org_custom.agent_ids
}

output "custom_webhook_urls" {
  description = "Webhook URLs from custom module"
  value       = module.engineering_org_custom.job_webhook_urls
}

output "custom_worker_queue_task_names" {
  description = "Worker queue task names for worker registration"
  value       = module.engineering_org_custom.worker_queue_task_names
}
