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

func TestEphemeralClientSecretCredentialEmpty(t *testing.T) {
	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_10_0),
		},
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testProtoV6ProviderFactoriesWithEcho(t, testNewTestCredentialFn(t)),
		Steps: []resource.TestStep{
			{
				Config: testEphemeralClientSecretCredentialEmptyConfig(),
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
						tfjsonpath.New("data").AtMapKey("client_secret"),
						knownvalue.StringExact("ze-secret"),
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

func TestEphemeralClientSecretCredential(t *testing.T) {
	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_10_0),
		},
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testProtoV6ProviderFactoriesWithEcho(t, testNewTestCredentialFn(t)),
		Steps: []resource.TestStep{
			{
				Config: testEphemeralClientSecretCredentialConfig(),
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
						tfjsonpath.New("data").AtMapKey("client_secret"),
						knownvalue.StringExact("ze-secret"),
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

func TestEphemeralClientSecretCredentialFailGetCredential(t *testing.T) {
	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_10_0),
		},
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testProtoV6ProviderFactoriesWithEcho(t, testNewGetCredentialFailureFn(t)),
		Steps: []resource.TestStep{
			{
				Config:      testEphemeralClientSecretCredentialEmptyConfig(),
				ExpectError: regexp.MustCompile(`ze-get-credential-fn-error`),
			},
		},
	})
}

func TestEphemeralClientSecretCredentialFailGetToken(t *testing.T) {
	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_10_0),
		},
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testProtoV6ProviderFactoriesWithEcho(t, testNewTestCredentialFailureFn(t)),
		Steps: []resource.TestStep{
			{
				Config: testEphemeralClientSecretCredentialConfigContinueOnError(),
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
						tfjsonpath.New("data").AtMapKey("client_secret"),
						knownvalue.StringExact("ze-secret"),
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

func TestEphemeralClientSecretCredentialGetTokenTimeout(t *testing.T) {
	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_10_0),
		},
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testProtoV6ProviderFactoriesWithEcho(t, testNewGetCredentialTimeoutFn(t, 50*time.Millisecond)),
		Steps: []resource.TestStep{
			{
				Config: `
ephemeral "azidentity_client_secret_credential" "this" {
	tenant_id     = "ze-tenant"
	client_id     = "ze-client"
	client_secret = "ze-secret"
    scopes        = ["ze-scope-1"]
	timeout 	  = "10ms"
}

provider "echo" {
  data = ephemeral.azidentity_client_secret_credential.this
}

resource "echo" "this" {}
`,
				ExpectError: regexp.MustCompile(`context deadline exceeded`),
			},
		},
	})
}

func testEphemeralClientSecretCredentialEmptyConfig() string {
	return `
ephemeral "azidentity_client_secret_credential" "this" {
	tenant_id     = "ze-tenant"
	client_id     = "ze-client"
	client_secret = "ze-secret"
    scopes        = ["ze-scope-1"]
}

provider "echo" {
  data = ephemeral.azidentity_client_secret_credential.this
}

resource "echo" "this" {}
`
}

func testEphemeralClientSecretCredentialConfig() string {
	return `
ephemeral "azidentity_client_secret_credential" "this" {
	tenant_id                    = "ze-tenant"
	client_id                    = "ze-client"
	client_secret                = "ze-secret"
	additionally_allowed_tenants = ["ze-additional-tenant-1", "ze-additional-tenant-2"]
	claims                       = "ze-claims"
	enable_cae                   = true
	scopes                       = ["ze-scope-1", "ze-scope-2"]
	continue_on_error            = true
	timeout                      = "1s"
}

provider "echo" {
  data = ephemeral.azidentity_client_secret_credential.this
}

resource "echo" "this" {}
`
}

func testEphemeralClientSecretCredentialConfigContinueOnError() string {
	return `
ephemeral "azidentity_client_secret_credential" "this" {
	tenant_id         = "ze-tenant"
	client_id         = "ze-client"
	client_secret     = "ze-secret"
	scopes            = ["ze-scope-1"]
	continue_on_error = true
}

provider "echo" {
  data = ephemeral.azidentity_client_secret_credential.this
}

resource "echo" "this" {}
`
}
