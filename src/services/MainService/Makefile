# Executables (local)
DOCKER_COMP = docker-compose

# Docker containers
APP_CONT = $(DOCKER_COMP) exec app

help: ## Outputs this help screen
	@grep -E '(^[a-zA-Z0-9_-]+:.*?##.*$$)|(^##)' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}{printf "\033[32m%-30s\033[0m %s\n", $$1, $$2}' | sed -e 's/\[32m##/[33m/'

## —— Docker
build: ## Builds the Docker images
	@$(DOCKER_COMP) build --pull

up: ## Start the docker hub in detached mode (no logs)
	@$(DOCKER_COMP) up --detach

start: build up ## Build and start the containers

down: ## Stop the docker hub
	@$(DOCKER_COMP) down --remove-orphans

prune: ## Prune the docker hub
	@docker image prune

logs: ## Show live logs
	@$(DOCKER_COMP) logs --tail=0 --follow

sh: ## Open a shell in the app container
	@$(DOCKER_COMP) exec app sh

bash: ## Open a bash in the app container
	@$(DOCKER_COMP) exec app bash

migration: ## Run migration
	@go run Cmd/main.go -migrate

coverage:
	@go test -coverprofile=coverage.out ./... -v
	@go tool cover -html=coverage.out -o coverage.html

drop: ## Drop database
	@go run Cmd/main.go -drop

es:
	@go run Cmd/main.go -esSetup

bulk:
	@go run Cmd/main.go -bulk

init_swagger: ## Init swagger
	@swag init --parseDependency --parseInternal

swagger: ## Generate swagger docs
	@rm -r docs
	@swag init --parseDependency --parseInternal

