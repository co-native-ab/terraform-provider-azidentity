terraform {
  required_providers {
    azidentity = {
      source = "registry.terraform.io/co-native-ab/azidentity"
    }
  }
}

provider "azidentity" {}

ephemeral "azidentity_client_secret_credential" "this" {
  tenant_id     = "00000000-0000-0000-0000-000000000000"
  client_id     = "00000000-0000-0000-0000-000000000000"
  client_secret = "supersecret"
  scopes        = ["https://management.azure.com/.default"]
}
