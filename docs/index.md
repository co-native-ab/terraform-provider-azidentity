---
page_title: "Provider: azidentity"
description: |-
  The azidentity provider is used to integrate Terraform with Azure Identity–based ephemeral resources.
---

# azidentity Provider

The **azidentity** provider surfaces ephemeral resources to acquire and manage short-lived Azure credentials at runtime. This lets you securely manage tokens, environment variables, or other secret data during a Terraform run, without persisting them in the state file.

> For an overview of ephemeral resources in Terraform, refer to the official [Ephemeral Resources documentation](https://developer.hashicorp.com/terraform/language/resources/ephemeral).

---

## Azure Identity for Go

This provider builds on Microsoft’s [Azure Identity Go SDK](https://github.com/Azure/azure-sdk-for-go/tree/main/sdk/azidentity), which supports various Microsoft Entra ID (formerly Azure Active Directory) authentication flows. Currently, the **azidentity** Terraform provider implements ephemeral resources for:

- **DefaultAzureCredential**: A multi-step credential that can use environment variables, managed identities (in some hosting environments), Azure CLI logins, or developer CLI.
- **ClientSecretCredential**: Authenticates a service principal via a client secret.
- **ClientAssertionCredential**: Authenticates a service principal via a JWT-based assertion.
- **AzureCLICredential**: Authenticates via a logged-in Azure CLI session.

Using ephemeral resources for these credential types ensures credentials and tokens are never stored in Terraform state.

---

## Example Usage

Below is a brief example showcasing how to acquire an ephemeral Azure CLI–based credential for use with both Azure DevOps and the AzureRM provider. Thanks to Terraform’s ephemeral resources, the token is never written to Terraform state:

```hcl
terraform {
  required_version = ">= 1.10.0"
  required_providers {
    azidentity = {
      source = "registry.terraform.io/co-native-ab/azidentity"
    }
    azuredevops = {
      source = "microsoft/azuredevops"
    }
    azurerm = {
      source = "hashicorp/azurerm"
    }
  }
}

provider "azidentity" {
  # No config needed for azidentity
}

ephemeral "azidentity_azure_cli_credential" "this" {
  # Acquire a short-lived credential from the Azure CLI
  scopes = ["499b84ac-1321-427f-aa17-267ca6975798/.default"]
}

provider "azuredevops" {
  # Use ephemeral token for Azure DevOps
  org_service_url       = var.azure_devops_org_url
  personal_access_token = ephemeral.azidentity_azure_cli_credential.this.access_token
}

data "azuredevops_projects" "this" {}

ephemeral "azidentity_azure_cli_account" "this" {
  # Retrieve subscription & tenant from the local Azure CLI
}

provider "azurerm" {
  features {}
  subscription_id = ephemeral.azidentity_azure_cli_account.this.subscription_id
}

resource "azurerm_resource_group" "this" {
  name     = "example-rg"
  location = "East US"
}
```

### Key Points

- **Short-Lived**: Ephemeral credentials and tokens are open only for the duration of a Terraform run.
- **No State Storage**: Ephemeral data is never written to the Terraform state file.
- **Integration**: You can chain ephemeral credentials across providers, or use them with local values, data sources, or other resources.

---

## Next Steps

- Review the [Ephemeral Resources documentation](https://developer.hashicorp.com/terraform/language/resources/ephemeral) to understand how ephemeral blocks work.
- Learn more about the underlying Azure Identity Go SDK in [azidentity’s GitHub repository](https://github.com/Azure/azure-sdk-for-go/tree/main/sdk/azidentity).
- Explore ephemeral resources for these available Azure Identity credential types (DefaultAzureCredential, ClientSecretCredential, ClientAssertionCredential, AzureCLICredential). Consult the provider’s resource documentation for usage details.

By using ephemeral resources, you can dynamically acquire secure credentials at runtime—reducing secrets exposure and improving security for your Terraform workflows.

