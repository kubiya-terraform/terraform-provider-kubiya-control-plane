# Examples

This directory contains example Terraform configurations for the Kubiya Control Plane provider.

## Prerequisites

1. Set up your API key and organization ID:
   ```shell
   export KUBIYA_CONTROL_PLANE_API_KEY=your_api_key_here
   export KUBIYA_CONTROL_PLANE_ORG_ID=your_organization_id_here
   ```

2. (Optional) Set the environment:
   ```shell
   export KUBIYA_CONTROL_PLANE_ENV=development  # or staging, production
   ```

## Running the Examples

1. Navigate to this directory:
   ```shell
   cd examples
   ```

2. Initialize Terraform:
   ```shell
   terraform init
   ```

3. Review the planned changes:
   ```shell
   terraform plan
   ```

4. Apply the configuration:
   ```shell
   terraform apply
   ```

## Available Examples

- `main.tf` - Basic provider configuration

More examples will be added as resources are implemented.
