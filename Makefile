# Makefile for brevity-go project

# Database connection string for PostgreSQL
DB_URL=postgres://postgres:postgres@localhost:5432/brevity?sslmode=disable

# Migration directory
MIGRATIONS_DIR=db/migrations

# Default target
.PHONY: all
all: help

# Help message
.PHONY: help
help:
	@echo "Available commands:"
	@echo "  make migrate-up      - Run all pending migrations"
	@echo "  make migrate-down    - Rollback the last migration"
	@echo "  make migrate-create  - Create a new migration (requires NAME=migration_name)"
	@echo "  make migrate-force   - Force migration version (requires VERSION=x)"
	@echo "  make migrate-version - Show current migration version"
	@echo "  make migrate-goto    - Migrate to specific version (requires VERSION=x)"
	@echo "  make migrate-drop    - Drop everything in the database"

# Run migrations up
.PHONY: migrate-up
migrate-up:
	@echo "Running migrations up..."
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" up

# Run migrations down (rollback)
.PHONY: migrate-down
migrate-down:
	@echo "Rolling back the last migration..."
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" down 1

# Create a new migration
.PHONY: migrate-create
migrate-create:
	@if [ -z "$(NAME)" ]; then \
		echo "Error: NAME is required. Use 'make migrate-create NAME=migration_name'"; \
		exit 1; \
	fi
	@echo "Creating migration $(NAME)..."
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(NAME)

# Force migration version
.PHONY: migrate-force
migrate-force:
	@if [ -z "$(VERSION)" ]; then \
		echo "Error: VERSION is required. Use 'make migrate-force VERSION=x'"; \
		exit 1; \
	fi
	@echo "Forcing migration version to $(VERSION)..."
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" force $(VERSION)

# Show current migration version
.PHONY: migrate-version
migrate-version:
	@echo "Current migration version:"
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" version

# Migrate to a specific version
.PHONY: migrate-goto
migrate-goto:
	@if [ -z "$(VERSION)" ]; then \
		echo "Error: VERSION is required. Use 'make migrate-goto VERSION=x'"; \
		exit 1; \
	fi
	@echo "Migrating to version $(VERSION)..."
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" goto $(VERSION)

# Drop everything in the database
.PHONY: migrate-drop
migrate-drop:
	@echo "WARNING: This will drop all tables in the database."
	@echo "Are you sure? [y/N]"
	@read -r CONFIRM; \
	if [ "$$CONFIRM" = "y" ] || [ "$$CONFIRM" = "Y" ]; then \
		echo "Dropping all tables..."; \
		migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" drop -f; \
	else \
		echo "Operation cancelled."; \
	fi

# Docker compose commands
.PHONY: up
up:
	docker-compose up -d

.PHONY: down
down:
	docker-compose down

.PHONY: restart
restart:
	docker-compose restart

.PHONY: logs
logs:
	docker-compose logs -f 