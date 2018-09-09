package main

import (
	"github.com/Azure/go-autorest/autorest/adal"
)

// AzureTokenAccessor Wrapper to grab Azure AAD token
type AzureTokenAccessor struct {
	environment *Environment
	oauthConfig *adal.OAuthConfig
}

// NewAzureTokenAccessor Create new AzureTokenAccessor
func NewAzureTokenAccessor(environment *Environment, tenantID string) *AzureTokenAccessor {
	oauthConfig, err := adal.NewOAuthConfig(environment.aadLoginURL, tenantID)
	if err != nil {
		logerrorf("Unable to initialize OAuth config: %v", err)
	}

	return &AzureTokenAccessor{
		environment: environment,
		oauthConfig: oauthConfig,
	}
}

// GetAccessToken Get access token to read or write custom metrics
func (tokenAccessor *AzureTokenAccessor) GetAccessToken(clientID string, clientSecret string, resourceName string) *adal.ServicePrincipalToken {
	token, err := adal.NewServicePrincipalToken(
		*tokenAccessor.oauthConfig,
		clientID,
		clientSecret,
		resourceName)

	if err != nil {
		logerrorf("Unable to get AAD token: %v", err)
	}

	return token
}
