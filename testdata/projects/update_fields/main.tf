
provider "controlplane" {}

variable "description" {
  type    = string
  default = "Original project description"
}

variable "visibility" {
  type    = string
  default = "org"
}

resource "controlplane_project" "test" {
  name        = "test-project-update"
  key         = "TPU"
  description = var.description
  visibility  = var.visibility
}

output "project_id" {
  value = controlplane_project.test.id
}

output "project_name" {
  value = controlplane_project.test.name
}

output "project_description" {
  value = controlplane_project.test.description
}

output "project_visibility" {
  value = controlplane_project.test.visibility
}

output "project_updated_at" {
  value = controlplane_project.test.updated_at
}
