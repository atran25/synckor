default: help

.PHONY: help
help: ## Show help for each of the Makefile recipes.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)


.PHONY: schema
schema: ## Generate the schema for the database
	sqlite3 db/sqlite/db.sqlite '.schema' > db/sqlite/schema.sql

.PHONY: tools
tools: ## Generate the oapi-codegen and sqlc code
	go generate ./tools/tools.go

.PHONY: check
check: ## Run the linters
	go mod tidy
	go vet ./...
	go fmt ./...

.PHONY: migrate
migrate: ## Run the migrations
	go run cmd/migrate/main.go

.PHONY: run
run: schema check tools migrate ## Run the application
	go run cmd/server/main.go