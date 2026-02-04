# RAG Knowledge Service

<p align="center">
  <strong>🚀 Enterprise-Grade Intelligent Document Processing & RAG System in Go</strong>
</p>

<p align="center">
  <a href="#features">Features</a> •
  <a href="#quick-start">Quick Start</a> •
  <a href="#architecture">Architecture</a> •
  <a href="#documentation">Documentation</a> •
  <a href="#api-reference">API Reference</a>
</p>

---

## 🌟 Overview

**RAG Knowledge Service** is a production-ready, microservices-based Retrieval-Augmented Generation (RAG) application built in Go that transforms your documents into an intelligent, searchable knowledge base. Built with enterprise requirements in mind, it processes 50+ file types, generates comprehensive summaries, and enables semantic search powered by Azure OpenAI and Pinecone.

### Key Capabilities

- 📊 **Multi-Format Processing**: Images, diagrams, documents, spreadsheets, code, and more
- 🤖 **AI-Powered Analysis**: Google Vision API for visual content, Azure OpenAI for understanding
- 🔍 **Semantic Search**: Find information using natural language queries
- 💬 **RAG Q&A**: Get accurate answers with source citations
- ⚡ **Incremental Processing**: Smart deduplication saves time and costs
- 🏗️ **Enterprise Architecture**: SOLID principles, microservices, comprehensive logging
- 🐹 **High Performance**: Built in Go for speed and efficiency

---

## ✨ Features

### 📚 Multi-Format Document Processing

| Category | Supported Formats |
|----------|------------------|
| **Images** | PNG, JPG, JPEG, SVG, GIF, BMP, WEBP |
| **Diagrams** | DrawIO, Excalidraw |
| **Documents** | DOCX, PDF, PPTX, ODT, TXT, MD |
| **Spreadsheets** | XLSX, XLS, CSV |
| **Structured** | JSON, GraphQL, YAML, XML, TOML |
| **Code** | Go, Python, JavaScript, TypeScript, Java, C/C++, Rust, SQL, and more |
| **Text** | Markdown, TXT, LOG, config files |

### 🧠 Intelligent Processing Pipeline

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   Scan      │────▶│   Extract   │────▶│   Analyze   │────▶│  Summarize  │
│   Files     │     │   Content   │     │   (Vision)  │     │   (LLM)     │
└─────────────┘     └─────────────┘     └─────────────┘     └─────────────┘
                                                                    │
                                                                    ▼
┌─────────────┐     ┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   Query     │◀────│   Search    │◀────│   Store     │◀────│   Embed     │
│   (RAG)     │     │  (Vector)   │     │  (Pinecone) │     │  (Azure AI) │
└─────────────┘     └─────────────┘     └─────────────┘     └─────────────┘
```

### 🔧 Enterprise Features

- **Smart Deduplication**: Automatically skip already-indexed documents
- **Incremental Updates**: Only process new or modified files
- **Comprehensive Logging**: Structured JSON logging with timing metrics
- **Health Monitoring**: Service health checks for all external dependencies
- **Error Recovery**: Graceful handling of failures with detailed error reporting
- **Configurable**: Environment-based configuration for easy deployment
- **Microservices**: Independent, scalable services with clear boundaries
- **SOLID Principles**: Clean architecture with dependency inversion

---

## 🚀 Quick Start

### Prerequisites

- Go 1.21 or higher
- Azure OpenAI account with API access
- Pinecone account (free tier works)
- Google Vision API key (optional, for image analysis)
- Redis (for caching)

### Installation

```bash
# Clone or navigate to project
cd /home/nadeeshame/go/rag-knowledge-service

# Copy environment configuration
cp .env.example .env

# Edit .env and add your API keys
nano .env

# Download dependencies
go mod download

# Build all services
make build

# Or build specific service
go build -o bin/rag-cli cmd/rag-cli/main.go
```

### Configuration

Edit the `.env` file with your credentials:

```bash
# Azure OpenAI
AZURE_OPENAI_API_KEY=your_key_here
AZURE_OPENAI_ENDPOINT=https://your-resource.openai.azure.com/
AZURE_OPENAI_EMBEDDINGS_DEPLOYMENT=text-embedding-ada-002
AZURE_OPENAI_CHAT_DEPLOYMENT=gpt-4

# Pinecone
PINECONE_API_KEY=your_pinecone_key
PINECONE_INDEX_NAME=repograph-ai-index

# Google Vision (optional)
GOOGLE_VISION_API_KEY=your_vision_key

# Data directory
DATA_DIRECTORY=./data/diagrams
```

### Run Services

#### Option 1: Docker Compose (Recommended)

```bash
# 1. Create .env file
cp .env.example .env
nano .env  # Add AZURE_OPENAI_API_KEY, PINECONE_API_KEY, etc.

# 2. Start services (automatically indexes documents on startup)
docker-compose up -d

# The orchestrator will automatically:
# - Scan documents in DATA_DIRECTORY
# - Extract content and analyze images
# - Generate embeddings via Azure OpenAI
# - Store vectors in Pinecone
# - Skip already-indexed documents (deduplication)

# View indexing progress
docker-compose logs -f orchestrator

# Check all services
docker-compose ps

# Stop services
docker-compose down
```

**Note**: On first startup, the orchestrator automatically indexes all documents in the configured `DATA_DIRECTORY` (default: `./data/diagrams`). This process runs in the background and can be monitored via logs.

#### Option 2: Manual Start

```bash
# Start Redis (using Docker)
docker run -d --name redis -p 6379:6379 redis:7

# Start each microservice in separate terminals
./bin/document-scanner
./bin/content-extractor
./bin/vision-service
./bin/summarization-service
./bin/embedding-service
./bin/vector-store
./bin/query-service
./bin/orchestrator
```

### Index Documents

**Automatic Indexing**: When you start the services with `docker-compose up -d`, the orchestrator automatically indexes all documents in the `DATA_DIRECTORY`.

**Manual Indexing** (optional - for CLI usage):

```bash
# Re-index all documents manually
./bin/rag-cli index

# Force reprocess all documents (ignores deduplication)
./bin/rag-cli index --force

# Index a specific directory
./bin/rag-cli index --directory ./my-docs
```

**Indexing Process**:
1. 🔍 Scans all files in directory
2. 📄 Extracts content based on file type
3. 👁️ Analyzes images/diagrams with Google Vision
4. 📝 Generates summaries with Azure OpenAI
5. 🧮 Creates embeddings (1536 dimensions)
6. 💾 Stores vectors + metadata in Pinecone
7. ⚡ Skips already-indexed files (hash-based deduplication)

### Query the Knowledge Base

```bash
# Ask a question
./bin/rag-cli query ask "What is the Choreo architecture?"

# Search documents
./bin/rag-cli query search "authentication flow"

# Interactive mode
./bin/rag-cli query interactive
```

---

## 🏗️ Architecture

### Microservices Design

RAG Knowledge Service follows a clean microservices architecture with clear separation of concerns:

```
┌────────────────────────────────────────────────────────────────┐
│                    CLI Layer (rag-cli)                   │
└────────────────────────────────────────────────────────────────┘
                                  │
                                  ▼
┌────────────────────────────────────────────────────────────────┐
│                    Orchestrator Service                         │
│         ┌─────────────────────────────────────────────┐        │
│         │  scan → extract → analyze → summarize →     │        │
│         │  chunk → embed → store → finalize           │        │
│         └─────────────────────────────────────────────┘        │
└────────────────────────────────────────────────────────────────┘
                                  │
          ┌───────────┬───────────┼───────────┬───────────┐
          ▼           ▼           ▼           ▼           ▼
     ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐
     │Document │ │ Vision  │ │Summarize│ │Embedding│ │ Vector  │
     │ Scanner │ │ Service │ │ Service │ │ Service │ │  Store  │
     └─────────┘ └─────────┘ └─────────┘ └─────────┘ └─────────┘
          │           │           │           │           │
          ▼           ▼           ▼           ▼           ▼
     ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐
     │File     │ │ Google  │ │ Azure   │ │ Azure   │ │Pinecone │
     │System   │ │ Vision  │ │ OpenAI  │ │ OpenAI  │ │         │
     └─────────┘ └─────────┘ └─────────┘ └─────────┘ └─────────┘
```

### Services

| Service | Port | Responsibility |
|---------|------|----------------|
| **Document Scanner** | 8081 | File discovery and metadata extraction |
| **Content Extractor** | 8082 | Multi-format content extraction |
| **Vision Service** | 8083 | Image and diagram analysis (Google Vision) |
| **Summarization Service** | 8084 | Content summarization (Azure OpenAI) |
| **Embedding Service** | 8085 | Generate embeddings (Azure OpenAI) |
| **Vector Store** | 8086 | Vector database operations (Pinecone) |
| **Query Service** | 8087 | RAG query handling |
| **Orchestrator** | 8088 | Workflow coordination |

### SOLID Principles

| Principle | Implementation |
|-----------|----------------|
| **S**ingle Responsibility | Each service handles one specific task |
| **O**pen/Closed | Easy to extend with new processors without modifying existing code |
| **L**iskov Substitution | All implementations follow interface contracts |
| **I**nterface Segregation | Small, focused interfaces (DocumentScanner, VisionAnalyzer, etc.) |
| **D**ependency Inversion | High-level modules depend on abstractions |

### Project Structure

```
rag-knowledge-service/
├── cmd/                        # Service entry points
│   ├── orchestrator/          # Orchestrator service
│   ├── document-scanner/      # Document scanner service
│   ├── content-extractor/     # Content extractor service
│   ├── vision-service/        # Vision service
│   ├── summarization-service/ # Summarization service
│   ├── embedding-service/     # Embedding service
│   ├── vector-store/          # Vector store service
│   ├── query-service/         # Query service
│   └── rag-cli/         # CLI application
│
├── internal/                   # Private application code
│   ├── domain/                # Domain layer
│   │   ├── models/           # Domain models
│   │   └── interfaces/       # Service interfaces (SOLID)
│   ├── orchestrator/         # Orchestrator implementation
│   ├── document-scanner/     # Scanner implementation
│   ├── content-extractor/    # Extractor implementation
│   │   └── processors/       # Format-specific processors
│   ├── vision-service/       # Vision service implementation
│   ├── summarization-service/# Summarization implementation
│   ├── embedding-service/    # Embedding implementation
│   ├── vector-store/         # Vector store implementation
│   ├── query-service/        # Query service implementation
│   ├── adapters/             # External service adapters
│   │   ├── azure/           # Azure OpenAI adapter
│   │   ├── google/          # Google Vision adapter
│   │   └── pinecone/        # Pinecone adapter
│   ├── config/              # Configuration management
│   ├── logger/              # Logging utilities
│   └── middleware/          # HTTP middleware
│
├── pkg/                       # Public libraries
│   ├── utils/                # Utility functions
│   └── health/               # Health checking
│
├── api/                       # API definitions
│   ├── proto/                # Protocol buffer definitions
│   └── openapi/              # OpenAPI specifications
│
├── configs/                   # Configuration files
│   └── config.yaml           # Default configuration
│
├── deployments/               # Deployment configurations
│   ├── docker/               # Dockerfiles
│   │   └── docker-compose.yml
│   └── kubernetes/           # Kubernetes manifests
│       ├── deployments/
│       ├── services/
│       └── ingress/
│
├── docs/                      # Documentation
│   ├── ARCHITECTURE.md       # Architecture guide
│   ├── API_REFERENCE.md      # API documentation
│   ├── DEPLOYMENT.md         # Deployment guide
│   └── DEVELOPMENT.md        # Development guide
│
├── scripts/                   # Utility scripts
│   ├── setup.sh              # Setup script
│   └── test.sh               # Test runner
│
├── tests/                     # Tests
│   ├── integration/          # Integration tests
│   └── e2e/                  # End-to-end tests
│
├── .github/                   # GitHub configuration
│   ├── workflows/            # GitHub Actions workflows
│   │   ├── ci.yml           # CI pipeline
│   │   ├── codeql.yml       # Security analysis
│   │   ├── docker.yml       # Docker build & push
│   │   ├── dependabot.yml   # Dependency updates
│   │   └── release.yml      # Release automation
│   └── dependabot.yml       # Dependabot configuration
│
├── credentials/               # Credentials (gitignored)
│   ├── README.md             # Security guidelines
│   └── *.json               # Service account keys
│
├── data/                      # Data directory
│   ├── diagrams/             # Sample diagrams
│   └── README.md             # Data organization guide
│
├── go.mod                     # Go module definition
├── go.sum                     # Dependency checksums
├── Makefile                   # Build automation
├── .env.example               # Environment template
├── .gitignore                 # Git ignore rules
└── README.md                  # This file
```

---

## 📖 Documentation

| Document | Description |
|----------|-------------|
| [Architecture Guide](docs/ARCHITECTURE.md) | Detailed architecture and design patterns |
| [API Reference](docs/API_REFERENCE.md) | Service APIs and interfaces |
| [Deployment Guide](docs/DEPLOYMENT.md) | Production deployment instructions |
| [Development Guide](docs/DEVELOPMENT.md) | Local development setup and guidelines |
| [CodeQL Troubleshooting](docs/CODEQL_TROUBLESHOOTING.md) | Fix CodeQL workflow issues |

---

## 🔄 CI/CD Pipeline

RAG Knowledge Service includes a comprehensive CI/CD pipeline using GitHub Actions.

### Workflows

| Workflow | Purpose |
|----------|---------|
| **CI** | Lint, test, build |
| **CodeQL** | Security analysis |
| **Docker** | Build & push images |
| **Dependabot** | Automated dependency updates |
| **Release** | Automated releases |

### Pipeline Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                         CI/CD Pipeline                           │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌─────────┐    ┌─────────┐    ┌─────────┐    ┌──────────────┐ │
│  │  Lint   │───▶│  Test   │───▶│Security │───▶│Build Docker  │ │
│  │(golangci│    │ (go test│    │ (CodeQL)│    │   Images     │ │
│  │  -lint) │    │coverage)│    │ (gosec) │    │              │ │
│  └─────────┘    └─────────┘    └─────────┘    └──────┬───────┘ │
│                                                        │         │
│                                          ┌─────────────┴──────┐  │
│                                          │                    │  │
│                                          ▼                    ▼  │
│                                 ┌──────────────┐     ┌──────────┐│
│                                 │Deploy Staging│     │Deploy Prod││
│                                 │(GCP/K8s)     │────▶│(GCP/K8s) ││
│                                 └──────────────┘     └──────────┘│
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

## 🐳 Docker Deployment

### Quick Start

```bash
# Build all images
docker-compose build

# Start services
docker-compose up -d

# View logs
docker-compose logs -f orchestrator

# Scale a service
docker-compose up -d --scale embedding-service=3

# Stop services
docker-compose down
```

### Individual Service

```bash
# Build
docker build -f deployments/docker/Dockerfile.orchestrator -t repograph-orchestrator .

# Run
docker run -d \
  --name orchestrator \
  -p 8088:8088 \
  --env-file .env \
  repograph-orchestrator
```

---

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package
go test ./internal/document-scanner/...

# Integration tests
go test -tags=integration ./tests/integration/...

# Benchmarks
go test -bench=. ./...
```

---

## 🛠️ Development

### Prerequisites

- Go 1.21+
- Docker & Docker Compose
- Make
- golangci-lint

### Setup

```bash
# Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Install development dependencies
make dev-setup

# Run linter
make lint

# Format code
make fmt

# Run tests
make test
```

### Adding a New Processor

1. Create processor in `internal/content-extractor/processors/`
2. Implement `ContentExtractor` interface
3. Register in `content-extractor/service.go`
4. Add tests
5. Update documentation

---

## 📊 Performance

- **Throughput**: Process 100+ documents/minute
- **Latency**: < 100ms for queries (excluding LLM)
- **Concurrency**: Handle 1000+ concurrent requests
- **Memory**: ~500MB per service
- **Storage**: Minimal local storage, uses Pinecone for vectors

---

## 🤝 Contributing

We welcome contributions! Please follow these guidelines:

1. Follow Go best practices and conventions
2. Implement interfaces for new services
3. Add comprehensive tests
4. Update documentation
5. Follow SOLID principles
6. Use conventional commits

### Development Workflow

```bash
# Create feature branch
git checkout -b feature/my-feature

# Make changes and test
make test

# Lint code
make lint

# Commit with conventional format
git commit -m "feat: add new feature"

# Push and create PR
git push origin feature/my-feature
```

---

## 📄 License

This project is provided as-is for educational and commercial use.

---

## 🙏 Acknowledgments

- **Azure OpenAI** - Embeddings and language models
- **Pinecone** - Vector database
- **Google Vision** - Image analysis
- **Gin** - HTTP web framework
- **Cobra** - CLI framework
- **Viper** - Configuration management
- **Zap** - Structured logging
- **GORM** - ORM library

---

<p align="center">
  Built with ❤️ in Go for intelligent document processing
</p>
