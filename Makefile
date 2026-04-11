.DEFAULT_GOAL := help

.PHONY: help bootstrap dev dev-frontend build down restart logs ps shell \
	migrate migrate-down generate sync-worker create-plugin config prod-config

##@ Setup
bootstrap: ## Create local env files for first run
	@./scripts/bootstrap.sh

##@ Development
dev: ## Start the dev stack (backend, db, redis)
	@./scripts/dev.sh up --build

dev-frontend: ## Start the dev stack including the optional frontend profile
	@DEV_WITH_FRONTEND=1 ./scripts/dev.sh up --build

build: ## Build the dev images
	@./scripts/dev.sh build

down: ## Stop and remove the dev stack
	@./scripts/dev.sh down

restart: ## Restart the dev stack
	@./scripts/dev.sh down
	@./scripts/dev.sh up --build

logs: ## Follow logs for the dev stack
	@./scripts/dev.sh logs -f

ps: ## List dev stack containers
	@./scripts/dev.sh ps

shell: ## Open a shell in the backend container
	@./scripts/dev.sh run --rm backend sh

##@ Backend
migrate: ## Run database migrations up
	@./scripts/migrate.sh up

migrate-down: ## Roll back one database migration
	@./scripts/migrate.sh down

generate: ## Run backend DB code generation
	@./scripts/generate.sh

sync-worker: ## Run the backend async worker process
	@./scripts/sync-worker.sh

create-plugin: ## Scaffold a local plugin folder, e.g. make create-plugin NAME=my-plugin TYPE=scaffolder
	@./scripts/create-plugin.sh \
		$(if $(NAME),--name "$(NAME)") \
		$(if $(TYPE),--type "$(TYPE)") \
		$(if $(SCOPE),--scope "$(SCOPE)") \
		$(if $(VERSION),--version "$(VERSION)") \
		$(if $(DESCRIPTION),--description "$(DESCRIPTION)") \
		$(if $(LANGUAGE),--language "$(LANGUAGE)") \
		$(if $(filter 1 true TRUE yes YES,$(FORCE)),--force) \
		$(ARGS)

##@ Compose
config: ## Render the development compose config
	@./scripts/dev.sh config

prod-config: ## Render the production compose config
	@COMPOSE_ENV=prod ./scripts/docker-build-and-run.sh config

##@ Help
help: ## Show help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z0-9_.-]+:.*?##/ { printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
