package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestEphemeralClientAssertionCredentialEmpty(t *testing.T) {
	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_10_0),
		},
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testProtoV6ProviderFactoriesWithEcho(t, testNewTestCredentialFn(t)),
		Steps: []resource.TestStep{
			{
				Config: testEphemeralClientAssertionCredentialEmptyConfig(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"echo.this",
						tfjsonpath.New("data").AtMapKey("access_token"),
						knownvalue.StringExact("ze-token"),
					),
					statecheck.ExpectKnownValue(
						"echo.this",
						tfjsonpath.New("data").AtMapKey("expires_on"),
						knownvalue.StringExact("2022-01-02T03:04:05Z"),
					),
					statecheck.ExpectKnownValue(
						"echo.this",
						tfjsonpath.New("data").AtMapKey("success"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"echo.this",
						tfjsonpath.New("data").AtMapKey("error"),
						knownvalue.Null(),
					),
					statecheck.ExpectKnownValue(
						"echo.this",
						tfjsonpath.New("data").AtMapKey("tenant_id"),
						knownvalue.StringExact("ze-tenant"),
					),
					statecheck.ExpectKnownValue(
						"echo.this",
						tfjsonpath.New("data").AtMapKey("client_id"),
						knownvalue.StringExact("ze-client"),
					),
					statecheck.ExpectKnownValue(
						"echo.this",
						tfjsonpath.New("data").AtMapKey("assertion"),
						knownvalue.StringExact("ze-assertion"),
					),
					statecheck.ExpectKnownValue(
						"echo.this",
						tfjsonpath.New("data").AtMapKey("additionally_allowed_tenants"),
						knownvalue.Null(),
					),
					statecheck.ExpectKnownValue(
						"echo.this",
						tfjsonpath.New("data").AtMapKey("claims"),
						knownvalue.Null(),
					),
					statecheck.ExpectKnownValue(
						"echo.this",
						tfjsonpath.New("data").AtMapKey("enable_cae"),
						knownvalue.Null(),
					),
					statecheck.ExpectKnownValue(
						"echo.this",
						tfjsonpath.New("data").AtMapKey("scopes"),
						knownvalue.SetExact([]knownvalue.Check{
							knownvalue.StringExact("ze-scope-1"),
						}),
					),
					statecheck.ExpectKnownValue(
						"echo.this",
						tfjsonpath.New("data").AtMapKey("continue_on_error"),
						knownvalue.Null(),
					),
				},
			},
		},
	})
}

func TestEphemeralClientAssertionCredential(t *testing.T) {
	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_10_0),
		},
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testProtoV6ProviderFactoriesWithEcho(t, testNewTestCredentialFn(t)),
		Steps: []resource.TestStep{
			{
				Config: testEphemeralClientAssertionCredentialConfig(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"echo.this",
						tfjsonpath.New("data").AtMapKey("access_token"),
						knownvalue.StringExact("ze-token"),
					),
					statecheck.ExpectKnownValue(
						"echo.this",
						tfjsonpath.New("data").AtMapKey("expires_on"),
						knownvalue.StringExact("2022-01-02T03:04:05Z"),
					),
					statecheck.ExpectKnownValue(
						"echo.this",
						tfjsonpath.New("data").AtMapKey("success"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"echo.this",
						tfjsonpath.New("data").AtMapKey("error"),
						knownvalue.Null(),
					),
					statecheck.ExpectKnownValue(
						"echo.this",
						tfjsonpath.New("data").AtMapKey("tenant_id"),
						knownvalue.StringExact("ze-tenant"),
					),
					statecheck.ExpectKnownValue(
						"echo.this",
						tfjsonpath.New("data").AtMapKey("client_id"),
						knownvalue.StringExact("ze-client"),
					),
					statecheck.ExpectKnownValue(
						"echo.this",
						tfjsonpath.New("data").AtMapKey("assertion"),
						knownvalue.StringExact("ze-assertion"),
					),
					statecheck.ExpectKnownValue(
						"echo.this",
						tfjsonpath.New("data").AtMapKey("additionally_allowed_tenants"),
						knownvalue.SetExact([]knownvalue.Check{
							knownvalue.StringExact("ze-additional-tenant-1"),
							knownvalue.StringExact("ze-additional-tenant-2"),
						}),
					),
					statecheck.ExpectKnownValue(
						"echo.this",
						tfjsonpath.New("data").AtMapKey("claims"),
						knownvalue.StringExact("ze-claims"),
					),
					statecheck.ExpectKnownValue(
						"echo.this",
						tfjsonpath.New("data").AtMapKey("enable_cae"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"echo.this",
						tfjsonpath.New("data").AtMapKey("scopes"),
						knownvalue.SetExact([]knownvalue.Check{
							knownvalue.StringExact("ze-scope-1"),
							knownvalue.StringExact("ze-scope-2"),
						}),
					),
					statecheck.ExpectKnownValue(
						"echo.this",
						tfjsonpath.New("data").AtMapKey("continue_on_error"),
						knownvalue.Bool(true),
					),
				},
			},
		},
	})
}

func TestEphemeralClientAssertionCredentialFailGetCredential(t *testing.T) {
	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_10_0),
		},
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testProtoV6ProviderFactoriesWithEcho(t, testNewGetCredentialFailureFn(t)),
		Steps: []resource.TestStep{
			{
				Config:      testEphemeralClientAssertionCredentialEmptyConfig(),
				ExpectError: regexp.MustCompile(`ze-get-credential-fn-error`),
			},
		},
	})
}

func TestEphemeralClientAssertionCredentialFailGetToken(t *testing.T) {
	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_10_0),
		},
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testProtoV6ProviderFactoriesWithEcho(t, testNewTestCredentialFailureFn(t)),
		Steps: []resource.TestStep{
			{
				Config: testEphemeralClientAssertionCredentialConfigContinueOnError(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"echo.this",
						tfjsonpath.New("data").AtMapKey("access_token"),
						knownvalue.Null(),
					),
					statecheck.ExpectKnownValue(
						"echo.this",
						tfjsonpath.New("data").AtMapKey("expires_on"),
						knownvalue.Null(),
					),
					statecheck.ExpectKnownValue(
						"echo.this",
						tfjsonpath.New("data").AtMapKey("success"),
						knownvalue.Bool(false),
					),
					statecheck.ExpectKnownValue(
						"echo.this",
						tfjsonpath.New("data").AtMapKey("error"),
						knownvalue.StringExact("ze-get-token-error"),
					),
					statecheck.ExpectKnownValue(
						"echo.this",
						tfjsonpath.New("data").AtMapKey("continue_on_error"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"echo.this",
						tfjsonpath.New("data").AtMapKey("tenant_id"),
						knownvalue.StringExact("ze-tenant"),
					),
					statecheck.ExpectKnownValue(
						"echo.this",
						tfjsonpath.New("data").AtMapKey("client_id"),
						knownvalue.StringExact("ze-client"),
					),
					statecheck.ExpectKnownValue(
						"echo.this",
						tfjsonpath.New("data").AtMapKey("assertion"),
						knownvalue.StringExact("ze-assertion"),
					),
					statecheck.ExpectKnownValue(
						"echo.this",
						tfjsonpath.New("data").AtMapKey("scopes"),
						knownvalue.SetExact([]knownvalue.Check{
							knownvalue.StringExact("ze-scope-1"),
						}),
					),
				},
			},
		},
	})
}

func testEphemeralClientAssertionCredentialEmptyConfig() string {
	return `
ephemeral "azidentity_client_assertion_credential" "this" {
	tenant_id = "ze-tenant"
	client_id = "ze-client"
	assertion = "ze-assertion"
    scopes    = ["ze-scope-1"]
}

provider "echo" {
  data = ephemeral.azidentity_client_assertion_credential.this
}

resource "echo" "this" {}
`
}

func testEphemeralClientAssertionCredentialConfig() string {
	return `
ephemeral "azidentity_client_assertion_credential" "this" {
	tenant_id                    = "ze-tenant"
	client_id                    = "ze-client"
	assertion                    = "ze-assertion"
	additionally_allowed_tenants = ["ze-additional-tenant-1", "ze-additional-tenant-2"]
	claims                       = "ze-claims"
	enable_cae                   = true
	scopes                       = ["ze-scope-1", "ze-scope-2"]
	continue_on_error            = true
}

provider "echo" {
  data = ephemeral.azidentity_client_assertion_credential.this
}

resource "echo" "this" {}
`
}

func testEphemeralClientAssertionCredentialConfigContinueOnError() string {
	return `
ephemeral "azidentity_client_assertion_credential" "this" {
	tenant_id         = "ze-tenant"
	client_id         = "ze-client"
	assertion         = "ze-assertion"
	scopes            = ["ze-scope-1"]
	continue_on_error = true
}

provider "echo" {
  data = ephemeral.azidentity_client_assertion_credential.this
}

resource "echo" "this" {}
`
}
