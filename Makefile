default: help

.PHONY: help
help: ## Show help for each of the Makefile recipes.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: test
test: ## Run the tests
	go test -v ./...

.PHONY: lint
lint: ## Run the linters
	golangci-lint run

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
run: migrate schema tools check ## Run the application
	go run cmd/server/main.go

.PHONY: docker-build
docker-build: ## Build the docker image
	docker build -t synckor .

.PHONY: docker-run
docker-run: ## Run the docker image
	docker run -it --network host -p 8050:8050 -e PORT=8050 -e REGISTRATION_ENABLED=True -e LITESTREAM_ACCESS_KEY_ID=minioadmin -e LITESTREAM_SECRET_ACCESS_KEY=minioadmin -e REPLICA_URL=s3://synckor-bkt.localhost:9000/db.sqlite synckor