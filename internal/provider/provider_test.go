package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/echoprovider"
)

func testNew(t *testing.T, credGetFn getCredentialFn) func() provider.Provider {
	t.Helper()

	return func() provider.Provider {
		return &azidentityProvider{
			version:   "test",
			credGetFn: credGetFn,
		}
	}
}

func testProtoV6ProviderFactoriesWithEcho(t *testing.T, credGetFn getCredentialFn) map[string]func() (tfprotov6.ProviderServer, error) {
	return map[string]func() (tfprotov6.ProviderServer, error){
		"azidentity": providerserver.NewProtocol6WithError(testNew(t, credGetFn)()),
		"echo":       echoprovider.NewProviderServer(),
	}
}
