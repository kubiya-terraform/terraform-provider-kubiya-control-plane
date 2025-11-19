terraform {
  required_providers {
    controlplane = {
      source = "kubiya/control-plane"
    }
  }
}

provider "controlplane" {}

variable "job_id" {
  type        = string
  description = "Job ID to import"
}

variable "job_name" {
  type        = string
  description = "Job name"
}

resource "controlplane_job" "imported" {
  name = var.job_name
}

output "imported_job_id" {
  value = controlplane_job.imported.id
}

output "imported_job_name" {
  value = controlplane_job.imported.name
}
