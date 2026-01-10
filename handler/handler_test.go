package handler_test

import (
	"testing"

	"github.com/martoc/kafka-auth-handler/handler"
	"github.com/stretchr/testify/assert"
)

func TestNewAuthHandler_GCP(t *testing.T) {
	t.Parallel()

	// When
	h := handler.NewAuthHandler("gcp", "europe-west2")

	// Then
	assert.NotNil(t, h)
	_, ok := h.(*handler.GCPAuthHandler)
	assert.True(t, ok, "Expected GCPAuthHandler for gcp provider")
}

func TestNewAuthHandler_AWS(t *testing.T) {
	t.Parallel()

	// When
	h := handler.NewAuthHandler("aws", "eu-central-1")

	// Then
	assert.NotNil(t, h)
	_, ok := h.(*handler.AWSAuthHandler)
	assert.True(t, ok, "Expected AWSAuthHandler for aws provider")
}

func TestNewAuthHandler_DefaultToGCP(t *testing.T) {
	t.Parallel()

	// When
	h := handler.NewAuthHandler("", "")

	// Then
	assert.NotNil(t, h)
	_, ok := h.(*handler.GCPAuthHandler)
	assert.True(t, ok, "Expected GCPAuthHandler as default")
}

func TestNewAuthHandler_UnknownProviderDefaultsToGCP(t *testing.T) {
	t.Parallel()

	// When
	h := handler.NewAuthHandler("unknown", "")

	// Then
	assert.NotNil(t, h)
	_, ok := h.(*handler.GCPAuthHandler)
	assert.True(t, ok, "Expected GCPAuthHandler for unknown provider")
}
