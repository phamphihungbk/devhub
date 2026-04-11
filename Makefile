.PHONY: set-permissions start restart npm npx run shell

set-permissions:
	@chmod +x ./scripts/docker-build-and-run.sh

##@ Operations
start: set-permissions ## Start the service
	@./scripts/docker-build-and-run.sh start $(ARGS)

build: set-permissions ## Build docker image
	@./scripts/docker-build-and-run.sh build

down: set-permissions ## Take down the service
	@./scripts/docker-build-and-run.sh down

restart: set-permissions ## Restart the service
	@./scripts/docker-build-and-run.sh restart

shell: ## Start a shell session
	@./scripts/docker-build-and-run.sh shell

##@ Help
help: ## Show help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z0-9_-]+:.*?##/ { printf "  \033[36m%-27s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)


# Catch any arguments not recognized by the Makefile
%:
	@: