# Kubiya Control Plane Terraform Provider Examples

This directory contains example Terraform configurations for using the Kubiya Control Plane Provider. Each subdirectory demonstrates a specific resource type or use case.

## Prerequisites

Before running any examples, ensure you have:

1. **Terraform installed** (version 1.0 or later)
2. **Kubiya Control Plane credentials**:
   - API Key
   - Organization ID

## Configuration

Set the required environment variables:

```bash
export KUBIYA_CONTROL_PLANE_API_KEY="your-api-key"
export KUBIYA_CONTROL_PLANE_ORG_ID="your-org-id"
# Optional: Override the default API URL
# export KUBIYA_CONTROL_PLANE_BASE_URL="https://custom-url.example.com"
```

## Examples Directory Structure

### Individual Resource Examples

Each directory contains examples for a specific resource type:

- **agent/** - Agent resource and data source examples
- **team/** - Team resource and data source examples
- **project/** - Project resource and data source examples
- **environment/** - Environment resource and data source examples
- **toolset/** - ToolSet resource and data source examples
- **policy/** - Policy resource and data source examples
- **worker/** - Worker resource examples
- **complete-setup/** - End-to-end example showing all resources working together

## Usage

Navigate to any example directory and run:

```bash
terraform init
terraform plan
terraform apply
```

See individual example directories for detailed documentation.
