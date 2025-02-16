terraform {
  required_version = ">= 1.10.0"
  required_providers {
    azurerm = {
      source = "hashicorp/azurerm"
    }
    azapi = {
      source = "Azure/azapi"
    }
    azidentity = {
      source = "co-native-ab/azidentity"
    }
    grafana = {
      source = "grafana/grafana"
    }
  }
}

provider "azidentity" {}

ephemeral "azidentity_azure_cli_account" "this" {}

provider "azurerm" {
  features {}
  subscription_id = ephemeral.azidentity_azure_cli_account.this.subscription_id
}

provider "azapi" {
  subscription_id = ephemeral.azidentity_azure_cli_account.this.subscription_id
}

resource "azurerm_resource_group" "this" {
  name     = "rg-lab-sc-grafana"
  location = "Sweden Central"
}

resource "azapi_resource" "grafana" {
  type                      = "Microsoft.Dashboard/grafana@2024-10-01"
  name                      = "grf-lab-sc-grafana"
  location                  = azurerm_resource_group.this.location
  parent_id                 = azurerm_resource_group.this.id
  schema_validation_enabled = false

  body = {
    sku = {
      name = "Essential"
    }
    properties = {
      apiKey                  = "Disabled"
      deterministicOutboundIP = "Disabled"
      grafanaConfigurations = {
        security = {
          csrfAlwaysCheck = true
        }
        smtp = {
          enabled = false
        }
        snapshots = {
          externalEnabled = false
        }
        users = {
          viewersCanEdit = false
        }
      }
      grafanaMajorVersion = "10"
      publicNetworkAccess = "Enabled"
      zoneRedundancy      = "Disabled" # Not supported in Sweden Central yet
    }
  }

  response_export_values = ["properties.endpoint"]
}

data "azurerm_client_config" "current" {}

resource "azurerm_role_assignment" "current_grafana_admin" {
  scope                = azapi_resource.grafana.id
  role_definition_name = "Grafana Admin"
  principal_id         = data.azurerm_client_config.current.object_id
}

ephemeral "azidentity_azure_cli_credential" "grafana" {
  depends_on = [azurerm_role_assignment.current_grafana_admin]
  scopes     = ["ce34e7e5-485f-4d76-964f-b3d2b16d1e4f/.default"] # Microsofts Grafana Application ID
}

provider "grafana" {
  url  = azapi_resource.grafana.output.properties.endpoint
  auth = ephemeral.azidentity_azure_cli_credential.grafana.access_token
}

data "grafana_organization_preferences" "this" {
  depends_on = [azurerm_role_assignment.current_grafana_admin]
  org_id     = "1"
}

data "grafana_data_source" "azure_monitor" {
  depends_on = [azurerm_role_assignment.current_grafana_admin]
  name       = "Azure Monitor"
}

output "grafana" {
  value = {
    organization_preferences  = data.grafana_organization_preferences.this
    azure_monitor_data_source = data.grafana_data_source.azure_monitor
  }
}
