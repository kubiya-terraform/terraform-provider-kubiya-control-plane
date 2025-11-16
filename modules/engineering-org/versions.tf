terraform {
  required_version = ">= 1.0"

  required_providers {
    controlplane = {
      source  = "kubiya/control-plane"
      version = ">= 0.1.0"
    }
  }
}
