terraform {
  required_providers {
    azidentity = {
      source = "registry.terraform.io/co-native-ab/azidentity"
    }
  }
}

provider "azidentity" {}

ephemeral "azidentity_azure_cli_credential" "this" {
  scopes = ["https://management.azure.com/.default"]
}
