package handler

import (
	"net/http"
)

// Provider constants for supported cloud providers.
const (
	ProviderGCP = "gcp"
	ProviderAWS = "aws"
)

// NewAuthHandler creates the appropriate auth handler based on provider.
// For GCP: uses Google Application Default Credentials.
// For AWS: uses IAM Roles for Service Accounts (IRSA) with the specified region.
// If provider is empty or unknown, defaults to GCP for backwards compatibility.
func NewAuthHandler(provider, region string) http.Handler {
	switch provider {
	case ProviderAWS:
		return NewAWSAuthHandlerBuilder().
			WithRegion(region).
			Build()
	default:
		return NewGCPAuthHandlerBuilder().Build()
	}
}
