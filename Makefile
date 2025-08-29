# Makefile for latency-exporter

# Variables
BINARY_NAME=latency-exporter
BINARY_PATH=./cmd/latency-exporter
BUILD_DIR=./build
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS=-ldflags "-X main.version=$(VERSION)"

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Default target
.PHONY: all
all: clean test build

# Build the binary
.PHONY: build
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(BINARY_PATH)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Build for multiple platforms
.PHONY: build-all
build-all: build-linux build-darwin build-windows

.PHONY: build-linux
build-linux:
	@echo "Building for Linux..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(BINARY_PATH)

.PHONY: build-darwin
build-darwin:
	@echo "Building for macOS..."
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(BINARY_PATH)
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(BINARY_PATH)

.PHONY: build-windows
build-windows:
	@echo "Building for Windows..."
	@mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(BINARY_PATH)

# Install dependencies
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html

# Run the application locally
.PHONY: run
run: build
	@echo "Running $(BINARY_NAME)..."
	LATENCY_PARSER_CONFIG_PATH=./config/example.yml $(BUILD_DIR)/$(BINARY_NAME)

# Docker targets using Chainguard images
.PHONY: docker-build
docker-build:
	@echo "Building Docker image with Chainguard base..."
	docker build -t $(BINARY_NAME):$(VERSION) .

.PHONY: docker-build-minimal
docker-build-minimal:
	@echo "Building minimal Docker image (HTTP only, no ICMP)..."
	docker build -f Dockerfile.minimal -t $(BINARY_NAME):$(VERSION)-minimal .

.PHONY: docker-run
docker-run: docker-build
	@echo "Running Docker container..."
	docker run --rm -p 8080:8080 -v $(PWD)/config:/var/latency-parser $(BINARY_NAME):$(VERSION)

.PHONY: docker-run-minimal
docker-run-minimal: docker-build-minimal
	@echo "Running minimal Docker container..."
	docker run --rm -p 8080:8080 $(BINARY_NAME):$(VERSION)-minimal

.PHONY: docker-scan
docker-scan: docker-build
	@echo "Scanning Docker image for vulnerabilities..."
	@which trivy > /dev/null || (echo "trivy not installed, install with: brew install trivy" && exit 1)
	trivy image $(BINARY_NAME):$(VERSION)

# Development helpers
.PHONY: dev
dev:
	@echo "Running in development mode..."
	LATENCY_PARSER_CONFIG_PATH=./config/example.yml $(GOCMD) run $(BINARY_PATH)

.PHONY: fmt
fmt:
	@echo "Formatting code..."
	$(GOCMD) fmt ./...

.PHONY: vet
vet:
	@echo "Running go vet..."
	$(GOCMD) vet ./...

.PHONY: lint
lint:
	@echo "Running golangci-lint..."
	@which golangci-lint > /dev/null || (echo "golangci-lint not installed, installing..." && $(GOCMD) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run

# Check all (format, vet, lint, test)
.PHONY: check
check: fmt vet lint test

# Install the binary to GOPATH/bin
.PHONY: install
install:
	@echo "Installing $(BINARY_NAME)..."
	$(GOCMD) install $(LDFLAGS) $(BINARY_PATH)

# GitHub Actions local testing
.PHONY: act-build
act-build:
	@echo "Running GitHub Actions build workflow locally..."
	@which act > /dev/null || (echo "act not installed, install with: brew install act" && exit 1)
	act -j test

.PHONY: act-docker
act-docker:
	@echo "Running GitHub Actions docker workflow locally..."
	@which act > /dev/null || (echo "act not installed, install with: brew install act" && exit 1)
	act -j build-and-push

# Show help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all           - Clean, test, and build"
	@echo "  build         - Build the binary"
	@echo "  build-all     - Build for all platforms"
	@echo "  build-linux   - Build for Linux"
	@echo "  build-darwin  - Build for macOS"
	@echo "  build-windows - Build for Windows"
	@echo "  deps          - Install dependencies"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage"
	@echo "  clean         - Clean build artifacts"
	@echo "  run           - Build and run locally"
	@echo "  dev           - Run in development mode"
	@echo "  docker-build  - Build Docker image with Chainguard"
	@echo "  docker-build-minimal - Build minimal Docker image (HTTP only)"
	@echo "  docker-run    - Build and run Docker container"
	@echo "  docker-run-minimal - Build and run minimal Docker container"
	@echo "  docker-scan   - Scan Docker image for vulnerabilities"
	@echo "  fmt           - Format code"
	@echo "  vet           - Run go vet"
	@echo "  lint          - Run golangci-lint"
	@echo "  check         - Run fmt, vet, lint, and test"
	@echo "  install       - Install binary to GOPATH/bin"
	@echo "  act-build     - Test GitHub Actions build workflow locally"
	@echo "  act-docker    - Test GitHub Actions docker workflow locally"
	@echo "  help          - Show this help"
