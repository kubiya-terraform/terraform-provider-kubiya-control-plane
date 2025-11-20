
provider "controlplane" {}

variable "worker_queue_id" {
  type        = string
  description = "Worker Queue ID to import"
}

variable "worker_queue_name" {
  type        = string
  description = "Worker Queue name"
}

resource "controlplane_worker_queue" "imported" {
  name = var.worker_queue_name
}

output "imported_worker_queue_id" {
  value = controlplane_worker_queue.imported.id
}

output "imported_worker_queue_name" {
  value = controlplane_worker_queue.imported.name
}
