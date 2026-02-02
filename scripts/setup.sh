#!/bin/bash

# RepoGraph Platform - Setup Script
# This script helps you set up the development environment

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Print functions
print_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

print_header() {
    echo ""
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${BLUE}  $1${NC}"
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo ""
}

# Check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Main setup function
main() {
    print_header "RepoGraph Platform Setup"

    # Check prerequisites
    print_info "Checking prerequisites..."

    if ! command_exists go; then
        print_error "Go is not installed. Please install Go 1.21 or higher."
        exit 1
    fi
    print_success "Go is installed ($(go version))"

    if ! command_exists docker; then
        print_warning "Docker is not installed. You'll need it for full functionality."
    else
        print_success "Docker is installed ($(docker --version))"
    fi

    if ! command_exists docker-compose; then
        print_warning "Docker Compose is not installed."
    else
        print_success "Docker Compose is installed ($(docker-compose --version))"
    fi

    # Create necessary directories
    print_header "Creating Directory Structure"
    mkdir -p data/diagrams
    mkdir -p credentials
    mkdir -p bin
    mkdir -p logs
    print_success "Directories created"

    # Copy environment file if it doesn't exist
    print_header "Configuration Setup"
    if [ ! -f .env ]; then
        cp .env.example .env
        print_success "Created .env file from .env.example"
        print_warning "Please edit .env file with your credentials"
    else
        print_info ".env file already exists"
    fi

    # Download Go dependencies
    print_header "Installing Go Dependencies"
    print_info "Downloading Go modules..."
    go mod download
    print_success "Go dependencies installed"

    # Install development tools
    print_header "Installing Development Tools"

    if ! command_exists golangci-lint; then
        print_info "Installing golangci-lint..."
        go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
        print_success "golangci-lint installed"
    else
        print_info "golangci-lint already installed"
    fi

    if ! command_exists goimports; then
        print_info "Installing goimports..."
        go install golang.org/x/tools/cmd/goimports@latest
        print_success "goimports installed"
    else
        print_info "goimports already installed"
    fi

    # Build services
    print_header "Building Services"
    print_info "Building all services (this may take a while)..."
    make build
    print_success "All services built successfully"

    # Start infrastructure with Docker
    if command_exists docker; then
        print_header "Starting Infrastructure"
        read -p "Do you want to start PostgreSQL and Redis with Docker? (y/n) " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            print_info "Starting PostgreSQL..."
            docker run -d --name repograph-postgres \
                -e POSTGRES_USER=repograph \
                -e POSTGRES_PASSWORD=repograph \
                -e POSTGRES_DB=repograph_db \
                -p 5432:5432 \
                postgres:15-alpine 2>/dev/null || print_warning "PostgreSQL container might already exist"

            print_info "Starting Redis..."
            docker run -d --name repograph-redis \
                -p 6379:6379 \
                redis:7-alpine 2>/dev/null || print_warning "Redis container might already exist"

            sleep 3
            print_success "Infrastructure started"
        fi
    fi

    # Run tests
    print_header "Running Tests"
    read -p "Do you want to run tests? (y/n) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        print_info "Running tests..."
        go test ./... -v
        print_success "Tests completed"
    fi

    # Summary
    print_header "Setup Complete! 🎉"
    echo ""
    echo "Next steps:"
    echo ""
    echo "1. Edit .env file with your API keys:"
    echo "   ${YELLOW}nano .env${NC}"
    echo ""
    echo "2. Start the services:"
    echo "   ${YELLOW}docker-compose -f deployments/docker/docker-compose.yml up -d${NC}"
    echo ""
    echo "   Or start individually:"
    echo "   ${YELLOW}./bin/orchestrator${NC}"
    echo ""
    echo "3. Index documents:"
    echo "   ${YELLOW}./bin/repograph-cli index --directory ./data/diagrams${NC}"
    echo ""
    echo "4. Query the knowledge base:"
    echo "   ${YELLOW}./bin/repograph-cli query ask \"What is Choreo?\"${NC}"
    echo ""
    echo "For more information, see:"
    echo "  - ${BLUE}README.md${NC}"
    echo "  - ${BLUE}docs/ARCHITECTURE.md${NC}"
    echo "  - ${BLUE}docs/DEPLOYMENT.md${NC}"
    echo ""
}

# Run main function
main
