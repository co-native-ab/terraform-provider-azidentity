---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "azidentity_http_request Ephemeral Resource - azidentity"
subcategory: ""
description: |-
  The azidentity_http_request resource performs HTTP requests within Terraform. This allows retrieval of external authentication tokens or metadata required for Terraform execution.
---

# azidentity_http_request (Ephemeral Resource)

The `azidentity_http_request` resource performs HTTP requests within Terraform. This allows retrieval of external authentication tokens or metadata required for Terraform execution.

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

ephemeral "azidentity_http_request" "this" {
  request_url    = "https://ipinfo.io/127.0.0.1"
  request_method = "GET"
  request_headers = {
    "Accept" = "application/json"
  }
}

locals {
  response_body = jsondecode(ephemeral.azidentity_http_request.this.response_body)
}

check "ip" {
  assert {
    condition     = local.response_body.ip == "127.0.0.1"
    error_message = "The IP address is not 127.0.0.1"
  }

  assert {
    condition     = local.response_body.bogon == true
    error_message = "The IP address is not a bogon IP address"
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `request_method` (String) The HTTP method to use for the request.
- `request_url` (String, Sensitive) The URL to send the HTTP request to.

### Optional

- `continue_on_error` (Boolean) ContinueOnError indicates whether to continue on error when the http request fails. The default is false.
- `request_body` (String, Sensitive) The body of the HTTP request. Defaults to an empty body.
- `request_headers` (Map of String, Sensitive) The headers to include in the HTTP request.
- `timeout` (String) Timeout sets the maximum time allowed for the request to complete, the string is a possibly signed sequence of decimal numbers, each with optional fraction and a unit suffix, such as '300ms', '1.5h' or '2h45m'. Valid time units are 'ns', 'us' (or 'µs'), 'ms', 's', 'm', 'h'. The default is 30 seconds ('30s').

### Read-Only

- `error` (String) Error message if the HTTP request failed.
- `response_body` (String, Sensitive) The body of the HTTP response.
- `response_headers` (Map of String, Sensitive) The headers of the HTTP response.
- `response_status_code` (Number) The status code of the HTTP response.
- `success` (Boolean) Indicates if the HTTP request was successful.
