# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial provider implementation
- Support for managing agents, teams, projects, and environments
- Data sources for read-only lookups of all resource types
- Full CRUD operations for all resources
- Import support for all resources
- Environment variable configuration (`KUBIYA_CONTROL_PLANE_*`)
- Comprehensive documentation and examples

### Features
- **Resources**
  - `kubiya_control_plane_agent` - Manage AI agents with Claude Code and default runtimes
  - `kubiya_control_plane_team` - Manage teams with toolset associations
  - `kubiya_control_plane_project` - Manage projects with policy controls
  - `kubiya_control_plane_environment` - Manage execution environments

- **Data Sources**
  - `kubiya_control_plane_agent` - Lookup existing agents
  - `kubiya_control_plane_team` - Lookup existing teams
  - `kubiya_control_plane_project` - Lookup existing projects
  - `kubiya_control_plane_environment` - Lookup existing environments

- **Configuration**
  - Multi-environment support (development, staging, production)
  - Organization-scoped authentication
  - Custom base URL support for self-hosted deployments

## [0.1.0] - TBD

### Added
- Initial release of Kubiya Control Plane Terraform Provider
- Basic resource management capabilities
- Authentication via API key and organization ID
- Support for local development and testing

[Unreleased]: https://github.com/kubiya-terraform/kubiya-control-plane/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/kubiya-terraform/kubiya-control-plane/releases/tag/v0.1.0
