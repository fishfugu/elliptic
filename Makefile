##@ BUILD

.PHONY: build-all
build-all: build-bigmath build-finitefield build-ecviz	## Build the project - all necessary components

.PHONY: build-bigmath
build-bigmath:	## Build bigmath executable
	@echo "Building bigmath executable..."
	GOOS=darwin GOARCH=amd64 go build -o bin/bigmath ./cmd/bigmath
	@echo "Build of bigmath complete."

.PHONY: build-finitefield
build-finitefield:	## Build finitefield executable
	@echo "Building finitefield executable..."
	GOOS=darwin GOARCH=amd64 go build -o bin/finitefield ./cmd/finitefield
	@echo "Build of finitefield complete."

.PHONY: build-ecviz
build-ecviz:	## Build Elliptic Curve Data Viz Tool
# TODO: Scaling / Translation: Adjust mapping of mathematical coords to screen coords to ensure curve fits in window / maintains aspect ratio
# TODO: Interactive: Zooming / panning to explore parts of curve
# Labeling: Optionally, labels / different colors to highlight properties / points on curve, such as inflection, zeros, etc.
	@echo "Building Elliptic Curve Data Viz Tool..."
	GOOS=darwin GOARCH=amd64 go build -o bin/ecviz ./cmd/ecviz
	@echo "Build finished"

##@ RUN

.PHONY: run-bigmath
run-bigmath: build-bigmath test-quiet	## Run bigmath binary after building it - ensuring latest build is executed - running tests first
	@echo
	@echo "************************************************************"
	@echo "✓✓✓✓✓ -- Bigmath built, and tested - continuing on to run..."
	@echo "************************************************************"
	@echo
	./bin/bigmath
	@echo "Run of bigmath complete."

.PHONY: run-finitefield
run-finitefield: build-finitefield test-quiet	## Run bigmath binary after building it - ensuring latest build is executed - running tests first
	@echo
	@echo "****************************************************************"
	@echo "✓✓✓✓✓ -- Finitefield built, and tested - continuing on to run..."
	@echo "****************************************************************"
	@echo
	./bin/finitefield
	@echo "Run of finitefield complete."

.PHONY: run-ecviz
run-ecviz: build-ecviz test-quiet	## Run Elliptic Curve Data Viz Tool (after doing a build)
	@echo
	@echo "****************************************************************"
	@echo "✓✓✓✓✓ -- EC Viz Tool built, and tested - continuing on to run..."
	@echo "****************************************************************"
	@echo
	./bin/ecviz &
	@echo "Elliptic Curve Data Viz Tool running"

##@ TEST and CLEAN

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
test-drive: help build-all run-bigmath run-finitefield run-ecviz test clean	## Run through all (appropriate) make file commands - just to take it for a test drive (check I haven't done stupidity)
	@echo "******************************************************************"
	@echo "✓✓✓✓✓ -- Seem to have got to end of test-dive without fatal errors"
	@echo "******************************************************************"

##@ HELPER

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

.DEFAULT_GOAL := help
.PHONY: help
help:	## Show this help
	@echo 'Usage: make <TARGETS>'
	@awk 'BEGIN {FS = ":.*##"; } /^[a-zA-Z0-9_-]+:.*?##/ { printf "  \033[36m%-30s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
