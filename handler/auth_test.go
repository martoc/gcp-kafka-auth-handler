package handler_test

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/martoc/gcp-kafka-auth-handler/handler"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var errExpectedError = errors.New("expected error")

func Test_NewAuthHandlerBuilder(t *testing.T) {
	t.Parallel()

	builder := handler.NewAuthHandlerBuilder()
	assert.NotNil(t, builder)
}

func TestAuthHandlerBuilder_Build(t *testing.T) {
	t.Parallel()

	// Given
	expectGoogleService := &handler.DefaultGoogleServiceImpl{}
	builder := handler.NewAuthHandlerBuilder()
	builder.WithGoogleService(expectGoogleService)

	// When
	service := builder.Build()

	// Then
	assert.NotNil(t, service)
}

func TestAuthHandlerBuilder_BuildWithGoogleServiceShouldCreateDefaultGoodleService(t *testing.T) {
	t.Parallel()

	// Given
	builder := handler.NewAuthHandlerBuilder()

	// When
	service := builder.Build()

	// Then
	assert.NotNil(t, service.GoogleService)
}

func TestAuthHandler_ServeHTTP(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Given
	request := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	expectedExpiry, _ := time.Parse(time.RFC3339, "2124-10-21T00:00:00Z")
	googleServiceMock := handler.NewMockGoogleService(ctrl)
	tokenSourceMock := handler.NewMockTokenSource(ctrl)
	googleServiceMock.EXPECT().FindDefaultCredentials(gomock.Any(), gomock.Any()).Return(&google.Credentials{
		ProjectID:   "project-id",
		TokenSource: tokenSourceMock,
		JSON:        []byte(`{"client_email": "client-email"}`),
	}, nil)
	tokenSourceMock.EXPECT().Token().Return(&oauth2.Token{
		AccessToken: "access-token",
		Expiry:      expectedExpiry,
	}, nil)

	builder := handler.NewAuthHandlerBuilder()
	builder.WithGoogleService(googleServiceMock)
	service := builder.Build()

	// When
	w := httptest.NewRecorder()
	service.ServeHTTP(w, request)
	resp := w.Result()

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	// Then
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
	assert.NotEmpty(t, body)

	var response map[string]interface{}
	_ = json.Unmarshal(body, &response)
	accessToken := response["access_token"]
	assert.NotEmpty(t, accessToken)
	assert.Equal(t, "Bearer", response["token_type"])
	assert.NotEmpty(t, response["expires_in"])

	accessTokenString, _ := accessToken.(string)
	tokenPart := strings.Split(accessTokenString, ".")

	assert.NotEmpty(t, tokenPart[0], "header is missing")
	assert.NotEmpty(t, tokenPart[1], "jwt is missing")
	assert.NotEmpty(t, tokenPart[2], "access-token is missing")

	var header map[string]string

	decodedJSON, _ := base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(tokenPart[0])
	_ = json.Unmarshal(decodedJSON, &header)
	assert.Equal(t, "JWT", header["typ"])
	assert.Equal(t, "GOOG_OAUTH2_TOKEN", header["alg"])

	var jwt map[string]interface{}

	decodedJSON, _ = base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(tokenPart[1])
	_ = json.Unmarshal(decodedJSON, &jwt)
	assert.NotEmpty(t, jwt["exp"])
	assert.NotEmpty(t, jwt["iat"])
	assert.Equal(t, "Google", jwt["iss"])
	assert.Equal(t, "client-email", jwt["sub"])
}

func TestAuthHandler_ServeHTTPErrorFeatchingGoogleCredentialsShouldReturnHttpInternalError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Given
	request := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	googleServiceMock := handler.NewMockGoogleService(ctrl)
	googleServiceMock.EXPECT().FindDefaultCredentials(gomock.Any(), gomock.Any()).Return(nil, errExpectedError)

	builder := handler.NewAuthHandlerBuilder()
	builder.WithGoogleService(googleServiceMock)
	service := builder.Build()

	// When
	w := httptest.NewRecorder()
	service.ServeHTTP(w, request)
	resp := w.Result()

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	// Then
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
	assert.Empty(t, body)
}

func TestAuthHandler_ServeHTTPGetTokenFails(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Given
	request := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	googleServiceMock := handler.NewMockGoogleService(ctrl)
	tokenSourceMock := handler.NewMockTokenSource(ctrl)
	googleServiceMock.EXPECT().FindDefaultCredentials(gomock.Any(), gomock.Any()).Return(&google.Credentials{
		ProjectID:   "project-id",
		TokenSource: tokenSourceMock,
		JSON:        []byte(`{"client_email": "client-email"}`),
	}, nil)
	tokenSourceMock.EXPECT().Token().Return(nil, errExpectedError)

	builder := handler.NewAuthHandlerBuilder()
	builder.WithGoogleService(googleServiceMock)
	service := builder.Build()

	// When
	w := httptest.NewRecorder()
	service.ServeHTTP(w, request)
	resp := w.Result()

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	// Then
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
	assert.Empty(t, body)
}
