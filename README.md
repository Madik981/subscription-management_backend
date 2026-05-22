# Subscription Management Backend (Microservices)

This project is split into two Go microservices for subscription management.

## Services

- accounts-service: authentication and user management (port 8080)
- billing-service: plans and billings (port 8081)

## Stack

- Go 1.25
- Gin (HTTP API)
- GORM (ORM)
- PostgreSQL
- golang-migrate (SQL migrations)
- resty v2 (internal HTTP calls)

## Repository Layout

- accounts-service/
- billing-service/
- docker-compose.yml
- Makefile
- .env.example

## Quick Start (Docker)

1) Create your env file (optional if you already have one):

```bash
cp .env.example .env
```

2) Start the stack:

```bash
docker compose up --build -d
```

3) Run migrations from the host:

```bash
make migrate-accounts-up ACCOUNTS_DB_URL=postgres://postgres:postgres@localhost:2345/subscription_management_accounts?sslmode=disable
make migrate-billing-up BILLING_DB_URL=postgres://postgres:postgres@localhost:2346/subscription_management_billing?sslmode=disable
```

4) Stop the stack:

```bash
docker compose down
```

Frontend:

- Next.js app source lives in `subscription-management_frontend/`
- Docker Compose runs `madik98/subscription-management_frontend` and serves it at `http://localhost:3000`
- The UI talks to accounts-service on `http://localhost:8080` and billing-service on `http://localhost:8081`
- Frontend can be started independently with `docker compose up -d frontend`; if backend APIs are down, the UI still opens and shows API status as unavailable.

## Environment Variables

The compose file loads values from .env. Most values are also available in .env.example.

Core settings:

- POSTGRES_USER
- POSTGRES_PASSWORD
- ACCOUNTS_DB_NAME
- BILLING_DB_NAME
- ACCOUNTS_PORT
- BILLING_PORT
- JWT_SECRET
- INTERNAL_TOKEN

Service connection strings (container network):

- ACCOUNTS_DATABASE_URL
- BILLING_DATABASE_URL
- BILLING_SERVICE_URL
- ACCOUNTS_SERVICE_URL

Migration URLs (host usage):

- ACCOUNTS_DB_URL
- BILLING_DB_URL

Note: Postgres ports are exposed as 2345 (accounts) and 2346 (billing) on the host, but services inside the Docker network connect to port 5432.

## Makefile Commands

- make test
- make test-accounts
- make test-billing
- make migrate-install
- make migrate-accounts-up
- make migrate-accounts-down
- make migrate-accounts-down1
- make migrate-accounts-version
- make migrate-accounts-force VERSION=1
- make migrate-accounts-create NAME=add_new_column
- make migrate-billing-up
- make migrate-billing-down
- make migrate-billing-down1
- make migrate-billing-version
- make migrate-billing-force VERSION=1
- make migrate-billing-create NAME=add_new_column

## Tests

Run all backend handler tests from the repository root:

```bash
make test
```


## API Summary

Accounts service:

- POST /auth/register
- POST /auth/login
- GET /auth/me
- POST /users
- GET /users
- GET /users/:id
- PATCH /users/:id
- DELETE /users/:id

Billing service:

- POST /plans
- GET /plans
- GET /plans/:id
- PATCH /plans/:id
- DELETE /plans/:id
- POST /billings
- GET /billings
- GET /billings/:id
- PATCH /billings/:id/pay
- PATCH /billings/:id/fail
- DELETE /billings/:id

## API Examples (JSON)

Accounts service

POST /auth/register

Request JSON:

```json
{
	"name": "John Student",
	"email": "john@example.com",
	"password": "12345678",
	"plan_id": 1
}
```

Response JSON:

```json
{
	"token": "<jwt_token>",
	"user": {
		"id": 1,
		"name": "John Student",
		"email": "john@example.com",
		"plan_id": 1,
		"is_active": true,
		"created_at": "2026-05-01T10:05:00Z",
		"updated_at": "2026-05-01T10:05:00Z"
	}
}
```

POST /auth/login

Request JSON:

```json
{
	"email": "john@example.com",
	"password": "12345678"
}
```

Response JSON:

```json
{
	"token": "<jwt_token>",
	"user": {
		"id": 1,
		"name": "John Student",
		"email": "john@example.com",
		"plan_id": 1,
		"is_active": true,
		"created_at": "2026-05-01T10:05:00Z",
		"updated_at": "2026-05-01T10:05:00Z"
	}
}
```

GET /auth/me

Response JSON:

```json
{
	"id": 1,
	"name": "John Student",
	"email": "john@example.com",
	"plan_id": 1,
	"is_active": true,
	"created_at": "2026-05-01T10:05:00Z",
	"updated_at": "2026-05-01T10:05:00Z"
}
```

POST /users

Request JSON:

```json
{
	"name": "Alex Brown",
	"email": "alex@example.com",
	"plan_id": 1,
	"is_active": true
}
```

Response JSON:

```json
{
	"id": 2,
	"name": "Alex Brown",
	"email": "alex@example.com",
	"plan_id": 1,
	"is_active": true,
	"created_at": "2026-05-01T10:10:00Z",
	"updated_at": "2026-05-01T10:10:00Z"
}
```

GET /users

Response JSON:

```json
[
	{
		"id": 2,
		"name": "Alex Brown",
		"email": "alex@example.com",
		"plan_id": 1,
		"is_active": true,
		"created_at": "2026-05-01T10:10:00Z",
		"updated_at": "2026-05-01T10:10:00Z"
	}
]
```

GET /users/:id

Response JSON:

```json
{
	"id": 2,
	"name": "Alex Brown",
	"email": "alex@example.com",
	"plan_id": 1,
	"is_active": true,
	"created_at": "2026-05-01T10:10:00Z",
	"updated_at": "2026-05-01T10:10:00Z"
}
```

PATCH /users/:id

Request JSON:

```json
{
	"is_active": false
}
```

Response JSON:

```json
{
	"id": 2,
	"name": "Alex Brown",
	"email": "alex@example.com",
	"plan_id": 1,
	"is_active": false,
	"created_at": "2026-05-01T10:10:00Z",
	"updated_at": "2026-05-01T10:15:00Z"
}
```

Billing service

POST /plans

Request JSON:

```json
{
	"name": "Pro",
	"description": "Pro monthly plan",
	"price": 29.99,
	"currency": "USD",
	"billing_cycle": "monthly"
}
```

Response JSON:

```json
{
	"id": 1,
	"name": "Pro",
	"description": "Pro monthly plan",
	"price": 29.99,
	"currency": "USD",
	"billing_cycle": "monthly",
	"created_at": "2026-05-01T10:00:00Z",
	"updated_at": "2026-05-01T10:00:00Z"
}
```

GET /plans

Response JSON:

```json
[
	{
		"id": 1,
		"name": "Pro",
		"description": "Pro monthly plan",
		"price": 29.99,
		"currency": "USD",
		"billing_cycle": "monthly",
		"created_at": "2026-05-01T10:00:00Z",
		"updated_at": "2026-05-01T10:00:00Z"
	}
]
```

GET /plans/:id

Response JSON:

```json
{
	"id": 1,
	"name": "Pro",
	"description": "Pro monthly plan",
	"price": 29.99,
	"currency": "USD",
	"billing_cycle": "monthly",
	"created_at": "2026-05-01T10:00:00Z",
	"updated_at": "2026-05-01T10:00:00Z"
}
```

PATCH /plans/:id

Request JSON:

```json
{
	"price": 34.99,
	"billing_cycle": "yearly"
}
```

Response JSON:

```json
{
	"id": 1,
	"name": "Pro",
	"description": "Pro monthly plan",
	"price": 34.99,
	"currency": "USD",
	"billing_cycle": "yearly",
	"created_at": "2026-05-01T10:00:00Z",
	"updated_at": "2026-05-01T10:20:00Z"
}
```

POST /billings

Request JSON:

```json
{
	"user_id": 2,
	"plan_id": 1,
	"amount": 29.99,
	"due_date": "2026-05-10T00:00:00Z",
	"description": "May subscription invoice"
}
```

Response JSON:

```json
{
	"id": 1,
	"user_id": 2,
	"plan_id": 1,
	"plan": {
		"id": 1,
		"name": "Pro",
		"description": "Pro monthly plan",
		"price": 34.99,
		"currency": "USD",
		"billing_cycle": "yearly",
		"created_at": "2026-05-01T10:00:00Z",
		"updated_at": "2026-05-01T10:20:00Z"
	},
	"amount": 29.99,
	"status": "pending",
	"due_date": "2026-05-10T00:00:00Z",
	"paid_at": null,
	"description": "May subscription invoice",
	"created_at": "2026-05-01T10:30:00Z",
	"updated_at": "2026-05-01T10:30:00Z"
}
```

GET /billings

Response JSON:

```json
[
	{
		"id": 1,
		"user_id": 2,
		"plan_id": 1,
		"plan": {
			"id": 1,
			"name": "Pro",
			"description": "Pro monthly plan",
			"price": 34.99,
			"currency": "USD",
			"billing_cycle": "yearly",
			"created_at": "2026-05-01T10:00:00Z",
			"updated_at": "2026-05-01T10:20:00Z"
		},
		"amount": 29.99,
		"status": "pending",
		"due_date": "2026-05-10T00:00:00Z",
		"paid_at": null,
		"description": "May subscription invoice",
		"created_at": "2026-05-01T10:30:00Z",
		"updated_at": "2026-05-01T10:30:00Z"
	}
]
```

GET /billings/:id

Response JSON:

```json
{
	"id": 1,
	"user_id": 2,
	"plan_id": 1,
	"plan": {
		"id": 1,
		"name": "Pro",
		"description": "Pro monthly plan",
		"price": 34.99,
		"currency": "USD",
		"billing_cycle": "yearly",
		"created_at": "2026-05-01T10:00:00Z",
		"updated_at": "2026-05-01T10:20:00Z"
	},
	"amount": 29.99,
	"status": "pending",
	"due_date": "2026-05-10T00:00:00Z",
	"paid_at": null,
	"description": "May subscription invoice",
	"created_at": "2026-05-01T10:30:00Z",
	"updated_at": "2026-05-01T10:30:00Z"
}
```

PATCH /billings/:id/pay

Response JSON:

```json
{
	"id": 1,
	"user_id": 2,
	"plan_id": 1,
	"plan": {
		"id": 1,
		"name": "Pro",
		"description": "Pro monthly plan",
		"price": 34.99,
		"currency": "USD",
		"billing_cycle": "yearly",
		"created_at": "2026-05-01T10:00:00Z",
		"updated_at": "2026-05-01T10:20:00Z"
	},
	"amount": 29.99,
	"status": "paid",
	"due_date": "2026-05-10T00:00:00Z",
	"paid_at": "2026-05-10T12:00:00Z",
	"description": "May subscription invoice",
	"created_at": "2026-05-01T10:30:00Z",
	"updated_at": "2026-05-10T12:00:00Z"
}
```

## Notes

- JWT must be sent as: Authorization: Bearer <token>
- Responses no longer embed related entities (no plan/user object nesting).
- Internal service calls use the shared INTERNAL_TOKEN via X-Internal-Token header.
