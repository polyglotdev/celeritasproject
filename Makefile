test: ## Run tests
	go test -v ./...

cover: ## Run tests with coverage
	@go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out

coverage: ## Display coverage
	@go tool cover ./...

cbuild: ## build the command line tool celeritas and copies it to myapp
	@go build -o ../myapp/celeritas ./cmd/cli

help: ## Display details on all commands
	@awk 'BEGIN {FS = ":.*?##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z0-9_-]+:.*?##/ { printf "  \033[36m%-25s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: help test cover coverage cbuild