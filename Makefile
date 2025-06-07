# Variables
BINARY_NAME=assetdownloader
OUTPUT_DIR=bin
VERSION=$(shell git describe --tags --always --dirty)
BUILD_DATE=$(shell date +%Y-%m-%dT%H:%M:%S)
LDFLAGS=-X 'main.Version=$(VERSION)' -X 'main.BuildDate=$(BUILD_DATE)'

# Default target
.PHONY: all
all: build

# Build program
.PHONY: build
build:
	@mkdir -p $(OUTPUT_DIR)
	@echo "Building $(BINARY_NAME) version $(VERSION)..."
	GOOS=windows go build -ldflags="$(LDFLAGS)" -o "$(OUTPUT_DIR)/$(BINARY_NAME)_windows.exe" ./cmd/main
	GOOS=darwin go build -ldflags="$(LDFLAGS)" -o "$(OUTPUT_DIR)/$(BINARY_NAME)_macos" ./cmd/main
	GOOS=linux go build -ldflags="$(LDFLAGS)" -o "$(OUTPUT_DIR)/$(BINARY_NAME)_linux" ./cmd/main
	@echo "Binaries saved to ./$(OUTPUT_DIR)"

# Clean up build artifacts
.PHONY: clean
clean:
	@echo "Cleaning..."
	@rm -f "./$(OUTPUT_DIR)/$(BINARY_NAME)_windows.exe"
	@rm -f "./$(OUTPUT_DIR)/$(BINARY_NAME)_macos"
	@rm -f "./$(OUTPUT_DIR)/$(BINARY_NAME)_linux"
	@echo "Clean complete!"

# Display current version
.PHONY: version
version:
	@echo "Version: $(VERSION)"
