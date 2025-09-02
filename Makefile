# Makefile for PCC-TimeBot Development

.PHONY: help build run test clean dev db-setup db-reset lint format

# Default target
help: ## Show this help message
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

# Build the application
build: ## Build the Go application
	go build -o bin/pcc-timebot .

# Run the application with hot reload
dev: ## Run the application in development mode with hot reload
	air

# Run the application normally
run: ## Run the application
	go run .

# Run tests
test: ## Run all tests
	go test ./...

# Run tests with verbose output
test-verbose: ## Run tests with verbose output
	go test -v ./...

# Run tests with coverage
test-coverage: ## Run tests with coverage report
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run database tests only
test-db: ## Run database tests only
	DATABASE_URL=postgresql://postgres@localhost:5432/timebot_test?sslmode=disable go test -v ./src/database

# Clean build artifacts
clean: ## Clean build artifacts
	rm -rf bin/
	rm -f coverage.out coverage.html

# Format code
format: ## Format Go code
	go fmt ./...
	goimports -w .

# Lint code
lint: ## Lint Go code
	golangci-lint run

# Install development dependencies
deps: ## Install/update Go dependencies
	go mod download
	go mod tidy

# Database setup
db-setup: ## Set up the development database
	@echo "Setting up development database..."
	psql -h localhost -U postgres -d postgres -f .devcontainer/init-db.sql

# Database reset
db-reset: ## Reset the development database
	@echo "Resetting development database..."
	psql -h localhost -U postgres -d postgres -c "DROP DATABASE IF EXISTS timebot_dev; DROP DATABASE IF EXISTS timebot_test;"
	make db-setup

# Install development tools
install-tools: ## Install Go development tools
	go install github.com/air-verse/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/go-delve/delve/cmd/dlv@latest

# Docker commands
docker-build: ## Build Docker image
	docker build -t pcc-timebot .

docker-run: ## Run Docker container
	docker run -p 8080:8080 pcc-timebot

# Development environment
dev-up: ## Start development environment
	docker-compose -f .devcontainer/docker-compose.yml up -d

dev-down: ## Stop development environment
	docker-compose -f .devcontainer/docker-compose.yml down

dev-logs: ## Show development environment logs
	docker-compose -f .devcontainer/docker-compose.yml logs -f
