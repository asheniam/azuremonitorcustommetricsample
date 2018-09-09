package main

import (
	"fmt"
	"strings"
)

// Environment This captures the environment configuration for an Azure cloud
type Environment struct {
	aadLoginURL          string
	ingestionURL         string
	ingestionURLTemplate string
	armURL               string
}

// PublicEnvironmentName This is the environment name for public Azure
const PublicEnvironmentName string = "Public"

var (
	// Public Azure
	publicAzureEnvironment = Environment{
		aadLoginURL:          "https://login.microsoftonline.com",
		ingestionURL:         "https://monitoring.azure.com/",
		ingestionURLTemplate: "https://%s.monitoring.azure.com/",
		armURL:               "https://management.azure.com/",
	}
)

func getCurrentEnvironment(environmentName string) *Environment {
	if strings.EqualFold(environmentName, PublicEnvironmentName) {
		return &publicAzureEnvironment
	}

	logerrorf("Unknown environment: %s", environmentName)
	return nil
}

func (environment Environment) getAadLoginURL(TenantID string) string {
	return fmt.Sprintf("%s/%s/oauth2/token", environment.aadLoginURL, TenantID)
}
