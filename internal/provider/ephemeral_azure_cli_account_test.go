package provider

import (
	"bytes"
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestEphemeralAzureCLIAccount(t *testing.T) {
	getRunCmdFn := func() runCommandFn {
		return func(ctx context.Context, stdout *bytes.Buffer, stderr *bytes.Buffer, extraEnv []string, name string, arg []string) error {
			t.Helper()
			fmt.Fprintf(stdout, `{"id":"00000000-0000-0000-0000-000000000001","tenantId":"00000000-0000-0000-0000-000000000010","foo":"bar"}`)
			return nil
		}
	}
	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_10_0),
		},
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testProtoV6ProviderFactoriesWithEchoRunCommand(t, getRunCmdFn()),
		Steps: []resource.TestStep{
			{
				Config: `
ephemeral "azidentity_azure_cli_account" "this" {}

provider "echo" {
  data = ephemeral.azidentity_azure_cli_account.this
}

resource "echo" "this" {}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"echo.this",
						tfjsonpath.New("data").AtMapKey("azure_config_dir"),
						knownvalue.Null(),
					),
					statecheck.ExpectKnownValue(
						"echo.this",
						tfjsonpath.New("data").AtMapKey("continue_on_error"),
						knownvalue.Null(),
					),
					statecheck.ExpectKnownValue(
						"echo.this",
						tfjsonpath.New("data").AtMapKey("subscription_id"),
						knownvalue.StringExact("00000000-0000-0000-0000-000000000001"),
					),
					statecheck.ExpectKnownValue(
						"echo.this",
						tfjsonpath.New("data").AtMapKey("tenant_id"),
						knownvalue.StringExact("00000000-0000-0000-0000-000000000010"),
					),
					statecheck.ExpectKnownValue(
						"echo.this",
						tfjsonpath.New("data").AtMapKey("json_result"),
						knownvalue.StringRegexp(regexp.MustCompile(`"foo":"bar"`)),
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
				},
			},
		},
	})
}

func TestEphemeralAzureCLIAccountFail(t *testing.T) {
	getRunCmdFn := func() runCommandFn {
		return func(ctx context.Context, stdout *bytes.Buffer, stderr *bytes.Buffer, extraEnv []string, name string, arg []string) error {
			t.Helper()
			return fmt.Errorf("ze-error")
		}
	}
	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_10_0),
		},
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testProtoV6ProviderFactoriesWithEchoRunCommand(t, getRunCmdFn()),
		Steps: []resource.TestStep{
			{
				Config: `
ephemeral "azidentity_azure_cli_account" "this" {}

provider "echo" {
  data = ephemeral.azidentity_azure_cli_account.this
}

resource "echo" "this" {}
`,
				ExpectError: regexp.MustCompile(`ze-error`),
			},
		},
	})
}

func TestEphemeralAzureCLIAccountFailContinueOnError(t *testing.T) {
	getRunCmdFn := func() runCommandFn {
		return func(ctx context.Context, stdout *bytes.Buffer, stderr *bytes.Buffer, extraEnv []string, name string, arg []string) error {
			t.Helper()
			return fmt.Errorf("ze-error")
		}
	}
	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_10_0),
		},
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testProtoV6ProviderFactoriesWithEchoRunCommand(t, getRunCmdFn()),
		Steps: []resource.TestStep{
			{
				Config: `
ephemeral "azidentity_azure_cli_account" "this" {
  continue_on_error = true
}

provider "echo" {
  data = ephemeral.azidentity_azure_cli_account.this
}

resource "echo" "this" {}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"echo.this",
						tfjsonpath.New("data").AtMapKey("azure_config_dir"),
						knownvalue.Null(),
					),
					statecheck.ExpectKnownValue(
						"echo.this",
						tfjsonpath.New("data").AtMapKey("continue_on_error"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"echo.this",
						tfjsonpath.New("data").AtMapKey("subscription_id"),
						knownvalue.Null(),
					),
					statecheck.ExpectKnownValue(
						"echo.this",
						tfjsonpath.New("data").AtMapKey("tenant_id"),
						knownvalue.Null(),
					),
					statecheck.ExpectKnownValue(
						"echo.this",
						tfjsonpath.New("data").AtMapKey("json_result"),
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
						knownvalue.StringExact("ze-error"),
					),
				},
			},
		},
	})
}

func TestEphemeralAzureCLIAccountFailStderr(t *testing.T) {
	getRunCmdFn := func() runCommandFn {
		return func(ctx context.Context, stdout *bytes.Buffer, stderr *bytes.Buffer, extraEnv []string, name string, arg []string) error {
			t.Helper()
			fmt.Fprintf(stderr, `ze-stderr-error`)
			return nil
		}
	}
	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_10_0),
		},
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testProtoV6ProviderFactoriesWithEchoRunCommand(t, getRunCmdFn()),
		Steps: []resource.TestStep{
			{
				Config: `
ephemeral "azidentity_azure_cli_account" "this" {}

provider "echo" {
  data = ephemeral.azidentity_azure_cli_account.this
}

resource "echo" "this" {}
`,
				ExpectError: regexp.MustCompile(`ze-stderr-error`),
			},
		},
	})
}

func TestEphemeralAzureCLIAccountFailJsonParse(t *testing.T) {
	getRunCmdFn := func() runCommandFn {
		return func(ctx context.Context, stdout *bytes.Buffer, stderr *bytes.Buffer, extraEnv []string, name string, arg []string) error {
			t.Helper()
			fmt.Fprintf(stdout, `ze-invalid-json`)
			return nil
		}
	}
	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_10_0),
		},
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testProtoV6ProviderFactoriesWithEchoRunCommand(t, getRunCmdFn()),
		Steps: []resource.TestStep{
			{
				Config: `
ephemeral "azidentity_azure_cli_account" "this" {}

provider "echo" {
  data = ephemeral.azidentity_azure_cli_account.this
}

resource "echo" "this" {}
`,
				ExpectError: regexp.MustCompile(`failed to unmarshal JSON`),
			},
		},
	})
}
