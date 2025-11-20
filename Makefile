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

migrate-up: ## Run database migrations up
	@echo "Running migrations..."
	@docker-compose exec -T postgres psql -U postgres -d worktrack < migrations/000001_create_users_table.up.sql
	@docker-compose exec -T postgres psql -U postgres -d worktrack < migrations/000002_create_tasks_table.up.sql
	@echo "Migrations completed"

migrate-down: ## Run database migrations down
	@echo "Rolling back migrations..."
	@docker-compose exec -T postgres psql -U postgres -d worktrack < migrations/000002_create_tasks_table.down.sql
	@docker-compose exec -T postgres psql -U postgres -d worktrack < migrations/000001_create_users_table.down.sql
	@echo "Rollback completed"

deps: ## Download dependencies
	go mod download
	go mod tidy

lint: ## Run linter
	golangci-lint run

fmt: ## Format code
	go fmt ./...
