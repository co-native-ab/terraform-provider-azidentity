terraform {
  required_version = ">= 1.10.0"
  required_providers {
    azidentity = {
      source = "co-native-ab/azidentity"
    }
  }
}

provider "azidentity" {}

ephemeral "azidentity_azure_cli_account" "this" {}
