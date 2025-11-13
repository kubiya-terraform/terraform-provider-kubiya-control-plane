# Kubiya Control Plane Provider Documentation

This directory contains the documentation for the Kubiya Control Plane Terraform Provider, formatted for the Terraform Registry.

## Documentation Structure

```
docs/
├── index.md                    # Provider overview and configuration
├── guides/
│   └── getting-started.md      # Getting started guide
├── resources/
│   ├── agent.md                # Agent resource documentation
│   ├── team.md                 # Team resource documentation
│   ├── project.md              # Project resource documentation
│   ├── environment.md          # Environment resource documentation
│   ├── skill.md                # Skill resource documentation
│   ├── policy.md               # Policy resource documentation
│   ├── worker_queue.md         # Worker Queue resource documentation
│   └── job.md                  # Job resource documentation
└── data-sources/
    ├── agent.md                # Agent data source documentation
    ├── team.md                 # Team data source documentation
    ├── project.md              # Project data source documentation
    ├── environment.md          # Environment data source documentation
    ├── skill.md                # Skill data source documentation
    ├── policy.md               # Policy data source documentation
    ├── worker_queue.md         # Worker Queue data source documentation
    ├── worker_queues.md        # Worker Queues (list) data source documentation
    ├── job.md                  # Job data source documentation
    └── jobs.md                 # Jobs (list) data source documentation
```

## Resources

The provider manages the following resources:

- **controlplane_agent** - AI agents with LLM configuration
- **controlplane_team** - Teams for organizing agents
- **controlplane_project** - Projects for grouping resources
- **controlplane_environment** - Execution environments
- **controlplane_skill** - Skills (filesystem, shell, docker, etc.)
- **controlplane_policy** - OPA Rego governance policies
- **controlplane_worker_queue** - Worker queue configuration and management
- **controlplane_job** - Scheduled, webhook-triggered, and manual jobs

## Data Sources

Corresponding data sources for resource lookup:

- **controlplane_agent** - Look up existing agents
- **controlplane_team** - Look up existing teams
- **controlplane_project** - Look up existing projects
- **controlplane_environment** - Look up existing environments
- **controlplane_skill** - Look up existing skills
- **controlplane_policy** - Look up existing policies
- **controlplane_worker_queue** - Look up a worker queue
- **controlplane_worker_queues** - List all worker queues in an environment
- **controlplane_job** - Look up a job
- **controlplane_jobs** - List all jobs

## Publishing to Terraform Registry

When publishing to the Terraform Registry, ensure:

1. Documentation follows the registry format (frontmatter with page_title, description)
2. Files are in the correct directories (resources/, data-sources/, guides/)
3. Examples are clear and complete
4. All required attributes are documented
5. Import statements are included

## Local Development

To preview documentation locally:

1. Use the Terraform Registry documentation preview tool
2. Or serve the markdown files with any markdown viewer
3. Check that all links work correctly

## Contributing

When adding new resources or data sources:

1. Create documentation in the appropriate directory
2. Follow the existing format and structure
3. Include practical examples
4. Document all attributes (required, optional, read-only)
5. Add import instructions
6. Update this README
