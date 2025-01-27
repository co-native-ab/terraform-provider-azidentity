terraform {
  required_version = ">= 1.10.0"
  required_providers {
    azidentity = {
      source = "registry.terraform.io/co-native-ab/azidentity"
    }
  }
}

provider "azidentity" {}

ephemeral "azidentity_environment_variable" "this" {
  key = "SYSTEM_ACCESSTOKEN"
}
