
provider "controlplane" {}

# Create environment for worker queue
resource "controlplane_environment" "test" {
  name = "test-environment-for-workers"
}

# Minimal worker queue (required fields only)
resource "controlplane_worker_queue" "minimal" {
  name           = "test-worker-queue-minimal"
  environment_id = controlplane_environment.test.id
}

# Data source test
data "controlplane_worker_queue" "minimal_lookup" {
  id = controlplane_worker_queue.minimal.id
}

# Outputs
output "worker_id" {
  value = controlplane_worker_queue.minimal.id
}

output "worker_name" {
  value = controlplane_worker_queue.minimal.name
}

output "worker_environment_id" {
  value = controlplane_worker_queue.minimal.environment_id
}

output "data_worker_name" {
  value = data.controlplane_worker_queue.minimal_lookup.name
}
