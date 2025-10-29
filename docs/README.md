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
│   ├── toolset.md              # ToolSet resource documentation
│   ├── policy.md               # Policy resource documentation
│   └── worker.md               # Worker resource documentation
└── data-sources/
    ├── agent.md                # Agent data source documentation
    ├── team.md                 # Team data source documentation
    ├── project.md              # Project data source documentation
    ├── environment.md          # Environment data source documentation
    ├── toolset.md              # ToolSet data source documentation
    └── policy.md               # Policy data source documentation
```

## Resources

The provider manages the following resources:

- **kubiya_control_plane_agent** - AI agents with LLM configuration
- **kubiya_control_plane_team** - Teams for organizing agents
- **kubiya_control_plane_project** - Projects for grouping resources
- **kubiya_control_plane_environment** - Execution environments
- **kubiya_control_plane_toolset** - Toolsets (filesystem, shell, docker, etc.)
- **kubiya_control_plane_policy** - OPA Rego governance policies
- **kubiya_control_plane_worker** - Worker registration

## Data Sources

Corresponding data sources for resource lookup:

- **kubiya_control_plane_agent** - Look up existing agents
- **kubiya_control_plane_team** - Look up existing teams
- **kubiya_control_plane_project** - Look up existing projects
- **kubiya_control_plane_environment** - Look up existing environments
- **kubiya_control_plane_toolset** - Look up existing toolsets
- **kubiya_control_plane_policy** - Look up existing policies

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
