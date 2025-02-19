---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "azidentity_client_assertion_credential Ephemeral Resource - azidentity"
subcategory: ""
description: |-
  The azidentity_client_assertion_credential resource supports authentication via a JWT assertion rather than a client secret. This is useful for scenarios where authentication tokens are issued dynamically or externally.
---

# azidentity_client_assertion_credential (Ephemeral Resource)

The `azidentity_client_assertion_credential` resource supports authentication via a **JWT assertion** rather than a client secret. This is useful for scenarios where authentication tokens are issued dynamically or externally.


## Example Usage

```terraform
terraform {
  required_version = ">= 1.10.0"
  required_providers {
    azidentity = {
      source = "co-native-ab/azidentity"
    }
  }
}

provider "azidentity" {}

ephemeral "azidentity_client_assertion_credential" "this" {
  tenant_id = "00000000-0000-0000-0000-000000000000"
  client_id = "00000000-0000-0000-0000-000000000000"
  assertion = "some-token"
  scopes    = ["https://management.azure.com/.default"]
}
```

## Azure DevOps Example

This is an example to show how to use the `azidentity_client_assertion_credential` resource with Azure DevOps terraform provider.

```terraform
terraform {
  required_version = ">= 1.10.0"
  required_providers {
    azuredevops = {
      source  = "microsoft/azuredevops"
      version = "1.6.0"
    }
    azidentity = {
      source = "co-native-ab/azidentity"
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
```

To use this, you will first need a Service Connection in Azure DevOps and use Workload identity federation.

An example pipeline could look like this (change `my-service-connection` to the name of your service connection):

```yaml
trigger:
  - main

pool:
  vmImage: ubuntu-latest

steps:
  - task: AzureCLI@2
    inputs:
      azureSubscription: "my-service-connection"
      scriptType: "bash"
      scriptLocation: "inlineScript"
      inlineScript: |
        set -e
        TEMP_DIR=$(mktemp -d)
        env --chdir=$TEMP_DIR curl --fail -L -o tenv_v4.1.0_Linux_x86_64.tar.gz https://github.com/tofuutils/tenv/releases/download/v4.1.0/tenv_v4.1.0_Linux_x86_64.tar.gz
        env --chdir=$TEMP_DIR tar xzvf tenv_v4.1.0_Linux_x86_64.tar.gz
        export PATH=$TEMP_DIR:$PATH
        tenv tf install 1.10.5
        tenv tf use 1.10.5
        export CI=true
        terraform init
        terraform plan
    env:
      SYSTEM_ACCESSTOKEN: $(System.AccessToken)
```


<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `assertion` (String, Sensitive) Assertion is a token (often JWT) assertion used to authenticate the client to the token service.
- `client_id` (String) ClientID is the application ID of the client.
- `scopes` (Set of String) Scopes contains the list of permission scopes required for the token. E.g. https://management.azure.com/.default for Azure Resource Manager or https://graph.microsoft.com/.default for Microsoft Graph.
- `tenant_id` (String) TenantID sets the default tenant for authentication via the Azure CLI and workload identity. Use 'organizations' or 'common' if you can't provide one but required to use one.

### Optional

- `additionally_allowed_tenants` (Set of String) AdditionallyAllowedTenants specifies tenants to which the credential may authenticate, in addition to TenantID. When TenantID is empty, this option has no effect and the credential will authenticate to any requested tenant. Add the wildcard value '*' to allow the credential to authenticate to any tenant. This value can also be set as a semicolon delimited list of tenants in the environment variable AZURE_ADDITIONALLY_ALLOWED_TENANTS. The default is an empty list.
- `claims` (String) Claims are any additional claims required for the token to satisfy a conditional access policy, such as a service may return in a claims challenge following an authorization failure. If a service returned the claims value base64 encoded, it must be decoded before setting this field. The default is an empty string.
- `cloud` (String) Cloud specifies a cloud for the client. The default is AzurePublic.
- `continue_on_error` (Boolean) ContinueOnError indicates whether to continue on error when acquiring a token. The default is false.
- `disable_instance_discovery` (Boolean) DisableInstanceDiscovery should be set true only by applications authenticating in disconnected clouds, or private clouds such as Azure Stack. It determines whether the credential requests Microsoft Entra instance metadata from https://login.microsoft.com before authenticating. Setting this to true will skip this request, making the application responsible for ensuring the configured authority is valid and trustworthy. The default is false.
- `enable_cae` (Boolean) EnableCAE indicates whether to enable Continuous Access Evaluation (CAE) for the requested token. When true, azidentity credentials request CAE tokens for resource APIs supporting CAE. Clients are responsible for handling CAE challenges. If a client that doesn't handle CAE challenges receives a CAE token, it may end up in a loop retrying an API call with a token that has been revoked due to CAE. The default is false.
- `timeout` (String) Timeout sets the maximum time allowed for the request to complete, the string is a possibly signed sequence of decimal numbers, each with optional fraction and a unit suffix, such as '300ms', '1.5h' or '2h45m'. Valid time units are 'ns', 'us' (or 'µs'), 'ms', 's', 'm', 'h'. The default is 30 seconds ('30s').

### Read-Only

- `access_token` (String, Sensitive) The issued access token.
- `error` (String) Error message if acquiring a token failed.
- `expires_on` (String) When the issued access token expires in RFC3339 format.
- `success` (Boolean) Indicates if a token was successfully acquired.