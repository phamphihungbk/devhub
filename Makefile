.DEFAULT_GOAL := help

.PHONY: help bootstrap generate-dev-cert setup-local-https \
	backend-up backend-down backend-watch frontend-up frontend-down frontend-watch \
	build-backend build-frontend up down logs logs-worker logs-runner logs-frontend ps shell \
	worker-up worker-down runner-up runner-down migrate migrate-down migrate-force new-migration seed-data generate sync-worker plugin-scan create-plugin config prod-config argocd-ui argocd-token minikube-registry

BACKEND_SERVICES := backend worker db redis devhub-registry gitea gitea-runner
FRONTEND_SERVICES := frontend nginx

##@ Setup
bootstrap: ## Create local env files for first run
	@./scripts/bootstrap.sh

generate-dev-cert: ## Generate local TLS certs, optionally with DOMAIN=devhub.local API_DOMAIN=api.devhub.local
	@./scripts/generate-dev-cert.sh $(if $(DOMAIN),$(DOMAIN),) $(if $(API_DOMAIN),$(API_DOMAIN),)

setup-local-https: ## Generate certs, update hosts, and trust local devhub/api certs (macOS)
	@./scripts/setup-local-https.sh $(if $(DOMAIN),$(DOMAIN),) $(if $(API_DOMAIN),$(API_DOMAIN),)

##@ Development
backend-up: ## Start backend services without UI or nginx
	@COMPOSE_PROFILES=ci ./scripts/dev.sh up --build $(BACKEND_SERVICES)

backend-watch: backend-up ## Start support services, then run the Go backend locally with Air hot reload
	@./scripts/backend-watch.sh

frontend-up: ## Start the UI stack; nginx will also bring up backend dependencies it needs
	@COMPOSE_PROFILES=ui ./scripts/dev.sh up --build $(FRONTEND_SERVICES)

frontend-watch: frontend-up ## Start and follow the frontend dev stack (frontend + nginx) for UI work
	@./scripts/frontend-watch.sh

build-backend: ## Build backend-side dev images without UI services
	@COMPOSE_PROFILES=ci ./scripts/dev.sh build $(BACKEND_SERVICES)

build-frontend: ## Build only the frontend UI services
	@COMPOSE_PROFILES=ui ./scripts/dev.sh build $(FRONTEND_SERVICES)

up: ## Start full stack services
	@COMPOSE_PROFILES=ui,ci ./scripts/dev.sh up --build

down: ## Stop and remove the dev stack
	@./scripts/dev.sh down

backend-down: ## Stop backend services without touching the UI profile
	@COMPOSE_PROFILES=ci ./scripts/dev.sh stop $(BACKEND_SERVICES)

frontend-down: ## Stop the frontend and nginx services
	@COMPOSE_PROFILES=ui ./scripts/dev.sh stop $(FRONTEND_SERVICES)

logs: ## Follow logs for the dev stack
	@./scripts/dev.sh logs -f

logs-worker: ## Follow logs for the worker service
	@./scripts/dev.sh logs -f worker

logs-frontend: ## Follow logs for the frontend and nginx services
	@COMPOSE_PROFILES=ui ./scripts/dev.sh logs -f frontend nginx

logs-runner: ## Follow logs for the Gitea Actions runner
	@COMPOSE_PROFILES=ci ./scripts/dev.sh logs -f gitea-runner

ps: ## List dev stack containers
	@./scripts/dev.sh ps

shell: ## Open a shell in the backend container
	@./scripts/dev.sh run --rm backend sh

worker-up: ## Start only the worker service and its dependencies
	@./scripts/dev.sh up -d worker

worker-down: ## Stop the worker service
	@./scripts/dev.sh stop worker

runner-up: ## Start the Gitea Actions runner
	@COMPOSE_PROFILES=ci ./scripts/dev.sh up -d gitea-runner

runner-down: ## Stop the Gitea Actions runner
	@COMPOSE_PROFILES=ci ./scripts/dev.sh stop gitea-runner

##@ Backend
migrate: ## Run database migrations up
	@./scripts/migrate.sh up

migrate-down: ## Roll back one database migration
	@./scripts/migrate.sh down

migrate-force: ## Force-set migration version; defaults to -1, or pass VERSION=<n>
	@FORCE_VERSION=$(if $(VERSION),$(VERSION),-1) ./scripts/migrate.sh

new-migration: ## Create a new migration pair, e.g. make new-migration NAME=add_services_table
	@test -n "$(NAME)" || (echo "Usage: make new-migration NAME=add_services_table" && exit 1)
	@./scripts/dev.sh run --rm backend go run . new-migration "$(NAME)"

seed-data: ## Seed RBAC roles, permissions, and role-permission mappings
	@./scripts/dev.sh run --rm backend go run . seed-data

generate: ## Run backend DB code generation
	@./scripts/generate.sh

sync-worker: ## Run the backend async worker process in a one-off worker container
	@./scripts/sync-worker.sh

plugin-scan: ## Scan plugin manifests and upsert them into the plugins table
	@./scripts/plugin-scan.sh $(ARGS)

argocd-ui: ## Install/apply Argo CD resources if needed, then port-forward the UI
	@./scripts/argocd.sh all

argocd-token: ## Generate and print an ARGOCD_AUTH_TOKEN export line
	@./scripts/argocd.sh token

minikube-registry: ## Recreate minikube with host.minikube.internal:5001 allowed as an insecure registry
	@./scripts/dev.sh up -d --build devhub-registry
	@minikube delete
	@minikube start --insecure-registry="host.minikube.internal:5001"
	@minikube ssh -- 'nc -vz host.minikube.internal 5001'

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
