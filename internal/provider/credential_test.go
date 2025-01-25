package provider

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

type testCredential struct {
	t *testing.T
}

func (c *testCredential) GetToken(ctx context.Context, options policy.TokenRequestOptions) (azcore.AccessToken, error) {
	c.t.Helper()

	return azcore.AccessToken{
		Token:     "ze-token",
		ExpiresOn: time.Date(2022, 1, 2, 3, 4, 5, 0, time.UTC),
	}, nil
}

var _ azcore.TokenCredential = (*testCredential)(nil)

func testNewTestCredentialFn(t *testing.T) getCredentialFn {
	t.Helper()

	return func(credType credentialType, cfg credentialConfig) (azcore.TokenCredential, error) {
		return &testCredential{
			t: t,
		}, nil
	}
}

func testNewGetCredentialFailureFn(t *testing.T) getCredentialFn {
	t.Helper()

	return func(credType credentialType, cfg credentialConfig) (azcore.TokenCredential, error) {
		return &testCredential{
			t: t,
		}, fmt.Errorf("ze-get-credential-fn-error")
	}
}

type testCredentialFailure struct {
	t *testing.T
}

func (c *testCredentialFailure) GetToken(ctx context.Context, options policy.TokenRequestOptions) (azcore.AccessToken, error) {
	c.t.Helper()

	return azcore.AccessToken{}, fmt.Errorf("ze-get-token-error")
}

var _ azcore.TokenCredential = (*testCredentialFailure)(nil)

func testNewTestCredentialFailureFn(t *testing.T) getCredentialFn {
	t.Helper()

	return func(credType credentialType, cfg credentialConfig) (azcore.TokenCredential, error) {
		return &testCredentialFailure{
			t: t,
		}, nil
	}
}
