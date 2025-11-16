# Engineering Organization Module

A comprehensive Terraform module for creating a complete engineering organization setup with the Kubiya Control Plane provider. This module creates and manages environments, projects, teams, agents, skills, policies, worker queues, and jobs in a declarative and scalable way.

## Features

- **Flexible Configuration**: Use `for_each` to create multiple instances of each resource type
- **Sensible Defaults**: Ready-to-use defaults for quick setup
- **Easy Extension**: Add or modify resources by updating variable maps
- **Complete Setup**: Creates all necessary resources for an engineering org
- **Resource Relationships**: Automatically handles dependencies between resources

## Resources Created

This module can create the following resources:

- **Environments**: Isolated execution environments for agents and workers
- **Projects**: Organizational units for grouping related work
- **Teams**: Groups of agents working together
- **Agents**: AI-powered automation agents
- **Skills**: Reusable capabilities for agents
- **Policies**: OPA Rego policies for governance and security
- **Worker Queues**: Queues for managing worker distribution
- **Jobs**: Scheduled, webhook-triggered, or manual tasks

## Usage

### Minimal Example (Using Defaults)

```hcl
module "engineering_org" {
  source = "../../modules/engineering-org"
}
```

This creates:
- 1 production environment
- 1 platform project
- 1 devops team
- 2 agents (deployer, monitor)
- 2 skills (shell, filesystem)
- 1 security policy
- 1 worker queue
- 1 health check job

### Custom Configuration

```hcl
module "engineering_org" {
  source = "../../modules/engineering-org"

  environments = {
    production = {
      description = "Production environment"
      settings = {
        region         = "us-east-1"
        max_workers    = 20
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
    staging = {
      description = "Staging environment"
      settings = {
        region         = "us-west-2"
        max_workers    = 10
        auto_scaling   = true
        retention_days = 30
      }
      execution_environment = {
        env_vars = {
          LOG_LEVEL = "debug"
          APP_ENV   = "staging"
        }
      }
    }
  }

  teams = {
    devops = {
      description = "DevOps and platform engineering team"
      runtime     = "default"
      configuration = {
        max_agents = 15
      }
      capabilities = ["deployment", "monitoring"]
    }
    sre = {
      description = "Site reliability engineering team"
      runtime     = "default"
      configuration = {
        max_agents = 10
      }
      capabilities = ["monitoring", "incident_response"]
    }
  }

  agents = {
    deployer = {
      description = "Production deployment agent"
      model_id    = "gpt-4"
      runtime     = "default"
      llm_config = {
        temperature = 0.3
        max_tokens  = 4000
      }
      capabilities = ["kubernetes_deploy", "helm_deploy"]
      team_name    = "devops"
    }
    monitor = {
      description = "Monitoring agent"
      model_id    = "gpt-4"
      runtime     = "default"
      capabilities = ["metrics_collection", "alerting"]
      team_name    = "sre"
    }
  }
}
```

### Adding More Resources

To add resources, simply add entries to the appropriate variable map:

```hcl
module "engineering_org" {
  source = "../../modules/engineering-org"

  # Add a new skill
  skills = {
    shell = {
      description = "Shell command execution"
      type        = "shell"
      enabled     = true
      configuration = {
        allowed_commands = ["kubectl", "helm", "aws"]
        timeout          = 300
      }
    }
    # Add your new skill here
    database = {
      description = "Database operations"
      type        = "database"
      enabled     = true
      configuration = {
        allowed_databases = ["postgresql", "mysql"]
        max_connections   = 10
      }
    }
  }
}
```

## Requirements

| Name | Version |
|------|---------|
| terraform | >= 1.0 |
| controlplane | >= 0.1.0 |

## Providers

| Name | Version |
|------|---------|
| controlplane | >= 0.1.0 |

## Inputs

### environments

Map of environments to create.

**Type:**
```hcl
map(object({
  description           = string
  settings              = optional(map(any), {})
  execution_environment = optional(map(any), {})
}))
```

**Default:** Creates a production environment with standard settings.

### projects

Map of projects to create.

**Type:**
```hcl
map(object({
  key         = string
  description = string
  settings    = optional(map(any), {})
}))
```

**Default:** Creates a platform project.

### teams

Map of teams to create.

**Type:**
```hcl
map(object({
  description   = string
  runtime       = optional(string, "default")
  configuration = optional(map(any), {})
  capabilities  = optional(list(string), [])
}))
```

**Default:** Creates a devops team.

### agents

Map of agents to create.

**Type:**
```hcl
map(object({
  description   = string
  model_id      = optional(string, "gpt-4")
  runtime       = optional(string, "default")
  llm_config    = optional(map(any), {})
  capabilities  = optional(list(string), [])
  configuration = optional(map(any), {})
  team_name     = optional(string, null)
}))
```

**Default:** Creates deployer and monitor agents.

### skills

Map of skills to create.

**Type:**
```hcl
map(object({
  description   = string
  type          = string
  enabled       = optional(bool, true)
  configuration = optional(map(any), {})
}))
```

**Default:** Creates shell and filesystem skills.

### policies

Map of policies to create.

**Type:**
```hcl
map(object({
  description    = string
  enabled        = optional(bool, true)
  policy_content = string
  tags           = optional(list(string), [])
}))
```

**Default:** Creates a security policy.

### worker_queues

Map of worker queues to create.

**Type:**
```hcl
map(object({
  environment_name   = string
  display_name       = string
  description        = string
  heartbeat_interval = optional(number, 60)
  max_workers        = optional(number, 10)
  tags               = optional(list(string), [])
  settings           = optional(map(string), {})
}))
```

**Default:** Creates a default queue for production.

### jobs

Map of jobs to create.

**Type:**
```hcl
map(object({
  description        = string
  enabled            = optional(bool, true)
  trigger_type       = string
  cron_schedule      = optional(string, null)
  cron_timezone      = optional(string, "UTC")
  planning_mode      = string
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
```

**Default:** Creates a daily health check job.

## Outputs

### Resource Collections

- `environments` - Map of created environments with details
- `projects` - Map of created projects with details
- `teams` - Map of created teams with details
- `agents` - Map of created agents with details
- `skills` - Map of created skills with details
- `policies` - Map of created policies with details
- `worker_queues` - Map of created worker queues with details
- `jobs` - Map of created jobs with details

### ID Maps

- `environment_ids` - Map of environment names to IDs
- `project_ids` - Map of project names to IDs
- `team_ids` - Map of team names to IDs
- `agent_ids` - Map of agent names to IDs
- `skill_ids` - Map of skill names to IDs
- `policy_ids` - Map of policy names to IDs
- `worker_queue_ids` - Map of worker queue names to IDs
- `job_ids` - Map of job names to IDs

### Special Outputs

- `job_webhook_urls` - Map of webhook-triggered job names to webhook URLs
- `worker_queue_task_names` - Map of worker queue names to task queue names (for worker registration)
- `summary` - Count of all created resources

## Examples

See the [examples/engineering-org](../../examples/engineering-org) directory for complete examples:

- **Defaults**: Minimal configuration using all defaults
- **Custom**: Extended configuration with multiple environments, teams, agents, and more

## Best Practices

1. **Start with defaults**: Begin with the default configuration and customize as needed
2. **Use meaningful names**: Resource keys in maps become part of the resource names
3. **Reference by name**: Use resource names (map keys) to create relationships between resources
4. **Environment separation**: Create separate environments for production, staging, and development
5. **Team organization**: Group agents into teams based on their function
6. **Policy enforcement**: Define policies early to ensure compliance from the start
7. **Worker queue strategy**: Create separate queues for different priorities or workload types

## Resource Relationships

The module automatically handles dependencies:

- Agents can reference teams using `team_name`
- Worker queues reference environments using `environment_name`
- Jobs reference agents or teams using `entity_name` and environments using `environment_name`

Example:
```hcl
agents = {
  deployer = {
    description = "Deployment agent"
    team_name   = "devops"  # References the devops team
    # ...
  }
}

jobs = {
  deploy_job = {
    entity_type  = "agent"
    entity_name  = "deployer"  # References the deployer agent
    environment_name = "production"  # References the production environment
    # ...
  }
}
```

## Extending Default Values

To extend default values while keeping some defaults:

```hcl
module "engineering_org" {
  source = "../../modules/engineering-org"

  # Keep default environments, projects, teams, skills, policies
  # Only customize agents
  agents = {
    deployer = {
      description = "My custom deployer"
      model_id    = "gpt-4"
      runtime     = "default"
      llm_config = {
        temperature = 0.2  # Lower temperature for more deterministic behavior
      }
      capabilities = ["kubernetes_deploy"]
      team_name    = "devops"
    }
    # Add more custom agents
    security_scanner = {
      description = "Security scanning agent"
      model_id    = "gpt-4"
      capabilities = ["security_scan", "vulnerability_assessment"]
      team_name    = "devops"
    }
  }
}
```

## License

See the provider's license for details.
