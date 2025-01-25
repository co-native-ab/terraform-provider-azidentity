terraform {
  required_providers {
    azidentity = {
      source = "registry.terraform.io/co-native-ab/azidentity"
    }
  }
}

provider "azidentity" {}

ephemeral "azidentity_default_credential" "this" {
  scopes = ["https://management.azure.com/.default"]
}
