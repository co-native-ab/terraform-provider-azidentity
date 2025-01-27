terraform {
  required_version = ">= 1.10.0"
  required_providers {
    azuredevops = {
      source  = "microsoft/azuredevops"
      version = "1.6.0"
    }
    azidentity = {
      source = "registry.terraform.io/co-native-ab/azidentity"
    }
  }
}

provider "azidentity" {}

ephemeral "azidentity_environment_variable" "system_oidcrequesturi" {
  key = "SYSTEM_OIDCREQUESTURI"
}

ephemeral "azidentity_environment_variable" "system_accesstoken" {
  key = "SYSTEM_ACCESSTOKEN"
}

ephemeral "azidentity_environment_variable" "azuresubscription_service_connection_id" {
  key = "AZURESUBSCRIPTION_SERVICE_CONNECTION_ID"
}

ephemeral "azidentity_environment_variable" "azuresubscription_client_id" {
  key = "AZURESUBSCRIPTION_CLIENT_ID"
}

ephemeral "azidentity_environment_variable" "azuresubscription_tenant_id" {
  key = "AZURESUBSCRIPTION_TENANT_ID"
}

ephemeral "azidentity_environment_variable" "system_collectionuri" {
  key = "SYSTEM_COLLECTIONURI"
}

check "environment_variables" {
  assert {
    condition     = ephemeral.azidentity_environment_variable.system_oidcrequesturi.value != null
    error_message = "SYSTEM_OIDCREQUESTURI environment variable is required"
  }

  assert {
    condition     = ephemeral.azidentity_environment_variable.system_accesstoken.value != null
    error_message = "SYSTEM_ACCESSTOKEN environment variable is required"
  }

  assert {
    condition     = ephemeral.azidentity_environment_variable.azuresubscription_service_connection_id.value != null
    error_message = "AZURESUBSCRIPTION_SERVICE_CONNECTION_ID environment variable is required"
  }

  assert {
    condition     = ephemeral.azidentity_environment_variable.azuresubscription_client_id.value != null
    error_message = "AZURESUBSCRIPTION_CLIENT_ID environment variable is required"
  }

  assert {
    condition     = ephemeral.azidentity_environment_variable.azuresubscription_tenant_id.value != null
    error_message = "AZURESUBSCRIPTION_TENANT_ID environment variable is required"
  }

  assert {
    condition     = ephemeral.azidentity_environment_variable.system_collectionuri.value != null
    error_message = "SYSTEM_COLLECTIONURI environment variable is required"
  }
}

ephemeral "azidentity_http_request" "azure_devops_token" {
  request_url    = "${ephemeral.azidentity_environment_variable.system_oidcrequesturi.value}?api-version=7.1&serviceConnectionId=${ephemeral.azidentity_environment_variable.azuresubscription_service_connection_id.value}"
  request_method = "POST"
  request_headers = {
    "Content-Length" = "0"
    "Content-Type"   = "application/json"
    "Authorization"  = "Bearer ${ephemeral.azidentity_environment_variable.system_accesstoken.value}"
  }
}

locals {
  azure_devops_app_id = "499b84ac-1321-427f-aa17-267ca6975798" # Azure DevOps Application ID in Entra
  azure_devops_jwt    = jsondecode(ephemeral.azidentity_http_request.azure_devops_token.response_body).oidcToken
}

ephemeral "azidentity_client_assertion_credential" "this" {
  tenant_id = ephemeral.azidentity_environment_variable.azuresubscription_tenant_id.value
  client_id = ephemeral.azidentity_environment_variable.azuresubscription_client_id.value
  assertion = local.azure_devops_jwt
  scopes    = ["${local.azure_devops_app_id}/.default"]
}

provider "azuredevops" {
  org_service_url       = ephemeral.azidentity_environment_variable.system_collectionuri.value
  personal_access_token = ephemeral.azidentity_client_assertion_credential.this.access_token
}

data "azuredevops_projects" "this" {}

output "projects" {
  value = data.azuredevops_projects.this
}
