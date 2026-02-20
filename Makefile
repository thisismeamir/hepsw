.PHONY: build test clean install dev fmt lint

BINARY_NAME=hepsw
BUILD_DIR=build
VERSION=$(shell git describe --tags --always --dirty)
LDFLAGS=-ldflags "-X main.Version=$(VERSION)"

build:
	@echo "Building $(BINARY_NAME)..."
	@export CGO_ENABLED=1
	@go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/hepsw

install: build
	@echo "Installing $(BINARY_NAME)..."
	@cp $(BUILD_DIR)/$(BINARY_NAME) $(GOPATH)/bin/

test:
	@echo "Running tests..."
	@go test -v -race -coverprofile=coverage.out ./...

test-integration:
	@echo "Running integration tests..."
	@go test -v -tags=integration ./tests/integration/...

clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out

fmt:
	@echo "Formatting code..."
	@go fmt ./...

lint:
	@echo "Linting code..."
	@golangci-lint run

dev:
	@echo "Running in development mode..."
	@go run ./cmd/hepsw

.DEFAULT_GOAL := build
