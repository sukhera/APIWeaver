# APIWeaver Makefile
# Following Go best practices and backend expert guidelines

# Variables
BINARY_NAME=apiweaver
BUILD_DIR=build
DIST_DIR=dist
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT_SHA=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S_UTC')
LDFLAGS=-ldflags "-X main.version=$(VERSION) -X main.commitSHA=$(COMMIT_SHA) -X main.buildTime=$(BUILD_TIME)"

# Go related variables
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOVET=$(GOCMD) vet
GOFMT=$(GOCMD) fmt
GOLINT=golangci-lint
GOSEC=gosec
MOCKERY=mockery

# Default target
.DEFAULT_GOAL := help

# Help target
.PHONY: help
help: ## Show this help message
	@echo "APIWeaver - Markdown to API Specification Parser"
	@echo ""
	@echo "Available targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Build targets
.PHONY: build
build: create-dirs ## Build the application
	@echo "Building APIWeaver..."
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/apiweaver

.PHONY: build-all
build-all: create-dirs ## Build for all platforms
	@echo "Building for all platforms..."
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/apiweaver
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-amd64 ./cmd/apiweaver
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-arm64 ./cmd/apiweaver
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-windows-amd64.exe ./cmd/apiweaver

.PHONY: install
install: ## Install the application
	@echo "Installing APIWeaver..."
	$(GOBUILD) $(LDFLAGS) -o $(GOPATH)/bin/$(BINARY_NAME) ./cmd/apiweaver

# Clean targets
.PHONY: clean
clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	rm -rf $(DIST_DIR)
	rm -f coverage.out
	rm -f coverage.html

# Test targets
.PHONY: test
test: ## Run tests
	@echo "Running tests..."
	$(GOTEST) -v ./internal/... ./pkg/... ./cmd/...

.PHONY: test-race
test-race: ## Run tests with race detection
	@echo "Running tests with race detection..."
	$(GOTEST) -race -v ./internal/... ./pkg/... ./cmd/...

.PHONY: test-coverage
test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	$(GOTEST) -coverprofile=coverage.out -covermode=atomic ./internal/... ./pkg/... ./cmd/...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

.PHONY: test-benchmark
test-benchmark: ## Run benchmarks
	@echo "Running benchmarks..."
	$(GOTEST) -bench=. -benchmem ./internal/... ./pkg/... ./cmd/...

# Mockery targets
.PHONY: generate-mocks
generate-mocks: ## Generate mocks using Mockery
	@echo "$(BLUE)Generating mocks...$(NC)"
	@go install github.com/vektra/mockery/v2@latest
	@rm -rf mocks/
	@mockery --all --inpackage --with-expecter=true
	@echo "$(GREEN)✓ Mocks generated!$(NC)"

.PHONY: test-with-mocks
test-with-mocks: ## Run tests with generated mocks
	@echo "$(BLUE)Running tests with mocks...$(NC)"
	@make generate-mocks
	@go test -v ./internal/domain/parser/...
	@go test -v ./internal/config/...
	@go test -v ./pkg/errors/...
	@echo "$(GREEN)✓ Tests with mocks completed!$(NC)"

.PHONY: clean-mocks
clean-mocks: ## Clean generated mocks
	@echo "Cleaning generated mocks..."
	@find . -name "*Mock*.go" -type f -delete
	@echo "✓ Mocks cleaned!"

# Code quality targets
.PHONY: fmt
fmt: ## Format code
	@echo "Formatting code..."
	$(GOFMT) ./internal/... ./pkg/... ./cmd/...

.PHONY: vet
vet: ## Run go vet
	@echo "Running go vet..."
	$(GOVET) ./internal/... ./pkg/... ./cmd/...

.PHONY: lint
lint: ## Run linter
	@echo "Running linter..."
	$(GOLINT) run ./internal/... ./pkg/... ./cmd/...

.PHONY: security
security: ## Run security scanner
	@echo "Running security scanner..."
	$(GOSEC) ./internal/... ./pkg/... ./cmd/...

.PHONY: mod-tidy
mod-tidy: ## Tidy go modules
	@echo "Tidying go modules..."
	$(GOMOD) tidy

.PHONY: mod-verify
mod-verify: ## Verify go modules
	@echo "Verifying go modules..."
	$(GOMOD) verify

# Pre-commit checks
.PHONY: pre-commit
pre-commit: ## Run all pre-commit checks
	@echo "Running pre-commit checks..."
	$(MAKE) fmt
	$(MAKE) vet
	$(MAKE) generate-mocks
	$(MAKE) test-race
	$(MAKE) lint
	$(MAKE) security
	$(MAKE) mod-tidy
	$(MAKE) mod-verify
	@echo "All pre-commit checks passed!"

# Development targets
.PHONY: dev
dev: ## Run in development mode
	@echo "Running in development mode..."
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/apiweaver
	./$(BUILD_DIR)/$(BINARY_NAME) -version

.PHONY: run
run: ## Run the application (requires -input flag)
	@echo "Running APIWeaver..."
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/apiweaver
	./$(BUILD_DIR)/$(BINARY_NAME) $(ARGS)

.PHONY: serve
serve: ## Start APIWeaver server in development mode with Docker
	@echo "Starting APIWeaver server with hot reload..."
	cd infra && make -f Makefile.docker docker-dev
	@echo ""
	@echo "APIWeaver server running with hot reload"
	@echo "  • API: http://localhost:8080"
	@echo "  • Health check: http://localhost:8080/api/v1/health"
	@echo "  • MongoDB Express: http://localhost:8081"
	@echo ""
	@echo "Press Ctrl+C to stop..."
	@trap 'cd infra && make -f Makefile.docker docker-stop' INT; cd infra && make -f Makefile.docker docker-dev-logs

# Documentation targets
.PHONY: docs
docs: ## Generate documentation
	@echo "Generating documentation..."
	godoc -http=:6060 &
	@echo "Documentation available at http://localhost:6060"

# Release targets
.PHONY: release
release: ## Create a release build
	@echo "Creating release build..."
	$(MAKE) clean
	$(MAKE) pre-commit
	$(MAKE) build-all
	@echo "Release build complete!"

# Docker targets
.PHONY: docker-build
docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t apiweaver:$(VERSION) .
	docker tag apiweaver:$(VERSION) apiweaver:latest

.PHONY: docker-run
docker-run: ## Run Docker container
	@echo "Running Docker container..."
	docker run --rm -it apiweaver:latest

.PHONY: docker-dev
docker-dev: ## Start development environment with Docker Compose
	@echo "Starting development environment..."
	cd infra && make -f Makefile.docker docker-dev

.PHONY: docker-dev-logs
docker-dev-logs: ## Show development logs
	cd infra && make -f Makefile.docker docker-dev-logs

.PHONY: docker-test
docker-test: ## Test backend in Docker environment
	cd infra && make -f Makefile.docker docker-test

.PHONY: docker-stop
docker-stop: ## Stop Docker development environment
	cd infra && make -f Makefile.docker docker-stop

.PHONY: docker-clean
docker-clean: ## Clean Docker resources
	cd infra && make -f Makefile.docker docker-clean

# Dependencies
.PHONY: deps
deps: ## Install dependencies
	@echo "Installing dependencies..."
	$(GOGET) -v -t -d ./...
	$(GOMOD) download

.PHONY: deps-update
deps-update: ## Update dependencies
	@echo "Updating dependencies..."
	$(GOGET) -u ./...
	$(GOMOD) tidy

# CI/CD targets
.PHONY: ci
ci: ## Run CI pipeline
	@echo "Running CI pipeline..."
	$(MAKE) pre-commit
	$(MAKE) test-coverage
	$(MAKE) build
	@echo "CI pipeline completed successfully!"

# Performance targets
.PHONY: profile
profile: ## Run profiling
	@echo "Running profiling..."
	$(GOTEST) -cpuprofile=cpu.prof -memprofile=mem.prof -bench=. ./...

.PHONY: profile-analyze
profile-analyze: ## Analyze profiling results
	@echo "Analyzing profiling results..."
	$(GOCMD) tool pprof cpu.prof
	$(GOCMD) tool pprof mem.prof

# Utility targets
.PHONY: version
version: ## Show version information
	@echo "Version: $(VERSION)"
	@echo "Commit: $(COMMIT_SHA)"
	@echo "Build Time: $(BUILD_TIME)"

.PHONY: check
check: ## Check if all tools are installed
	@echo "Checking required tools..."
	@command -v $(GOCMD) >/dev/null 2>&1 || { echo "Go is not installed"; exit 1; }
	@command -v $(GOLINT) >/dev/null 2>&1 || { echo "golangci-lint is not installed"; exit 1; }
	@command -v $(GOSEC) >/dev/null 2>&1 || { echo "gosec is not installed"; exit 1; }
	@echo "All required tools are installed!"

# Install development tools
.PHONY: install-tools
install-tools: ## Install development tools
	@echo "Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
	go install github.com/vektra/mockery/v2@latest
	go install golang.org/x/tools/cmd/godoc@latest
	go install golang.org/x/tools/cmd/benchcmp@latest
	@echo "Development tools installed!"

# Create directories
.PHONY: create-dirs
create-dirs:
	mkdir -p $(BUILD_DIR)
	mkdir -p $(DIST_DIR)

# Frontend targets
.PHONY: build-frontend dev-frontend test-frontend lint-frontend type-check-frontend
build-frontend: ## Build frontend for production
	@echo "🏗️  Building frontend for production..."
	@cd web && ./scripts/embed-build.sh

dev-frontend: ## Start frontend development server
	@echo "🚀 Starting frontend development server..."
	@cd web && npm run dev

test-frontend: ## Run frontend tests
	@echo "🧪 Running frontend tests..."
	@cd web && npm run test:run

lint-frontend: ## Lint frontend code
	@echo "🔍 Linting frontend code..."
	@cd web && npm run lint

type-check-frontend: ## Type check frontend code
	@echo "🔍 Type checking frontend code..."
	@cd web && npm run type-check

# TODO: Add Go embedding targets when web server is implemented
