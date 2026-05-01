# Accounts Service

This service handles authentication and user management.

## Environment

- ACCOUNTS_PORT (default: 8080)
- ACCOUNTS_DATABASE_URL (preferred)
- DATABASE_URL (fallback)
- JWT_SECRET (shared with billing service)
- BILLING_SERVICE_URL (default: http://localhost:8081)
- INTERNAL_TOKEN (shared internal token)

## Run

```bash
go run ./
```

## Migrations

Use the root Makefile targets:

```bash
make migrate-accounts-up
```
