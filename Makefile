##@ BUILD

.PHONY: build-all
build-all: build-bigmath build-finitefield build-ecvis build-cli	## Build the project - all necessary components

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

.PHONY: build-ecvis
build-ecvis:	## Build Elliptic Curve Data Vis Tool
# TODO: Scaling / Translation: Adjust mapping of mathematical coords to screen coords to ensure curve fits in window / maintains aspect ratio
# TODO: Interactive: Zooming / panning to explore parts of curve
# Labeling: Optionally, labels / different colors to highlight properties / points on curve, such as inflection, zeros, etc.
	@echo "Building Elliptic Curve Data Vis Tool..."
	GOOS=darwin GOARCH=amd64 go build -o bin/ecvis ./cmd/ecvis
	@echo "Build finished"

.PHONY: build-cli
build-cli:	## Build ellipticcurvecli executable
	@echo "Building ellipticcurvecli executable..."
	GOOS=darwin GOARCH=amd64 go build -o bin/ellipticcurvecli ./cmd/ellipticcurvecli
	@echo "Build of ellipticcurvecli complete."

.PHONY: build-pns
build-pns:	## Build prime nummber system (psumsys - pns) executable
	@echo "Building prime nummber system (psumsys - pns) executable..."
	GOOS=darwin GOARCH=amd64 go build -o bin/pnumsys ./cmd/pnumsys
	@echo "Build of prime nummber system (psumsys - pns) complete."

##@ RUN

.PHONY: run-bigmath
run-bigmath: build-bigmath test-quiet	## Run bigmath binary - ensuring latest build is executed - running tests first
	@echo
	@echo "************************************************************"
	@echo "✓✓✓✓✓ -- Bigmath built, and tested - continuing on to run..."
	@echo "************************************************************"
	@echo
	./bin/bigmath
	@echo
	@echo "Run of bigmath complete."

.PHONY: run-finitefield
run-finitefield: build-finitefield test-fast	## Run finitefield binary - ensuring latest build is executed - running tests first
	@echo
	@echo "****************************************************************"
	@echo "✓✓✓✓✓ -- Finitefield built, and tested - continuing on to run..."
	@echo "****************************************************************"
	@echo
	./bin/finitefield
	@echo
	@echo "Run of finitefield complete."

.PHONY: run-ecvis
run-ecvis: build-ecvis test-fast	## Run Elliptic Curve Data Vis Tool - ensuring latest build is executed - running tests first
	@echo
	@echo "****************************************************************"
	@echo "✓✓✓✓✓ -- EC Vis Tool built, and tested - continuing on to run..."
	@echo "****************************************************************"
	@echo
	./bin/ecvis &
	@echo
	@echo "Elliptic Curve Data Vis Tool running"

.PHONY: run-cli
run-cli: build-cli test-fast	## Run Elliptic Curve CLI - ensuring latest build is executed - running tests first
	@echo
	@echo "***********************************************************************"
	@echo "✓✓✓✓✓ -- Elliptic Curve CLI built, and tested - continuing on to run..."
	@echo "***********************************************************************"
	@echo
	./bin/ellipticcurvecli
	@echo
	@echo "Run of Elliptic Curve CLI (without visualise) complete"

##@ TEST and CLEAN

.PHONY: test
test:	## Run unit tests for all packages under pkg
	@echo "Running tests..."
	go clean -testcache
	go test -v -timeout=10m -bench=. ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	open coverage.html

.PHONY: test-performance
test-performance:	## run unit tests and extract duration / avg perforamcnce stats
	go test -v -timeout=10m -bench=. -count=1 ./pkg/ellipticcurve/... | go run ignorescripts/parse_test_stats/parse_test_stats.go 

.PHONY: clean
clean:	## Remove binaries and any temporary files
	@echo "Cleaning up..."
	rm -rf bin
	@echo "Clean complete."

.PHONY: test-drive
test-drive: help build-all run-bigmath run-finitefield run-ecvis test clean	## Run through all (appropriate) make file commands - just to take it for a test drive (check I haven't done stupidity)
	@echo "******************************************************************"
	@echo "✓✓✓✓✓ -- Seem to have got to end of test-dive without fatal errors"
	@echo "******************************************************************"

##@ HELPER

## Helper make targets - not to be run as part of test-drive
## They either replicate other stuff or are just inappropriate for running without thinking about it

.PHONY: test-fast
test-fast:	## Run unit tests for all packages under pkg - but quietly - quits at first error
	go clean -testcache
	go test -failfast -timeout=10m ./pkg/...

.PHONY: test-quiet
test-quiet:	## Run unit tests for all packages under pkg - but quietly - quits at first error
	go clean -testcache
	go test -failfast -timeout=10m -bench=. ./pkg/...

.PHONY: test-verbose
test-verbose:	## Run unit tests for all packages under pkg - in verbose mode
# TODO: implement verbose mode more conistently
	@echo "Running tests in verbose mode..."
	go clean -testcache
	go test -v -timeout=10m -bench=. ./pkg/... -coverprofile=coverage.out -args
	go tool cover -html=coverage.out -o coverage.html
	open coverage.html

.DEFAULT_GOAL := help
.PHONY: help
help:	## Show this help
	@echo 'Usage: make <TARGETS>'
	@awk 'BEGIN {FS = ":.*##"; } /^[a-zA-Z0-9_-]+:.*?##/ { printf "  \033[36m%-30s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
