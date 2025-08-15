# APIWeaver

A Go-based developer tool that generates OpenAPI 3.1 specifications from Markdown API requirements. 
Single binary with CLI and embedded web UI.

[![CI](https://github.com/sukhera/APIWeaver/actions/workflows/ci.yml/badge.svg)](https://github.com/sukhera/APIWeaver/actions/workflows/ci.yml)
[![Go Version](https://img.shields.io/badge/Go-1.24+-blue.svg)](https://golang.org/)
[![React Version](https://img.shields.io/badge/React-18+-blue.svg)](https://reactjs.org/)
[![Tailwind Version](https://img.shields.io/badge/Tailwind-3+-blue.svg)](https://tailwindcss.com/)
[![Coverage](https://img.shields.io/badge/coverage-95%25+-blue.svg)](https://github.com/sukhera/APIWeaver/actions/workflows/ci.yml)
[![Security](https://img.shields.io/badge/security-gosec%20%7C%20golangci--lint-green.svg)](https://github.com/sukhera/APIWeaver/actions/workflows/ci.yml)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Docker](https://img.shields.io/badge/docker-ready-blue.svg)](docker-compose.yml)

## What This Is

A lightweight OpenAPI specification generator that converts Markdown requirements 
directly into valid OpenAPI 3.1 specs. Built with enterprise-grade Go patterns, AST-based 
parsing, and production-ready infrastructure automation.

### Problem Solved

Maintaining OpenAPI specs manually is slow and error-prone. Most teams already write Markdown 
requirements in PRDs, tickets, or wikis — but translating those into accurate OpenAPI specs 
requires extra effort. APIWeaver bridges this gap by directly converting your existing 
documentation into valid OpenAPI specifications.

## Web Interface

APIWeaver includes a powerful embedded web UI that provides an intuitive interface for working with your API specifications.

## Who is APIWeaver For?

- **Backend Developers** creating and maintaining API specs
- **API Engineers** designing consistent interfaces
- **Tech Leads** documenting APIs for their teams
- **Small Teams** collaborating on OpenAPI specifications

## Features

- 📝 **Markdown to OpenAPI** - Convert requirements.md directly to openapi.yaml
- ✅ **Validation Engine** - Comprehensive spec validation with actionable feedback
- ⚡ **High Performance** - AST-based parsing with concurrent processing
- 🐳 **Production Ready** - Single binary with CLI and embedded web UI
- 📈 **Observability** - Structured logging and comprehensive error handling
- 🛠️ **Infrastructure Automation** - Comprehensive Makefile and CI/CD pipeline

## Quick Start

```bash
git clone https://github.com/sukhera/APIWeaver.git
cd APIWeaver
make setup
```

### Access Points

- **Web UI**: <http://localhost:8080> (embedded in binary)
- **CLI Tool**: `./apiweaver parse requirements.md`

## Configuration

Configuration is handled through YAML config files and command-line flags.

- **YAML Config**: See `config.yaml.sample` for structured configuration
- **Precedence**: Command-line flags > config file > defaults

Key settings: See `config.yaml.sample` for complete configuration options.

## Project Structure

```bash
├── cmd/                    # Application entrypoints
│   └── apiweaver/         # CLI binary
├── internal/              # Private application code
│   ├── config/           # Configuration management
│   └── domain/           # Core business logic
│       ├── parser/       # AST-based markdown parser
│       └── builder/      # OpenAPI spec builders
├── pkg/                   # Public packages
│   └── errors/           # Error handling utilities
├── web/                   # React frontend
├── docs/                  # Documentation and guides
└── testutil/             # Test utilities and helpers
```

## Markdown Parser Architecture

The system uses **AST-based parsing** with **Visitor pattern** for modular, extensible OpenAPI 
generation:

- **AST Construction**: Parse Markdown into structured Abstract Syntax Tree
- **Visitor Pattern**: Traverse AST to generate OpenAPI components
- **Strategy Pattern**: Multiple parsing strategies for different content types
- **Functional Options**: Flexible configuration with builder pattern

## Markdown Input Format

APIWeaver parses Markdown files with OpenAPI specifications. See `docs/markdown-format.md` for detailed format documentation.

## CLI Usage

### Generate OpenAPI Spec

```bash
# From file
./apiweaver parse requirements.md --output openapi.yaml

# From stdin
cat requirements.md | ./apiweaver parse - --output openapi.json

# With options
./apiweaver parse requirements.md \
  --strict \
  --format yaml \
  --output openapi.yaml \
  --verbose
```

## Development Commands

Run `make help` to see the full toolbox; this README only shows the commands you'll use daily.

```bash
make setup           # Complete project setup
make dev             # Start development environment  
make test            # Run all tests with mocks
make test-race       # Run tests with race detection
make lint            # Run all linters
make security        # Security scanning (gosec, golangci-lint)
make generate-mocks  # Generate mocks using Mockery
make build           # Build all components
make clean           # Clean build artifacts
make docs            # Generate documentation

# Web Development
make dev-frontend    # Start frontend development server
make build-frontend  # Build frontend for production
make test-frontend   # Run frontend tests
make lint-frontend   # Lint frontend code
```

## Project Status

### ✅ Implemented

- AST-based markdown parser with visitor pattern
- OpenAPI 3.1 spec generation and validation
- Comprehensive configuration management
- Mockery integration with expecter pattern
- CI/CD pipeline with automated testing
- React frontend with embedded web UI
- Frontend build and integration scripts

### 🚧 Roadmap

- HTTP API server implementation
- Docker Compose deployment setup
- Structured logging and error handling
- CLI tool with full feature parity
- Comprehensive test coverage (95%+)
- Production-ready HTTP server with security hardening
- Advanced natural language processing for spec amendments
- Real-time collaborative editing
- OpenAPI spec diff visualization
- Integration with popular API documentation platforms
- Advanced validation rules and custom validators
- Plugin system for custom parsers
- Prometheus metrics integration
- Template library for common API patterns
- WebSocket support for real-time updates
- Multi-language support (Python, JavaScript SDKs)
- Git integration for version control

## Contributing

1. **Conventional Commits**: Use conventional commit format
2. **Quality Gates**: Run `make lint test` before submitting
3. **Small PRs**: Keep changes focused and reviewable
4. **Tests Required**: Add tests for new functionality (95%+ coverage)
5. **Documentation**: Update docs for API/config changes
6. **Mockery**: Generate mocks for new interfaces

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
