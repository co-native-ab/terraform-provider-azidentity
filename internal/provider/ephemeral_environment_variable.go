// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ ephemeral.EphemeralResource = &ephemeralEnvironmentVariable{}

func newEphemeralEnvironmentVariable() ephemeral.EphemeralResource {
	return &ephemeralEnvironmentVariable{}
}

type ephemeralEnvironmentVariable struct{}

type ephemeralEnvironmentVariableModel struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

func (r *ephemeralEnvironmentVariable) Metadata(_ context.Context, req ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_environment_variable"
}

func (r *ephemeralEnvironmentVariable) Schema(ctx context.Context, _ ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Environment Variable Ephemeral Resource",
		Attributes: map[string]schema.Attribute{
			"key": schema.StringAttribute{
				MarkdownDescription: "The key of the environment variable to get.",
				Required:            true,
			},
			"value": schema.StringAttribute{
				MarkdownDescription: "The value of the environment variable. Returns `null` if the environment variable is not set.",
				Sensitive:           true,
				Computed:            true,
			},
		},
	}
}

func (r *ephemeralEnvironmentVariable) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	var data ephemeralEnvironmentVariableModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	key := data.Key.ValueString()
	value, ok := os.LookupEnv(key)

	data.Value = types.StringNull()
	if ok {
		data.Value = types.StringValue(value)
	}

	resp.Diagnostics.Append(resp.Result.Set(ctx, &data)...)
}
