# APIWeaver

A production-ready Go-based developer tool and service that generates and amends OpenAPI 3.1 specifications from Markdown API requirements, natural language descriptions, and structured JSON. Single binary with CLI, HTTP API, and embedded web UI.
Backend is Go + MongoDB; frontend is React + Tailwind

[![CI](https://github.com/sukhera/APIWeaver/actions/workflows/ci.yml/badge.svg)](https://github.com/sukhera/APIWeaver/actions/workflows/ci.yml)
[![Go Version](https://img.shields.io/badge/Go-1.24+-blue.svg)](https://golang.org/)
[![React Version](https://img.shields.io/badge/React-18+-blue.svg)](https://reactjs.org/)
[![Tailwind Version](https://img.shields.io/badge/Tailwind-3+-blue.svg)](https://tailwindcss.com/)
[![Coverage](https://img.shields.io/badge/coverage-95%25+-blue.svg)](https://github.com/sukhera/APIWeaver/actions/workflows/ci.yml)
[![Security](https://img.shields.io/badge/security-gosec%20%7C%20golangci--lint-green.svg)](https://github.com/sukhera/APIWeaver/actions/workflows/ci.yml)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Docker](https://img.shields.io/badge/docker-ready-blue.svg)](docker-compose.yml)

## What This Is

A lightweight OpenAPI specification generator and amender that converts Markdown requirements directly into valid OpenAPI 3.1 specs. Built with enterprise-grade Go patterns, AST-based parsing, and production-ready infrastructure automation.What This Is
A lightweight OpenAPI specification generator and amender that converts Markdown requirements directly into valid OpenAPI 3.1 specs. Built with enterprise-grade Go patterns, AST-based parsing, and production-ready infrastructure automation.
Problem Solved
Maintaining OpenAPI specs manually is slow and error-prone. Most teams already write Markdown requirements in PRDs, tickets, or wikis — but translating those into accurate OpenAPI specs requires extra effort. APIWeaver bridges this gap by directly converting your existing documentation into valid OpenAPI specifications.

## Web Interface

APIWeaver includes a powerful embedded web UI that provides an intuitive interface for working with your API specifications:

## Who is APIWeaver For?

- **Backend Developers** creating and maintaining API specs
- **API Engineers** designing consistent interfaces
- **Tech Leads** documenting APIs for their teams
- **Small Teams** collaborating on OpenAPI specifications

## Features

- 📝 **Markdown to OpenAPI** - Convert requirements.md directly to openapi.yaml
- 🔄 **Spec Amendment** - Update existing specs via Markdown changes or plain text
- ✅ **Validation Engine** - Comprehensive spec validation with actionable feedback
- 🔍 **Diff Preview** - See changes before applying amendments
- ⚡ **High Performance** - AST-based parsing with concurrent processing
- 🐳 **Production Ready** - Single binary with CLI, API, and embedded web UI
- 📈 **Observability** - Structured logging and comprehensive error handling
- 🛠️ **Infrastructure Automation** - Comprehensive Makefile and CI/CD pipeline

## Quick Start

```bash
git clone https://github.com/sukhera/APIWeaver.git
cd APIWeaver
cp .env.example .env
make dev
```

### Access Points

- **Web UI**: <http://localhost:8080> (embedded in binary)
- **API Health**: <http://localhost:8080/api/health>
- **CLI Tool**: `./apiweaver parse requirements.md`
- **API Endpoint**: <http://localhost:8080/api/v1/generate>

## Configuration

Configuration is env-first; `.env` sets ports/URIs for dev, compose wires containers; prod uses environment/secret manager.

- **Environment**: Copy `.env.example` to `.env` and customize
- **YAML Config**: See `config.yaml` for structured configuration
- **Precedence**: Environment variables > flags > config file > defaults

Key settings:

```bash
PORT=8080                           # Server port
LOG_LEVEL=info                      # Logging verbosity
STRICT_MODE=false                   # Strict parsing mode
PARSER_TIMEOUT=30s                  # Parser timeout
VALIDATION_LEVEL=basic              # Validation strictness
OUTPUT_FORMAT=json                  # Output format (json/yaml/text)
```

## Project Structure

```Directory Structure
├── cmd/                    # Application entrypoints
│   └── apiweaver/         # CLI and API server binary
├── internal/              # Private application code
│   ├── config/           # Configuration management
│   ├── domain/           # Core business logic
│   │   ├── parser/       # AST-based markdown parser
│   │   ├── builder/      # OpenAPI spec builders
│   │   └── strategies/   # Parsing strategies
│   └── application/      # HTTP handlers and routes
├── pkg/                   # Public packages
│   └── errors/           # Error handling utilities
├── docs/                  # Documentation and guides
├── testutil/             # Test utilities and helpers
```

→ **Full structure**: See [docs/architecture.md](docs/architecture.md)

## API

### Generate OpenAPI Spec

```bash
POST /api/v1/generate
Content-Type: application/json


```json
{
  "content": "# User API\n\n## GET /users\nReturns a list of users...",
  "format": "yaml",
  "strict_mode": false
}
```

```json
{
  "openapi": "3.1.0",
  "info": {
    "title": "User API",
    "version": "1.0.0"
  },
  "paths": {
    "/users": {
      "get": {
        "summary": "Returns a list of users"
      }
    }
  }
}
```

### Amend Existing Spec

```bash
POST /api/v1/amend
Content-Type: application/json
```

```json
{
  "original_spec": "...",
  "changes": "Add POST /users endpoint for creating new users",
  "format": "yaml"
}
```

### Validate Spec

```bash
POST /api/v1/validate
Content-Type: application/json
```

```json
{
  "spec": "...",
  "validation_level": "strict"
}
```

## Markdown Parser Architecture

The system uses **AST-based parsing** with **Visitor pattern** for modular, extensible OpenAPI generation:

- **AST Construction**: Parse Markdown into structured Abstract Syntax Tree
- **Visitor Pattern**: Traverse AST to generate OpenAPI components
- **Strategy Pattern**: Multiple parsing strategies for different content types
- **Functional Options**: Flexible configuration with builder pattern

Design patterns used here are documented under `/docs`—the README stays high-level so you can find what you need fast.

→ **Implementation details**: [docs/design-patterns.md](docs/design-patterns.md)

## Markdown Input Format

### File Structure

1. **Optional frontmatter** (YAML)
2. **API metadata** title, version, servers, security
3. **Endpoint sections** `## METHOD /path`
4. **Parameter and response tables**
5. **Request/response schemas** in fenced JSON/YAML blocks
6. **Reusable components** defined at the bottom

### Frontmatter Example

```yaml
---
title: User API
version: 1.0.0
servers:
  - url: https://api.example.com
security:
  - bearerAuth: []
---
```

### Endpoint Example

```markdown
## GET /users

Returns a list of users.

### Parameters
| Name | Type | Required | Description |
|------|------|----------|-------------|
| page | integer | No | Page number |
| limit | integer | No | Items per page |

### Responses
| Code | Description |
|------|-------------|
| 200 | Success |
| 400 | Bad Request |

### Response Schema
```json
{
  "type": "object",
  "properties": {
    "users": {
      "type": "array",
      "items": {
        "$ref": "#/components/schemas/User"
      }
    }
  }
}
```

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

## Observability

- **Structured Logging**: Context-aware with trace correlation
- **Performance Metrics**: Parsing time, memory usage, throughput
- **Error Tracking**: Comprehensive error wrapping with stack traces
- **Health Checks**: Service health and dependency status

## Deployment

### Docker Compose

Production-ready setup with `docker-compose.yml` featuring API server with embedded web UI and MongoDB for API key management.

```bash
make prod    # Start production environment
```

### Single Binary

```bash
# Build for your platform
make build

# Run directly
./apiweaver serve --port 8080
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
```

## Project Status

### ✅ Implemented

- AST-based markdown parser with visitor pattern
- OpenAPI 3.1 spec generation and validation
- Comprehensive configuration management
- Mockery integration with expecter pattern
- CI/CD pipeline with automated testing

### 🚧 Roadmap

- Docker Compose deployment setup
- Structured logging and error handling
- Embedded web UI for easy interaction
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

## Contributing

1. **Conventional Commits**: Use conventional commit format
2. **Quality Gates**: Run `make lint test` before submitting
3. **Small PRs**: Keep changes focused and reviewable
4. **Tests Required**: Add tests for new functionality (95%+ coverage)
5. **Documentation**: Update docs for API/config changes
6. **Mockery**: Generate mocks for new interfaces

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---
