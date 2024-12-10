# Description: Makefile for Go projects
SHA := $(shell git rev-parse --short HEAD)
BIN_NAME := $(shell basename `pwd`)
PACKAGES = ./...
TARGET := ./target
GOPATH := $(shell go env GOPATH)
GOOS := $(shell go env GOOS)
OS_LIST := darwin linux
ARCH_LIST := amd64 arm64
BIN_VERSION := $(shell echo $$TAG_VERSION)
ifeq ($(strip $(BIN_VERSION)),)
BIN_VERSION := 1.0.0
endif

.PHONY: all
all: clean init build run-integration-tests ## Run all targets

.PHONY: clean
clean: ## Clean the binary
	@echo "==> Cleaning..."
	rm -rf $(TARGET)

.PHONY: init
init: install ## Installing binaries
	@echo "==> Initialising..."
	go version

.PHONY: test
test: ## Run tests
	@mkdir -p $(TARGET)
	go test -coverprofile=$(TARGET)/coverage.out $(PACKAGES)
	go tool cover -html=$(TARGET)/coverage.out -o $(TARGET)/coverage.html

.PHONY: build
build: clean check lint test ## Cross platform build the binary
	@echo "==> Building..."
	@mkdir -p $(TARGET)/builds
	for GOOS in $(OS_LIST); do \
		for GOARCH in $(ARCH_LIST); do \
			GOOS=$$GOOS GOARCH=$$GOARCH CGO_ENABLED=0 go build -tags musl -ldflags "-s -w \
					-X github.com/martoc/$(BIN_NAME)/cmd.CLIVersion=$(BIN_VERSION) \
					-X github.com/martoc/$(BIN_NAME)/cmd.CLIOs=$$GOOS \
					-X github.com/martoc/$(BIN_NAME)/cmd.CLIArch=$$GOARCH \
					-X github.com/martoc/$(BIN_NAME)/cmd.CLISha=$(SHA) \
					-o $(TARGET)/builds/$(BIN_NAME)-$$GOOS-$$GOARCH main.go ; \
			chmod 755 $(TARGET)/builds/$(BIN_NAME)-$$GOOS-$$GOARCH ; \
		done ; \
	done

.PHONY: run-integration-tests
run-integration-tests: ## Run integration tests
	@echo "==> Running integration tests..."
	./integration-tests/run.sh

.PHONY: docs
docs: ## Run docs
	@echo "==> Running docs..."
	godoc -http=:6060

.PHONY: format
format: ## Run format files
	@echo "==> Running format..."
	go mod tidy
	go fmt $(PACKAGES)
	$(GOPATH)/bin/gofumpt -d .

.PHONY: generate
generate: ## Run source code generation
	@echo "==> Generating source files..."
	go generate $(PACKAGES)

.PHONY: install
install: ## Install development dependencies
	@echo "==> Installing dependencies..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.60.1
	go install github.com/golang/mock/mockgen@v1.6.0
	go install mvdan.cc/gofumpt@v0.6.0
	go install golang.org/x/tools/cmd/godoc@v0.12.0

.PHONY: lint
lint: ## Run linter
	@echo "==> Running linter..."
	$(GOPATH)/bin/golangci-lint run --timeout=5m $(PACKAGES)

.PHONY: check
check: ## Run checks
	@echo "==> Running checks..."
	go mod verify
	go vet -all $(PACKAGES)
