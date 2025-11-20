
provider "controlplane" {}

variable "project_id" {
  type        = string
  description = "Project ID to import"
}

variable "project_name" {
  type        = string
  description = "Project name"
}

variable "project_key" {
  type        = string
  description = "Project key"
}

resource "controlplane_project" "imported" {
  name = var.project_name
  key  = var.project_key
}

output "imported_project_id" {
  value = controlplane_project.imported.id
}

output "imported_project_name" {
  value = controlplane_project.imported.name
}

output "imported_project_key" {
  value = controlplane_project.imported.key
}
