# azidentity: Ephemeral Azure Identity Credentials for Terraform

[![Terraform Registry](https://img.shields.io/badge/Terraform-Registry-blue.svg)](https://registry.terraform.io/providers/co-native-ab/azidentity)
[![License](https://img.shields.io/github/license/co-native-ab/terraform-provider-azidentity)](https://opensource.org/licenses/MIT)
[![GitHub Stars](https://img.shields.io/github/stars/co-native-ab/terraform-provider-azidentity?style=social)](https://github.com/co-native-ab/terraform-provider-azidentity)

The **azidentity** Terraform provider enables secure, ephemeral authentication for Azure by dynamically acquiring short-lived credentials at runtime. It supports a range of authentication methods without persisting secrets in the Terraform state file.

## 🌟 Features

- **Ephemeral credentials**: Acquire Azure authentication tokens dynamically during Terraform runs.
- **No secrets in state**: Tokens are never written to Terraform state, improving security.
- **Multiple credential types**: Supports DefaultAzureCredential, ClientSecretCredential, ClientAssertionCredential, AzureCLICredential, and more.
- **Seamless integration**: Works with AzureRM, Azure DevOps, and other Terraform providers.

---

## 🚀 Quick Start

### Install the Provider

```hcl
terraform {
  required_providers {
    azidentity = {
      source  = "co-native-ab/azidentity"
    }
  }
}

provider "azidentity" {}
```

### Acquire a Token Using Default Credentials

```hcl
ephemeral "azidentity_default_credential" "this" {
  scopes = ["https://management.azure.com/.default"]
}
```

### Use a Client Secret Credential

```hcl
ephemeral "azidentity_client_secret_credential" "this" {
  tenant_id     = "your-tenant-id"
  client_id     = "your-client-id"
  client_secret = "your-client-secret"
  scopes        = ["https://management.azure.com/.default"]
}
```

---

## 🔍 Supported Credential Types

| Credential Type               | Description                                                          |
| ----------------------------- | -------------------------------------------------------------------- |
| **DefaultAzureCredential**    | Uses environment variables, managed identities, or Azure CLI logins. |
| **ClientSecretCredential**    | Authenticates a service principal using a client secret.             |
| **ClientAssertionCredential** | Authenticates a service principal with a JWT assertion.              |
| **AzureCLICredential**        | Uses an active Azure CLI session.                                    |
| **HTTP Request**              | Performs HTTP request.                                               |
| **Environment Variable**      | Reads value from environment variables.                              |

---

## 📖 Documentation

- **[Provider Documentation](https://registry.terraform.io/providers/co-native-ab/azidentity/latest/docs)**
- **[Ephemeral Resources Overview](https://developer.hashicorp.com/terraform/language/resources/ephemeral)**
- **[Microsoft Azure Identity SDK](https://github.com/Azure/azure-sdk-for-go/tree/main/sdk/azidentity)**

---

## 💡 Why Use Ephemeral Credentials?

🔒 **Enhanced Security**: Credentials are used only during execution and never stored.

🛠 **Easy Integration**: Works with Terraform’s native ephemeral resource framework.

🚀 **Zero Configuration**: Default credentials auto-detect authentication methods.

---

## 🎯 Contributing

We welcome contributions! To get started:

1. Clone the repository: `git clone https://github.com/co-native-ab/terraform-provider-azidentity.git`
2. Install dependencies and build: `make build`
3. Run tests: `make test`

Create an issue and a pull request.

---

## 📜 License

This project is licensed under the [MIT License](LICENSE).

---

## 📢 Stay Updated

- **GitHub**: [co-native-ab/terraform-provider-azidentity](https://github.com/co-native-ab/terraform-provider-azidentity)
- **Terraform Registry**: [azidentity Provider](https://registry.terraform.io/providers/co-native-ab/azidentity)
- **Issues & Discussions**: [Open an Issue](https://github.com/co-native-ab/terraform-provider-azidentity/issues)

Happy Terraforming! 🚀
