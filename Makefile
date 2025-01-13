.DEFAULT_GOAL := help

help: ## Displays all the available commands
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
.PHONY: help

app: ## Runs the application
	@go run cmd/app/main.go
.PHONY: app

migrate: ## Migrates the database to the most recent version
	@go run cmd/migrate/main.go -command up
.PHONY: migrate

migrate-down: ## Roll back the migration version by 1
	@go run cmd/migrate/main.go -command down
.PHONY: migrate-down

migrate-reset: ## Roll back all the migrations
	@go run cmd/migrate/main.go -command reset
.PHONY: migrate-reset

format: ## Runs go fmt and go vet
	@go fmt ./... && go vet -v ./...
.PHONY: format

test: ## Runs tests
	@go test \
		-shuffle=on \
		-count=1 \
		-short \
		-timeout=5m \
		./...
.PHONY: test