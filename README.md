# PCC-TimeBot

A payroll management application for Poulsen Concrete Company. A React single-page
app lets employees submit time entries and lets admins manage payroll, backed by a
secured Go JSON API. The compiled SPA is embedded into the Go binary and served
same-origin in production.

## Features

- Mobile-friendly web app for time entry submission
- Passwordless sign-in via email magic links (no passwords)
- Admin dashboard for employee and timecard management
- Email notifications for time card submissions
- Admin-provisioned accounts (no public signup)
- PostgreSQL for persistence

## Prerequisites

- Go 1.24+
- Node.js 20+ and npm

## Getting Started

1. **Clone the repository**

   ```bash
   git clone <repository-url>
   cd PCC-TimeBot
   ```

2. **Set up environment variables**

   ```bash
   cp .env.example .env
   # edit .env with real values
   ```

3. **Initialize the database**

   ```bash
   psql "$DATABASE_URL" -f scripts/init-db.sql
   ```

4. **Build the SPA and server**

   ```bash
   npm run build
   go build -o bin/web-api ./cmd/web-api/
   ```

5. **Create the first admin** (accounts are admin-provisioned; no public signup)

   ```bash
   go run ./cmd/bootstrap-admin -first Jane -last Doe -email "$ADMIN_EMAIL"
   ```

6. **Run the server**

   ```bash
   bin/web-api
   # App is at http://localhost:$PORT (default 8080)
   ```

   The server auto-loads `.env` if present. On Heroku, the file won't exist and
   config vars are injected by the platform instead.

## Local Development with Hot Reload

Run the Go API and the Vite dev server in separate terminals:

```bash
# Terminal 1 – start Go API (auto-loads .env)
bin/web-api

# Terminal 2 – Vite dev server with hot reload, proxies /api to Go
npm run dev
# App is at http://localhost:5173
```

When running Vite separately, uncomment `CORS_ALLOWED_ORIGINS` in `.env.local`
so the Go server accepts requests from the Vite origin.

## Dev Login (local testing only)

When `APP_ENV=dev`, submitting the login form skips the email entirely
and logs you in immediately. Just type your email and hit **Send sign-in link** —
no email, no curl, no browser console needed.

This behaviour is **not active in production** (`APP_ENV=prod`).

## Testing

```bash
go test ./...
```

All tests are unit tests — no database required.

## Admin CLI

`timebot-service` is an admin-only CLI for operations that don't fit in the web UI:

```bash
# Manually add time for an employee (interactive prompts)
go run ./cmd/timebot-service -addTime

# Update an employee record (interactive prompts)
go run ./cmd/timebot-service -updateEmployee

# Send time cards for a pay period
go run ./cmd/timebot-service -sendTimeCards -startDate 2026-06-09 -endDate 2026-06-22

# Send time cards using the last recorded period end date
go run ./cmd/timebot-service -sendTimeCards -useLastPeriod
```

## Project Structure

```
├── cmd/
│   ├── web-api/            # Production HTTP server (JSON API + embedded SPA)
│   ├── bootstrap-admin/    # Seed/promote the first admin employee
│   ├── timebot-service/    # Admin CLI (add time, send timecards, etc.)
│   └── ...                 # Other admin CLIs (add-time, update-employee, etc.)
├── internal/
│   ├── auth/               # Magic-link tokens, session middleware
│   ├── handlers/           # HTTP handlers (auth, timecards, admin)
│   ├── admin/              # Business logic for timecard building/sending
│   ├── email/              # SMTP helpers including magic-link email
│   └── middleware/         # CORS, request logger
├── shared/
│   ├── database/           # DB access layer (employees, sessions, payroll)
│   ├── timecalc/           # Hours calculation
│   ├── timecard/           # Timecard model and HTML rendering
│   └── helpers/            # Shared utilities
├── web/
│   ├── embed.go            # go:embed declaration for app/dist
│   └── app/                # Vite + React + TypeScript SPA
│       └── src/            # App source (pages, components, API client)
└── scripts/
    └── init-db.sql         # Full schema (run this to (re)initialize the DB)
```

## Deployment (Heroku)

Configured in `app.json` using two buildpacks:

1. `heroku/nodejs` — runs `heroku-postbuild`, which builds the SPA into `web/app/dist`
2. `heroku/go` — builds `./cmd/web-api`, embedding `web/app/dist` into the binary

`Procfile` runs `web: bin/web-api`. One binary serves the API (`/api/v1`) and the SPA.

**Required environment variables on Heroku:**

| Variable | Description |
|---|---|
| `DATABASE_URL` | PostgreSQL connection string (Heroku Postgres add-on) |
| `SMTP_USERNAME` | SMTP username (also receives payroll summaries) |
| `SMTP_PASSWORD` | SMTP password |
| `APP_BASE_URL` | Public app URL, used to build magic-link URLs |
| `APP_ENV` | Set to `prod` (enables `Secure` cookies) |

**Optional / reserved environment variables:**

| Variable | Description |
|---|---|
| `CORS_ALLOWED_ORIGINS` | Local development only, when the SPA is served from a different origin such as Vite |
| `TWILIO_ACCOUNT_SID` | Reserved for future supervisor SMS notifications on time submission |
| `TWILIO_AUTH_TOKEN` | Reserved for future supervisor SMS notifications on time submission |
| `TWILIO_PHONE_NUMBER` | Reserved sender number for future supervisor SMS notifications |

Twilio is intentionally **not wired into the current web submission flow on this branch**,
but the placeholders remain so the future SMS notification work has a clear home.

