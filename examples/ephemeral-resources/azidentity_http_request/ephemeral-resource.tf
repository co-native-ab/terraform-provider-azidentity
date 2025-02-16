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
