# PCC-TimeBot

A Go-based payroll management application for Poulsen Concrete Company, enabling employees to submit time entries and supervisors to manage payroll data.

## Features

- Employee time entry submission
- Supervisor payroll management
- Email notifications for time card submissions
- SMS notifications for reminders
- Administrative functions for employee management
- PostgreSQL database for data persistence

## Development

### Prerequisites

- [Docker Desktop](https://www.docker.com/products/docker-desktop)
- [VS Code](https://code.visualstudio.com/)
- [Dev Containers extension](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers)

### Getting Started

1. **Open in Dev Container**

   ```bash
   # Clone the repository
   git clone <repository-url>
   cd PCC-TimeBot

   # Open in VS Code and reopen in container
   code .
   # Use Command Palette: "Dev Containers: Reopen in Container"
   ```

2. **Set up environment variables**

   ```bash
   cp .devcontainer/.env.example .env
   # Edit .env with your configuration values
   ```

3. **Start development**
   ```bash
   make dev  # Starts with hot reload
   ```

### Available Commands

```bash
make help           # Show all available commands
make dev            # Start with hot reload
make build          # Build the application
make test           # Run all tests
make test-db        # Run database tests only
make lint           # Run code linting
make format         # Format code
make db-setup       # Set up development database
make db-reset       # Reset development database
```

### Project Structure

```
├── src/
│   ├── admin/          # Administrative functions
│   ├── constants/      # Application constants
│   ├── database/       # Database operations
│   ├── email/          # Email functionality
│   ├── helpers/        # Utility functions
│   ├── timecalc/       # Time calculation logic
│   └── timecard/       # Time card processing
├── .devcontainer/      # Development container configuration
├── scripts/            # Database and setup scripts
└── main.go            # Application entry point
```

### Testing

The project includes comprehensive tests for all major components:

```bash
# Run all tests
make test

# Run specific test packages
go test ./src/database
go test ./src/admin
go test ./src/timecard
```

Tests automatically use a separate test database to avoid affecting development data.

### Database

The application uses PostgreSQL with separate databases for development and testing:

- **Development**: `timebot_dev` (port 5432)
- **Testing**: `timebot_test` (port 5433)

Database schema and sample data are automatically set up in the dev container.

## Deployment

The application is designed to be deployed as a single binary:

```bash
# Build for production
make build

# The binary will be created as ./pcc-timebot
```

## Environment Variables

Key environment variables:

- `DATABASE_URL`: PostgreSQL connection string
- `SMTP_USERNAME`: Email service username
- `SMTP_PASSWORD`: Email service password
- `TWILIO_ACCOUNT_SID`: Twilio account SID for SMS
- `TWILIO_AUTH_TOKEN`: Twilio auth token
- `PORT`: Application port (default: 8080)

## Contributing

1. Make changes in the dev container environment
2. Run tests to ensure functionality
3. Follow Go formatting standards (`make format`)
4. Ensure code passes linting (`make lint`)
5. Submit pull requests for review
