terraform {
  required_providers {
    controlplane = {
      source = "kubiya/control-plane"
    }
  }
}

provider "controlplane" {}

variable "description" {
  type    = string
  default = "Original worker_queue description"
}

resource "controlplane_worker_queue" "test" {
  name        = "test-worker-queue-update"
  description = var.description
}

output "worker_queue_id" {
  value = controlplane_worker_queue.test.id
}

output "worker_queue_name" {
  value = controlplane_worker_queue.test.name
}

output "worker_queue_description" {
  value = controlplane_worker_queue.test.description
}

output "worker_queue_updated_at" {
  value = controlplane_worker_queue.test.updated_at
}
