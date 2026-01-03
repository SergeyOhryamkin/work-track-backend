.PHONY: help run build test clean docker-build docker-up docker-down migrate-up migrate-down migrate-create

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

run: ## Run the application locally
	go run cmd/api/main.go

build: ## Build the application
	go build -o bin/api cmd/api/main.go

test: ## Run tests
	go test -v ./...

test-coverage: ## Run tests with coverage
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

clean: ## Clean build artifacts
	rm -rf bin/
	rm -f coverage.out coverage.html

docker-build: ## Build Docker image
	docker build -t worktrack-backend:latest .

docker-up: ## Start services with Docker Compose
	docker-compose up -d

docker-down: ## Stop services with Docker Compose
	docker-compose down

docker-logs: ## View Docker Compose logs
	docker-compose logs -f

migrate-up: ## Run database migrations up (not needed for SQLite as it runs on start)
	@echo "Migrations are handled automatically on application startup"

migrate-down: ## Run database migrations down (not implemented for SQLite)
	@echo "Manual rollback not implemented for SQLite"

db-reset: ## Reset the SQLite database (Deletes the .db file)
	@echo "Warning: This will delete all data. Are you sure? [y/N] " && read ans && [ $${ans:-N} = y ]
	rm -f worktrack.db worktrack.db-shm worktrack.db-wal
	rm -f data/worktrack.db data/worktrack.db-shm data/worktrack.db-wal
	@echo "Database deleted. It will be recreated on next startup."

deps: ## Download dependencies
	go mod download
	go mod tidy

lint: ## Run linter
	golangci-lint run

fmt: ## Format code
	go fmt ./...
