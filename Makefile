DB_URL ?= postgres://postgres:postgres@localhost:5432/subscription_management?sslmode=disable
MIGRATIONS_DIR ?= migrations
MIGRATE_BIN ?= migrate

.PHONY: migrate-install migrate-up migrate-down migrate-down1 migrate-force migrate-version migrate-create

migrate-install:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

migrate-up:
	$(MIGRATE_BIN) -path $(MIGRATIONS_DIR) -database "$(DB_URL)" up

migrate-down:
	$(MIGRATE_BIN) -path $(MIGRATIONS_DIR) -database "$(DB_URL)" down

migrate-down1:
	$(MIGRATE_BIN) -path $(MIGRATIONS_DIR) -database "$(DB_URL)" down 1

migrate-force:
	$(MIGRATE_BIN) -path $(MIGRATIONS_DIR) -database "$(DB_URL)" force $(VERSION)

migrate-version:
	$(MIGRATE_BIN) -path $(MIGRATIONS_DIR) -database "$(DB_URL)" version

migrate-create:
	$(MIGRATE_BIN) create -ext sql -dir $(MIGRATIONS_DIR) -seq $(NAME)
