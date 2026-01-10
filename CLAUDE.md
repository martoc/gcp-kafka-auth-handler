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

This is a Kafka authentication handler CLI that provides OAuth2 tokens for Kafka clients authenticating against Google Cloud Platform (GCP) or Amazon Web Services (AWS).

### Structure

- `main.go` - Entry point, delegates to cmd package
- `cmd/` - Cobra CLI commands (serve, version)
- `handler/` - HTTP handlers for cloud provider authentication
  - `handler.go` - Factory function to create provider-specific handlers
  - `gcp.go` - GCP handler using Google Application Default Credentials
  - `aws.go` - AWS handler using MSK IAM SASL signer
  - `server.go` - HTTP server startup logic

### How It Works

The `serve` command starts an HTTP server on port 14293. Based on the `PROVIDER` environment variable, it creates the appropriate handler:

#### GCP Provider (default)
1. Fetches GCP default credentials using `google.FindDefaultCredentials`
2. Constructs a JWT-like token with header, claims, and the GCP access token
3. Returns a JSON response with `access_token`, `token_type`, and `expires_in`

#### AWS Provider
1. Uses AWS MSK IAM SASL signer to generate authentication tokens
2. Requires `REGION` environment variable
3. Constructs a JWT-like token compatible with OAUTHBEARER mechanism
4. Returns a JSON response with `access_token`, `token_type`, and `expires_in`

### Environment Variables

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `PROVIDER` | No | `gcp` | Cloud provider: `gcp` or `aws` |
| `REGION` | Yes (AWS) | - | AWS region for MSK IAM token generation |

### Code Generation

Mocks are generated using mockgen. Regenerate with:
```bash
make generate
```

## Code Style

Follow the [Google Go Style Guide](https://google.github.io/styleguide/go/guide). Use British English spelling (e.g., "marshalling" not "marshaling"). The linter enforces UK locale for misspell checks.
