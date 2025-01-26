package provider

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/lestrrat-go/jwx/v3/jwa"
	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/lestrrat-go/jwx/v3/jwt"
)

func testGetJWK(t *testing.T) (jwk.Key, jwk.Key) {
	t.Helper()

	ecdsaKey, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		t.Fatalf("failed to generate new ECDSA privatre key: %s", err)
	}

	key, err := jwk.Import(ecdsaKey)
	if err != nil {
		t.Fatalf("failed to import ECDSA private key: %s", err)
	}

	if _, ok := key.(jwk.ECDSAPrivateKey); !ok {
		t.Fatalf("expected jwk.ECDSAPrivateKey, got %T", key)
	}

	thumbprint, err := key.Thumbprint(crypto.SHA256)
	if err != nil {
		t.Fatalf("failed to compute thumbprint: %s", err)
	}

	keyID := fmt.Sprintf("%x", thumbprint)
	err = key.Set(jwk.KeyIDKey, keyID)
	if err != nil {
		t.Fatalf("failed to set key ID: %s", err)
	}

	pubKey, err := jwk.Import(ecdsaKey.PublicKey)
	if err != nil {
		t.Fatalf("failed to create public key: %s", err)
	}

	if _, ok := pubKey.(jwk.ECDSAPublicKey); !ok {
		t.Fatalf("expected jwk.ECDSAPublicKey, got %T", pubKey)
	}

	err = pubKey.Set(jwk.KeyIDKey, keyID)
	if err != nil {
		t.Fatalf("failed to set key ID: %s", err)
	}

	err = pubKey.Set(jwk.AlgorithmKey, jwa.ES384())
	if err != nil {
		t.Fatalf("failed to set algorithm: %s", err)
	}

	return key, pubKey
}

func testGetJWT(t *testing.T) string {
	t.Helper()

	key, _ := testGetJWK(t)

	token, err := jwt.NewBuilder().
		Issuer("ze-issuer").
		Audience([]string{"ze-audience"}).
		Subject("ze-subject").
		Claim("ze-claim", "ze-value").
		Expiration(time.Now().Add(10 * time.Second)).
		IssuedAt(time.Now()).
		NotBefore(time.Now()).
		Build()
	if err != nil {
		t.Fatalf("failed to build JWT: %s", err)
	}

	signedToken, err := jwt.Sign(token, jwt.WithKey(jwa.ES384(), key))
	if err != nil {
		t.Fatalf("failed to sign JWT: %s", err)
	}

	return string(signedToken)
}

func TestFunctionUnsafeParseJWT(t *testing.T) {
	tokenStr := testGetJWT(t)
	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_8_0),
		},
		ProtoV6ProviderFactories: testProtoV6ProviderFactoriesWithEcho(t, testNewTestCredentialFn(t)),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
				variable "token_string" {
					type    = string
					default = "%s"
			    }
				
				output "parsed_token_string" {
				    value = provider::azidentity::unsafe_parse_jwt(var.token_string)
				}
				
				output "audience" {
					value = jsondecode(provider::azidentity::unsafe_parse_jwt(var.token_string)).aud[0]
				}
				
				output "issuer" {
					value = jsondecode(provider::azidentity::unsafe_parse_jwt(var.token_string)).iss
				}

				output "subject" {
					value = jsondecode(provider::azidentity::unsafe_parse_jwt(var.token_string)).sub
				}

				output "ze_claim" {
					value = jsondecode(provider::azidentity::unsafe_parse_jwt(var.token_string))["ze-claim"]
				}
				`, tokenStr),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownOutputValue(
						"parsed_token_string",
						knownvalue.StringRegexp(regexp.MustCompile(`"aud":\["ze-audience"\]`)),
					),
					statecheck.ExpectKnownOutputValue(
						"audience",
						knownvalue.StringExact("ze-audience"),
					),
					statecheck.ExpectKnownOutputValue(
						"issuer",
						knownvalue.StringExact("ze-issuer"),
					),
					statecheck.ExpectKnownOutputValue(
						"subject",
						knownvalue.StringExact("ze-subject"),
					),
					statecheck.ExpectKnownOutputValue(
						"ze_claim",
						knownvalue.StringExact("ze-value"),
					),
				},
			},
		},
	})
}

func TestFunctionUnsafeParseJWT_Null(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_8_0),
		},
		ProtoV6ProviderFactories: testProtoV6ProviderFactoriesWithEcho(t, testNewTestCredentialFn(t)),
		Steps: []resource.TestStep{
			{
				Config: `
				output "test" {
					value = provider::azidentity::unsafe_parse_jwt(null)
				}
				`,
				ExpectError: regexp.MustCompile(`argument must not be null`),
			},
		},
	})
}

func TestFunctionUnsafeParseJWT_Invalid(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_8_0),
		},
		ProtoV6ProviderFactories: testProtoV6ProviderFactoriesWithEcho(t, testNewTestCredentialFn(t)),
		Steps: []resource.TestStep{
			{
				Config: `
				output "test" {
					value = provider::azidentity::unsafe_parse_jwt("ze-invalid-token")
				}
				`,
				ExpectError: regexp.MustCompile(`failed to parse string: unknown payload type`),
			},
		},
	})
}
