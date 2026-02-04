# Development Guide

## Table of Contents
- [Getting Started](#getting-started)
- [Project Structure](#project-structure)
- [Development Workflow](#development-workflow)
- [Coding Standards](#coding-standards)
- [Testing](#testing)
- [Debugging](#debugging)
- [Adding New Features](#adding-new-features)
- [Git Workflow](#git-workflow)

---

## Getting Started

### Quick Start

```bash
# Clone the repository
cd /home/nadeeshame/go/rag-knowledge-service

# Run setup script
./scripts/setup.sh

# Start development
make dev
```

### Manual Setup

```bash
# Install dependencies
go mod download

# Install development tools
make install-tools

# Build all services
make build

# Run tests
make test
```

---

## Project Structure

```
rag-knowledge-service/
├── cmd/                      # Service entry points
│   ├── orchestrator/        # Main orchestrator
│   ├── document-scanner/    # File scanning service
│   ├── content-extractor/   # Content extraction
│   ├── vision-service/      # Image analysis
│   ├── summarization-service/
│   ├── embedding-service/
│   ├── vector-store/
│   ├── query-service/
│   └── repograph-cli/       # CLI application
│
├── internal/                # Private application code
│   ├── domain/             # Domain layer (SOLID)
│   │   ├── models/        # Business entities
│   │   └── interfaces/    # Service contracts
│   ├── orchestrator/      # Workflow coordination
│   ├── document-scanner/  # Scanner implementation
│   ├── content-extractor/ # Extractor implementation
│   │   └── processors/   # File processors
│   ├── adapters/         # External service adapters
│   ├── config/           # Configuration management
│   ├── logger/           # Logging utilities
│   └── middleware/       # HTTP middleware
│
├── pkg/                   # Public libraries
│   ├── utils/           # Utility functions
│   └── health/          # Health checking
│
├── docs/                # Documentation
├── tests/               # Tests
├── scripts/             # Utility scripts
└── deployments/         # Deployment configs
```

### Package Naming Conventions

- **cmd/**: Main applications
- **internal/**: Private application code
- **pkg/**: Public, reusable libraries
- **api/**: API definitions (protobuf, OpenAPI)
- **configs/**: Configuration files
- **deployments/**: Infrastructure configs

---

## Development Workflow

### 1. Create Feature Branch

```bash
git checkout -b feature/my-feature
```

### 2. Implement Feature

Follow TDD approach:
1. Write tests first
2. Implement functionality
3. Refactor

### 3. Run Tests

```bash
# Run all tests
make test

# Run specific package
go test ./internal/document-scanner/...

# With coverage
make test-coverage
```

### 4. Format and Lint

```bash
# Format code
make fmt

# Run linter
make lint

# Run all checks
make check
```

### 5. Commit Changes

```bash
# Conventional commit format
git commit -m "feat: add new document processor"
git commit -m "fix: resolve parsing error"
git commit -m "docs: update API documentation"
```

### 6. Push and Create PR

```bash
git push origin feature/my-feature
# Create PR on GitHub
```

---

## Coding Standards

### Go Best Practices

1. **Follow Effective Go**: https://golang.org/doc/effective_go
2. **Use gofmt**: Code must be formatted
3. **Error handling**: Always check errors
4. **Naming conventions**: 
   - Packages: lowercase, single word
   - Interfaces: noun or adjective + "er" (e.g., Reader, Scanner)
   - Variables: camelCase
   - Constants: MixedCaps or ALL_CAPS for exported

### Code Style

```go
// Good
func ProcessDocument(ctx context.Context, doc *models.Document) error {
    if doc == nil {
        return fmt.Errorf("document is nil")
    }
    
    // Implementation
    return nil
}

// Bad
func processDocument(doc *models.Document) {
    // No context
    // No error handling
    // Not exported when it should be
}
```

### Interface Design

```go
// Good: Small, focused interface
type DocumentScanner interface {
    ScanDirectory(ctx context.Context, dir string) ([]*models.FileMetadata, error)
}

// Bad: Large, unfocused interface
type DocumentManager interface {
    Scan(dir string) []File
    Process(file File) Result
    Store(result Result) bool
    Query(q string) Answer
    // Too many responsibilities
}
```

### Error Handling

```go
// Good
if err := processFile(file); err != nil {
    logger.Error("failed to process file",
        zap.String("file", file),
        zap.Error(err))
    return fmt.Errorf("processing file %s: %w", file, err)
}

// Bad
err := processFile(file)
// No error handling
```

### Logging

```go
// Good
logger.Info("processing document",
    zap.String("document_id", doc.ID.String()),
    zap.String("file_name", doc.FileName),
    zap.Int64("file_size", doc.FileSize))

// Bad
log.Println("Processing document:", doc.ID)
```

---

## Testing

### Unit Tests

```go
func TestDocumentScanner_ScanDirectory(t *testing.T) {
    scanner := NewDocumentScanner()
    
    tests := []struct {
        name    string
        dir     string
        want    int
        wantErr bool
    }{
        {
            name:    "valid directory",
            dir:     "./testdata",
            want:    5,
            wantErr: false,
        },
        {
            name:    "invalid directory",
            dir:     "./nonexistent",
            want:    0,
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            files, err := scanner.ScanDirectory(context.Background(), tt.dir)
            if (err != nil) != tt.wantErr {
                t.Errorf("ScanDirectory() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if len(files) != tt.want {
                t.Errorf("ScanDirectory() got %v files, want %v", len(files), tt.want)
            }
        })
    }
}
```

### Integration Tests

```go
//go:build integration
// +build integration

func TestDocumentProcessing_Integration(t *testing.T) {
    // Setup test environment
    db := setupTestDB(t)
    defer db.Close()
    
    // Run integration test
    // ...
}
```

Run integration tests:
```bash
go test -tags=integration ./tests/integration/...
```

### Test Coverage

```bash
# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# Or use make
make test-coverage
```

### Benchmarks

```go
func BenchmarkProcessDocument(b *testing.B) {
    doc := &models.Document{
        Content: strings.Repeat("test ", 1000),
    }
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        processDocument(doc)
    }
}
```

Run benchmarks:
```bash
go test -bench=. ./...
```

---

## Debugging

### Using Delve

```bash
# Install delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug a service
dlv debug ./cmd/orchestrator

# Debug tests
dlv test ./internal/document-scanner
```

### Logging for Debugging

```go
logger.Debug("detailed debug info",
    zap.Any("data", complexObject),
    zap.Stack("stack"))
```

### VS Code Debug Configuration

Create `.vscode/launch.json`:

```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Debug Orchestrator",
            "type": "go",
            "request": "launch",
            "mode": "debug",
            "program": "${workspaceFolder}/cmd/orchestrator",
            "env": {
                "LOG_LEVEL": "debug"
            }
        }
    ]
}
```

---

## Adding New Features

### Adding a New File Processor

1. **Create processor file**:
```go
// internal/content-extractor/processors/my_processor.go
package processors

type MyProcessor struct{}

func (p *MyProcessor) Extract(ctx context.Context, filePath string) (string, error) {
    // Implementation
    return content, nil
}

func (p *MyProcessor) SupportsFileType(fileType string) bool {
    return fileType == "mytype"
}
```

2. **Register processor**:
```go
// internal/content-extractor/service.go
func NewContentExtractor() *ContentExtractor {
    processors := []processors.Processor{
        &processors.ImageProcessor{},
        &processors.MyProcessor{}, // Add here
    }
    return &ContentExtractor{processors: processors}
}
```

3. **Add tests**:
```go
func TestMyProcessor_Extract(t *testing.T) {
    // Test implementation
}
```

### Adding a New Service

1. Create service directory in `cmd/` and `internal/`
2. Implement service interface
3. Add Dockerfile in `deployments/docker/`
4. Add to docker-compose.yml
5. Update documentation

### Adding a New API Endpoint

```go
// cmd/orchestrator/main.go
v1.POST("/new-endpoint", newEndpointHandler)

func newEndpointHandler(c *gin.Context) {
    // Implementation
}
```

---

## Git Workflow

### Branching Strategy

- `main`: Production-ready code
- `develop`: Development branch
- `feature/*`: New features
- `fix/*`: Bug fixes
- `hotfix/*`: Emergency fixes

### Commit Messages

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
feat: add new feature
fix: resolve bug
docs: update documentation
style: format code
refactor: refactor code
test: add tests
chore: update dependencies
```

### Pull Request Process

1. Create feature branch
2. Implement feature with tests
3. Run `make check`
4. Push and create PR
5. Request reviews
6. Address feedback
7. Merge after approval

---

## Environment Setup

### Required Tools

```bash
# Go
go version  # 1.21+

# Development tools
golangci-lint --version
goimports --version

# Docker
docker --version
docker-compose --version

# Optional
kubectl version  # For Kubernetes
helm version     # For Helm charts
```

### IDE Setup

#### VS Code Extensions
- Go (golang.go)
- Go Test Explorer
- Docker
- YAML
- EditorConfig

#### GoLand
- Enable Go modules
- Configure Go tools
- Set up code style

---

## Useful Commands

```bash
# Development
make dev          # Start development environment
make build        # Build all services
make test         # Run tests
make lint         # Run linter
make fmt          # Format code

# Docker
make docker-build # Build Docker images
make docker-up    # Start with Docker Compose
make docker-down  # Stop containers

# Cleaning
make clean        # Clean build artifacts
go clean -modcache # Clean module cache
```

---

## Resources

- [Effective Go](https://golang.org/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)
- [SOLID Principles in Go](https://dave.cheney.net/2016/08/20/solid-go-design)

---

*Last Updated: February 2026*
