package provider

import (
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/echoprovider"
)

func testNew(t *testing.T, getCredFn getCredentialFn) func() provider.Provider {
	t.Helper()

	return func() provider.Provider {
		return &azidentityProvider{
			version:    "test",
			getCredFn:  getCredFn,
			httpClient: &http.Client{},
		}
	}
}

func testNewRunCommand(t *testing.T, runCmdFn runCommandFn) func() provider.Provider {
	t.Helper()

	return func() provider.Provider {
		return &azidentityProvider{
			version:  "test",
			runCmdFn: runCmdFn,
		}
	}
}

func testProtoV6ProviderFactoriesWithEcho(t *testing.T, getCredFn getCredentialFn) map[string]func() (tfprotov6.ProviderServer, error) {
	return map[string]func() (tfprotov6.ProviderServer, error){
		"azidentity": providerserver.NewProtocol6WithError(testNew(t, getCredFn)()),
		"echo":       echoprovider.NewProviderServer(),
	}
}

func testProtoV6ProviderFactoriesWithEchoRunCommand(t *testing.T, runCmdFn runCommandFn) map[string]func() (tfprotov6.ProviderServer, error) {
	return map[string]func() (tfprotov6.ProviderServer, error){
		"azidentity": providerserver.NewProtocol6WithError(testNewRunCommand(t, runCmdFn)()),
		"echo":       echoprovider.NewProviderServer(),
	}
}
