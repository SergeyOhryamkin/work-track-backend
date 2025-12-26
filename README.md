# Work Track Backend

A robust Go backend for a work tracking application, featuring JWT authentication, PostgreSQL integration, and a clean architecture.

## Features

- **Authentication**: Secure user registration and login with JWT tokens
- **Work Tracking**: Track daily work items, hours, shifts, and special conditions (emergency/holiday calls)
- **RESTful API**: Clean, resource-oriented API design
- **PostgreSQL**: Reliable data persistence with `pgx` driver
- **Docker**: Containerized for easy development and deployment
- **CI/CD**: Automated testing and build pipeline with GitHub Actions

## Tech Stack

- **Language**: Go 1.21+
- **Router**: [Chi](https://github.com/go-chi/chi)
- **Database**: PostgreSQL
- **Driver**: [pgx](https://github.com/jackc/pgx)
- **Auth**: JWT (JSON Web Tokens)
- **Config**: [godotenv](https://github.com/joho/godotenv)
- **Testing**: Go standard library `testing` package

## Project Structure

```
backend/
├── cmd/
│   └── api/            # Application entry point
├── internal/
│   ├── config/         # Configuration management
│   ├── database/       # Database connection
│   ├── handler/        # HTTP handlers (controllers)
│   ├── middleware/     # HTTP middleware (auth, logging, CORS)
│   ├── models/         # Data models
│   ├── repository/     # Data access layer
│   ├── service/        # Business logic layer
│   └── util/           # Utility functions
├── migrations/         # Database migrations
└── ...
```

## Quick Start

1. **Clone the repository**
2. **Set up environment variables**
   ```bash
   cp .env.example .env
   ```
3. **Run with Docker**
   ```bash
   make docker-up
   ```
4. **Run locally**
   ```bash
   make run
   ```

## API Documentation

See [API_DOCUMENTATION.md](API_DOCUMENTATION.md) for detailed endpoint descriptions.

## Deployment

See [DEPLOYMENT.md](DEPLOYMENT.md) for deployment instructions.
# Test CI/CD
