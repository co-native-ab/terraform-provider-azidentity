terraform {
  required_version = ">= 1.10.0"
  required_providers {
    azidentity = {
      source = "registry.terraform.io/co-native-ab/azidentity"
    }
    azurerm = {
      source = "hashicorp/azurerm"
    }
  }
}

provider "azidentity" {}

ephemeral "azidentity_azure_cli_account" "this" {}

provider "azurerm" {
  features {}
  subscription_id = ephemeral.azidentity_azure_cli_account.this.subscription_id
}

resource "azurerm_resource_group" "this" {
  name     = "rg-example"
  location = "West Europe"
}
