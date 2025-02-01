// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ ephemeral.EphemeralResource = &ephemeralAzureCLIAccount{}

func newEphemeralAzureCLIAccount() ephemeral.EphemeralResource {
	return &ephemeralAzureCLIAccount{}
}

type ephemeralAzureCLIAccount struct{}

type ephemeralAzureCLIAccountModel struct {
	AzureConfigDir  types.String `tfsdk:"azure_config_dir"`
	ContinueOnError types.Bool   `tfsdk:"continue_on_error"`
	SubscriptionID  types.String `tfsdk:"subscription_id"`
	TenantID        types.String `tfsdk:"tenant_id"`
	JsonResult      types.String `tfsdk:"json_result"`
	Success         types.Bool   `tfsdk:"success"`
	Error           types.String `tfsdk:"error"`
}

func (r *ephemeralAzureCLIAccount) Metadata(_ context.Context, req ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_azure_cli_account"
}

func (r *ephemeralAzureCLIAccount) Schema(ctx context.Context, _ ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Azure CLI Account Ephemeral Resource, which executes the `az account show` command to retrieve the subscription ID and tenant ID.",
		Attributes: map[string]schema.Attribute{
			"azure_config_dir": schema.StringAttribute{
				MarkdownDescription: "The directory where the Azure CLI configuration is stored. Default to not being set.",
				Optional:            true,
			},
			"continue_on_error": schema.BoolAttribute{
				MarkdownDescription: "ContinueOnError indicates whether to continue on error when the http request fails. The default is false.",
				Optional:            true,
			},
			"subscription_id": schema.StringAttribute{
				MarkdownDescription: "The subscription ID of the Azure account.",
				Computed:            true,
			},
			"tenant_id": schema.StringAttribute{
				MarkdownDescription: "The tenant ID of the Azure account.",
				Computed:            true,
			},
			"json_result": schema.StringAttribute{
				MarkdownDescription: "The JSON result of the Azure CLI account show command.",
				Computed:            true,
			},
			"success": schema.BoolAttribute{
				MarkdownDescription: "Indicates if a token was successfully acquired.",
				Computed:            true,
			},
			"error": schema.StringAttribute{
				MarkdownDescription: "Error message if acquiring a token failed.",
				Computed:            true,
			},
		},
	}
}

func (r *ephemeralAzureCLIAccount) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	var data ephemeralAzureCLIAccountModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	cmd := exec.CommandContext(ctx, "az", []string{"account", "show", "--output", "json"}...)
	if data.AzureConfigDir.ValueString() != "" {
		cmd.Env = append(cmd.Env, fmt.Sprintf("AZURE_CONFIG_DIR=%s", data.AzureConfigDir.ValueString()))
	}

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	err := cmd.Run()
	if err != nil {
		if data.ContinueOnError.ValueBool() {
			data.Error = types.StringValue(err.Error())
			data.Success = types.BoolValue(false)
			resp.Diagnostics.Append(resp.Result.Set(ctx, &data)...)
			return
		}

		resp.Diagnostics.AddError("Failed to run Azure CLI account show command", err.Error())
		return
	}

	if stderrBuf.Len() > 0 {
		if data.ContinueOnError.ValueBool() {
			data.Error = types.StringValue(stderrBuf.String())
			data.Success = types.BoolValue(false)
			resp.Diagnostics.Append(resp.Result.Set(ctx, &data)...)
			return
		}

		resp.Diagnostics.AddError("Failed to run Azure CLI account show command", stderrBuf.String())
		return
	}

	compacted, err := compactJSON(stdoutBuf.String())
	if err != nil {
		if data.ContinueOnError.ValueBool() {
			data.Error = types.StringValue(err.Error())
			data.Success = types.BoolValue(false)
			resp.Diagnostics.Append(resp.Result.Set(ctx, &data)...)
			return
		}

		resp.Diagnostics.AddError("Failed to compact JSON", err.Error())
		return
	}

	var account struct {
		SubscriptionID string `json:"id"`
		TenantID       string `json:"tenantId"`
	}

	err = json.Unmarshal([]byte(compacted), &account)
	if err != nil {
		if data.ContinueOnError.ValueBool() {
			data.Error = types.StringValue(err.Error())
			data.Success = types.BoolValue(false)
			resp.Diagnostics.Append(resp.Result.Set(ctx, &data)...)
			return
		}

		resp.Diagnostics.AddError("Failed to unmarshal JSON", err.Error())
		return
	}

	data.JsonResult = types.StringValue(compacted)
	data.SubscriptionID = types.StringValue(account.SubscriptionID)
	data.TenantID = types.StringValue(account.TenantID)
	data.Success = types.BoolValue(true)

	tflog.Debug(ctx, fmt.Sprintf("Azure CLI account succeeded:\njson_result=%s\nsubscription_id=%s\ntenant_id=%s\n", compacted, account.SubscriptionID, account.TenantID))

	resp.Diagnostics.Append(resp.Result.Set(ctx, &data)...)
}

func compactJSON(input string) (string, error) {
	var data interface{}
	if err := json.Unmarshal([]byte(input), &data); err != nil {
		return "", fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	out, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %w", err)
	}

	return string(out), nil
}
