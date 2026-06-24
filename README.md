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

### Building

```bash
go build -o bin/timebot-service ./cmd/timebot-service/
```

### CLI Usage

The `timebot-service` binary serves as both the web server and the admin CLI, controlled by flags:

```bash
# Run the web server (used by Heroku via Procfile)
bin/timebot-service -heroku

# Manually add or update time for an employee (interactive prompts)
bin/timebot-service -addTime

# Update an employee record (interactive prompts)
bin/timebot-service -updateEmployee

# Send time cards for a specific pay period
bin/timebot-service -sendTimeCards -startDate 2026-06-09 -endDate 2026-06-22

# Send time cards using the last recorded period end date (auto-calculates next 7-day period)
bin/timebot-service -sendTimeCards -useLastPeriod
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
