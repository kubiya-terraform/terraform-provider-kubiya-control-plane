# Getting Started with Control Plane Provider

This guide will walk you through setting up and using the Control Plane Terraform provider.

## Prerequisites

- Terraform >= 1.0
- A Control Plane API key
- Go >= 1.22 (for local development)

## Installation

### For Local Development

1. Clone the repository:
   ```shell
   git clone <repository-url>
   cd kubiya-control-plane-terraform-provider
   ```

2. Build and install the provider:
   ```shell
   make install
   ```

3. Configure Terraform to use the local provider by adding this to your `~/.terraformrc`:
   ```hcl
   provider_installation {
     dev_overrides {
       "hashicorp.com/kubiya/control-plane" = "/path/to/your/GOPATH/bin"
     }
     direct {}
   }
   ```

## Authentication

Set your API key as an environment variable:

```shell
export CONTROL_PLANE_API_KEY=your_api_key_here
```

Optionally, set the environment:

```shell
export CONTROL_PLANE_ENV=production  # or staging
```

## Basic Usage

Create a `main.tf` file:

```terraform
terraform {
  required_providers {
    control_plane = {
      source = "hashicorp.com/kubiya/control-plane"
    }
  }
}

provider "control_plane" {}

# Add your resources here
```

Initialize and apply:

```shell
terraform init
terraform plan
terraform apply
```

## Next Steps

- Explore the [examples](../../examples) directory for more configurations
- Check out the [resources documentation](../resources) for available resources
- Read about [data sources](../data-sources) to query existing infrastructure
