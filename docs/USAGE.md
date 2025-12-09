# Usage

## Installation

### From Source

```bash
make build
```

Binaries are built for darwin/linux on amd64/arm64 in `./target/builds/`.

### Docker

```bash
docker pull martoc/gcp-kafka-auth-handler:latest
```

## Running the Server

```bash
# From binary
./target/builds/gcp-kafka-auth-handler-darwin-arm64 serve

# From Docker
docker run -p 14293:14293 martoc/gcp-kafka-auth-handler:latest
```

The server listens on port 14293.

## API

### Get OAuth2 Token

```bash
curl http://localhost:14293/
```

Response:

```json
{
  "access_token": "<header>.<jwt>.<gcp-access-token>",
  "token_type": "Bearer",
  "expires_in": 3600
}
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
