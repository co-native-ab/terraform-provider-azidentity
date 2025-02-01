// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ ephemeral.EphemeralResource = &ephemeralHttpRequest{}

func newEphemeralHttpRequest() ephemeral.EphemeralResource {
	return &ephemeralHttpRequest{}
}

type ephemeralHttpRequest struct {
	httpClient *http.Client
}

type ephemeralHttpRequestModel struct {
	RequestURL         types.String `tfsdk:"request_url"`
	RequestMethod      types.String `tfsdk:"request_method"`
	RequestBody        types.String `tfsdk:"request_body"`
	RequestHeaders     types.Map    `tfsdk:"request_headers"`
	ResponseBody       types.String `tfsdk:"response_body"`
	ResponseHeaders    types.Map    `tfsdk:"response_headers"`
	ResponseStatusCode types.Int32  `tfsdk:"response_status_code"`
	ContinueOnError    types.Bool   `tfsdk:"continue_on_error"`
	Timeout            types.String `tfsdk:"timeout"`
	Success            types.Bool   `tfsdk:"success"`
	Error              types.String `tfsdk:"error"`
}

func (r *ephemeralHttpRequest) Metadata(_ context.Context, req ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_http_request"
}

func (r *ephemeralHttpRequest) Schema(ctx context.Context, _ ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "HTTP Request Ephemeral Resource",
		Attributes: map[string]schema.Attribute{
			"request_url": schema.StringAttribute{
				MarkdownDescription: "The URL to send the HTTP request to.",
				Required:            true,
				Sensitive:           true,
			},
			"request_method": schema.StringAttribute{
				MarkdownDescription: "The HTTP method to use for the request.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						http.MethodConnect,
						http.MethodDelete,
						http.MethodGet,
						http.MethodHead,
						http.MethodOptions,
						http.MethodPatch,
						http.MethodPost,
						http.MethodPut,
						http.MethodTrace,
					),
				},
			},
			"request_body": schema.StringAttribute{
				MarkdownDescription: "The body of the HTTP request. Defaults to an empty body.",
				Optional:            true,
				Sensitive:           true,
			},
			"request_headers": schema.MapAttribute{
				MarkdownDescription: "The headers to include in the HTTP request.",
				ElementType:         types.StringType,
				Optional:            true,
				Sensitive:           true,
			},
			"continue_on_error": schema.BoolAttribute{
				MarkdownDescription: "ContinueOnError indicates whether to continue on error when the http request fails. The default is false.",
				Optional:            true,
			},
			"timeout": schema.StringAttribute{
				MarkdownDescription: "Timeout sets the maximum time allowed for the request to complete, the string is a possibly signed sequence of decimal numbers, each with optional fraction and a unit suffix, such as '300ms', '1.5h' or '2h45m'. Valid time units are 'ns', 'us' (or 'µs'), 'ms', 's', 'm', 'h'. The default is 30 seconds ('30s').",
				Optional:            true,
			},
			"response_body": schema.StringAttribute{
				MarkdownDescription: "The body of the HTTP response.",
				Sensitive:           true,
				Computed:            true,
			},
			"response_headers": schema.MapAttribute{
				MarkdownDescription: "The headers of the HTTP response.",
				ElementType:         types.StringType,
				Sensitive:           true,
				Computed:            true,
			},
			"response_status_code": schema.Int32Attribute{
				MarkdownDescription: "The status code of the HTTP response.",
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

func (p *ephemeralHttpRequest) Configure(ctx context.Context, req ephemeral.ConfigureRequest, resp *ephemeral.ConfigureResponse) {
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

	if provider.httpClient == nil {
		resp.Diagnostics.AddError("HTTP Client is not set", "HTTP Client is required to send HTTP requests")
		return
	}

	p.httpClient = provider.httpClient
}

func (r *ephemeralHttpRequest) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	var data ephemeralHttpRequestModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqTimeout := parseTimeout(ctx, data.Timeout)
	reqCtx, cancel := context.WithTimeout(ctx, reqTimeout)
	defer cancel()

	reqMethod := data.RequestMethod.ValueString()
	reqUrl := data.RequestURL.ValueString()
	var reqBody io.ReadCloser = http.NoBody
	if !data.RequestBody.IsNull() {
		reqBody = io.NopCloser(strings.NewReader(data.RequestBody.ValueString()))
	}

	httpReq, err := http.NewRequestWithContext(reqCtx, reqMethod, reqUrl, reqBody)
	if err != nil {
		if data.ContinueOnError.ValueBool() {
			data.Error = types.StringValue(err.Error())
			data.Success = types.BoolValue(false)
			resp.Diagnostics.Append(resp.Result.Set(ctx, &data)...)
			return
		}

		resp.Diagnostics.AddError("Failed to create HTTP request", err.Error())
		return
	}

	for k, v := range data.RequestHeaders.Elements() {
		if v.IsNull() {
			continue
		}

		vv, ok := v.(types.String)
		if !ok {
			continue
		}

		httpReq.Header.Set(k, vv.ValueString())
	}

	httpRes, err := r.httpClient.Do(httpReq)
	if err != nil {
		if data.ContinueOnError.ValueBool() {
			data.Error = types.StringValue(err.Error())
			data.Success = types.BoolValue(false)
			resp.Diagnostics.Append(resp.Result.Set(ctx, &data)...)
			return
		}

		resp.Diagnostics.AddError("Failed to send HTTP request", err.Error())
		return
	}

	defer httpRes.Body.Close()

	resBody, err := io.ReadAll(httpRes.Body)
	if err != nil {
		if data.ContinueOnError.ValueBool() {
			data.Error = types.StringValue(err.Error())
			data.Success = types.BoolValue(false)
			resp.Diagnostics.Append(resp.Result.Set(ctx, &data)...)
			return
		}

		resp.Diagnostics.AddError("Failed to read HTTP response body", err.Error())
		return
	}

	resHeaders := map[string]attr.Value{}
	for k, v := range httpRes.Header {
		resHeaders[k] = types.StringValue(strings.Join(v, ", "))
	}

	tflog.Debug(ctx, fmt.Sprintf("Body: %s", resBody))

	data.ResponseBody = types.StringValue(string(resBody))
	responseHeaders, diag := types.MapValue(types.StringType, resHeaders)
	if diag.HasError() {
		if data.ContinueOnError.ValueBool() {
			data.Error = types.StringValue("Failed to set response headers")
			data.Success = types.BoolValue(false)
			resp.Diagnostics.Append(resp.Result.Set(ctx, &data)...)
			return
		}

		resp.Diagnostics.Append(diag...)
		return
	}
	data.ResponseHeaders = responseHeaders
	data.ResponseStatusCode = types.Int32Value(int32(httpRes.StatusCode))
	data.Success = types.BoolValue(true)

	resp.Diagnostics.Append(resp.Result.Set(ctx, &data)...)
}
