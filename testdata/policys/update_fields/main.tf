
provider "controlplane" {}

variable "description" {
  type    = string
  default = "Original policy description"
}

resource "controlplane_policy" "test" {
  name        = "test-policy-update"
  description = var.description
}

output "policy_id" {
  value = controlplane_policy.test.id
}

output "policy_name" {
  value = controlplane_policy.test.name
}

output "policy_description" {
  value = controlplane_policy.test.description
}

output "policy_updated_at" {
  value = controlplane_policy.test.updated_at
}
