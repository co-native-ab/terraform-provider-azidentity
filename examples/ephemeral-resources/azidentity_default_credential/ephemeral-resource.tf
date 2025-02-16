terraform {
  required_version = ">= 1.10.0"
  required_providers {
    azidentity = {
      source = "co-native-ab/azidentity"
    }
  }
}

provider "azidentity" {}

ephemeral "azidentity_default_credential" "this" {
  scopes = ["https://management.azure.com/.default"]
}
