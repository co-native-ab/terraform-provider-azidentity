package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ ephemeral.EphemeralResource = &ephemeralAzureCLICredential{}

func newEphemeralAzureCLICredential() ephemeral.EphemeralResource {
	return &ephemeralAzureCLICredential{}
}

type ephemeralAzureCLICredential struct {
	getCredFn getCredentialFn
}

type ephemeralAzureCLICredentialModel struct {
	TenantID                   types.String `tfsdk:"tenant_id"`
	SubscriptionID             types.String `tfsdk:"subscription_id"`
	AdditionallyAllowedTenants types.Set    `tfsdk:"additionally_allowed_tenants"`
	Claims                     types.String `tfsdk:"claims"`
	EnableCAE                  types.Bool   `tfsdk:"enable_cae"`
	Scopes                     types.Set    `tfsdk:"scopes"`
	ContinueOnError            types.Bool   `tfsdk:"continue_on_error"`
	Timeout                    types.String `tfsdk:"timeout"`
	AccessToken                types.String `tfsdk:"access_token"`
	ExpiresOn                  types.String `tfsdk:"expires_on"`
	Success                    types.Bool   `tfsdk:"success"`
	Error                      types.String `tfsdk:"error"`
}

func (r *ephemeralAzureCLICredentialModel) newCredentialConfig(ctx context.Context) credentialConfig {
	return credentialConfig{
		TenantID:                   r.TenantID.ValueString(),
		SubscriptionID:             r.SubscriptionID.ValueString(),
		AdditionallyAllowedTenants: typesSetToStringSlice(r.AdditionallyAllowedTenants),
		Claims:                     r.Claims.ValueString(),
		Scopes:                     typesSetToStringSlice(r.Scopes),
		ContinueOnError:            r.ContinueOnError.ValueBool(),
		Timeout:                    parseTimeout(ctx, r.Timeout),
	}
}

func (r *ephemeralAzureCLICredential) Metadata(_ context.Context, req ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_azure_cli_credential"
}

func (r *ephemeralAzureCLICredential) Schema(ctx context.Context, _ ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The `azidentity_azure_cli_credential` resource provides authentication using an active **Azure CLI session**. This allows Terraform to acquire tokens from the CLI without requiring stored credentials.",
		Attributes: map[string]schema.Attribute{
			"tenant_id": schema.StringAttribute{
				MarkdownDescription: "TenantID sets the default tenant for authentication via the Azure CLI and workload identity. The default is empty, use 'organizations' or 'common' if you can't provide one but required to use one.",
				Optional:            true,
			},
			"subscription_id": schema.StringAttribute{
				MarkdownDescription: "SubscriptionID is the ID (or name) of a subscription. Set this to acquire tokens for an account other than the Azure CLI's current account. The default is empty.",
				Optional:            true,
			},
			"additionally_allowed_tenants": schema.SetAttribute{
				MarkdownDescription: "AdditionallyAllowedTenants specifies tenants to which the credential may authenticate, in addition to TenantID. When TenantID is empty, this option has no effect and the credential will authenticate to any requested tenant. Add the wildcard value '*' to allow the credential to authenticate to any tenant. This value can also be set as a semicolon delimited list of tenants in the environment variable AZURE_ADDITIONALLY_ALLOWED_TENANTS. The default is an empty list.",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"claims": schema.StringAttribute{
				MarkdownDescription: "Claims are any additional claims required for the token to satisfy a conditional access policy, such as a service may return in a claims challenge following an authorization failure. If a service returned the claims value base64 encoded, it must be decoded before setting this field. The default is an empty string.",
				Optional:            true,
			},
			"enable_cae": schema.BoolAttribute{
				MarkdownDescription: "EnableCAE indicates whether to enable Continuous Access Evaluation (CAE) for the requested token. When true, azidentity credentials request CAE tokens for resource APIs supporting CAE. Clients are responsible for handling CAE challenges. If a client that doesn't handle CAE challenges receives a CAE token, it may end up in a loop retrying an API call with a token that has been revoked due to CAE. The default is false.",
				Optional:            true,
			},
			"scopes": schema.SetAttribute{
				MarkdownDescription: "Scopes contains the list of permission scopes required for the token. E.g. https://management.azure.com/.default for Azure Resource Manager or https://graph.microsoft.com/.default for Microsoft Graph.",
				Required:            true,
				ElementType:         types.StringType,
			},
			"continue_on_error": schema.BoolAttribute{
				MarkdownDescription: "ContinueOnError indicates whether to continue on error when acquiring a token. The default is false.",
				Optional:            true,
			},
			"timeout": schema.StringAttribute{
				MarkdownDescription: "Timeout sets the maximum time allowed for the request to complete, the string is a possibly signed sequence of decimal numbers, each with optional fraction and a unit suffix, such as '300ms', '1.5h' or '2h45m'. Valid time units are 'ns', 'us' (or 'µs'), 'ms', 's', 'm', 'h'. The default is 30 seconds ('30s').",
				Optional:            true,
			},
			"access_token": schema.StringAttribute{
				MarkdownDescription: "The issued access token.",
				Computed:            true,
				Sensitive:           true,
			},
			"expires_on": schema.StringAttribute{
				MarkdownDescription: "When the issued access token expires in RFC3339 format.",
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

func (p *ephemeralAzureCLICredential) Configure(ctx context.Context, req ephemeral.ConfigureRequest, resp *ephemeral.ConfigureResponse) {
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

	p.getCredFn = provider.getCredFn
}

func (r *ephemeralAzureCLICredential) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	var data ephemeralAzureCLICredentialModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	cfg := data.newCredentialConfig(ctx)
	token, errSummary, err := getToken(ctx, azureCLICredential, r.getCredFn, cfg)
	if err != nil && cfg.ContinueOnError {
		data.Error = types.StringValue(err.Error())
		data.Success = types.BoolValue(false)
		resp.Diagnostics.Append(resp.Result.Set(ctx, &data)...)
		return
	}

	if err != nil {
		resp.Diagnostics.AddError(errSummary, err.Error())
		return
	}

	data.AccessToken = types.StringValue(token.Token)
	data.ExpiresOn = types.StringValue(token.ExpiresOn.Format(time.RFC3339))
	data.Success = types.BoolValue(true)

	resp.Diagnostics.Append(resp.Result.Set(ctx, &data)...)
}
