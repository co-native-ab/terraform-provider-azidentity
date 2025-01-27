package provider

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestEphemeralHttpRequestGet(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET method, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"ze-key": "ze-value"}`)) // nolint:errcheck
	}))

	defer server.Close()

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_10_0),
		},
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testProtoV6ProviderFactoriesWithEcho(t, testNewTestCredentialFn(t)),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
ephemeral "azidentity_http_request" "this" {
	request_url    = "%s"
	request_method = "GET"
}

provider "echo" {
  data = ephemeral.azidentity_http_request.this
}

resource "echo" "this" {}
`, server.URL),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"echo.this",
						tfjsonpath.New("data").AtMapKey("request_url"),
						knownvalue.StringExact(server.URL),
					),
					statecheck.ExpectKnownValue(
						"echo.this",
						tfjsonpath.New("data").AtMapKey("request_method"),
						knownvalue.StringExact("GET"),
					),
					statecheck.ExpectKnownValue(
						"echo.this",
						tfjsonpath.New("data").AtMapKey("response_body"),
						knownvalue.StringExact(`{"ze-key": "ze-value"}`),
					),
					statecheck.ExpectKnownValue(
						"echo.this",
						tfjsonpath.New("data").AtMapKey("response_status_code"),
						knownvalue.Int32Exact(http.StatusOK),
					),
					statecheck.ExpectKnownValue(
						"echo.this",
						tfjsonpath.New("data").AtMapKey("response_headers").AtMapKey("Content-Type"),
						knownvalue.StringExact("application/json"),
					),
					statecheck.ExpectKnownValue(
						"echo.this",
						tfjsonpath.New("data").AtMapKey("success"),
						knownvalue.Bool(true),
					),
				},
			},
		},
	})
}

func TestEphemeralHttpRequestPost(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST method, got %s", r.Method)
		}

		if r.Header.Get("Content-Type") != "application/json" {
			t.Fatalf("expected Content-Type header to be application/json, got %s", r.Header.Get("Content-Type"))
		}

		reqBody, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("failed to read request body: %v", err)
		}
		defer r.Body.Close()

		if string(reqBody) != `{"foo":"bar"}` {
			t.Fatalf("expected request body to be {\"foo\":\"bar\"}, got %s", string(reqBody))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"ze-key": "ze-value"}`)) // nolint:errcheck
	}))

	defer server.Close()

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_10_0),
		},
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testProtoV6ProviderFactoriesWithEcho(t, testNewTestCredentialFn(t)),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
ephemeral "azidentity_http_request" "this" {
	request_url     = "%s"
	request_method  = "POST"
	request_body    = jsonencode({foo = "bar"})
	request_headers = {
		"Content-Type" = "application/json"
	}
}

provider "echo" {
  data = ephemeral.azidentity_http_request.this
}

resource "echo" "this" {}
`, server.URL),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"echo.this",
						tfjsonpath.New("data").AtMapKey("request_url"),
						knownvalue.StringExact(server.URL),
					),
					statecheck.ExpectKnownValue(
						"echo.this",
						tfjsonpath.New("data").AtMapKey("request_method"),
						knownvalue.StringExact("POST"),
					),
					statecheck.ExpectKnownValue(
						"echo.this",
						tfjsonpath.New("data").AtMapKey("request_body"),
						knownvalue.StringExact(`{"foo":"bar"}`),
					),
					statecheck.ExpectKnownValue(
						"echo.this",
						tfjsonpath.New("data").AtMapKey("request_headers").AtMapKey("Content-Type"),
						knownvalue.StringExact("application/json"),
					),
					statecheck.ExpectKnownValue(
						"echo.this",
						tfjsonpath.New("data").AtMapKey("response_body"),
						knownvalue.StringExact(`{"ze-key": "ze-value"}`),
					),
					statecheck.ExpectKnownValue(
						"echo.this",
						tfjsonpath.New("data").AtMapKey("response_status_code"),
						knownvalue.Int32Exact(http.StatusOK),
					),
					statecheck.ExpectKnownValue(
						"echo.this",
						tfjsonpath.New("data").AtMapKey("response_headers").AtMapKey("Content-Type"),
						knownvalue.StringExact("application/json"),
					),
					statecheck.ExpectKnownValue(
						"echo.this",
						tfjsonpath.New("data").AtMapKey("success"),
						knownvalue.Bool(true),
					),
				},
			},
		},
	})
}

func TestEphemeralHttpRequestGetFail(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	serverURL := server.URL
	server.Close()

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_10_0),
		},
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testProtoV6ProviderFactoriesWithEcho(t, testNewTestCredentialFn(t)),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
ephemeral "azidentity_http_request" "this" {
	request_url    = "%s"
	request_method = "GET"
}

provider "echo" {
  data = ephemeral.azidentity_http_request.this
}

resource "echo" "this" {}
`, serverURL),
				ExpectError: regexp.MustCompile(`connection refused`),
			},
		},
	})
}

func TestEphemeralHttpRequestGetFailContinueOnError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	serverURL := server.URL
	server.Close()

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_10_0),
		},
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testProtoV6ProviderFactoriesWithEcho(t, testNewTestCredentialFn(t)),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
ephemeral "azidentity_http_request" "this" {
	request_url       = "%s"
	request_method    = "GET"
	continue_on_error = true
}

provider "echo" {
  data = ephemeral.azidentity_http_request.this
}

resource "echo" "this" {}
`, serverURL),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"echo.this",
						tfjsonpath.New("data").AtMapKey("request_url"),
						knownvalue.StringExact(server.URL),
					),
					statecheck.ExpectKnownValue(
						"echo.this",
						tfjsonpath.New("data").AtMapKey("request_method"),
						knownvalue.StringExact("GET"),
					),
					statecheck.ExpectKnownValue(
						"echo.this",
						tfjsonpath.New("data").AtMapKey("response_body"),
						knownvalue.Null(),
					),
					statecheck.ExpectKnownValue(
						"echo.this",
						tfjsonpath.New("data").AtMapKey("response_status_code"),
						knownvalue.Null(),
					),
					statecheck.ExpectKnownValue(
						"echo.this",
						tfjsonpath.New("data").AtMapKey("response_headers"),
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
						knownvalue.StringRegexp(regexp.MustCompile(`connection refused`)),
					),
				},
			},
		},
	})
}
