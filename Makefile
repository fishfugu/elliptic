.DEFAULT_GOAL := help

.PHONY: help
help: ## Show this help
	@echo
	@echo 'Usage: make <TARGETS>'
	@echo
	@awk 'BEGIN {FS = ":.*##"; } /^[a-zA-Z0-9_-]+:.*?##/ { printf "  \033[36m%-30s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
	@echo

.PHONY: all
all: build ## Default target for no arguments

.PHONY: build
build: build-bigmath ## Build the project - all necessary components

.PHONY: build-bigmath
build-bigmath: ## Build bigmath executable
	@echo "Building bigmath executable..."
	go build -o bin/bigmath ./cmd/bigmath
	@echo "Build complete."

.PHONY: test
test: ## Run unit tests for all packages under pkg
	@echo "Running tests..."
	go test -v ./pkg/...

.PHONY: clean
clean: ## Remove binaries and any temporary files
	@echo "Cleaning up..."
	rm -rf bin
	@echo "Clean complete."

.PHONY: run
run: build-bigmath ## Run bigmath binary after building it - ensuring latest build is executed
	@echo "Running bigmath..."
	./bin/bigmath
