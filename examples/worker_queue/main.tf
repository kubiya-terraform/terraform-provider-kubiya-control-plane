terraform {
  required_providers {
    controlplane = {
      source = "kubiya/control-plane"
    }
  }
}

provider "controlplane" {
  # Configuration is via environment variables:
  # KUBIYA_CONTROL_PLANE_API_KEY
  # KUBIYA_CONTROL_PLANE_BASE_URL (optional, defaults to https://control-plane.kubiya.ai)
}

# First, create or reference an environment
resource "controlplane_environment" "example" {
  name         = "production"
  display_name = "Production Environment"
  description  = "Main production environment"
}

# Create a worker queue
resource "controlplane_worker_queue" "example" {
  environment_id     = controlplane_environment.example.id
  name               = "default-queue"
  display_name       = "Default Worker Queue"
  description        = "Main worker queue for production"
  heartbeat_interval = 60 # Default changed from 30 to 60 (lightweight heartbeats)
  max_workers        = 10

  tags = ["production", "primary"]

  settings = {
    region = "us-east-1"
    tier   = "production"
  }
}

# Data source example - fetch a worker queue by ID
data "controlplane_worker_queue" "example" {
  id = controlplane_worker_queue.example.id
}

# Data source example - list all worker queues in an environment
data "controlplane_worker_queues" "all" {
  environment_id = controlplane_environment.example.id
}

# Output worker queue information
output "queue_id" {
  value       = controlplane_worker_queue.example.id
  description = "The ID of the worker queue"
}

output "queue_name" {
  value       = controlplane_worker_queue.example.name
  description = "The name of the worker queue"
}

output "task_queue_name" {
  value       = controlplane_worker_queue.example.task_queue_name
  description = "The Temporal task queue name"
}

output "active_workers" {
  value       = controlplane_worker_queue.example.active_workers
  description = "Number of active workers in the queue"
}

output "all_queues_count" {
  value       = length(data.controlplane_worker_queues.all.queues)
  description = "Total number of worker queues in the environment"
}

# Note: Worker queues organize and manage workers within an environment.
# Workers connect to the control plane and register with a specific queue.
# Use the task_queue_name when starting workers to connect to this queue.
