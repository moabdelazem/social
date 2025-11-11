include .env

# Configuration
MIGRATIONS_DIR ?= ./cmd/migrate/migrations

# Development
.PHONY: watch
watch:
	@air

# Build
.PHONY: build
build:
	@echo "Building the binary in the bin dir"
	@go build -o ./bin/main ./cmd/api

# Database Seeding
.PHONY: seed
seed:
	@echo "Seeding the database..."
	@go run cmd/seed/main.go

# Database Migrations
.PHONY: migrate-up
migrate-up:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_ADDR)" up

.PHONY: migrate-down
migrate-down:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_ADDR)" down

.PHONY: migrate-create
migrate-create:
ifndef NAME
	$(error NAME is not set. Usage: make migrate-create NAME=create_users_table)
endif
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(NAME)

.PHONY: migrate-version
migrate-version:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_ADDR)" version

.PHONY: migrate-force
migrate-force:
ifndef VERSION
	$(error VERSION is not set. Usage: make migrate-force VERSION=1)
endif
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_ADDR)" force $(VERSION)
