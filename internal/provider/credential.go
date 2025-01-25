package provider

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type credentialType string

const (
	defaultCredential      credentialType = "DefaultCredential"
	azureCLICredential     credentialType = "AzureCLICredential"
	clientSecretCredential credentialType = "ClientSecretCredential"
)

type credentialConfig struct {
	CloudConfig                cloud.Configuration `json:"cloud_config"`
	TenantID                   string              `json:"tenant_id"`
	ClientID                   string              `json:"client_id"`
	ClientSecret               string              `json:"client_secret"`
	SubscriptionID             string              `json:"subscription_id"`
	AdditionallyAllowedTenants []string            `json:"additionally_allowed_tenants"`
	DisableInstanceDiscovery   bool                `json:"disable_instance_discovery"`
	Claims                     string              `json:"claims"`
	EnableCAE                  bool                `json:"enable_cae"`
	Scopes                     []string            `json:"scopes"`
	ContinueOnError            bool                `json:"continue_on_error"`
}

type getCredentialFn func(credType credentialType, cfg credentialConfig) (azcore.TokenCredential, error)

func newGetcredentialFn() getCredentialFn {
	return func(credType credentialType, cfg credentialConfig) (azcore.TokenCredential, error) {
		switch credType {
		case defaultCredential:
			return newDefaultAzureCredential(cfg)
		case azureCLICredential:
			return newAzureCLICredential(cfg)
		case clientSecretCredential:
			return newClientSecretCredential(cfg)
		default:
			return nil, fmt.Errorf("unsupported credential type: %s", credType)
		}
	}
}

func newDefaultAzureCredential(cfg credentialConfig) (azcore.TokenCredential, error) {
	options := &azidentity.DefaultAzureCredentialOptions{
		TenantID:                   cfg.TenantID,
		AdditionallyAllowedTenants: cfg.AdditionallyAllowedTenants,
		DisableInstanceDiscovery:   cfg.DisableInstanceDiscovery,
		ClientOptions: azcore.ClientOptions{
			Cloud: cfg.CloudConfig,
		},
	}

	return azidentity.NewDefaultAzureCredential(options)
}

func newAzureCLICredential(cfg credentialConfig) (azcore.TokenCredential, error) {
	options := &azidentity.AzureCLICredentialOptions{
		AdditionallyAllowedTenants: cfg.AdditionallyAllowedTenants,
		Subscription:               cfg.SubscriptionID,
		TenantID:                   cfg.TenantID,
	}
	return azidentity.NewAzureCLICredential(options)
}

func newClientSecretCredential(cfg credentialConfig) (azcore.TokenCredential, error) {
	options := &azidentity.ClientSecretCredentialOptions{
		AdditionallyAllowedTenants: cfg.AdditionallyAllowedTenants,
		DisableInstanceDiscovery:   cfg.DisableInstanceDiscovery,
		ClientOptions: azcore.ClientOptions{
			Cloud: cfg.CloudConfig,
		},
	}

	tenantID := cfg.TenantID
	clientID := cfg.ClientID
	clientSecret := cfg.ClientSecret

	return azidentity.NewClientSecretCredential(tenantID, clientID, clientSecret, options)
}

func getToken(ctx context.Context, credType credentialType, getCredFn getCredentialFn, cfg credentialConfig) (azcore.AccessToken, string, error) {
	cred, err := getCredFn(credType, cfg)
	if err != nil {
		return azcore.AccessToken{}, "Error creating credential", err
	}

	tokenOpts := newTokenRequestOptions(cfg)
	token, err := cred.GetToken(ctx, tokenOpts)
	if err != nil {
		return azcore.AccessToken{}, "Error getting token", err
	}

	return token, "", nil
}

func newTokenRequestOptions(cfg credentialConfig) policy.TokenRequestOptions {
	return policy.TokenRequestOptions{
		Claims:    cfg.Claims,
		EnableCAE: cfg.EnableCAE,
		Scopes:    cfg.Scopes,
		TenantID:  cfg.TenantID,
	}
}

func typesSetToStringSlice(input types.Set) []string {
	result := []string{}
	for _, v := range input.Elements() {
		vv, ok := v.(basetypes.StringValue)
		if !ok {
			continue
		}
		result = append(result, vv.ValueString())
	}

	return result
}

func getCloudConfig(input string) cloud.Configuration {
	switch input {
	case "AzurePublic":
		return cloud.AzurePublic
	case "AzureChina":
		return cloud.AzureChina
	case "AzureGovernment":
		return cloud.AzureGovernment
	default:
		return cloud.AzurePublic
	}
}
