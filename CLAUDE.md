# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build Commands

```bash
make build          # Full build: clean, check, lint, test, cross-compile binaries
make lint           # Run golangci-lint
make test           # Run tests with coverage
make install        # Install development dependencies
make generate       # Regenerate mocks
make format         # Format code with gofmt and gofumpt
make run-integration-tests  # Run integration tests
```

Run a single test:
```bash
go test -v -run TestName ./handler/...
```

## Architecture

This is a GCP Kafka authentication handler CLI that provides OAuth2 tokens for Kafka clients authenticating against GCP.

### Structure

- `main.go` - Entry point, delegates to cmd package
- `cmd/` - Cobra CLI commands (serve, version)
- `handler/` - HTTP handler that exchanges GCP credentials for Kafka-compatible OAuth2 tokens

### How It Works

The `serve` command starts an HTTP server on port 14293. When a request is received, the `AuthHandler`:
1. Fetches GCP default credentials using `google.FindDefaultCredentials`
2. Constructs a JWT-like token with header, claims, and the GCP access token
3. Returns a JSON response with `access_token`, `token_type`, and `expires_in`

### Code Generation

Mocks are generated using mockgen. Regenerate with:
```bash
make generate
```

## Code Style

Follow the [Google Go Style Guide](https://google.github.io/styleguide/go/guide). Use British English spelling (e.g., "marshalling" not "marshaling"). The linter enforces UK locale for misspell checks.
