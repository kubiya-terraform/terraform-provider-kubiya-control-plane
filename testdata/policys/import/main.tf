
provider "controlplane" {}

variable "policy_id" {
  type        = string
  description = "Policy ID to import"
}

variable "policy_name" {
  type        = string
  description = "Policy name"
}

resource "controlplane_policy" "imported" {
  name = var.policy_name
}

output "imported_policy_id" {
  value = controlplane_policy.imported.id
}

output "imported_policy_name" {
  value = controlplane_policy.imported.name
}
