# Variables
BINARY_NAME=authula-playground
BUILD_DIR=build
SRC_DIR=.

.PHONY: all install build run clean test deps fmt lint help

# Default target
all: build

# Development setup
setup: install # Setup development environment
	@echo "Setting up development environment..."
	@curl -sSfL https://golangci-lint.run/install.sh | sh -s v2.13.2
	@go install github.com/air-verse/air@latest
	@echo "Development environment setup complete!"

# Build the application
build: ## Build the package (library)
	@echo "Building $(BINARY_NAME)..."
	@go build ./...
	@echo "Build complete!"

# Run the application
run:
	go run $(SRC_DIR)

# Clean build artifacts
clean:
	rm -rf $(BUILD_DIR)

# Run tests
test:
	go test ./...

# Install dependencies
deps:
	go mod tidy
	go mod download

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run

# Help
help:
	@echo "Available targets:"
	@echo "  build    - Build the application"
	@echo "  run      - Run the application"
	@echo "  test     - Run tests"
	@echo "  clean    - Clean build artifacts"
	@echo "  deps     - Install dependencies"
	@echo "  fmt      - Format code"
	@echo "  lint     - Lint code"
	@echo "  help     - Show this help message"
