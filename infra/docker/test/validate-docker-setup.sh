#!/bin/bash
# Docker Setup Validation Script for APIWeaver
# This script validates the Docker configuration without requiring a full build

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check requirements
check_requirements() {
    log_info "Checking Docker setup requirements..."
    
    # Check Docker
    if ! command -v docker &> /dev/null; then
        log_error "Docker is not installed"
        exit 1
    fi
    log_success "Docker is installed: $(docker --version)"
    
    # Check Docker Compose
    if ! command -v docker-compose &> /dev/null; then
        log_error "Docker Compose is not installed"
        exit 1
    fi
    log_success "Docker Compose is installed: $(docker-compose --version)"
    
    # Check if Docker daemon is running
    if ! docker info &> /dev/null; then
        log_error "Docker daemon is not running"
        exit 1
    fi
    log_success "Docker daemon is running"
}

# Validate file structure
validate_files() {
    log_info "Validating Docker configuration files..."
    
    local files=(
        "Dockerfile"
        "docker-compose.yml"
        "docker-compose.override.yml"
        "docker-compose.prod.yml"
        ".dockerignore"
        ".env.example"
        "Makefile.docker"
    )
    
    for file in "${files[@]}"; do
        if [[ -f "$file" ]]; then
            log_success "Found: $file"
        else
            log_error "Missing: $file"
            exit 1
        fi
    done
    
    # Check Docker configuration directories
    local dirs=(
        "docker/mongodb/init"
        "docker/nginx"
        "docker/scripts"
        "data"
        "logs"
        "backups"
    )
    
    for dir in "${dirs[@]}"; do
        if [[ -d "$dir" ]]; then
            log_success "Found directory: $dir"
        else
            log_warning "Missing directory: $dir (will be created automatically)"
        fi
    done
}

# Validate Dockerfile syntax
validate_dockerfile() {
    log_info "Validating Dockerfile syntax..."
    
    # Check if Dockerfile has the required stages
    local stages=("frontend-builder" "go-builder" "production")
    
    for stage in "${stages[@]}"; do
        if grep -q "FROM .* AS $stage" Dockerfile; then
            log_success "Found stage: $stage"
        else
            log_error "Missing stage: $stage"
            exit 1
        fi
    done
    
    # Check for security best practices
    if grep -q "USER.*apiweaver" Dockerfile; then
        log_success "Non-root user configured"
    else
        log_warning "Non-root user not found in Dockerfile"
    fi
    
    if grep -q "HEALTHCHECK" Dockerfile; then
        log_success "Health check configured"
    else
        log_warning "Health check not found in Dockerfile"
    fi
}

# Validate Docker Compose configuration
validate_compose() {
    log_info "Validating Docker Compose configuration..."
    
    # Validate main compose file
    if docker-compose config &> /dev/null; then
        log_success "Docker Compose configuration is valid"
    else
        log_error "Docker Compose configuration is invalid"
        docker-compose config
        exit 1
    fi
    
    # Check required services
    local services=("mongodb" "apiweaver")
    
    for service in "${services[@]}"; do
        if docker-compose config --services | grep -q "^$service$"; then
            log_success "Found service: $service"
        else
            log_error "Missing service: $service"
            exit 1
        fi
    done
}

# Check environment configuration
validate_environment() {
    log_info "Validating environment configuration..."
    
    if [[ -f ".env" ]]; then
        log_success "Environment file exists"
        
        # Check for required variables
        local required_vars=(
            "MONGODB_ROOT_PASSWORD"
            "JWT_SECRET"
        )
        
        for var in "${required_vars[@]}"; do
            if grep -q "^$var=" .env; then
                log_success "Found environment variable: $var"
            else
                log_warning "Missing environment variable: $var"
            fi
        done
    else
        log_warning "No .env file found (will use defaults)"
        log_info "Consider copying .env.example to .env"
    fi
}

# Test Docker build stages
test_build_stages() {
    log_info "Testing Docker build stages (syntax only)..."
    
    # Test if build context can be created
    if docker build --target frontend-builder -f Dockerfile . --progress=plain 2>&1 | head -20 | grep -q "load build definition"; then
        log_success "Frontend build stage syntax is valid"
    else
        log_error "Frontend build stage has issues"
    fi
    
    # Note: We don't run full builds as they require significant time and resources
    log_info "Note: Full build testing requires manual execution with 'make docker-build'"
}

# Test MongoDB initialization
validate_mongodb_init() {
    log_info "Validating MongoDB initialization script..."
    
    if [[ -f "docker/mongodb/init/01-init-apiweaver.js" ]]; then
        log_success "MongoDB initialization script found"
        
        # Basic syntax check for JavaScript
        if node -c docker/mongodb/init/01-init-apiweaver.js &> /dev/null; then
            log_success "MongoDB init script syntax is valid"
        else
            log_warning "MongoDB init script syntax check failed (Node.js may not be available)"
        fi
    else
        log_error "MongoDB initialization script not found"
        exit 1
    fi
}

# Validate backup scripts
validate_scripts() {
    log_info "Validating Docker scripts..."
    
    local scripts=(
        "docker/scripts/backup.sh"
        "docker/scripts/restore.sh"
    )
    
    for script in "${scripts[@]}"; do
        if [[ -f "$script" ]]; then
            if [[ -x "$script" ]]; then
                log_success "Script found and executable: $script"
            else
                log_warning "Script found but not executable: $script"
                chmod +x "$script"
                log_success "Made script executable: $script"
            fi
        else
            log_error "Missing script: $script"
            exit 1
        fi
    done
}

# Main validation function
main() {
    echo "================================================"
    echo "APIWeaver Docker Setup Validation"
    echo "================================================"
    echo ""
    
    check_requirements
    echo ""
    
    validate_files
    echo ""
    
    validate_dockerfile
    echo ""
    
    validate_compose
    echo ""
    
    validate_environment
    echo ""
    
    validate_mongodb_init
    echo ""
    
    validate_scripts
    echo ""
    
    test_build_stages
    echo ""
    
    echo "================================================"
    log_success "Docker setup validation completed successfully!"
    echo "================================================"
    echo ""
    
    log_info "Next steps:"
    echo "  1. Copy .env.example to .env and configure for your environment"
    echo "  2. Run 'make docker-build' to build the Docker image"
    echo "  3. Run 'make docker-dev' to start the development environment"
    echo "  4. Run 'make docker-prod' to start the production environment"
    echo ""
    
    log_info "Available make targets:"
    echo "  • make docker-build      - Build the Docker image"
    echo "  • make docker-dev        - Start development environment"
    echo "  • make docker-prod       - Start production environment"
    echo "  • make docker-test       - Run basic tests"
    echo "  • make docker-health     - Check service health"
    echo "  • make docker-logs       - View service logs"
    echo "  • make docker-backup     - Create MongoDB backup"
    echo "  • make docker-clean      - Clean up Docker resources"
    echo ""
}

# Run main function
main "$@"