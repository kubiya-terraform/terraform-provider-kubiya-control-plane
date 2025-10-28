# Control Plane Provider

The Control Plane provider is used to interact with Kubiya Control Plane resources.

## Example Usage

```terraform
terraform {
  required_providers {
    control_plane = {
      source = "hashicorp.com/kubiya/control-plane"
    }
  }
}

provider "control_plane" {
  # Configuration is read from environment variables:
  # - CONTROL_PLANE_API_KEY (required)
  # - CONTROL_PLANE_ENV (optional, defaults to "production")
}
```

## Authentication

The provider requires authentication using an API key. You can set the API key using the `CONTROL_PLANE_API_KEY` environment variable:

```shell
export CONTROL_PLANE_API_KEY=your_api_key_here
```

Optionally, you can specify the environment (production or staging):

```shell
export CONTROL_PLANE_ENV=production
```

## Resources

Resources will be documented as they are implemented.

## Data Sources

Data sources will be documented as they are implemented.
