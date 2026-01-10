package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-msk-iam-sasl-signer-go/signer"
)

const millisecondsPerSecond = 1000

//go:generate ${GOPATH}/bin/mockgen -source=aws.go -destination=aws_mock.go -package=handler

// AWSTokenGenerator provides AWS MSK IAM token generation.
type AWSTokenGenerator interface {
	GenerateAuthToken(ctx context.Context, region string) (string, int64, error)
}

// DefaultAWSTokenGenerator uses the aws-msk-iam-sasl-signer-go library.
type DefaultAWSTokenGenerator struct{}

// GenerateAuthToken generates an AWS MSK IAM authentication token.
//
//nolint:gocritic // Named results conflict with nonamedreturns rule.
func (g *DefaultAWSTokenGenerator) GenerateAuthToken(ctx context.Context, region string) (string, int64, error) {
	return signer.GenerateAuthToken(ctx, region)
}

// AWSAuthHandler handles AWS MSK IAM token requests for Kafka authentication.
type AWSAuthHandler struct {
	TokenGenerator AWSTokenGenerator
	Region         string
}

// AWSAuthHandlerBuilder builds AWSAuthHandler instances.
type AWSAuthHandlerBuilder struct {
	tokenGenerator AWSTokenGenerator
	region         string
}

// NewAWSAuthHandlerBuilder creates a new AWSAuthHandlerBuilder.
func NewAWSAuthHandlerBuilder() *AWSAuthHandlerBuilder {
	return &AWSAuthHandlerBuilder{}
}

// WithTokenGenerator sets the token generator for testing.
func (b *AWSAuthHandlerBuilder) WithTokenGenerator(generator AWSTokenGenerator) *AWSAuthHandlerBuilder {
	b.tokenGenerator = generator

	return b
}

// WithRegion sets the AWS region for token generation.
func (b *AWSAuthHandlerBuilder) WithRegion(region string) *AWSAuthHandlerBuilder {
	b.region = region

	return b
}

// Build creates the AWSAuthHandler.
func (b *AWSAuthHandlerBuilder) Build() *AWSAuthHandler {
	if b.tokenGenerator == nil {
		b.tokenGenerator = &DefaultAWSTokenGenerator{}
	}

	return &AWSAuthHandler{
		TokenGenerator: b.tokenGenerator,
		Region:         b.region,
	}
}

func (h *AWSAuthHandler) buildMessage(token string, expiryMs int64) ([]byte, error) {
	// Calculate expires_in in seconds from now.
	now := time.Now().UnixMilli()
	expiresInSeconds := (expiryMs - now) / millisecondsPerSecond

	if expiresInSeconds < 0 {
		expiresInSeconds = 0
	}

	// Build JWT-like structure matching GCP format for librdkafka compatibility.
	header := `{"typ": "JWT", "alg": "AWS_MSK_IAM"}`

	claims := fmt.Sprintf(`{"exp": %d, "iss": "AWS", "iat": %d}`,
		expiryMs/millisecondsPerSecond, time.Now().Unix())

	// The MSK token from the signer is already in the correct format.
	fullAccessToken := fmt.Sprintf("%s.%s.%s",
		b64Encode(header),
		b64Encode(claims),
		b64Encode(token))

	message := map[string]interface{}{
		"access_token": fullAccessToken,
		"token_type":   "Bearer",
		"expires_in":   int(expiresInSeconds),
	}

	return json.Marshal(message)
}

// ServeHTTP handles the HTTP request for AWS MSK IAM tokens.
func (h *AWSAuthHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	log.Println("Received request: ", request.Method, request.URL)

	ctx := request.Context()

	writer.Header().Set("Content-Type", "application/json")

	if h.Region == "" {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Println("REGION environment variable is required for MSK IAM authentication")

		return
	}

	// Generate MSK IAM auth token.
	token, expiryMs, err := h.TokenGenerator.GenerateAuthToken(ctx, h.Region)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Printf("Failed to generate MSK IAM auth token: %v", err)

		return
	}

	// Build OAUTHBEARER-compatible response.
	message, err := h.buildMessage(token, expiryMs)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Printf("Failed to build token response: %v", err)

		return
	}

	_, err = writer.Write(message)
	if err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}
