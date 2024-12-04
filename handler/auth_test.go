package handler_test

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

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
	googleServiceMock := handler.NewMockGoogleService(ctrl)
	tokenSourceMock := handler.NewMockTokenSource(ctrl)
	googleServiceMock.EXPECT().FindDefaultCredentials(gomock.Any(), gomock.Any()).Return(&google.Credentials{
		ProjectID:   "project-id",
		TokenSource: tokenSourceMock,
	}, nil)
	tokenSourceMock.EXPECT().Token().Return(&oauth2.Token{}, nil)

	builder := handler.NewAuthHandlerBuilder()
	builder.WithGoogleService(googleServiceMock)
	service := builder.Build()

	// When
	w := httptest.NewRecorder()
	service.ServeHTTP(w, nil)
	resp := w.Result()

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	// Then
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
	assert.NotEmpty(t, body)
}

func TestAuthHandler_ServeHTTPErrorFeatchingGoogleCredentialsShouldReturnHttpInternalError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Given
	googleServiceMock := handler.NewMockGoogleService(ctrl)
	googleServiceMock.EXPECT().FindDefaultCredentials(gomock.Any(), gomock.Any()).Return(nil, errExpectedError)

	builder := handler.NewAuthHandlerBuilder()
	builder.WithGoogleService(googleServiceMock)
	service := builder.Build()

	// When
	w := httptest.NewRecorder()
	service.ServeHTTP(w, nil)
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
	googleServiceMock := handler.NewMockGoogleService(ctrl)
	tokenSourceMock := handler.NewMockTokenSource(ctrl)
	googleServiceMock.EXPECT().FindDefaultCredentials(gomock.Any(), gomock.Any()).Return(&google.Credentials{
		ProjectID:   "project-id",
		TokenSource: tokenSourceMock,
	}, nil)
	tokenSourceMock.EXPECT().Token().Return(nil, errExpectedError)

	builder := handler.NewAuthHandlerBuilder()
	builder.WithGoogleService(googleServiceMock)
	service := builder.Build()

	// When
	w := httptest.NewRecorder()
	service.ServeHTTP(w, nil)
	resp := w.Result()

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	// Then
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
	assert.Empty(t, body)
}
