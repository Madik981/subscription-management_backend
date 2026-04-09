# Subscription Management Backend

Hi! This is my student backend project for subscription management.
I built it with Go, Gin, GORM, and PostgreSQL.
The API lets you manage plans, users, and billings.

## Stack

- Go 1.25
- Gin (HTTP API)
- GORM (ORM)
- PostgreSQL (database)

## Project Structure

- main.go
- models/
  - plan.go
  - user.go
  - billing.go
- handlers/
  - plan_handler.go
  - user_handler.go
  - billing_handler.go

## Environment Variables

- PORT: server port (default: 8080)
- DATABASE_URL: PostgreSQL connection string
- JWT_SECRET: secret key for signing JWT tokens (default: super-secret-key)

Example DATABASE_URL format:

host=localhost user=postgres password=postgres dbname=subscription_management port=5432 sslmode=disable TimeZone=UTC

## Data Models

### Plan

- id
- name
- description
- price
- currency
- billing_cycle
- created_at
- updated_at

### User

- id
- name
- email
- password (stored in DB, hidden in JSON responses)
- plan_id
- plan
- is_active
- created_at
- updated_at

### Billing

- id
- user_id
- user
- plan_id
- plan
- amount
- status
- due_date
- paid_at
- description
- created_at
- updated_at

Status values:

- pending
- paid
- failed

## API Endpoints

### Authentication

- POST /auth/register
- POST /auth/login
- GET /auth/me

Register request JSON:

```json
{
  "name": "John Student",
  "email": "john@example.com",
  "password": "12345678",
  "plan_id": 1
}
```

Register response JSON:

```json
{
  "token": "<jwt_token>",
  "user": {
    "id": 1,
    "name": "John Student",
    "email": "john@example.com",
    "plan_id": 1,
    "plan": {
      "id": 1,
      "name": "Pro",
      "description": "Pro monthly plan",
      "price": 29.99,
      "currency": "USD",
      "billing_cycle": "monthly",
      "created_at": "2026-04-09T10:00:00Z",
      "updated_at": "2026-04-09T10:00:00Z"
    },
    "is_active": true,
    "created_at": "2026-04-09T10:05:00Z",
    "updated_at": "2026-04-09T10:05:00Z"
  }
}
```

Login request JSON:

```json
{
  "email": "john@example.com",
  "password": "12345678"
}
```

Login response JSON:

```json
{
  "token": "<jwt_token>",
  "user": {
    "id": 1,
    "name": "John Student",
    "email": "john@example.com",
    "plan_id": 1,
    "is_active": true,
    "created_at": "2026-04-09T10:05:00Z",
    "updated_at": "2026-04-09T10:05:00Z"
  }
}
```

Get current user response JSON (`GET /auth/me`):

```json
{
  "id": 1,
  "name": "John Student",
  "email": "john@example.com",
  "plan_id": 1,
  "is_active": true,
  "created_at": "2026-04-09T10:05:00Z",
  "updated_at": "2026-04-09T10:05:00Z"
}
```

Important: all `/plans`, `/users`, and `/billings` endpoints are protected now.
Send JWT in `Authorization` header as:

`Bearer <jwt_token>`

### Health

- GET /health

Response JSON:

```json
{
  "status": "ok"
}
```

### Plans

- POST /plans
- GET /plans
- GET /plans/:id
- PATCH /plans/:id

Create plan request JSON:

```json
{
  "name": "Pro",
  "description": "Pro monthly plan",
  "price": 29.99,
  "currency": "USD",
  "billing_cycle": "monthly"
}
```

Create plan response JSON:

```json
{
  "id": 1,
  "name": "Pro",
  "description": "Pro monthly plan",
  "price": 29.99,
  "currency": "USD",
  "billing_cycle": "monthly",
  "created_at": "2026-04-03T10:00:00Z",
  "updated_at": "2026-04-03T10:00:00Z"
}
```

Update plan request JSON (partial):

```json
{
  "price": 34.99,
  "billing_cycle": "yearly"
}
```

List plans response JSON:

```json
[
  {
    "id": 1,
    "name": "Pro",
    "description": "Pro monthly plan",
    "price": 29.99,
    "currency": "USD",
    "billing_cycle": "monthly",
    "created_at": "2026-04-03T10:00:00Z",
    "updated_at": "2026-04-03T10:00:00Z"
  }
]
```

### Users

- POST /users
- GET /users
- GET /users/:id
- PATCH /users/:id

Create user request JSON:

```json
{
  "name": "Alex Brown",
  "email": "alex@example.com",
  "plan_id": 1,
  "is_active": true
}
```

Create user response JSON:

```json
{
  "id": 1,
  "name": "Alex Brown",
  "email": "alex@example.com",
  "plan_id": 1,
  "plan": {
    "id": 1,
    "name": "Pro",
    "description": "Pro monthly plan",
    "price": 29.99,
    "currency": "USD",
    "billing_cycle": "monthly",
    "created_at": "2026-04-03T10:00:00Z",
    "updated_at": "2026-04-03T10:00:00Z"
  },
  "is_active": true,
  "created_at": "2026-04-03T10:05:00Z",
  "updated_at": "2026-04-03T10:05:00Z"
}
```

Update user request JSON (partial):

```json
{
  "is_active": false
}
```

List users response JSON:

```json
[
  {
    "id": 1,
    "name": "Alex Brown",
    "email": "alex@example.com",
    "plan_id": 1,
    "plan": {
      "id": 1,
      "name": "Pro",
      "description": "Pro monthly plan",
      "price": 29.99,
      "currency": "USD",
      "billing_cycle": "monthly",
      "created_at": "2026-04-03T10:00:00Z",
      "updated_at": "2026-04-03T10:00:00Z"
    },
    "is_active": true,
    "created_at": "2026-04-03T10:05:00Z",
    "updated_at": "2026-04-03T10:05:00Z"
  }
]
```

### Billings

- POST /billings
- GET /billings
- GET /billings/:id
- PATCH /billings/:id/pay

Create billing request JSON:

```json
{
  "user_id": 1,
  "plan_id": 1,
  "amount": 29.99,
  "due_date": "2026-05-01T00:00:00Z",
  "description": "May subscription invoice"
}
```

Note: due_date must be RFC3339 format.

Create billing response JSON:

```json
{
  "id": 1,
  "user_id": 1,
  "user": {
    "id": 1,
    "name": "Alex Brown",
    "email": "alex@example.com",
    "plan_id": 1,
    "is_active": true,
    "created_at": "2026-04-03T10:05:00Z",
    "updated_at": "2026-04-03T10:05:00Z"
  },
  "plan_id": 1,
  "plan": {
    "id": 1,
    "name": "Pro",
    "description": "Pro monthly plan",
    "price": 29.99,
    "currency": "USD",
    "billing_cycle": "monthly",
    "created_at": "2026-04-03T10:00:00Z",
    "updated_at": "2026-04-03T10:00:00Z"
  },
  "amount": 29.99,
  "status": "pending",
  "due_date": "2026-05-01T00:00:00Z",
  "paid_at": null,
  "description": "May subscription invoice",
  "created_at": "2026-04-03T10:10:00Z",
  "updated_at": "2026-04-03T10:10:00Z"
}
```

Pay billing response JSON:

```json
{
  "id": 1,
  "user_id": 1,
  "plan_id": 1,
  "amount": 29.99,
  "status": "paid",
  "due_date": "2026-05-01T00:00:00Z",
  "paid_at": "2026-04-03T10:15:00Z",
  "description": "May subscription invoice",
  "created_at": "2026-04-03T10:10:00Z",
  "updated_at": "2026-04-03T10:15:00Z"
}
```

List billings response JSON:

```json
[
  {
    "id": 1,
    "user_id": 1,
    "plan_id": 1,
    "amount": 29.99,
    "status": "pending",
    "due_date": "2026-05-01T00:00:00Z",
    "paid_at": null,
    "description": "May subscription invoice",
    "created_at": "2026-04-03T10:10:00Z",
    "updated_at": "2026-04-03T10:10:00Z"
  }
]
```

## Error Response Format

When request data is invalid or record is not found, API returns this shape:

```json
{
  "error": "message text"
}
```

## Notes

- PATCH endpoints support partial updates.
- Users can be created with or without a plan.
- If billing amount is not sent, it is taken from the selected plan price.
- For this stage, passwords are checked as plain text (bcrypt dependency is added but not used yet).
