# ============================================================================
# Environments
# ============================================================================

variable "environments" {
  description = "Map of environments to create"
  type = map(object({
    description           = string
    settings              = optional(map(any), {})
    execution_environment = optional(map(any), {})
  }))
  default = {
    production = {
      description = "Production environment"
      settings = {
        region         = "us-east-1"
        max_workers    = 10
        auto_scaling   = true
        retention_days = 90
      }
      execution_environment = {
        env_vars = {
          LOG_LEVEL = "info"
          APP_ENV   = "production"
        }
      }
    }
  }
}

# ============================================================================
# Projects
# ============================================================================

variable "projects" {
  description = "Map of projects to create"
  type = map(object({
    key         = string
    description = string
    settings    = optional(map(any), {})
  }))
  default = {
    platform = {
      key         = "PLAT"
      description = "Platform engineering project"
      settings = {
        owner       = "platform-team"
        cost_center = "engineering"
      }
    }
  }
}

# ============================================================================
# Teams
# ============================================================================

variable "teams" {
  description = "Map of teams to create"
  type = map(object({
    description   = string
    runtime       = optional(string, "default")
    configuration = optional(map(any), {})
    capabilities  = optional(list(string), [])
  }))
  default = {
    devops = {
      description = "DevOps and platform engineering team"
      runtime     = "default"
      configuration = {
        max_agents        = 10
        enable_monitoring = true
      }
      capabilities = ["deployment", "monitoring", "incident_response"]
    }
  }
}

# ============================================================================
# Skills
# ============================================================================

variable "skills" {
  description = "Map of skills to create"
  type = map(object({
    description   = string
    type          = string
    enabled       = optional(bool, true)
    configuration = optional(map(any), {})
  }))
  default = {
    shell = {
      description = "Shell command execution"
      type        = "shell"
      enabled     = true
      configuration = {
        allowed_commands = ["kubectl", "helm", "aws", "terraform"]
        timeout          = 300
      }
    }
    filesystem = {
      description = "File system operations"
      type        = "file_system"
      enabled     = true
      configuration = {
        allowed_paths = ["/app/configs", "/app/data"]
        max_file_size = 10485760 # 10MB
        operations    = ["read", "write", "list"]
      }
    }
  }
}

# ============================================================================
# Policies
# ============================================================================

variable "policies" {
  description = "Map of policies to create"
  type = map(object({
    description    = string
    enabled        = optional(bool, true)
    policy_content = string
    tags           = optional(list(string), [])
  }))
  default = {
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
  }
}

# ============================================================================
# Agents
# ============================================================================

variable "agents" {
  description = "Map of agents to create"
  type = map(object({
    description   = string
    model_id      = optional(string, "gpt-4")
    runtime       = optional(string, "default")
    llm_config    = optional(map(any), {})
    capabilities  = optional(list(string), [])
    configuration = optional(map(any), {})
    team_name     = optional(string, null)
  }))
  default = {
    deployer = {
      description = "Deployment agent"
      model_id    = "gpt-4"
      runtime     = "default"
      llm_config = {
        temperature = 0.3
        max_tokens  = 4000
      }
      capabilities = ["kubernetes_deploy", "helm_deploy", "rollback"]
      configuration = {
        max_retries     = 3
        timeout         = 600
        approval_needed = true
      }
      team_name = "devops"
    }
    monitor = {
      description = "Monitoring and alerting agent"
      model_id    = "gpt-4"
      runtime     = "default"
      llm_config = {
        temperature = 0.5
        max_tokens  = 2000
      }
      capabilities = ["metrics_collection", "alerting", "log_analysis"]
      configuration = {
        check_interval = 60
        alert_channels = ["slack", "pagerduty"]
      }
      team_name = "devops"
    }
  }
}

# ============================================================================
# Worker Queues
# ============================================================================

variable "worker_queues" {
  description = "Map of worker queues to create"
  type = map(object({
    environment_name   = string
    display_name       = string
    description        = string
    heartbeat_interval = optional(number, 60)
    max_workers        = optional(number, 10)
    tags               = optional(list(string), [])
    settings           = optional(map(string), {})
  }))
  default = {
    default = {
      environment_name   = "production"
      display_name       = "Default Queue"
      description        = "Default worker queue for production"
      heartbeat_interval = 60
      max_workers        = 10
      tags               = ["production", "primary"]
      settings = {
        region = "us-east-1"
        tier   = "production"
      }
    }
  }
}

# ============================================================================
# Jobs
# ============================================================================

variable "jobs" {
  description = "Map of jobs to create"
  type = map(object({
    description        = string
    enabled            = optional(bool, true)
    trigger_type       = string # "cron", "webhook", or "manual"
    cron_schedule      = optional(string, null)
    cron_timezone      = optional(string, "UTC")
    planning_mode      = string # "predefined_agent", "predefined_team", or "on_the_fly"
    entity_type        = optional(string, null)
    entity_name        = optional(string, null)
    prompt_template    = string
    system_prompt      = optional(string, null)
    executor_type      = optional(string, "auto")
    environment_name   = optional(string, null)
    execution_env_vars = optional(map(string), {})
    execution_secrets  = optional(list(string), [])
    config             = optional(map(any), null)
  }))
  default = {
    health_check = {
      description   = "Daily health check"
      enabled       = true
      trigger_type  = "cron"
      cron_schedule = "0 9 * * *" # 9 AM UTC daily
      cron_timezone = "UTC"
      planning_mode = "predefined_agent"
      entity_type   = "agent"
      entity_name   = "monitor"
      prompt_template = "Run daily health check for all services"
      system_prompt = "Check the health of all production services and report any issues"
      executor_type = "auto"
      execution_env_vars = {
        CHECK_TYPE       = "comprehensive"
        ALERT_ON_FAILURE = "true"
      }
    }
  }
}
