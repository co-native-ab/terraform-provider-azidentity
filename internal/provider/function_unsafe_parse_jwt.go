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
		Summary:             "Parse JWT without validating anything or verifying the signature, only use for the purpose of checking claims but not for anything security related. Outputs a JSON string that can be decoded with jsondecode().",
		MarkdownDescription: "Parse JWT without validating anything or verifying the signature, only use for the purpose of checking claims but not for anything security related. Outputs a JSON string that can be decoded with jsondecode().",
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
