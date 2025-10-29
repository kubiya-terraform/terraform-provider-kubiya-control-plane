terraform {
  required_providers {
    kubiya_control_plane = {
      source = "kubiya/control-plane"
    }
  }
}

provider "kubiya_control_plane" {
  # Configuration via environment variables:
  # KUBIYA_CONTROL_PLANE_API_KEY
  # KUBIYA_CONTROL_PLANE_ORG_ID
}

# Test worker resource
resource "kubiya_control_plane_worker" "test" {
  name        = "test-worker"
  description = "Test worker for automated testing"

  configuration = jsonencode({
    max_concurrent_tasks = 5
    timeout              = 300
  })
}

output "worker_id" {
  value = kubiya_control_plane_worker.test.id
}

output "worker_name" {
  value = kubiya_control_plane_worker.test.name
}

output "worker_status" {
  value = kubiya_control_plane_worker.test.status
}
