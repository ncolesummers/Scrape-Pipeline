# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build and Run Commands

```bash
# Build the application
make build

# Run the application with default configuration
make run

# Run all tests with formatted output
make test

# Run tests with standard Go output
make test-standard

# Run linters
make lint

# Run linters with auto-fix
make lint-fix

# Clean build artifacts and rebuild
make clean build

# Create default config.yaml from example
make config
```

## Project Architecture

This project is a high-performance, modular web scraping and RAG (Retrieval-Augmented Generation) system built in Go. It follows a pipeline architecture with the following modules:

1. **Web Scraping Module** (`pkg/scraper`, `internal/scraper`): Discovers and crawls blog URLs, manages rate limiting. Implements the `Scraper` interface.

2. **Content Extraction Module** (`pkg/extractor`, `internal/extractor`): Separates main article content from boilerplate HTML using go-readability. Implements the `Extractor` interface.

3. **Text Normalization Module**: Standardizes text formatting. Implements the `Normalizer` interface.

4. **Document Chunking Module**: Splits documents into appropriate-sized chunks. Implements the `Chunker` interface.

5. **Quality Control Module**: Detects low-quality or duplicate content. Implements the `QualityControl` interface.

6. **Embedding Service Module**: Converts text chunks to vector embeddings. Implements the `Embedder` interface.

7. **Vector Storage Module**: Stores and indexes vector embeddings. Implements the `VectorStorage` interface.

8. **Observability Module**: Collects metrics, traces, and logs. Implements the `Observer` interface.

All module interfaces are defined in `internal/models/interfaces.go`.

## Project Structure

- `cmd/`: Application entry points
- `pkg/`: Public packages that can be used by external applications
- `internal/`: Private packages that are internal to this application
  - `config/`: Configuration loading and management
  - `extractor/`: Content extraction from HTML
  - `models/`: Data models and interfaces
  - `scraper/`: Web scraping implementation
- `docs/`: Documentation files

## Key Dependencies

- github.com/gocolly/colly/v2: For web scraping
- github.com/go-shiori/go-readability: For content extraction
- golang.org/x/net: For HTML/DOM processing
- gopkg.in/yaml.v3: For configuration parsing

## Configuration

The application uses a YAML configuration file (`config.yaml`) with settings for all pipeline components. Create it from the example:

```bash
make config
```

Key configuration sections include:
- Scraping parameters (rate limits, concurrency, user agents)
- Content extraction rules
- Document chunking strategies
- Quality control thresholds
- Embedding service configuration
- Vector storage options

## Project Management

The project uses GitHub Projects v2 for task tracking. The project number is **2**.

For CLI commands to manage the project, refer to:
```bash
# View GitHub Projects CLI guide
cat docs/research/github_projects_cli.md
```