ACCOUNTS_DB_URL ?= postgres://postgres:postgres@localhost:5432/subscription_management_accounts?sslmode=disable
BILLING_DB_URL ?= postgres://postgres:postgres@localhost:5432/subscription_management_billing?sslmode=disable

ACCOUNTS_MIGRATIONS_DIR ?= accounts-service/migrations
BILLING_MIGRATIONS_DIR ?= billing-service/migrations
MIGRATE_BIN ?= migrate

.PHONY: migrate-install \
	migrate-accounts-up migrate-accounts-down migrate-accounts-down1 migrate-accounts-force migrate-accounts-version migrate-accounts-create \
	migrate-billing-up migrate-billing-down migrate-billing-down1 migrate-billing-force migrate-billing-version migrate-billing-create

migrate-install:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

migrate-accounts-up:
	$(MIGRATE_BIN) -path $(ACCOUNTS_MIGRATIONS_DIR) -database "$(ACCOUNTS_DB_URL)" up

migrate-accounts-down:
	$(MIGRATE_BIN) -path $(ACCOUNTS_MIGRATIONS_DIR) -database "$(ACCOUNTS_DB_URL)" down

migrate-accounts-down1:
	$(MIGRATE_BIN) -path $(ACCOUNTS_MIGRATIONS_DIR) -database "$(ACCOUNTS_DB_URL)" down 1

migrate-accounts-force:
	$(MIGRATE_BIN) -path $(ACCOUNTS_MIGRATIONS_DIR) -database "$(ACCOUNTS_DB_URL)" force $(VERSION)

migrate-accounts-version:
	$(MIGRATE_BIN) -path $(ACCOUNTS_MIGRATIONS_DIR) -database "$(ACCOUNTS_DB_URL)" version

migrate-accounts-create:
	$(MIGRATE_BIN) create -ext sql -dir $(ACCOUNTS_MIGRATIONS_DIR) -seq $(NAME)

migrate-billing-up:
	$(MIGRATE_BIN) -path $(BILLING_MIGRATIONS_DIR) -database "$(BILLING_DB_URL)" up

migrate-billing-down:
	$(MIGRATE_BIN) -path $(BILLING_MIGRATIONS_DIR) -database "$(BILLING_DB_URL)" down

migrate-billing-down1:
	$(MIGRATE_BIN) -path $(BILLING_MIGRATIONS_DIR) -database "$(BILLING_DB_URL)" down 1

migrate-billing-force:
	$(MIGRATE_BIN) -path $(BILLING_MIGRATIONS_DIR) -database "$(BILLING_DB_URL)" force $(VERSION)

migrate-billing-version:
	$(MIGRATE_BIN) -path $(BILLING_MIGRATIONS_DIR) -database "$(BILLING_DB_URL)" version

migrate-billing-create:
	$(MIGRATE_BIN) create -ext sql -dir $(BILLING_MIGRATIONS_DIR) -seq $(NAME)
