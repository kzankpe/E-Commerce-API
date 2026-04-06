# Makefile
.PHONY: help setup dev-up dev-down db-migrate db-seed clean build run test lint fmt

# Variables
BINARY_NAME=ecommerce-api
GO=go
GOFLAGS=-v

help:
	@echo "Available commands:"
	@echo "  make setup          - Setup development environment"
	@echo "  make dev-up         - Start Docker containers"
	@echo "  make dev-down       - Stop Docker containers"
	@echo "  make db-migrate     - Run database migrations"
	@echo "  make db-seed        - Seed database with sample data"
	@echo "  make run            - Run the application"
	@echo "  make build          - Build the application"
	@echo "  make test           - Run tests"
	@echo "  make clean          - Clean build artifacts"
	@echo "  make lint           - Run linter"
	@echo "  make fmt            - Format code"

setup:
	@echo "Setting up development environment..."
	@cp .env.example .env
	@$(GO) mod download
	@$(GO) mod tidy
	@echo "Setup complete! Run 'make dev-up' to start containers"

dev-up:
	@echo "Starting Docker containers..."
	docker-compose up -d
	@echo "Waiting for services to be ready..."
	@sleep 5
	@echo "Services started successfully!"
	@docker-compose ps

dev-down:
	@echo "Stopping Docker containers..."
	docker-compose down

dev-logs:
	docker-compose logs -f

db-migrate:
	@echo "Running database migrations..."
	@psql -h localhost -U postgres -d ecommerce_db -f internal/infrastructure/database/migrations/001_create_users_table.sql
	@psql -h localhost -U postgres -d ecommerce_db -f internal/infrastructure/database/migrations/002_create_token_blacklist_table.sql
	@echo "Migrations completed!"

db-seed:
	@echo "Seeding database..."
	@$(GO) run scripts/db_seed.go
	@echo "Database seeded!"

build:
	@echo "Building $(BINARY_NAME)..."
	@$(GO) build $(GOFLAGS) -o bin/$(BINARY_NAME) ./cmd/server

run: build
	@echo "Running $(BINARY_NAME)..."
	@./bin/$(BINARY_NAME)

dev:
	@echo "Running in development mode with hot reload..."
	@which air > /dev/null || go install github.com/cosmtrek/air@latest
	@air

test:
	@echo "Running tests..."
	@$(GO) test -v -cover -timeout 30s ./...

test-coverage:
	@echo "Running tests with coverage..."
	@$(GO) test -v -coverprofile=coverage.out ./...
	@$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

lint:
	@echo "Running linter..."
	@which golangci-lint > /dev/null || go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@golangci-lint run ./...

fmt:
	@echo "Formatting code..."
	@$(GO) fmt ./...

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@$(GO) clean
	@echo "Clean complete!"

deps:
	@echo "Downloading dependencies..."
	@$(GO) mod download
	@$(GO) mod tidy

install-tools:
	@echo "Installing development tools..."
	@go install github.com/cosmtrek/air@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@echo "Tools installed!"
