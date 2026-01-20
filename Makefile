# PCC-TimeBot Multi-Service Makefile

.PHONY: help build-all test-all clean deps

# Default target
help:
	@echo "Available targets:"
	@echo "  build-all     - Build all services and CLI tools (current platform)"
	@echo "  build-all-cross - Build all services and CLI tools (all platforms)"
	@echo "  test-all      - Run tests for all services"
	@echo "  clean         - Clean build artifacts"
	@echo "  deps          - Download dependencies for all services"
	@echo ""
	@echo "Individual service targets:"
	@echo "  build-timebot - Build timebot service"
	@echo "  build-api     - Build web API service"
	@echo "  build-cli     - Build all CLI tools (current platform)"
	@echo "  build-cli-cross - Build all CLI tools (all platforms)"
	@echo "  test-timebot  - Test timebot service"
	@echo "  test-api      - Test web API service"
	@echo ""
	@echo "CLI tools (current platform):"
	@echo "  build-add-time        - Build add-time CLI"
	@echo "  build-send-timecards  - Build send-timecards CLI"
	@echo "  build-update-employee - Build update-employee CLI"
	@echo ""
	@echo "CLI tools (cross-platform):"
	@echo "  build-add-time-cross        - Build add-time CLI for all platforms"
	@echo "  build-send-timecards-cross  - Build send-timecards CLI for all platforms"
	@echo "  build-update-employee-cross - Build update-employee CLI for all platforms"

# Build all services
build-all: build-timebot build-api build-cli-cross

build-timebot:
	@echo "Building timebot service..."
	mkdir -p bin/darwin-arm64
	cd cmd/timebot-service && go build -o ../../bin/darwin-arm64/timebot-service .

build-api:
	@echo "Building web API service..."
	mkdir -p bin/darwin-arm64
	cd cmd/web-api && go build -o ../../bin/darwin-arm64/web-api .

# Build CLI tools
build-cli: build-add-time build-send-timecards build-update-employee

# Cross-platform builds for all CLI tools
build-cli-cross: build-add-time-cross build-send-timecards-cross build-update-employee-cross

build-add-time:
	@echo "Building add-time CLI..."
	mkdir -p bin/darwin-arm64
	cd cmd/add-time && go build -o ../../bin/darwin-arm64/add-time .

build-add-time-cross:
	@echo "Building add-time CLI for all platforms..."
	mkdir -p bin/darwin-arm64 bin/linux-amd64 bin/linux-arm64 bin/windows-amd64
	cd cmd/add-time && GOOS=darwin GOARCH=arm64 go build -o ../../bin/darwin-arm64/add-time .
	cd cmd/add-time && GOOS=linux GOARCH=amd64 go build -o ../../bin/linux-amd64/add-time .
	cd cmd/add-time && GOOS=linux GOARCH=arm64 go build -o ../../bin/linux-arm64/add-time .
	cd cmd/add-time && GOOS=windows GOARCH=amd64 go build -o ../../bin/windows-amd64/add-time.exe .

build-send-timecards:
	@echo "Building send-timecards CLI..."
	mkdir -p bin/darwin-arm64
	cd cmd/send-timecards && go build -o ../../bin/darwin-arm64/send-timecards .

build-send-timecards-cross:
	@echo "Building send-timecards CLI for all platforms..."
	mkdir -p bin/darwin-arm64 bin/linux-amd64 bin/linux-arm64 bin/windows-amd64
	cd cmd/send-timecards && GOOS=darwin GOARCH=arm64 go build -o ../../bin/darwin-arm64/send-timecards .
	cd cmd/send-timecards && GOOS=linux GOARCH=amd64 go build -o ../../bin/linux-amd64/send-timecards .
	cd cmd/send-timecards && GOOS=linux GOARCH=arm64 go build -o ../../bin/linux-arm64/send-timecards .
	cd cmd/send-timecards && GOOS=windows GOARCH=amd64 go build -o ../../bin/windows-amd64/send-timecards.exe .

build-update-employee:
	@echo "Building update-employee CLI..."
	mkdir -p bin/darwin-arm64
	cd cmd/update-employee && go build -o ../../bin/darwin-arm64/update-employee .

build-update-employee-cross:
	@echo "Building update-employee CLI for all platforms..."
	mkdir -p bin/darwin-arm64 bin/linux-amd64 bin/linux-arm64 bin/windows-amd64
	cd cmd/update-employee && GOOS=darwin GOARCH=arm64 go build -o ../../bin/darwin-arm64/update-employee .
	cd cmd/update-employee && GOOS=linux GOARCH=amd64 go build -o ../../bin/linux-amd64/update-employee .
	cd cmd/update-employee && GOOS=linux GOARCH=arm64 go build -o ../../bin/linux-arm64/update-employee .
	cd cmd/update-employee && GOOS=windows GOARCH=amd64 go build -o ../../bin/windows-amd64/update-employee.exe .

# Test all services
test-all: test-shared test-timebot test-api

test-shared:
	@echo "Testing shared packages..."
	cd shared && go test ./...

test-timebot:
	@echo "Testing timebot service..."
	go test ./cmd/timebot-service/... ./internal/admin/... ./internal/email/...

test-api:
	@echo "Testing web API service..."
	go test ./cmd/web-api/... ./internal/handlers/... ./internal/middleware/...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/darwin-arm64/* bin/linux-amd64/* bin/linux-arm64/* bin/windows-amd64/*
	rm -rf tmp/
	go clean ./...

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	go mod download && go mod tidy

# Create bin directories
bin-dirs:
	mkdir -p bin/darwin-arm64 bin/linux-amd64 bin/linux-arm64 bin/windows-amd64

# Build targets depend on bin directories
build-timebot: bin-dirs
build-api: bin-dirs

# Development shortcuts
dev-timebot:
	@echo "Running timebot service in development mode..."
	cd cmd/timebot-service && go run .

dev-api:
	@echo "Running web API service in development mode..."
	cd cmd/web-api && go run .

# Air live reload (if installed)
air-timebot:
	@echo "Starting timebot service with Air live reload..."
	cd cmd/timebot-service && air

air-api:
	@echo "Starting web API service with Air live reload..."
	cd cmd/web-api && air
