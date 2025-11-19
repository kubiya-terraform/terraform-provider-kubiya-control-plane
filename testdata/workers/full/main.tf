terraform {
  required_providers {
    controlplane = {
      source = "kubiya/control-plane"
    }
  }
}

provider "controlplane" {}

# Create environment for worker queue
resource "controlplane_environment" "test" {
  name = "test-environment-for-full-workers"
}

# Full worker queue with all optional fields
resource "controlplane_worker_queue" "full" {
  name           = "test-worker-queue-full"
  description    = "Comprehensive test worker queue with all fields configured"
  environment_id = controlplane_environment.test.id

  max_workers         = 10
  heartbeat_interval  = 60
  status              = "active"

  tags = ["test", "comprehensive", "full-config"]
}

# Data source test
data "controlplane_worker_queue" "full_lookup" {
  id = controlplane_worker_queue.full.id
}

# Outputs
output "worker_id" {
  value = controlplane_worker_queue.full.id
}

output "worker_name" {
  value = controlplane_worker_queue.full.name
}

output "worker_description" {
  value = controlplane_worker_queue.full.description
}

output "worker_max_workers" {
  value = controlplane_worker_queue.full.max_workers
}

output "worker_heartbeat_interval" {
  value = controlplane_worker_queue.full.heartbeat_interval
}

output "data_worker_description" {
  value = data.controlplane_worker_queue.full_lookup.description
}

output "data_worker_max_workers" {
  value = data.controlplane_worker_queue.full_lookup.max_workers
}
