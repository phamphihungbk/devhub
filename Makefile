.DEFAULT_GOAL := help

.PHONY: help bootstrap generate-dev-cert setup-local-https dev dev-ui build build-ui down restart logs logs-worker ps shell \
	worker-up worker-down migrate migrate-down generate sync-worker create-plugin config prod-config argocd-ui argocd-token

CORE_SERVICES := backend worker nginx db redis gitea

##@ Setup
bootstrap: ## Create local env files for first run
	@./scripts/bootstrap.sh

generate-dev-cert: ## Generate local TLS certs, optionally with DOMAIN=devhub.local API_DOMAIN=api.devhub.local
	@./scripts/generate-dev-cert.sh $(if $(DOMAIN),$(DOMAIN),) $(if $(API_DOMAIN),$(API_DOMAIN),)

setup-local-https: ## Generate certs, update hosts, and trust local devhub/api certs (macOS)
	@./scripts/setup-local-https.sh $(if $(DOMAIN),$(DOMAIN),) $(if $(API_DOMAIN),$(API_DOMAIN),)

##@ Development
dev: ## Start the dev stack without UI (backend, worker, db, redis, gitea)
	@./scripts/dev.sh up --build $(CORE_SERVICES)

dev-ui: ## Start the full dev stack including frontend and nginx
	@COMPOSE_PROFILES=ui ./scripts/dev.sh up --build

build: ## Build the dev images without UI services
	@./scripts/dev.sh build $(CORE_SERVICES)

build-ui: ## Build the full dev stack including frontend and nginx
	@COMPOSE_PROFILES=ui ./scripts/dev.sh build

down: ## Stop and remove the dev stack
	@./scripts/dev.sh down

restart: ## Restart the dev stack
	@./scripts/dev.sh down
	@./scripts/dev.sh up --build $(CORE_SERVICES)

logs: ## Follow logs for the dev stack
	@./scripts/dev.sh logs -f

logs-worker: ## Follow logs for the worker service
	@./scripts/dev.sh logs -f worker

ps: ## List dev stack containers
	@./scripts/dev.sh ps

shell: ## Open a shell in the backend container
	@./scripts/dev.sh run --rm backend sh

worker-up: ## Start only the worker service and its dependencies
	@./scripts/dev.sh up -d worker

worker-down: ## Stop the worker service
	@./scripts/dev.sh stop worker

##@ Backend
migrate: ## Run database migrations up
	@./scripts/migrate.sh up

migrate-down: ## Roll back one database migration
	@./scripts/migrate.sh down

generate: ## Run backend DB code generation
	@./scripts/generate.sh

sync-worker: ## Run the backend async worker process in a one-off worker container
	@./scripts/sync-worker.sh

argocd-ui: ## Install/apply Argo CD resources if needed, then port-forward the UI
	@./scripts/argocd.sh all

argocd-token: ## Generate and print an ARGOCD_AUTH_TOKEN export line
	@./scripts/argocd.sh token

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
