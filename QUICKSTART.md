# Quick Start Guide - RepoGraph Platform

## Prerequisites Check

Before starting, ensure you have:
- [ ] Go 1.21+ installed (`go version`)
- [ ] Docker installed (`docker --version`)
- [ ] Docker Compose installed (`docker-compose --version`)
- [ ] Azure OpenAI API key
- [ ] Pinecone API key
- [ ] (Optional) Google Vision API key

## Step-by-Step Setup

### 1. Configure Environment Variables

```bash
# Copy the example environment file
cp .env.example .env

# Edit with your credentials
nano .env
```

**Required Variables**:
```bash
AZURE_OPENAI_API_KEY=your_azure_key_here
AZURE_OPENAI_ENDPOINT=https://your-resource.openai.azure.com/
PINECONE_API_KEY=your_pinecone_key_here
PINECONE_INDEX_NAME=repograph-platform
```

### 2. Start Infrastructure Only (Redis)

If you want to develop locally without Docker for services:

```bash
# Start just Redis
docker-compose up -d redis

# Check it's running
docker-compose ps
```

### 3. Start All Services with Docker

```bash
# Build and start everything
docker-compose up -d

# This will:
# - Build all 8 microservices
# - Start Redis
# - Create the network
# - Set up health checks
```

**⚠️ IMPORTANT - Automatic Indexing Status**:

The infrastructure is ready, but the **automatic indexing workflow needs to be implemented**. 

**What's Ready**:
- ✅ All microservices running
- ✅ API endpoints available
- ✅ Configuration system in place
- ✅ Azure OpenAI & Pinecone clients ready

**What Needs Implementation** (See `cmd/orchestrator/main.go` TODO):
1. Document scanner integration
2. Content extraction pipeline
3. Vision API calls for images
4. Summary generation via Azure OpenAI
5. Embedding generation (1536 dims)
6. Pinecone storage with metadata
7. Hash-based deduplication

**Planned Workflow** (when implemented):
```
DATA_DIRECTORY (./data/diagrams)
         ↓
   Document Scanner → Extract Content → Analyze Images
         ↓                    ↓                ↓
    File Metadata      Text Content      Visual Analysis
         ↓                    ↓                ↓
         └────────────────────┴────────────────┘
                              ↓
                    Generate Summary (Azure OpenAI)
                              ↓
                    Create Embeddings (1536 dims)
                              ↓
                    Store in Pinecone (vectors + metadata)
                              ↓
                        ✅ Ready to Query
```

**Current Workaround** - Use API to trigger processing:
```bash
# Once implementation is complete, you can trigger indexing via API
curl -X POST http://localhost:8088/api/v1/process/directory \
  -H "Content-Type: application/json" \
  -d '{"directory": "./data/diagrams", "force_reprocess": false}'
```

# View logs
docker-compose logs -f

# View specific service logs
docker-compose logs -f orchestrator
```

### 4. Verify Services are Running

```bash
# Check all containers
docker-compose ps

# Should show:
# - repograph-redis  
# - repograph-document-scanner
# - repograph-content-extractor
# - repograph-vision-service
# - repograph-summarization-service
# - repograph-embedding-service
# - repograph-vector-store
# - repograph-query-service
# - repograph-orchestrator

# Test health endpoints
curl http://localhost:8088/health  # Orchestrator
curl http://localhost:8081/health  # Document Scanner
curl http://localhost:8082/health  # Content Extractor
```

### 5. Build CLI (Optional - for local development)

```bash
# Build the CLI tool
go build -o bin/repograph-cli cmd/repograph-cli/main.go

# Or use Make
make build-cli
```

### 6. Document Indexing

**📋 Implementation Status**: The indexing workflow is **designed but not yet implemented**.

**Architecture is Ready**:
- ✅ All microservices are running and accessible
- ✅ API endpoints defined (ready for implementation)
- ✅ Configuration system complete
- ✅ Pinecone integration planned

**Planned Indexing Workflow**:

When implemented, the system will:

1. **🔍 Scan**: Document Scanner service reads all files from `DATA_DIRECTORY`
2. **📄 Extract**: Content Extractor service processes each file type
3. **👁️ Analyze**: Vision Service analyzes images/diagrams (Google Vision API)
4. **📝 Summarize**: Summarization Service generates summaries (Azure OpenAI GPT-4)
5. **🧮 Embed**: Embedding Service creates 1536-dimension vectors (Azure OpenAI)
6. **💾 Store**: Vector Store service saves to Pinecone with metadata
7. **⚡ Deduplicate**: Skip files already indexed (hash-based)

**Data Flow**:
```
./data/diagrams/*.pdf, *.png, *.xlsx...
         ↓
  [Document Scanner: 8081]
         ↓
  [Content Extractor: 8082] → [Vision Service: 8083]
         ↓                            ↓
     Text Content              Image Analysis
         ↓                            ↓
         └────────────┬───────────────┘
                      ↓
         [Summarization: 8084] (Azure OpenAI GPT-4)
                      ↓
            "This document describes..."
                      ↓
         [Embedding: 8085] (Azure OpenAI text-embedding-ada-002)
                      ↓
         [0.123, -0.456, 0.789, ...] (1536 dims)
                      ↓
         [Vector Store: 8086]
                      ↓
         🌐 PINECONE VECTOR DATABASE
              {
                id: "doc-123-chunk-0",
                values: [embedding],
                metadata: {
                  file_name: "architecture.pdf",
                  content: "text...",
                  summary: "...",
                  file_hash: "sha256:..."
                }
              }
```

**Current Implementation Tasks** (See TODO in code):

```bash
# Check orchestrator implementation status
grep -n "TODO" cmd/orchestrator/main.go

# You'll see placeholders for:
# - Document processing workflow
# - Service integration
# - Error handling
```

**To Implement**:
1. Connect orchestrator to document-scanner service
2. Implement file-by-file processing loop
3. Call content-extractor for each file
4. Call vision-service for images
5. Call summarization-service
6. Call embedding-service
7. Call vector-store to save in Pinecone
8. Add deduplication logic (check file hash in Pinecone metadata)

### 7. Query the System

```bash
# Using CLI
./bin/repograph-cli query ask "What is in the documents?"

# Or via API
curl -X POST http://localhost:8087/api/v1/query \
  -H "Content-Type: application/json" \
  -d '{"text": "What is in the documents?", "top_k": 5}'
```

## Common Commands

### Docker Compose Management

```bash
# Start services
docker-compose up -d

# Stop services
docker-compose down

# Rebuild a specific service
docker-compose up -d --build orchestrator

# View logs (all services)
docker-compose logs -f

# View logs (specific service)
docker-compose logs -f orchestrator

# Restart a service
docker-compose restart orchestrator

# Stop a specific service
docker-compose stop orchestrator

# Remove all containers and volumes
docker-compose down -v

# Scale a service
docker-compose up -d --scale embedding-service=3
```

### Local Development (without Docker)

```bash
# 1. Start infrastructure
docker-compose up -d redis

# 2. Build all services
make build

# 3. Start services manually in separate terminals

# Terminal 1
./bin/document-scanner

# Terminal 2
./bin/content-extractor

# Terminal 3
./bin/vision-service

# Terminal 4
./bin/summarization-service

# Terminal 5
./bin/embedding-service

# Terminal 6
./bin/vector-store

# Terminal 7
./bin/query-service

# Terminal 8
./bin/orchestrator
```

## Troubleshooting

### Issue: "no configuration file provided: not found"

**Solution**: Make sure you're in the project root directory where `docker-compose.yml` exists:
```bash
cd /home/nadeeshame/go/repograph_platform
docker-compose up -d
```

### Issue: Services won't start

**Check**: 
1. .env file exists and has valid API keys
2. Ports 8081-8088, 6379 are not in use
3. Docker daemon is running

```bash
# Check if ports are available
lsof -i :8088
lsof -i :6379

# Check Docker
docker ps
systemctl status docker  # or: sudo service docker status
```

### Issue: Build fails

**Solution**: Update dependencies
```bash
go mod tidy
go mod download
```

### Issue: Can't connect to services

**Check**:
1. Services are healthy: `docker-compose ps`
2. Logs for errors: `docker-compose logs -f`
3. Network connectivity: `docker network ls`

### Issue: Permission denied for credentials

**Solution**: 
```bash
# Make sure credentials directory exists and has correct permissions
chmod 755 credentials
chmod 644 credentials/*.json
```

## Next Steps

1. ✅ Services running? → Read [ARCHITECTURE.md](docs/ARCHITECTURE.md)
2. ✅ Want to develop? → Read [DEVELOPMENT.md](docs/DEVELOPMENT.md)
3. ✅ Ready to deploy? → Read [DEPLOYMENT.md](docs/DEPLOYMENT.md)
4. ✅ Need API details? → Read [API_REFERENCE.md](docs/API_REFERENCE.md)

## Health Check Dashboard

Create a simple script to check all services:

```bash
#!/bin/bash
# Save as check-health.sh

echo "Checking RepoGraph Platform Services..."
echo "======================================"

services=(
  "8088:Orchestrator"
  "8081:Document Scanner"
  "8082:Content Extractor"
  "8083:Vision Service"
  "8084:Summarization"
  "8085:Embedding"
  "8086:Vector Store"
  "8087:Query Service"
)

for service in "${services[@]}"; do
  port=$(echo $service | cut -d: -f1)
  name=$(echo $service | cut -d: -f2)
  
  if curl -s http://localhost:$port/health > /dev/null; then
    echo "✅ $name (port $port): Healthy"
  else
    echo "❌ $name (port $port): Not responding"
  fi
done
```

Run it:
```bash
chmod +x check-health.sh
./check-health.sh
```

## Useful Aliases

Add to your `~/.bashrc` or `~/.zshrc`:

```bash
# RepoGraph Platform shortcuts
alias rg-up='docker-compose up -d'
alias rg-down='docker-compose down'
alias rg-logs='docker-compose logs -f'
alias rg-ps='docker-compose ps'
alias rg-build='docker-compose up -d --build'
alias rg-health='curl http://localhost:8088/health | jq'
```

Then reload: `source ~/.bashrc`

---

**Need help?** Check the documentation or logs:
- Documentation: `./docs/`
- Logs: `docker-compose logs -f`
- Status: `docker-compose ps`

Happy building! 🚀
