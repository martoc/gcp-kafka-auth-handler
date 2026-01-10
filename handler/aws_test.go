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

	"github.com/martoc/kafka-auth-handler/handler"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

var errTokenGenerationFailed = errors.New("token generation failed")

func TestNewAWSAuthHandlerBuilder(t *testing.T) {
	t.Parallel()

	builder := handler.NewAWSAuthHandlerBuilder()
	assert.NotNil(t, builder)
}

func TestAWSAuthHandlerBuilder_Build(t *testing.T) {
	t.Parallel()

	// Given
	builder := handler.NewAWSAuthHandlerBuilder()
	builder.WithRegion("eu-central-1")

	// When
	service := builder.Build()

	// Then
	assert.NotNil(t, service)
	assert.Equal(t, "eu-central-1", service.Region)
}

func TestAWSAuthHandlerBuilder_BuildWithTokenGenerator(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Given
	mockGenerator := handler.NewMockAWSTokenGenerator(ctrl)
	builder := handler.NewAWSAuthHandlerBuilder()
	builder.WithTokenGenerator(mockGenerator)
	builder.WithRegion("eu-central-1")

	// When
	service := builder.Build()

	// Then
	assert.NotNil(t, service)
	assert.NotNil(t, service.TokenGenerator)
}

func TestAWSAuthHandlerBuilder_BuildWithDefaultTokenGenerator(t *testing.T) {
	t.Parallel()

	// Given
	builder := handler.NewAWSAuthHandlerBuilder()

	// When
	service := builder.Build()

	// Then
	assert.NotNil(t, service.TokenGenerator)
}

func TestAWSAuthHandler_ServeHTTP(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Given
	expectedExpiryMs := time.Now().Add(15 * time.Minute).UnixMilli()
	mockGenerator := handler.NewMockAWSTokenGenerator(ctrl)
	mockGenerator.EXPECT().GenerateAuthToken(gomock.Any(), "eu-central-1").Return("test-token", expectedExpiryMs, nil)

	builder := handler.NewAWSAuthHandlerBuilder()
	builder.WithTokenGenerator(mockGenerator)
	builder.WithRegion("eu-central-1")
	service := builder.Build()

	// When
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/token", http.NoBody)
	service.ServeHTTP(w, req)
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
	assert.NotEmpty(t, tokenPart[1], "claims is missing")
	assert.NotEmpty(t, tokenPart[2], "token is missing")

	var header map[string]string

	decodedJSON, _ := base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(tokenPart[0])
	_ = json.Unmarshal(decodedJSON, &header)
	assert.Equal(t, "JWT", header["typ"])
	assert.Equal(t, "AWS_MSK_IAM", header["alg"])

	var claims map[string]interface{}

	decodedJSON, _ = base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(tokenPart[1])
	_ = json.Unmarshal(decodedJSON, &claims)
	assert.NotEmpty(t, claims["exp"])
	assert.NotEmpty(t, claims["iat"])
	assert.Equal(t, "AWS", claims["iss"])
}

func TestAWSAuthHandler_ServeHTTPMissingRegion(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Given
	mockGenerator := handler.NewMockAWSTokenGenerator(ctrl)
	builder := handler.NewAWSAuthHandlerBuilder()
	builder.WithTokenGenerator(mockGenerator)
	// No region set
	service := builder.Build()

	// When
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/token", http.NoBody)
	service.ServeHTTP(w, req)
	resp := w.Result()

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	// Then
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
	assert.Empty(t, body)
}

func TestAWSAuthHandler_ServeHTTPTokenGenerationFails(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Given
	mockGenerator := handler.NewMockAWSTokenGenerator(ctrl)
	mockGenerator.EXPECT().GenerateAuthToken(gomock.Any(), "eu-central-1").Return("", int64(0), errTokenGenerationFailed)

	builder := handler.NewAWSAuthHandlerBuilder()
	builder.WithTokenGenerator(mockGenerator)
	builder.WithRegion("eu-central-1")
	service := builder.Build()

	// When
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/token", http.NoBody)
	service.ServeHTTP(w, req)
	resp := w.Result()

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	// Then
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
	assert.Empty(t, body)
}

func TestAWSAuthHandler_ServeHTTPExpiredToken(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Given - token already expired
	expiredExpiryMs := time.Now().Add(-5 * time.Minute).UnixMilli()
	mockGenerator := handler.NewMockAWSTokenGenerator(ctrl)
	mockGenerator.EXPECT().GenerateAuthToken(gomock.Any(), "eu-central-1").Return("expired-token", expiredExpiryMs, nil)

	builder := handler.NewAWSAuthHandlerBuilder()
	builder.WithTokenGenerator(mockGenerator)
	builder.WithRegion("eu-central-1")
	service := builder.Build()

	// When
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/token", http.NoBody)
	service.ServeHTTP(w, req)
	resp := w.Result()

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	// Then
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response map[string]interface{}

	_ = json.Unmarshal(body, &response)
	// expires_in should be 0 for expired tokens
	assert.Equal(t, float64(0), response["expires_in"])
}
