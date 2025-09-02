# PCC-TimeBot Development Container

This directory contains the development container configuration for the PCC-TimeBot project, providing a consistent and isolated development environment.

## What's Included

### Development Environment

- **Go 1.24**: Latest Go runtime and tools
- **PostgreSQL 15**: Database server with separate development and test databases
- **VS Code Extensions**: Pre-configured extensions for Go development
- **Development Tools**: Air (hot reload), golangci-lint, delve debugger, and more

### Services

- **app**: Main application container with Go development environment
- **db**: PostgreSQL development database (port 5432)
- **test-db**: PostgreSQL test database (port 5433)

## Quick Start

### Prerequisites

- [Docker Desktop](https://www.docker.com/products/docker-desktop)
- [VS Code](https://code.visualstudio.com/)
- [Dev Containers extension](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers)

### Getting Started

1. **Open in Dev Container**

   ```bash
   # From VS Code Command Palette (Ctrl/Cmd + Shift + P)
   # Run: "Dev Containers: Open Folder in Container"
   # Select the project root directory
   ```

2. **Or use the command line**

   ```bash
   # Clone and open
   git clone <repository-url>
   cd PCC-TimeBot
   code .
   # Then use "Dev Containers: Reopen in Container"
   ```

3. **Set up environment variables**

   ```bash
   cp .devcontainer/.env.example .env
   # Edit .env with your actual configuration values
   ```

4. **Initialize the database**

   ```bash
   make db-setup
   ```

5. **Start development**
   ```bash
   make dev  # Starts with hot reload
   # or
   make run  # Normal start
   ```

## Available Commands

Use the provided Makefile for common development tasks:

```bash
make help           # Show all available commands
make dev            # Start with hot reload
make test           # Run all tests
make test-db        # Run database tests only
make lint           # Run code linting
make format         # Format code
make db-reset       # Reset development database
```

## Database Access

### Development Database

- **Host**: localhost
- **Port**: 5432
- **Database**: timebot_dev
- **Username**: postgres
- **Password**: (none required)

### Test Database

- **Host**: localhost
- **Port**: 5432
- **Database**: timebot_test
- **Username**: postgres
- **Password**: (none required)

### Sample Data

The development database comes pre-populated with sample data:

- Admin user (admin/admin)
- Test employees (john doe, jane smith)
- Sample payroll entries

## Configuration

### Environment Variables

Key environment variables for development:

```bash
DATABASE_URL=postgresql://postgres@localhost:5432/timebot_dev?sslmode=disable
TEST_DATABASE_URL=postgresql://postgres@localhost:5432/timebot_test?sslmode=disable
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password
TWILIO_ACCOUNT_SID=your-twilio-sid
TWILIO_AUTH_TOKEN=your-twilio-token
PORT=8080
```

### VS Code Settings

The dev container comes with pre-configured VS Code settings for:

- Go development tools
- Formatting and linting
- Testing integration
- Debugging configuration

## Testing

### Running Tests

```bash
# All tests
make test

# Database tests only (uses test database)
make test-db

# Tests with coverage
make test-coverage
```

### Test Database

Tests automatically use the separate test database to avoid affecting development data.

## Debugging

The dev container includes the Delve debugger. You can:

1. Set breakpoints in VS Code
2. Use the Run and Debug panel
3. Or debug from the command line with `dlv`

## Troubleshooting

### Container Issues

```bash
# Rebuild container
# From VS Code: "Dev Containers: Rebuild Container"

# Or manually
docker-compose -f .devcontainer/docker-compose.yml down
docker-compose -f .devcontainer/docker-compose.yml build --no-cache
```

### Database Issues

````bash
# Reset database
make db-reset

# Check database connection
```bash
# Check database connection
psql -h localhost -U postgres -d timebot_dev
````

````

### Go Module Issues
```bash
# Clean and reinstall dependencies
go clean -modcache
go mod download
go mod tidy
````

## Ports

The following ports are forwarded to your local machine:

- **8080**: Application server
- **5432**: PostgreSQL database (both dev and test databases)

## File Structure

```
.devcontainer/
├── devcontainer.json      # Main dev container configuration
├── docker-compose.yml     # Services configuration
├── Dockerfile            # Development environment setup
├── init-db.sql           # Database initialization
├── .env.example          # Environment variables template
├── Makefile              # Development commands
├── .air.toml             # Hot reload configuration
└── README.md             # This file
```
