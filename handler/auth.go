package handler

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"golang.org/x/oauth2/google"
)

//go:generate ${GOPATH}/bin/mockgen -source=auth.go -destination=auth_mock.go -package=handler
//go:generate ${GOPATH}/bin/mockgen -destination=oauth2_mock.go -package=handler golang.org/x/oauth2 TokenSource
type GoogleService interface {
	FindDefaultCredentials(ctx context.Context, scopes ...string) (*google.Credentials, error)
}

type DefaultGoogleServiceImpl struct {
	GoogleService
}

func (s *DefaultGoogleServiceImpl) FindDefaultCredentials(ctx context.Context, scopes ...string) (*google.Credentials, error) {
	return google.FindDefaultCredentials(ctx, scopes...)
}

type AuthHandler struct {
	GoogleService GoogleService
}

type AuthHandlerBuilder struct {
	googleService GoogleService
}

func NewAuthHandlerBuilder() *AuthHandlerBuilder {
	return &AuthHandlerBuilder{}
}

func (b *AuthHandlerBuilder) WithGoogleService(googleService GoogleService) *AuthHandlerBuilder {
	b.googleService = googleService

	return b
}

func (b *AuthHandlerBuilder) Build() *AuthHandler {
	if b.googleService == nil {
		b.googleService = &DefaultGoogleServiceImpl{}
	}

	return &AuthHandler{
		GoogleService: b.googleService,
	}
}

func (*AuthHandler) buildMessage(googleCreds *google.Credentials) ([]byte, error) {
	tokenSource, err := googleCreds.TokenSource.Token()
	if err != nil {
		return nil, err
	}

	header := `{"typ": "JWT", "alg": "GOOG_OAUTH2_TOKEN"}`

	jwt := fmt.Sprintf(`{"exp": %d, "iss": "%q", "iat": %d, "sub": "%q"}`,
		time.Now().Add(time.Hour).Unix(), "Google", time.Now().Unix(), googleCreds.ProjectID)

	fullAccessToken := fmt.Sprintf("%s.%s.%s", b64Encode(header), b64Encode(jwt), b64Encode(tokenSource.AccessToken))

	expirySeconds := int(time.Until(tokenSource.Expiry).Seconds())

	message := map[string]interface{}{
		"access_token": fullAccessToken,
		"token_type":   "Bearer",
		"expires_in":   expirySeconds,
	}

	return json.Marshal(message)
}

func (h *AuthHandler) ServeHTTP(writer http.ResponseWriter, _ *http.Request) {
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
	encoded := base64.URLEncoding.EncodeToString([]byte(source))

	return strings.TrimRight(encoded, "=")
}
