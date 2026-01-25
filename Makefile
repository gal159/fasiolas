.PHONY: help build run test clean migrate-up migrate-down migrate-create docker-up docker-down

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the application
	go build -o bin/server cmd/server/main.go

run: ## Run the application
	go run cmd/server/main.go

test: ## Run tests
	go test -v -cover ./...

clean: ## Clean build artifacts
	rm -rf bin/
	go clean

migrate-up: ## Run database migrations up
	migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/fasiolas_game?sslmode=disable" up

migrate-down: ## Run database migrations down
	migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/fasiolas_game?sslmode=disable" down

migrate-create: ## Create a new migration (usage: make migrate-create NAME=migration_name)
	migrate create -ext sql -dir migrations -seq $(NAME)

docker-up: ## Start Docker containers
	docker-compose up -d

docker-down: ## Stop Docker containers
	docker-compose down

docker-logs: ## Show Docker logs
	docker-compose logs -f

install-tools: ## Install development tools
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	go install github.com/swaggo/swag/cmd/swag@latest

swagger: ## Generate Swagger documentation
	swag init -g cmd/server/main.go -o docs

