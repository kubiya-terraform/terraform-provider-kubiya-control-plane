terraform {
  required_providers {
    controlplane = {
      source = "kubiya/control-plane"
    }
  }
}

provider "controlplane" {}

# Create environment for worker queues
resource "controlplane_environment" "test" {
  name = "test-environment-for-workers-ds"
}

# Minimal worker queue (required fields only)
resource "controlplane_worker_queue" "minimal" {
  name           = "test-worker-queue-minimal"
  environment_id = controlplane_environment.test.id
}

# Full worker queue with all optional fields
resource "controlplane_worker_queue" "full" {
  name           = "test-worker-queue-full-ds"
  description    = "Comprehensive test worker queue with all fields configured"
  environment_id = controlplane_environment.test.id

  max_workers        = 10
  heartbeat_interval = 60
  status             = "active"

  tags = ["test", "comprehensive", "full-config"]
}

# Inactive worker queue for status testing
resource "controlplane_worker_queue" "inactive" {
  name           = "test-worker-queue-inactive-ds"
  description    = "Inactive worker queue for status testing"
  environment_id = controlplane_environment.test.id
  status         = "inactive"
}

# Data sources
data "controlplane_worker_queue" "minimal_lookup" {
  id = controlplane_worker_queue.minimal.id
}

data "controlplane_worker_queue" "full_lookup" {
  id = controlplane_worker_queue.full.id
}

data "controlplane_worker_queue" "inactive_lookup" {
  id = controlplane_worker_queue.inactive.id
}

# List data source - filter by environment
data "controlplane_worker_queues" "test_env" {
  environment_id = controlplane_environment.test.id
}

# Outputs for tests
output "data_minimal_name" {
  value = data.controlplane_worker_queue.minimal_lookup.name
}

output "data_minimal_environment_id" {
  value = data.controlplane_worker_queue.minimal_lookup.environment_id
}

output "data_full_description" {
  value = data.controlplane_worker_queue.full_lookup.description
}

output "data_full_max_workers" {
  value = tostring(data.controlplane_worker_queue.full_lookup.max_workers)
}

output "data_full_tags" {
  value = jsonencode(data.controlplane_worker_queue.full_lookup.tags)
}

output "data_inactive_status" {
  value = data.controlplane_worker_queue.inactive_lookup.status
}

output "test_env_queues_count" {
  value = length(data.controlplane_worker_queues.test_env.queues)
}

output "test_env_queues_list" {
  value = jsonencode([for q in data.controlplane_worker_queues.test_env.queues : q.name])
}
