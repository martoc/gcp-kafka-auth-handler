# Usage

## Using as a Library

The primary purpose of this module is to be used as a library in any Go HTTP server. The `handler` package provides an `http.Handler` implementation that can be mounted on any route.

### Installation

```bash
go get github.com/martoc/kafka-auth-handler
```

### Basic Integration

```go
package main

import (
    "log"
    "net/http"
    "os"

    "github.com/martoc/kafka-auth-handler/handler"
)

func main() {
    // Create the auth handler based on provider
    provider := os.Getenv("PROVIDER") // "gcp" or "aws"
    region := os.Getenv("REGION")     // Required for AWS
    authHandler := handler.NewAuthHandler(provider, region)

    // Mount on your preferred route
    http.Handle("/oauth/token", authHandler)

    // Or use with your preferred router (e.g., gorilla/mux, chi, gin)
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

### Provider-Specific Handlers

#### GCP Handler

```go
authHandler := handler.NewGCPAuthHandlerBuilder().Build()
```

With custom Google service for testing:

```go
authHandler := handler.NewGCPAuthHandlerBuilder().
    WithGoogleService(myCustomGoogleService).
    Build()
```

#### AWS Handler

```go
authHandler := handler.NewAWSAuthHandlerBuilder().
    WithRegion("eu-central-1").
    Build()
```

With custom token generator for testing:

```go
authHandler := handler.NewAWSAuthHandlerBuilder().
    WithRegion("eu-central-1").
    WithTokenGenerator(myMockTokenGenerator).
    Build()
```

### Integration with Popular Routers

#### gorilla/mux

```go
import "github.com/gorilla/mux"

r := mux.NewRouter()
r.Handle("/oauth/token", handler.NewAuthHandler("gcp", ""))
```

#### chi

```go
import "github.com/go-chi/chi/v5"

r := chi.NewRouter()
r.Handle("/oauth/token", handler.NewAuthHandler("aws", "eu-central-1"))
```

#### gin

```go
import "github.com/gin-gonic/gin"

r := gin.Default()
authHandler := handler.NewAuthHandler("gcp", "")
r.GET("/oauth/token", gin.WrapH(authHandler))
```

### Response Format

Both GCP and AWS handlers return a JSON response with the following structure:

```json
{
  "access_token": "<header>.<claims>.<token>",
  "token_type": "Bearer",
  "expires_in": 3600
}
```

The `access_token` is a JWT-like token composed of:
- Base64-encoded header with type and algorithm
- Base64-encoded claims (exp, iss, iat, sub)
- Base64-encoded cloud provider access token

#### GCP Token Structure

| Field | Value |
|-------|-------|
| `alg` | `GOOG_OAUTH2_TOKEN` |
| `iss` | `Google` |
| `sub` | GCP service account email |

#### AWS Token Structure

| Field | Value |
|-------|-------|
| `alg` | `AWS_MSK_IAM` |
| `iss` | `AWS` |

## Standalone Server

The module also includes a standalone server for quick deployment.

### Environment Variables

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `PROVIDER` | No | `gcp` | Cloud provider: `gcp` or `aws` |
| `REGION` | Yes (AWS) | - | AWS region for MSK IAM token generation |

### From Source

```bash
make build

# GCP (default)
./target/builds/kafka-auth-handler-darwin-arm64 serve

# AWS
PROVIDER=aws REGION=eu-central-1 ./target/builds/kafka-auth-handler-darwin-arm64 serve
```

Binaries are built for darwin/linux on amd64/arm64 in `./target/builds/`.

### Docker

```bash
# GCP
docker pull martoc/kafka-auth-handler:latest
docker run -p 14293:14293 martoc/kafka-auth-handler:latest

# AWS
docker run -p 14293:14293 \
  -e PROVIDER=aws \
  -e REGION=eu-central-1 \
  martoc/kafka-auth-handler:latest
```

The standalone server listens on port 14293.

## API

### Get OAuth2 Token

```bash
curl http://localhost:14293/
```

## AWS MSK Configuration

When using AWS MSK with IAM authentication:

1. Use port **9098** for IAM authentication on your MSK bootstrap servers
2. Set `KAFKA_SECURITY_PROTOCOL=SASL_SSL`
3. Set `KAFKA_SASL_MECHANISM=OAUTHBEARER`
4. Configure your Kafka client to use `http://localhost:14293/` as the token endpoint
5. Ensure your pod/service has proper IAM permissions (via IRSA on EKS)

## Development

### Build

```bash
make build
```

### Run Tests

```bash
make test
```

### Run Linter

```bash
make lint
```

### Run Integration Tests

```bash
make run-integration-tests
```

### Install Dependencies

```bash
make install
```

### Regenerate Mocks

```bash
make generate
```

### Format Code

```bash
make format
```
