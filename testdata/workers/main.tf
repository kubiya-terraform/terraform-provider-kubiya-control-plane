terraform {
  required_providers {
    controlplane = {
      source = "kubiya/control-plane"
    }
  }
}

provider "controlplane" {
  # Configuration via environment variables:
  # KUBIYA_CONTROL_PLANE_API_KEY
  # KUBIYA_CONTROL_PLANE_BASE_URL (optional, defaults to https://control-plane.kubiya.ai)
}

# Test worker resource
# Note: Workers typically self-register at runtime, but can be pre-registered
resource "controlplane_worker" "test" {
  environment_name = "default"
  hostname         = "test-worker-01"

  # Worker metadata
  metadata = jsonencode({
    region     = "us-east-1"
    datacenter = "dc1"
    capacity   = "test"
    tags = {
      environment = "test"
      purpose     = "automated-testing"
    }
  })
}

output "worker_id" {
  value = controlplane_worker.test.id
}

output "worker_status" {
  value = controlplane_worker.test.status
}

output "worker_registered_at" {
  value = controlplane_worker.test.registered_at
}

# Note: Workers are runtime entities that self-manage their lifecycle.
# The worker resource is primarily for registration and discovery.
# Workers will connect to the control plane using the environment's worker token.
