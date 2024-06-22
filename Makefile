.DEFAULT_GOAL := help

.PHONY: help
help:	## Show this help
	@echo 'Usage: make <TARGETS>'
	@awk 'BEGIN {FS = ":.*##"; } /^[a-zA-Z0-9_-]+:.*?##/ { printf "  \033[36m%-30s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: build-bigmath
build-bigmath:	## Build bigmath executable
	@echo "Building bigmath executable..."
	go build -o bin/bigmath ./cmd/bigmath
	@echo "Build complete."

.PHONY: build-all
build-all: build-bigmath	## Build the project - all necessary components

.PHONY: run-bigmath
run-bigmath: build-bigmath test-quiet	## Run bigmath binary after building it - ensuring latest build is executed - running tests first
	@echo
	@echo "************************************************************"
	@echo "✓✓✓✓✓ -- Bigmath built, and tested - continuing on to run..."
	@echo "************************************************************"
	@echo
	./bin/bigmath
	@echo "Run complete."

.PHONY: test
test:	## Run unit tests for all packages under pkg
	@echo "Running tests..."
	go clean -testcache
	go test -v ./pkg/... -coverprofile=coverage.out
	go tool cover -html=coverage.out

.PHONY: clean
clean:	## Remove binaries and any temporary files
	@echo "Cleaning up..."
	rm -rf bin
	@echo "Clean complete."

.PHONY: test-drive
test-drive: help build-bigmath build-all run-bigmath test clean	## Run through all (appropriate) make file commands - just to take it for a test drive (check I haven't done stupidity)
	@echo "******************************************************************"
	@echo "✓✓✓✓✓ -- Seem to have got to end of test-dive without fatal errors"
	@echo "******************************************************************"

## Helper make targets - not to be run as part of test-drive
## They either replicate other stuff or are just inappropriate for running without thinking about it

.PHONY: test-quiet
test-quiet:	## Run unit tests for all packages under pkg - but quietly - quits at first error
	go clean -testcache
	go test -failfast ./pkg/...

.PHONY: test-verbose
test-verbose:	## Run unit tests for all packages under pkg - in verbose mode
# TODO: implement verbose mode more conistently
	@echo "Running tests in verbose mode..."
	go clean -testcache
	go test -v ./pkg/... -coverprofile=coverage.out -args -verbose
	go tool cover -html=coverage.out