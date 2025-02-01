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

type runCommandFn func(ctx context.Context, stdout *bytes.Buffer, stderr *bytes.Buffer, extraEnv []string, name string, arg []string) error

func newRunCommandFn() runCommandFn {
	return func(ctx context.Context, stdout *bytes.Buffer, stderr *bytes.Buffer, extraEnv []string, name string, arg []string) error {
		cmd := exec.CommandContext(ctx, name, arg...)
		if len(extraEnv) > 0 {
			cmd.Env = append(cmd.Env, extraEnv...)
		}

		cmd.Stdout = stdout
		cmd.Stderr = stderr

		return cmd.Run()
	}
}

type ephemeralAzureCLIAccount struct {
	runCmdFn runCommandFn
}

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
				MarkdownDescription: "Indicates whether the Azure CLI account show command succeeded.",
				Computed:            true,
			},
			"error": schema.StringAttribute{
				MarkdownDescription: "Error message if the Azure CLI account show command failed.",
				Computed:            true,
			},
		},
	}
}

func (p *ephemeralAzureCLIAccount) Configure(ctx context.Context, req ephemeral.ConfigureRequest, resp *ephemeral.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	provider, ok := req.ProviderData.(*azidentityProvider)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected ProviderData Type",
			fmt.Sprintf("Expected *azidentityProvider, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	if provider.runCmdFn == nil {
		resp.Diagnostics.AddError("RunCommandFn is not set", "RunCommandFn is required to run the Azure CLI account show command")
		return
	}

	p.runCmdFn = provider.runCmdFn
}

func (r *ephemeralAzureCLIAccount) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	var data ephemeralAzureCLIAccountModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	extraEnv := []string{}
	if data.AzureConfigDir.ValueString() != "" {
		extraEnv = append(extraEnv, fmt.Sprintf("AZURE_CONFIG_DIR=%s", data.AzureConfigDir.ValueString()))
	}
	var stdoutBuf, stderrBuf bytes.Buffer
	executableName := "az"
	executableArgs := []string{"account", "show", "--output", "json"}

	err := r.runCmdFn(ctx, &stdoutBuf, &stderrBuf, extraEnv, executableName, executableArgs)
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
