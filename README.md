# Web Scraping and RAG System Pipeline

A high-performance, modular web scraping and RAG (Retrieval-Augmented Generation) system pipeline built in Go.

## Project Overview

This project implements a robust and scalable pipeline for web scraping blog content, processing the extracted data, and preparing it for use in a Retrieval-Augmented Generation (RAG) system. The system is designed to create a local knowledge base of thousands of blog posts that can be efficiently queried for relevant information.

### Key Features

- High-performance web scraping using Colly
- Intelligent content extraction from HTML
- Advanced text normalization and document chunking
- Quality control with duplicate detection
- Vector embedding generation and storage
- Comprehensive observability with metrics, logging, and tracing

## Architecture

The system follows a modular architecture with the following components:

1. Web Scraping Module
2. Content Extraction Module
3. Text Normalization Module
4. Document Chunking Module
5. Quality Control Module
6. Embedding Service Module
7. Vector Storage Module
8. Observability Module

## Getting Started

### Prerequisites

- Go 1.22 or higher
- Git
- golangci-lint (for development)
- gotestsum (for improved test output formatting)

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/ncolesummers/scrape-pipeline.git
   cd scrape-pipeline
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Install development tools:
   ```bash
   # Install golangci-lint
   # macOS with Homebrew
   brew install golangci-lint
   
   # Go install
   go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
   
   # Install gotestsum for better test output
   go install gotest.tools/gotestsum@latest
   
   # See https://golangci-lint.run/usage/install/ for more options
   ```

4. Build the application:
   ```bash
   go build -o scrape-pipeline ./cmd/scrape-pipeline
   ```

   Or use the provided Makefile:
   ```bash
   make build
   ```

### Usage

Run the application with the default configuration:
```bash
./scrape-pipeline
```

Or specify a custom configuration file:
```bash
./scrape-pipeline -config=custom-config.yaml
```

Using the Makefile:
```bash
make run
```

## Makefile Commands

The project includes a Makefile to simplify common development tasks:

| Command | Description |
|---------|-------------|
| `make build` | Build the application binary in the `bin` directory |
| `make clean` | Remove build artifacts |
| `make test` | Run all tests with formatted output (uses gotestsum) |
| `make test-standard` | Run tests with standard Go output (without gotestsum) |
| `make lint` | Run linters (requires golangci-lint) |
| `make lint-fix` | Run linters with auto-fix for some issues |
| `make run` | Build and run the application with the default configuration |
| `make config` | Create a default configuration file if none exists |
| `make help` | Show available Makefile commands |

Example:
```bash
# Build the application
make build

# Run all tests
make test

# Clean build artifacts and rebuild
make clean build

# Run the application with default settings
make run
```

## Configuration

The application is configured using a YAML file. See the example configuration file `config.yaml.example` for more details.

To create a default configuration:
```bash
make config
```

## Continuous Integration

This project uses GitHub Actions for continuous integration to ensure code quality and functionality. The following CI workflows are set up:

### Main CI Workflow

The main CI workflow (`ci.yml`) runs on every push to `main` and on pull requests. It performs the following steps:
- Builds the application
- Runs all tests using `gotestsum`
- Archives test results as artifacts
- Runs code linting with `golangci-lint`

### Makefile CI Workflow

The Makefile CI workflow (`makefile-ci.yml`) ensures that the Makefile commands work as expected. It runs:
- `make build` - Builds the application
- `make test` - Runs all tests
- `make lint` - Lints the code

### Release Workflow

The Release workflow (`release.yml`) automates the creation of release artifacts when a new version tag is pushed:
- Triggered when a tag matching the pattern `v*` is pushed (e.g., `v1.0.0`)
- Builds binaries for multiple platforms (Linux, macOS Intel/ARM, Windows)
- Creates release archives (.tar.gz for Linux/macOS, .zip for Windows)
- Automatically creates a GitHub Release with the artifacts
- Includes the configuration example file with the release

To create a release, tag your commit and push the tag:
```bash
git tag v1.0.0
git push origin v1.0.0
```

The CI ensures that all changes maintain the quality standards of the project and don't introduce regressions.

## Development

### Project Structure

- `cmd/`: Application entry points
- `pkg/`: Public packages that can be used by external applications
- `internal/`: Private packages that are internal to this application
- `docs/`: Documentation files

### Running Tests

```bash
go test ./...
```

Or using the Makefile for improved test output:
```bash
make test
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- [Colly](https://github.com/gocolly/colly) for web scraping
- [go-readability](https://github.com/go-shiori/go-readability) for content extraction 