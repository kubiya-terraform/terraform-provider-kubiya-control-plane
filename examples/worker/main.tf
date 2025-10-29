terraform {
  required_providers {
    kubiya_control_plane = {
      source = "kubiya/control-plane"
    }
  }
}

provider "kubiya_control_plane" {
  # Configuration is via environment variables:
  # KUBIYA_CONTROL_PLANE_API_KEY
  # KUBIYA_CONTROL_PLANE_ORG_ID
  # KUBIYA_CONTROL_PLANE_BASE_URL (optional, defaults to https://control-plane.kubiya.ai)
}

# Register a worker
# Note: Workers typically self-register at runtime, but can be pre-registered
resource "kubiya_control_plane_worker" "example" {
  environment_name = "production"
  hostname         = "worker-node-01"

  # Worker metadata
  metadata = jsonencode({
    region     = "us-east-1"
    datacenter = "dc1"
    capacity   = "high"
    tags = {
      environment = "production"
      team        = "platform"
    }
  })
}

# Output worker information
output "worker_id" {
  value       = kubiya_control_plane_worker.example.id
  description = "The ID of the registered worker"
}

output "worker_status" {
  value       = kubiya_control_plane_worker.example.status
  description = "The current status of the worker"
}

output "worker_registered_at" {
  value       = kubiya_control_plane_worker.example.registered_at
  description = "When the worker was registered"
}

# Note: Workers are runtime entities that self-manage their lifecycle.
# The worker resource is primarily for registration and discovery.
# Workers will connect to the control plane using the environment's worker token.
