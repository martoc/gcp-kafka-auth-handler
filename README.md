[![build](https://github.com/martoc/gcp-kafka-auth-handler/actions/workflows/main.yml/badge.svg)](https://github.com/martoc/gcp-kafka-auth-handler/actions/workflows/main.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/martoc/gcp-kafka-auth-handler)](https://goreportcard.com/report/github.com/martoc/gcp-kafka-auth-handler)
![Go Version](https://img.shields.io/github/go-mod-go-version/martoc/gcp-kafka-auth-handler/main)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

# gcp-kafka-auth-handler

A lightweight HTTP server that provides OAuth2 tokens for Kafka clients authenticating against Google Cloud Platform.

## Overview

This handler exchanges GCP default credentials for Kafka-compatible OAuth2 tokens. It runs as a sidecar or local service, listening on port 14293 and returning JWT-formatted access tokens.

## Quick Start

```bash
# Build
make build

# Run the server
./target/builds/gcp-kafka-auth-handler-darwin-arm64 serve
```

## Docker

```bash
docker pull martoc/gcp-kafka-auth-handler:latest
docker run -p 14293:14293 martoc/gcp-kafka-auth-handler:latest
```

## Documentation

- [Usage Guide](./docs/USAGE.md)
- [Code Style](./docs/CODESTYLE.md)
