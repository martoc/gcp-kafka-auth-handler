package handler

import (
	"log"
	"net/http"
	"strconv"
	"time"
)

// StartServer starts the HTTP server with provider-based authentication.
// provider: "gcp" or "aws" (defaults to "gcp" if empty).
// region: AWS region for MSK IAM authentication (required for AWS, ignored for GCP).
// port: The port to listen on.
func StartServer(provider, region string, port int) {
	addr := ":" + strconv.Itoa(port)
	log.Printf("Starting server listening at %s with provider %s", addr, provider)

	server := &http.Server{
		Addr:        addr,
		ReadTimeout: 5 * time.Second,
		Handler:     NewAuthHandler(provider, region),
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
