.PHONY: all build clean test test-standard lint run help lint-fix

# Variables
BINARY_NAME=scrape-pipeline
BUILD_DIR=./bin
MAIN_PATH=./cmd/scrape-pipeline
CONFIG_FILE=config.yaml

all: clean build test

help:
	@echo "Available commands:"
	@echo "  make build   - Build the application"
	@echo "  make clean   - Remove build artifacts"
	@echo "  make test    - Run tests with formatted output and summary"
	@echo "  make test-standard - Run tests with standard Go output"
	@echo "  make lint    - Run linters"
	@echo "  make lint-fix - Run linters with auto-fix for some issues"
	@echo "  make run     - Run the application with default configuration"
	@echo "  make help    - Show this help message"

build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)

test:
	@echo "Running tests with formatted output..."
	@if command -v gotestsum > /dev/null; then \
		gotestsum --format pkgname-and-test-fails --format-hide-empty-pkg --junitfile=test-report.xml -- ./...; \
	else \
		echo "gotestsum not installed. Using standard test output."; \
		echo "To install gotestsum: go install gotest.tools/gotestsum@latest"; \
		$(MAKE) test-standard; \
	fi

test-standard:
	@echo "Running tests with standard output..."
	@go test -v ./...

lint:
	@echo "Running linters..."
	@go vet ./...
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Please install with:"; \
		echo "  brew install golangci-lint    # macOS with Homebrew"; \
		echo "  go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest  # Go install"; \
		echo "See https://golangci-lint.run/usage/install/ for more options"; \
		exit 1; \
	fi

lint-fix:
	@echo "Running linters with auto-fix..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run --fix; \
	else \
		echo "golangci-lint not installed. Please install with:"; \
		echo "  brew install golangci-lint    # macOS with Homebrew"; \
		echo "  go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest  # Go install"; \
		echo "See https://golangci-lint.run/usage/install/ for more options"; \
		exit 1; \
	fi

run: build
	@echo "Running $(BINARY_NAME)..."
	@$(BUILD_DIR)/$(BINARY_NAME) -config=$(CONFIG_FILE)

# Create default config file if it doesn't exist
config:
	@if [ ! -f $(CONFIG_FILE) ]; then \
		echo "Creating default configuration file..."; \
		cp config.yaml.example $(CONFIG_FILE); \
	else \
		echo "Configuration file already exists."; \
	fi 