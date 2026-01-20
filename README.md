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

- Go (see `go.mod` for the version)

Optional (only if you want to run integration tests against Postgres):

- PostgreSQL

### Getting Started

1. **Clone the repository**

   ```bash
   git clone <repository-url>
   cd PCC-TimeBot
   ```

2. **Set up environment variables**

   Copy `.env.example` to `.env` and edit values as needed.

3. **Run tests**

   ```bash
   make test-all
   ```

### Available Commands

```bash
make help           # Show all available commands
make build-all       # Build the application binaries
make test-all        # Run unit tests (no database required)
```

### Project Structure

```
├── cmd/                # Binaries (CLIs + services)
├── internal/           # App internals (handlers, admin, email, middleware)
├── shared/             # Shared packages (database, timecard, helpers, etc.)
├── services/           # Dockerfiles for services (optional)
└── timebot-mobile/     # Mobile app
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

By default, `go test ./...` / `make test-all` runs unit tests only (no database required).

Database-backed tests are marked as integration tests and can be run explicitly with:

```bash
go test -tags=integration ./...
```

### Database

The application uses PostgreSQL for persistence in production (e.g. Heroku Postgres).
Local Postgres is only required if you want to run the integration tests.

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
