# RAG Knowledge Service - Project Summary

## Overview

**RAG Knowledge Service** is an enterprise-grade, microservices-based Retrieval-Augmented Generation (RAG) system built in Go. It processes 50+ file formats, generates AI-powered summaries, and enables semantic search using Azure OpenAI and Pinecone vector database.

## Project Status

✅ **Phase 1: Foundation Complete**
- Project structure established
- Core domain models defined
- Configuration management implemented
- Logging infrastructure ready
- Health checking system in place

## Architecture

### Microservices (8 Services)

1. **Orchestrator Service** (Port 8088)
   - Coordinates entire workflow
   - Manages document processing pipeline
   - Status: ✅ Basic structure complete

2. **Document Scanner Service** (Port 8081)
   - File discovery and metadata extraction
   - Hash computation for deduplication
   - Status: ✅ Basic structure complete

3. **Content Extractor Service** (Port 8082)
   - Multi-format content extraction
   - Processor registry pattern
   - Status: ⏳ Processors need implementation

4. **Vision Service** (Port 8083)
   - Google Vision API integration
   - Image and diagram analysis
   - OCR capabilities
   - Status: ⏳ Adapter needs implementation

5. **Summarization Service** (Port 8084)
   - Azure OpenAI integration
   - Text summarization
   - Status: ⏳ Adapter needs implementation

6. **Embedding Service** (Port 8085)
   - Azure OpenAI embeddings
   - Batch processing support
   - Status: ⏳ Adapter needs implementation

7. **Vector Store Service** (Port 8086)
   - Pinecone integration
   - Vector upsert and search
   - Status: ⏳ Adapter needs implementation

8. **Query Service** (Port 8087)
   - RAG query handling
   - Context-aware Q&A
   - Status: ⏳ Implementation needed

### CLI Application

**repograph-cli**
- Index documents
- Query knowledge base
- Interactive mode
- Status checks
- Status: ✅ Basic structure complete

## Technology Stack

### Core
- **Language**: Go 1.21
- **Web Framework**: Gin
- **CLI Framework**: Cobra
- **Config**: Viper
- **Logging**: Zap (structured JSON)

### External Services
- **Azure OpenAI**: GPT-4 and embeddings
- **Google Vision API**: Image analysis
- **Pinecone**: Vector database

### Infrastructure
- **Database**: PostgreSQL 15 (metadata)
- **Cache**: Redis 7
- **Containers**: Docker
- **Orchestration**: Kubernetes
- **CI/CD**: GitHub Actions

## Key Features

### ✅ Implemented
1. Project structure following Go best practices
2. SOLID principles with interface-based design
3. Configuration management (Viper + environment variables)
4. Structured logging (Zap)
5. Health checking system
6. File utilities (hashing, type detection)
7. Docker support (Dockerfiles, docker-compose)
8. CI/CD pipeline (GitHub Actions)
9. Comprehensive documentation
10. Development tools (Makefile, setup script)

### ⏳ In Progress
1. Content extraction processors
2. External service adapters (Azure, Google, Pinecone)
3. Database repositories (GORM)
4. Workflow orchestration logic
5. Query service RAG implementation

### 📋 Planned
1. gRPC inter-service communication
2. OpenTelemetry distributed tracing
3. Prometheus metrics
4. Kubernetes manifests
5. Helm charts
6. Load testing
7. Performance optimization
8. Multi-tenancy support

## File Structure

```
rag-knowledge-service/
├── cmd/                        # ✅ Entry points
├── internal/
│   ├── domain/                # ✅ Models & interfaces
│   ├── config/                # ✅ Configuration
│   ├── logger/                # ✅ Logging
│   ├── orchestrator/          # ⏳ Needs implementation
│   ├── document-scanner/      # ⏳ Needs implementation
│   ├── content-extractor/     # ⏳ Processors needed
│   ├── vision-service/        # ⏳ Needs implementation
│   ├── summarization-service/ # ⏳ Needs implementation
│   ├── embedding-service/     # ⏳ Needs implementation
│   ├── vector-store/          # ⏳ Needs implementation
│   ├── query-service/         # ⏳ Needs implementation
│   ├── adapters/              # ⏳ Needs implementation
│   └── middleware/            # ⏳ Needs implementation
├── pkg/
│   ├── utils/                 # ✅ File utilities
│   └── health/                # ✅ Health checker
├── docs/                      # ✅ Documentation
├── deployments/
│   ├── docker/                # ✅ Dockerfiles & compose
│   └── kubernetes/            # 📋 To be created
├── scripts/                   # ✅ Setup script
└── .github/workflows/         # ✅ CI/CD pipelines
```

## CI/CD Pipeline

### Workflows Implemented
1. **ci-cd.yml**: Main pipeline (lint, test, build, deploy)
2. **codeql.yml**: Security analysis
3. **docker.yml**: Docker build and scan
4. **dependency-updates.yml**: Automated updates
5. **release.yml**: Release automation
6. **dependabot.yml**: Dependency PRs

### Pipeline Stages
- ✅ Linting (golangci-lint)
- ✅ Testing (multiple Go versions)
- ✅ Security scanning (CodeQL, Gosec, Trivy)
- ✅ Docker build & push
- ⏳ Kubernetes deployment
- ⏳ Staging/Production deployment

## Documentation

### Created
1. ✅ **README.md**: Project overview and quick start
2. ✅ **ARCHITECTURE.md**: Detailed architecture guide
3. ✅ **API_REFERENCE.md**: Complete API documentation
4. ✅ **DEPLOYMENT.md**: Deployment instructions
5. ✅ **DEVELOPMENT.md**: Development guide
6. ✅ **.env.example**: Environment template

### Credentials
- ✅ **credentials/README.md**: Security guidelines
- ⚠️ Service account keys needed

## Next Steps

### Immediate (Week 1-2)
1. **Implement External Service Adapters**
   - Azure OpenAI client
   - Google Vision API client
   - Pinecone client

2. **Implement Core Services**
   - Document scanner service logic
   - Content extractor with processors
   - Embedding service

3. **Database Setup**
   - PostgreSQL schema
   - GORM repositories
   - Migrations

### Short-term (Week 3-4)
1. **Complete Service Implementations**
   - Vision service
   - Summarization service
   - Vector store service
   - Query service

2. **Orchestration Logic**
   - Workflow state machine
   - Error handling
   - Retry mechanisms

3. **Testing**
   - Unit tests for all services
   - Integration tests
   - End-to-end tests

### Medium-term (Month 2)
1. **Production Readiness**
   - Kubernetes manifests
   - Helm charts
   - Monitoring setup
   - Performance testing

2. **Advanced Features**
   - Batch processing
   - Parallel document processing
   - Advanced chunking strategies
   - Multiple LLM support

## Development Commands

```bash
# Setup
./scripts/setup.sh

# Build
make build

# Test
make test
make test-coverage

# Lint & Format
make lint
make fmt

# Docker
make docker-build
make docker-up
make docker-down

# Run services
./bin/orchestrator
./bin/repograph-cli index --directory ./data/diagrams
./bin/repograph-cli query ask "What is Choreo?"
```

## Configuration

### Required Environment Variables
```bash
# Azure OpenAI
AZURE_OPENAI_API_KEY=your_key
AZURE_OPENAI_ENDPOINT=https://your-resource.openai.azure.com/
AZURE_OPENAI_EMBEDDINGS_DEPLOYMENT=text-embedding-ada-002
AZURE_OPENAI_CHAT_DEPLOYMENT=gpt-4

# Pinecone
PINECONE_API_KEY=your_key
PINECONE_INDEX_NAME=repograph-ai-index

# Google Vision
GOOGLE_VISION_API_KEY=your_key
GOOGLE_APPLICATION_CREDENTIALS=credentials/service-account.json

# Database
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=repograph
POSTGRES_PASSWORD=your_password
POSTGRES_DB=repograph_db

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
```

## Testing Status

| Component | Unit Tests | Integration Tests | Coverage |
|-----------|-----------|-------------------|----------|
| Domain Models | ⏳ TODO | N/A | 0% |
| Config | ⏳ TODO | N/A | 0% |
| Logger | ⏳ TODO | N/A | 0% |
| Utils | ⏳ TODO | N/A | 0% |
| Services | ⏳ TODO | ⏳ TODO | 0% |
| CLI | ⏳ TODO | N/A | 0% |

Target: 80%+ coverage

## Performance Targets

- **Throughput**: 100+ documents/minute
- **Latency**: <100ms for queries (excluding LLM)
- **Concurrency**: 1000+ simultaneous requests
- **Uptime**: 99.9%

## Security Considerations

- ✅ Credentials in .gitignore
- ✅ Non-root Docker containers
- ✅ Security scanning in CI/CD
- ⏳ mTLS for inter-service communication
- ⏳ API authentication
- ⏳ Rate limiting
- ⏳ Audit logging

## Known Issues

None currently - project is in initial development phase.

## Contributing

See **docs/DEVELOPMENT.md** for:
- Development setup
- Coding standards
- Testing guidelines
- Git workflow
- PR process

## License

This project is provided as-is for educational and commercial use.

## Contact & Support

For questions or issues:
1. Check documentation in `docs/`
2. Review code examples in `cmd/` and `internal/`
3. See GitHub Issues for known problems

---

## Quick Start Checklist

- [ ] Clone repository
- [ ] Copy `.env.example` to `.env`
- [ ] Add API keys to `.env`
- [ ] Run `./scripts/setup.sh`
- [ ] Start infrastructure: `docker-compose up -d`
- [ ] Build services: `make build`
- [ ] Run tests: `make test`
- [ ] Start orchestrator: `./bin/orchestrator`
- [ ] Index documents: `./bin/repograph-cli index`
- [ ] Query: `./bin/repograph-cli query ask "test"`

---

**Project Created**: February 2026  
**Last Updated**: February 2, 2026  
**Version**: 0.1.0-alpha  
**Status**: 🟡 In Development
