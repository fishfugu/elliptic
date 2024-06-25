.PHONY: build
build:	## Build Elliptic Curve Data Viz Tool
# TODO: Scaling / Translation: Adjust mapping of mathematical coords to screen coords to ensure curve fits in window / maintains aspect ratio
# TODO: Interactive: Zooming / panning to explore parts of curve
# Labeling: Optionally, labels / different colors to highlight properties / points on curve, such as inflection, zeros, etc.
	@echo "Building Elliptic Curve Data Viz Tool..."
	GOOS=darwin GOARCH=amd64 go build -o ../../bin/ecviz
	@echo "Build finished"

.PHONY: run
run: build	## Run Elliptic Curve Data Viz Tool (after doing a build)
	@echo "Running Elliptic Curve Data Viz Tool..."
	../../bin/ecviz &
	@echo "Elliptic Curve Data Viz Tool running"

.PHONY: clean
ecviz-clean:	## Clean Elliptic Curve Data Viz Tool binary
	@echo "Cleaning Elliptic Curve Data Viz Tool binary..."
	rm ../../bin/ecviz
	@echo "Elliptic Curve Data Viz Tool binary gone"

.DEFAULT_GOAL := help
.PHONY: help
help:	## Show this help
	@echo 'Usage: make <TARGETS>'
	@awk 'BEGIN {FS = ":.*##"; } /^[a-zA-Z0-9_-]+:.*?##/ { printf "  \033[36m%-30s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
