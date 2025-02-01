package provider

import (
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestEphemeralDefaultCredentialEmpty(t *testing.T) {
	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_10_0),
		},
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testProtoV6ProviderFactoriesWithEcho(t, testNewTestCredentialFn(t)),
		Steps: []resource.TestStep{
			{
				Config: testEphemeralDefaultCredentialEmptyConfig(),
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
						tfjsonpath.New("data").AtMapKey("cloud"),
						knownvalue.Null(),
					),
					statecheck.ExpectKnownValue(
						"echo.this",
						tfjsonpath.New("data").AtMapKey("tenant_id"),
						knownvalue.Null(),
					),
					statecheck.ExpectKnownValue(
						"echo.this",
						tfjsonpath.New("data").AtMapKey("additionally_allowed_tenants"),
						knownvalue.Null(),
					),
					statecheck.ExpectKnownValue(
						"echo.this",
						tfjsonpath.New("data").AtMapKey("disable_instance_discovery"),
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
					statecheck.ExpectKnownValue(
						"echo.this",
						tfjsonpath.New("data").AtMapKey("timeout"),
						knownvalue.Null(),
					),
				},
			},
		},
	})
}

func TestEphemeralDefaultCredential(t *testing.T) {
	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_10_0),
		},
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testProtoV6ProviderFactoriesWithEcho(t, testNewTestCredentialFn(t)),
		Steps: []resource.TestStep{
			{
				Config: testEphemeralDefaultCredentialConfig(),
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
						tfjsonpath.New("data").AtMapKey("cloud"),
						knownvalue.StringExact("AzureGovernment"),
					),
					statecheck.ExpectKnownValue(
						"echo.this",
						tfjsonpath.New("data").AtMapKey("tenant_id"),
						knownvalue.StringExact("ze-tenant"),
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
						tfjsonpath.New("data").AtMapKey("disable_instance_discovery"),
						knownvalue.Bool(true),
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
					statecheck.ExpectKnownValue(
						"echo.this",
						tfjsonpath.New("data").AtMapKey("timeout"),
						knownvalue.StringExact("1s"),
					),
				},
			},
		},
	})
}

func TestEphemeralDefaultCredentialFailGetCredential(t *testing.T) {
	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_10_0),
		},
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testProtoV6ProviderFactoriesWithEcho(t, testNewGetCredentialFailureFn(t)),
		Steps: []resource.TestStep{
			{
				Config:      testEphemeralDefaultCredentialEmptyConfig(),
				ExpectError: regexp.MustCompile(`ze-get-credential-fn-error`),
			},
		},
	})
}

func TestEphemeralDefaultCredentialFailGetToken(t *testing.T) {
	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_10_0),
		},
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testProtoV6ProviderFactoriesWithEcho(t, testNewTestCredentialFailureFn(t)),
		Steps: []resource.TestStep{
			{
				Config: testEphemeralDefaultCredentialConfigContinueOnError(),
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

func TestEphemeralDefaultCredentialGetTokenTimeout(t *testing.T) {
	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_10_0),
		},
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testProtoV6ProviderFactoriesWithEcho(t, testNewGetCredentialTimeoutFn(t, 50*time.Millisecond)),
		Steps: []resource.TestStep{
			{
				Config: `
ephemeral "azidentity_default_credential" "this" {
	scopes  = ["ze-scope-1"]
	timeout = "10ms"
}

provider "echo" {
  data = ephemeral.azidentity_default_credential.this
}

resource "echo" "this" {}
`,
				ExpectError: regexp.MustCompile(`context deadline exceeded`),
			},
		},
	})
}

func testEphemeralDefaultCredentialEmptyConfig() string {
	return `
ephemeral "azidentity_default_credential" "this" {
	scopes = ["ze-scope-1"]
}

provider "echo" {
  data = ephemeral.azidentity_default_credential.this
}

resource "echo" "this" {}
`
}

func testEphemeralDefaultCredentialConfig() string {
	return `
ephemeral "azidentity_default_credential" "this" {
	cloud                        = "AzureGovernment"
	tenant_id                    = "ze-tenant"
	additionally_allowed_tenants = ["ze-additional-tenant-1", "ze-additional-tenant-2"]
	disable_instance_discovery   = true
	claims                       = "ze-claims"
	enable_cae                   = true
	scopes                       = ["ze-scope-1", "ze-scope-2"]
	continue_on_error            = true
	timeout                      = "1s"
}

provider "echo" {
  data = ephemeral.azidentity_default_credential.this
}

resource "echo" "this" {}
`
}

func testEphemeralDefaultCredentialConfigContinueOnError() string {
	return `
ephemeral "azidentity_default_credential" "this" {
	scopes            = ["ze-scope-1"]
	continue_on_error = true
}

provider "echo" {
  data = ephemeral.azidentity_default_credential.this
}

resource "echo" "this" {}
`
}
