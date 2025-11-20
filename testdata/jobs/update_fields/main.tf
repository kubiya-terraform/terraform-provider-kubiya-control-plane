
provider "controlplane" {}

variable "description" {
  type    = string
  default = "Original job description"
}

resource "controlplane_job" "test" {
  name        = "test-job-update"
  description = var.description
}

output "job_id" {
  value = controlplane_job.test.id
}

output "job_name" {
  value = controlplane_job.test.name
}

output "job_description" {
  value = controlplane_job.test.description
}

output "job_updated_at" {
  value = controlplane_job.test.updated_at
}
