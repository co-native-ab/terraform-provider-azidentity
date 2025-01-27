package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestEphemeralEnvironmentVariable(t *testing.T) {
	t.Setenv("ZE_KEY", "ze-value")

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_10_0),
		},
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testProtoV6ProviderFactoriesWithEcho(t, testNewTestCredentialFn(t)),
		Steps: []resource.TestStep{
			{
				Config: `
ephemeral "azidentity_environment_variable" "this" {
	key = "ZE_KEY"
}

provider "echo" {
  data = ephemeral.azidentity_environment_variable.this
}

resource "echo" "this" {}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"echo.this",
						tfjsonpath.New("data").AtMapKey("key"),
						knownvalue.StringExact("ZE_KEY"),
					),
					statecheck.ExpectKnownValue(
						"echo.this",
						tfjsonpath.New("data").AtMapKey("value"),
						knownvalue.StringExact("ze-value"),
					),
				},
			},
		},
	})
}

func TestEphemeralEnvironmentVariableEmpty(t *testing.T) {
	t.Setenv("ZE_KEY", "")

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_10_0),
		},
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testProtoV6ProviderFactoriesWithEcho(t, testNewTestCredentialFn(t)),
		Steps: []resource.TestStep{
			{
				Config: `
ephemeral "azidentity_environment_variable" "this" {
	key = "ZE_KEY"
}

provider "echo" {
  data = ephemeral.azidentity_environment_variable.this
}

resource "echo" "this" {}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"echo.this",
						tfjsonpath.New("data").AtMapKey("key"),
						knownvalue.StringExact("ZE_KEY"),
					),
					statecheck.ExpectKnownValue(
						"echo.this",
						tfjsonpath.New("data").AtMapKey("value"),
						knownvalue.StringExact(""),
					),
				},
			},
		},
	})
}

func TestEphemeralEnvironmentVariableUnset(t *testing.T) {
	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_10_0),
		},
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testProtoV6ProviderFactoriesWithEcho(t, testNewTestCredentialFn(t)),
		Steps: []resource.TestStep{
			{
				Config: `
ephemeral "azidentity_environment_variable" "this" {
	key = "ZE_EMPTY_KEY"
}

provider "echo" {
  data = ephemeral.azidentity_environment_variable.this
}

resource "echo" "this" {}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"echo.this",
						tfjsonpath.New("data").AtMapKey("key"),
						knownvalue.StringExact("ZE_EMPTY_KEY"),
					),
					statecheck.ExpectKnownValue(
						"echo.this",
						tfjsonpath.New("data").AtMapKey("value"),
						knownvalue.Null(),
					),
				},
			},
		},
	})
}
