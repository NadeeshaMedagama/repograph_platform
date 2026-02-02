.PHONY: help build clean test lint fmt run dev docker-build docker-up docker-down install-tools

# Variables
GO := go
GOFLAGS := -v
BINARY_DIR := bin
DOCKER_COMPOSE := docker-compose

# Service names
SERVICES := orchestrator document-scanner content-extractor vision-service summarization-service embedding-service vector-store query-service repograph-cli

# Colors for output
COLOR_RESET := \033[0m
COLOR_BOLD := \033[1m
COLOR_GREEN := \033[32m
COLOR_YELLOW := \033[33m
COLOR_BLUE := \033[34m

help: ## Show this help message
	@echo "$(COLOR_BOLD)RepoGraph Platform - Makefile Commands$(COLOR_RESET)"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "$(COLOR_GREEN)%-20s$(COLOR_RESET) %s\n", $$1, $$2}'

all: clean build test ## Clean, build, and test all services

build: ## Build all services
	@echo "$(COLOR_BLUE)Building all services...$(COLOR_RESET)"
	@mkdir -p $(BINARY_DIR)
	@for service in $(SERVICES); do \
		echo "$(COLOR_YELLOW)Building $$service...$(COLOR_RESET)"; \
		$(GO) build $(GOFLAGS) -o $(BINARY_DIR)/$$service ./cmd/$$service; \
	done
	@echo "$(COLOR_GREEN)✓ Build complete$(COLOR_RESET)"

build-orchestrator: ## Build orchestrator service
	@echo "$(COLOR_BLUE)Building orchestrator...$(COLOR_RESET)"
	@mkdir -p $(BINARY_DIR)
	@$(GO) build $(GOFLAGS) -o $(BINARY_DIR)/orchestrator ./cmd/orchestrator

build-cli: ## Build CLI
	@echo "$(COLOR_BLUE)Building CLI...$(COLOR_RESET)"
	@mkdir -p $(BINARY_DIR)
	@$(GO) build $(GOFLAGS) -o $(BINARY_DIR)/repograph-cli ./cmd/repograph-cli

clean: ## Clean build artifacts
	@echo "$(COLOR_YELLOW)Cleaning...$(COLOR_RESET)"
	@rm -rf $(BINARY_DIR)
	@$(GO) clean
	@echo "$(COLOR_GREEN)✓ Clean complete$(COLOR_RESET)"

test: ## Run tests
	@echo "$(COLOR_BLUE)Running tests...$(COLOR_RESET)"
	@$(GO) test -v -race -coverprofile=coverage.out ./...
	@echo "$(COLOR_GREEN)✓ Tests complete$(COLOR_RESET)"

test-coverage: test ## Run tests with coverage report
	@echo "$(COLOR_BLUE)Generating coverage report...$(COLOR_RESET)"
	@$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "$(COLOR_GREEN)✓ Coverage report generated: coverage.html$(COLOR_RESET)"

test-integration: ## Run integration tests
	@echo "$(COLOR_BLUE)Running integration tests...$(COLOR_RESET)"
	@$(GO) test -v -tags=integration ./tests/integration/...

lint: ## Run linter
	@echo "$(COLOR_BLUE)Running linter...$(COLOR_RESET)"
	@golangci-lint run ./...
	@echo "$(COLOR_GREEN)✓ Lint complete$(COLOR_RESET)"

fmt: ## Format code
	@echo "$(COLOR_BLUE)Formatting code...$(COLOR_RESET)"
	@gofmt -s -w .
	@goimports -w .
	@echo "$(COLOR_GREEN)✓ Format complete$(COLOR_RESET)"

vet: ## Run go vet
	@echo "$(COLOR_BLUE)Running go vet...$(COLOR_RESET)"
	@$(GO) vet ./...

tidy: ## Tidy and verify modules
	@echo "$(COLOR_BLUE)Tidying modules...$(COLOR_RESET)"
	@$(GO) mod tidy
	@$(GO) mod verify

download: ## Download dependencies
	@echo "$(COLOR_BLUE)Downloading dependencies...$(COLOR_RESET)"
	@$(GO) mod download

install-tools: ## Install development tools
	@echo "$(COLOR_BLUE)Installing development tools...$(COLOR_RESET)"
	@$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@$(GO) install golang.org/x/tools/cmd/goimports@latest
	@echo "$(COLOR_GREEN)✓ Tools installed$(COLOR_RESET)"

docker-build: ## Build Docker images
	@echo "$(COLOR_BLUE)Building Docker images...$(COLOR_RESET)"
	@$(DOCKER_COMPOSE) -f deployments/docker/docker-compose.yml build
	@echo "$(COLOR_GREEN)✓ Docker build complete$(COLOR_RESET)"

docker-up: ## Start services with Docker Compose
	@echo "$(COLOR_BLUE)Starting services...$(COLOR_RESET)"
	@$(DOCKER_COMPOSE) -f deployments/docker/docker-compose.yml up -d
	@echo "$(COLOR_GREEN)✓ Services started$(COLOR_RESET)"

docker-down: ## Stop services
	@echo "$(COLOR_YELLOW)Stopping services...$(COLOR_RESET)"
	@$(DOCKER_COMPOSE) -f deployments/docker/docker-compose.yml down
	@echo "$(COLOR_GREEN)✓ Services stopped$(COLOR_RESET)"

docker-logs: ## View Docker logs
	@$(DOCKER_COMPOSE) -f deployments/docker/docker-compose.yml logs -f

docker-ps: ## List running containers
	@$(DOCKER_COMPOSE) -f deployments/docker/docker-compose.yml ps

run-orchestrator: build-orchestrator ## Run orchestrator service
	@echo "$(COLOR_BLUE)Starting orchestrator...$(COLOR_RESET)"
	@./$(BINARY_DIR)/orchestrator

run-cli: build-cli ## Run CLI
	@./$(BINARY_DIR)/repograph-cli

dev-setup: install-tools download ## Setup development environment
	@echo "$(COLOR_GREEN)✓ Development environment setup complete$(COLOR_RESET)"

check: lint test ## Run linter and tests

benchmark: ## Run benchmarks
	@echo "$(COLOR_BLUE)Running benchmarks...$(COLOR_RESET)"
	@$(GO) test -bench=. -benchmem ./...

security: ## Run security checks
	@echo "$(COLOR_BLUE)Running security checks...$(COLOR_RESET)"
	@gosec ./...

proto-gen: ## Generate protobuf code
	@echo "$(COLOR_BLUE)Generating protobuf code...$(COLOR_RESET)"
	@protoc --go_out=. --go-grpc_out=. api/proto/*.proto

.DEFAULT_GOAL := help
