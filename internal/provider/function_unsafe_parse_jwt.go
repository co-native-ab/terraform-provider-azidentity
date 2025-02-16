package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/lestrrat-go/jwx/v3/jwt"
)

var (
	_ function.Function = functionUnsafeParseJWT{}
)

func newFunctionUnsafeParseJWT() function.Function {
	return functionUnsafeParseJWT{}
}

type functionUnsafeParseJWT struct{}

func (r functionUnsafeParseJWT) Metadata(_ context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "unsafe_parse_jwt"
}

func (r functionUnsafeParseJWT) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		MarkdownDescription: "The `unsafe_parse_jwt` function parses a JSON Web Token (JWT) without validating its signature or verifying its authenticity. This function is useful for extracting and inspecting claims from a JWT but should **never** be used for security-sensitive operations.",
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:                "jwt",
				MarkdownDescription: "The JWT to parse.",
			},
		},
		Return: function.StringReturn{},
	}
}

func (r functionUnsafeParseJWT) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var data string

	resp.Error = function.ConcatFuncErrors(req.Arguments.Get(ctx, &data))

	if resp.Error != nil {
		return
	}

	parsedToken, err := jwt.ParseString(data, jwt.WithVerify(false), jwt.WithValidate(false))
	if err != nil {
		resp.Error = function.ConcatFuncErrors(resp.Error, function.NewFuncError(fmt.Sprintf("failed to parse JWT: %s", err)))
		return
	}

	tokenBytes, err := json.Marshal(parsedToken)
	if err != nil {
		resp.Error = function.ConcatFuncErrors(resp.Error, function.NewFuncError(fmt.Sprintf("failed to marshal token to JSON: %s", err)))
		return
	}

	resp.Error = function.ConcatFuncErrors(resp.Result.Set(ctx, string(tokenBytes)))
}
