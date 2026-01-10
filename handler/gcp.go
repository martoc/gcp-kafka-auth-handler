package handler

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/oauth2/google"
)

//go:generate ${GOPATH}/bin/mockgen -source=gcp.go -destination=gcp_mock.go -package=handler
//go:generate ${GOPATH}/bin/mockgen -destination=oauth2_mock.go -package=handler golang.org/x/oauth2 TokenSource

// GoogleService provides Google credential operations.
type GoogleService interface {
	FindDefaultCredentials(ctx context.Context, scopes ...string) (*google.Credentials, error)
}

// DefaultGoogleServiceImpl is the default implementation using Google's OAuth2 library.
type DefaultGoogleServiceImpl struct {
	GoogleService
}

// FindDefaultCredentials retrieves Google Application Default Credentials.
func (s *DefaultGoogleServiceImpl) FindDefaultCredentials(ctx context.Context, scopes ...string) (*google.Credentials, error) {
	return google.FindDefaultCredentials(ctx, scopes...)
}

// GCPAuthHandler handles GCP OAuth token requests for Kafka authentication.
type GCPAuthHandler struct {
	GoogleService GoogleService
}

// GCPAuthHandlerBuilder builds GCPAuthHandler instances.
type GCPAuthHandlerBuilder struct {
	googleService GoogleService
}

// NewGCPAuthHandlerBuilder creates a new GCPAuthHandlerBuilder.
func NewGCPAuthHandlerBuilder() *GCPAuthHandlerBuilder {
	return &GCPAuthHandlerBuilder{}
}

// WithGoogleService sets the GoogleService for testing.
func (b *GCPAuthHandlerBuilder) WithGoogleService(googleService GoogleService) *GCPAuthHandlerBuilder {
	b.googleService = googleService

	return b
}

// Build creates the GCPAuthHandler.
func (b *GCPAuthHandlerBuilder) Build() *GCPAuthHandler {
	if b.googleService == nil {
		b.googleService = &DefaultGoogleServiceImpl{}
	}

	return &GCPAuthHandler{
		GoogleService: b.googleService,
	}
}

func (*GCPAuthHandler) buildMessage(googleCreds *google.Credentials) ([]byte, error) {
	tokenSource, err := googleCreds.TokenSource.Token()
	if err != nil {
		return nil, err
	}

	header := `{"typ": "JWT", "alg": "GOOG_OAUTH2_TOKEN"}`

	var rawCredentials map[string]string

	err = json.Unmarshal(googleCreds.JSON, &rawCredentials)
	if err != nil {
		return nil, err
	}

	jwt := fmt.Sprintf(`{"exp": %d, "iss": "Google", "iat": %d, "sub": %q}`,
		tokenSource.Expiry.Unix(), time.Now().Unix(), rawCredentials["client_email"])

	fullAccessToken := fmt.Sprintf("%s.%s.%s", b64Encode(header), b64Encode(jwt), b64Encode(tokenSource.AccessToken))

	expirySeconds := int(time.Until(tokenSource.Expiry).Seconds())

	message := map[string]interface{}{
		"access_token": fullAccessToken,
		"token_type":   "Bearer",
		"expires_in":   expirySeconds,
	}

	return json.Marshal(message)
}

// ServeHTTP handles the HTTP request for GCP OAuth tokens.
func (h *GCPAuthHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	log.Println("Received request: ", request.Method, request.URL)

	ctx := context.Background()

	writer.Header().Set("Content-Type", "application/json")

	creds, err := h.GoogleService.FindDefaultCredentials(ctx, "https://www.googleapis.com/auth/cloud-platform") //nolint:contextcheck
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Println(err)

		return
	}

	message, err := h.buildMessage(creds)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Println(err)

		return
	}

	_, err = writer.Write(message)
	if err != nil {
		log.Println(err)

		return
	}
}

func b64Encode(source string) string {
	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString([]byte(source))
}

// AuthHandler is an alias for GCPAuthHandler.
//
// Deprecated: Use GCPAuthHandler and NewGCPAuthHandlerBuilder instead.
type AuthHandler = GCPAuthHandler

// NewAuthHandlerBuilder creates a new GCPAuthHandlerBuilder.
//
// Deprecated: Use NewGCPAuthHandlerBuilder instead.
func NewAuthHandlerBuilder() *GCPAuthHandlerBuilder {
	return NewGCPAuthHandlerBuilder()
}
