# Contributing to the Web Scraping and RAG System Pipeline

Thank you for considering contributing to this project! This document provides guidelines and instructions for contributing to our high-performance, modular web scraping and RAG system pipeline built in Go.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Workflow](#development-workflow)
- [Pull Request Process](#pull-request-process)
- [Coding Standards](#coding-standards)
- [Testing Guidelines](#testing-guidelines)
- [Documentation](#documentation)
- [Issue Reporting](#issue-reporting)
- [Communication](#communication)

## Code of Conduct

This project adheres to a Code of Conduct that all contributors are expected to follow. By participating, you are expected to uphold this code. Please report unacceptable behavior to the project maintainers.

## Getting Started

### Prerequisites

Before you begin contributing, ensure you have the following installed:

- Go 1.22 or higher
- Git
- golangci-lint (for linting)
- gotestsum (for improved test output formatting)

### Setting Up Your Development Environment

1. Fork the repository on GitHub.
2. Clone your forked repository:
   ```bash
   git clone https://github.com/YOUR_USERNAME/scrape-pipeline.git
   cd scrape-pipeline
   ```
3. Add the original repository as an upstream remote:
   ```bash
   git remote add upstream https://github.com/ncolesummers/scrape-pipeline.git
   ```
4. Install dependencies:
   ```bash
   go mod tidy
   ```
5. Install development tools:
   ```bash
   # Install golangci-lint
   # macOS with Homebrew
   brew install golangci-lint
   
   # Go install
   go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
   
   # Install gotestsum for better test output
   go install gotest.tools/gotestsum@latest
   ```

## Development Workflow

We follow a feature branch workflow:

1. Sync your fork with the upstream repository:
   ```bash
   git fetch upstream
   git checkout main
   git merge upstream/main
   ```

2. Create a new branch for your feature or bug fix:
   ```bash
   git checkout -b feature/your-feature-name
   ```
   
   Use prefixes like `feature/`, `bugfix/`, `docs/`, or `test/` to categorize your branches.

3. Make your changes, following our [Coding Standards](#coding-standards).

4. Run tests and linters locally:
   ```bash
   make test
   make lint
   ```

5. Commit your changes with clear, descriptive commit messages:
   ```bash
   git commit -m "feat: add new feature X" 
   ```
   
   We follow [Conventional Commits](https://www.conventionalcommits.org/) for commit messages.

6. Push to your fork:
   ```bash
   git push origin feature/your-feature-name
   ```

7. Create a Pull Request from your fork to the original repository.

## Pull Request Process

1. Ensure your PR includes a clear description of the changes and their purpose.
2. Link the PR to any relevant issues.
3. Make sure all CI checks pass.
4. Update documentation as necessary.
5. Respond to review comments and make requested changes.
6. A maintainer will merge your PR when it's ready.

## Coding Standards

### Core Development Principles

- **Simplicity First**: Prefer simple, readable solutions over complex ones. Clear code that others can understand and maintain is a priority.
- **Incremental Development**: Make small, focused changes with each commit. This enables easier review, testing, and rollback if needed.
- **Modularity**: Implement each component with clean interfaces that allow for component swapping and system evolution.
- **Performance Focus**: Optimize for Go's strengths in concurrency and efficient resource usage, especially for high-volume processing.
- **Ethical Compliance**: Ensure all scraping respects robots.txt, implements appropriate rate limiting, and follows legal requirements.

### Go Style Guidelines

- Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments) and [Effective Go](https://golang.org/doc/effective_go) guidelines.
- Run `golangci-lint` on your code to catch common issues:
  ```bash
  make lint
  ```
- Keep functions focused and under 50 lines where possible.
- Use meaningful variable and function names that reflect their purpose.
- Comment complex algorithms and business logic, not obvious code.
- Handle errors explicitly and avoid panic in production code.
- Ensure all exported functions, types, and interfaces are properly documented.

## Testing Guidelines

- Write tests for all new code. Aim for at least 80% test coverage.
- Create both unit tests and integration tests as appropriate.
- Test edge cases and error conditions.
- Use table-driven tests when testing multiple inputs/outputs.
- Run the test suite before submitting a PR:
  ```bash
  make test
  ```
- Include performance benchmarks for performance-critical code:
  ```bash
  go test -bench=. ./...
  ```

## Documentation

- Update documentation when adding new features or changing existing functionality.
- Document all exported functions, types, and interfaces according to Go standards.
- Keep the README.md updated with any user-facing changes.
- Add examples to help users understand how to use new features.

## Issue Reporting

When reporting issues, please use one of our issue templates and include:

- A clear and descriptive title
- A detailed description of the issue
- Steps to reproduce the behavior
- Expected behavior
- Actual behavior
- Environment details (OS, Go version, etc.)
- Any relevant logs or screenshots

## Communication

- For quick questions or discussions, open a Discussion on GitHub.
- For bug reports or feature requests, open an Issue.
- For complex changes, consider discussing in an Issue before submitting a PR.

## License

By contributing to this project, you agree that your contributions will be licensed under the project's license (see LICENSE file).

Thank you for contributing to the Web Scraping and RAG System Pipeline! 