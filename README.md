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

### What is `timebot-service`?

`timebot-service` is the core deployable binary. It runs in two modes:

**Server mode** — receives inbound SMS webhooks from Twilio, parses time submissions, writes hours to the database, and sends confirmation texts to the employee, their supervisor, and the owner.

**Admin CLI mode** — interactive command-line tools for managing employee records and sending payroll emails, controlled by flags (see CLI Usage below). The same operations are also available as standalone binaries under `cmd/`.

### CLI Usage

The `timebot-service` binary serves as both the web server and the admin CLI, controlled by flags:

```bash
# Run the web server (used by Heroku via Procfile)
bin/timebot-service -heroku

# Add a new employee (interactive prompts, phone validation)
bin/timebot-service -addEmployee

# Add or correct time for an employee (interactive prompts)
bin/timebot-service -addTime

# Update an employee record field (interactive prompts)
bin/timebot-service -updateEmployee

# Send time cards for a specific pay period
bin/timebot-service -sendTimeCards -startDate 2026-06-09 -endDate 2026-06-22

# Send time cards using the last recorded period end date (auto-calculates next 7-day period)
bin/timebot-service -sendTimeCards -useLastPeriod
```

The same operations (plus additional ones) are available as standalone binaries:

```bash
# Add a new employee (interactive prompts, phone validation)
go run ./cmd/add-employee/

# Add or correct time for an employee (interactive prompts)
go run ./cmd/add-time/

# Update an employee record field (interactive prompts)
go run ./cmd/update-employee/

# Send time cards for a pay period (same flags as timebot-service -sendTimeCards)
go run ./cmd/send-timecards/ [-start YYYY-MM-DD -end YYYY-MM-DD | -lastperiod]

# Run the Monday timecard scheduler (sends automatically if today is Monday)
go run ./cmd/timecard-scheduler/
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
go test ./...

# Run specific test packages
go test ./shared/database/...
go test ./internal/admin/...
go test ./shared/timecard/...
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

Key environment variables (see `.env.example` for the full list):

- `DATABASE_URL`: PostgreSQL connection string
- `SMTP_USERNAME`: Gmail address used to send emails
- `SMTP_PASSWORD`: Gmail App Password (not your regular password)
- `TWILIO_ACCOUNT_SID`: Twilio account SID
- `TWILIO_AUTH_TOKEN`: Twilio auth token
- `TWILIO_PHONE`: Twilio phone number (10 digits, no `+1`)
- `ADMIN_EMAIL`: Receives the weekly payroll summary
- `BOOKKEEPER_EMAIL`: Also receives the weekly payroll summary
- `PORT`: Server port (default: 8080)

## Contributing

1. Make changes in the dev container environment
2. Run tests to ensure functionality
3. Follow Go formatting standards (`make format`)
4. Ensure code passes linting (`make lint`)
5. Submit pull requests for review
