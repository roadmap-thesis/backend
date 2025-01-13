.DEFAULT_GOAL := help

help: ## Displays all the available commands
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

app: ## Runs the application
	go run cmd/app/main.go

migrate: ## Migrates the database to the most recent version
	go run cmd/migrate/main.go -command up

migrate-down: ## Roll back the migration version by 1
	go run cmd/migrate/main.go -command down

migrate-reset: ## Roll back all the migrations
	go run cmd/migrate/main.go -command reset

format: ## Runs go fmt and go vet
	go fmt ./... && go vet ./...

test: ## Runs tests
	go test -v ./...

bench: ## Runs benchmark
	go test -bench=. ./... -benchmem