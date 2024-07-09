help: ## This help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

run-sample-echo: ## Run sample with Echo Http Framework
	@go run cmd/sample-echo/main.go

run-sample-gin: ## Run sample with Gin Http Framework
	@go run cmd/sample-gin/main.go

run-sample-fiber: ## Run sample with Fiber Http Framework
	@go run cmd/sample-fiber/main.go