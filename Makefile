.PHONY: build run test clean install help

# Build the application
build:
	@echo "Building application..."
	go build -o ai-interviewer ./cmd/server

# Run the application
run:
	@echo "Starting server..."
	go run ./cmd/server

# Run all tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
test-cover:
	@echo "Running tests with coverage..."
	go test -cover ./...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f ai-interviewer
	rm -f ai_interviewer.db

# Install dependencies
install:
	@echo "Installing dependencies..."
	go mod download

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Run linter
lint:
	@echo "Running linter..."
	go vet ./...

# Help
help:
	@echo "Available commands:"
	@echo "  make build       - Build the application"
	@echo "  make run         - Run the application"
	@echo "  make test        - Run all tests"
	@echo "  make test-cover  - Run tests with coverage"
	@echo "  make clean       - Clean build artifacts"
	@echo "  make install     - Install dependencies"
	@echo "  make fmt         - Format code"
	@echo "  make lint        - Run linter"
	@echo "  make help        - Show this help message"
