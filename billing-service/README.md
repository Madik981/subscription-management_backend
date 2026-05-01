# Billing Service

This service manages plans and billings.

## Environment

- BILLING_PORT (default: 8081)
- BILLING_DATABASE_URL (preferred)
- DATABASE_URL (fallback)
- JWT_SECRET (shared with accounts service)
- ACCOUNTS_SERVICE_URL (default: http://localhost:8080)
- INTERNAL_TOKEN (shared internal token)

## Run

```bash
go run ./
```

## Migrations

Use the root Makefile targets:

```bash
make migrate-billing-up
```
