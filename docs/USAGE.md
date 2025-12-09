# Usage

## Using as a Library

The primary purpose of this module is to be used as a library in any Go HTTP server. The `handler` package provides an `http.Handler` implementation that can be mounted on any route.

### Installation

```bash
go get github.com/martoc/gcp-kafka-auth-handler
```

### Basic Integration

```go
package main

import (
    "log"
    "net/http"

    "github.com/martoc/gcp-kafka-auth-handler/handler"
)

func main() {
    // Create the auth handler
    authHandler := handler.NewAuthHandlerBuilder().Build()

    // Mount on your preferred route
    http.Handle("/oauth/token", authHandler)

    // Or use with your preferred router (e.g., gorilla/mux, chi, gin)
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

### With Custom Google Service

You can provide a custom `GoogleService` implementation for testing or custom credential handling:

```go
authHandler := handler.NewAuthHandlerBuilder().
    WithGoogleService(myCustomGoogleService).
    Build()
```

### Integration with Popular Routers

#### gorilla/mux

```go
import "github.com/gorilla/mux"

r := mux.NewRouter()
r.Handle("/oauth/token", handler.NewAuthHandlerBuilder().Build())
```

#### chi

```go
import "github.com/go-chi/chi/v5"

r := chi.NewRouter()
r.Handle("/oauth/token", handler.NewAuthHandlerBuilder().Build())
```

#### gin

```go
import "github.com/gin-gonic/gin"

r := gin.Default()
authHandler := handler.NewAuthHandlerBuilder().Build()
r.GET("/oauth/token", gin.WrapH(authHandler))
```

### Response Format

The handler returns a JSON response with the following structure:

```json
{
  "access_token": "<header>.<jwt>.<gcp-access-token>",
  "token_type": "Bearer",
  "expires_in": 3600
}
```

The `access_token` is a JWT-like token composed of:
- Base64-encoded header with type and algorithm
- Base64-encoded JWT claims (exp, iss, iat, sub)
- Base64-encoded GCP access token

## Standalone Server

The module also includes a standalone server for quick deployment.

### From Source

```bash
make build
./target/builds/gcp-kafka-auth-handler-darwin-arm64 serve
```

Binaries are built for darwin/linux on amd64/arm64 in `./target/builds/`.

### Docker

```bash
docker pull martoc/gcp-kafka-auth-handler:latest
docker run -p 14293:14293 martoc/gcp-kafka-auth-handler:latest
```

The standalone server listens on port 14293.

## API

### Get OAuth2 Token

```bash
curl http://localhost:14293/
```

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
