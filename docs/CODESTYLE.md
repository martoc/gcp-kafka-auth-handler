# Code Style

This project follows the [Google Go Style Guide](https://google.github.io/styleguide/go/guide).

## General Guidelines

- Use British English spelling in code, comments, and documentation (e.g., "marshalling" not "marshaling")
- The linter enforces UK locale for misspell checks
- Follow Effective Go guidelines
- Prefer interfaces for abstractions

## Project Structure

```
.
├── cmd/           # CLI commands (Cobra)
├── handler/       # HTTP handlers and authentication logic
├── docs/          # Documentation
└── integration-tests/  # Integration test scripts
```

## Naming Conventions

- Use descriptive names for functions and variables
- Use CamelCase for exported identifiers
- Use camelCase for unexported identifiers
- Interface names should describe behaviour (e.g., `GoogleService`, `AWSTokenGenerator`)

## Error Handling

- Use idiomatic Go error handling patterns
- Return errors to callers rather than logging and continuing
- Log errors at the appropriate level

## Testing

- Unit tests must be written using the `testing` package
- Assertions using the `testify` library
- Mocking using `gomock` from uber
- Add `go:generate` comments for mock generation

## Code Generation

Mocks are generated using mockgen. Each interface has a corresponding mock file:

- `gcp.go` → `gcp_mock.go`
- `aws.go` → `aws_mock.go`

Regenerate mocks with:
```bash
make generate
```

## Formatting

Code formatting is handled by:
- `gofmt` - Standard Go formatting
- `gofumpt` - Stricter formatting rules

Run formatting with:
```bash
make format
```

## Linting

The project uses `golangci-lint` with multiple linters enabled. Run linting with:
```bash
make lint
```

## Design Patterns

- Builder Pattern for handler construction (e.g., `NewGCPAuthHandlerBuilder`, `NewAWSAuthHandlerBuilder`)
- Factory Pattern for handler creation (e.g., `NewAuthHandler`)
- Interface-based abstractions for testability
