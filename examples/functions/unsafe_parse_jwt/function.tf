terraform {
  required_version = ">= 1.8.0"
  required_providers {
    azidentity = {
      source = "co-native-ab/azidentity"
    }
  }
}

locals {
  # Just a sample generated using the code in function_unsafe_parse_jwt_test.go
  jwt = "eyJhbGciOiJFUzM4NCIsImtpZCI6ImQyMWJlMGVmMzQzNzg5YzM2ZWEwYzBmNjlmOWRiZDRiM2JmYWU0ZmQ2OGUzM2E5NWM5ZWE2Y2RhZGE3MDlkMGQiLCJ0eXAiOiJKV1QifQ.eyJhdWQiOlsiemUtYXVkaWVuY2UiXSwiZXhwIjoxNzM3OTI2NzEwLCJpYXQiOjE3Mzc5MjY3MDAsImlzcyI6InplLWlzc3VlciIsIm5iZiI6MTczNzkyNjcwMCwic3ViIjoiemUtc3ViamVjdCIsInplLWNsYWltIjoiemUtdmFsdWUifQ.D26iVBChqPqMoJO3AP29yJDixflGY3KqszhsspOVbh3g1lMwZpNDFrskXE_TsA_JozMFDFnuz2g7iSidL9Hr64gwG8oRKa_LxnAyJNvii76lalu_c0PtvfCZFu2aRqt5"
}

output "claims_string_json" {
  description = "The claims of the JWT as a JSON string"
  value       = "JSON string: ${provider::azidentity::unsafe_parse_jwt(local.jwt)}"
}

output "claims_map" {
  description = "The claims of the JWT as a map"
  value       = jsondecode(provider::azidentity::unsafe_parse_jwt(local.jwt))
}

output "audience" {
  description = "The audience of the JWT"
  value       = jsondecode(provider::azidentity::unsafe_parse_jwt(local.jwt)).aud[0]
}
