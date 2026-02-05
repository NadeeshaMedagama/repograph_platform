# Deployment Guide

## Table of Contents
- [Prerequisites](#prerequisites)
- [Local Development](#local-development)
- [Docker Deployment](#docker-deployment)
- [Kubernetes Deployment](#kubernetes-deployment)
- [Cloud Deployment](#cloud-deployment)
- [Configuration](#configuration)
- [Monitoring](#monitoring)
- [Troubleshooting](#troubleshooting)

---

## Prerequisites

### Required
- Go 1.21 or higher
- Docker 20.10+
- Docker Compose 2.0+
- PostgreSQL 15+
- Redis 7+

### For Production
- Kubernetes 1.25+
- Helm 3.0+
- kubectl configured

### External Services
- Azure OpenAI account with:
  - GPT-4 deployment
  - text-embedding-ada-002 deployment
- Pinecone account and API key
- Google Cloud Platform account (for Vision API)

---

## Local Development

### 1. Clone and Setup

```bash
cd /home/nadeeshame/go/rag-knowledge-service

# Copy environment file
cp .env.example .env

# Edit with your credentials
nano .env
```

### 2. Install Dependencies

```bash
# Download Go dependencies
go mod download

# Install development tools
make install-tools
```

### 3. Start Infrastructure

```bash
# Start PostgreSQL and Redis
docker run -d --name postgres \
  -e POSTGRES_PASSWORD=repograph \
  -e POSTGRES_DB=repograph_db \
  -p 5432:5432 \
  postgres:15-alpine

docker run -d --name redis \
  -p 6379:6379 \
  redis:7-alpine
```

### 4. Build Services

```bash
# Build all services
make build

# Or build individually
make build-orchestrator
make build-cli
```

### 5. Run Services

Open separate terminals for each service:

```bash
# Terminal 1: Document Scanner
./bin/document-scanner

# Terminal 2: Content Extractor
./bin/content-extractor

# Terminal 3: Vision Service
./bin/vision-service

# Terminal 4: Summarization Service
./bin/summarization-service

# Terminal 5: Embedding Service
./bin/embedding-service

# Terminal 6: Vector Store
./bin/vector-store

# Terminal 7: Query Service
./bin/query-service

# Terminal 8: Orchestrator
./bin/orchestrator
```

### 6. Test Installation

```bash
# Check health
./bin/repograph-cli health

# Index sample documents
./bin/repograph-cli index --directory ./data/diagrams

# Query
./bin/repograph-cli query ask "What is Choreo?"
```

---

## Docker Deployment

### Using Docker Compose (Recommended for Development)

```bash
# Build all images
docker-compose -f deployments/docker/docker-compose.yml build

# Start all services
docker-compose -f deployments/docker/docker-compose.yml up -d

# View logs
docker-compose -f deployments/docker/docker-compose.yml logs -f

# Check status
docker-compose -f deployments/docker/docker-compose.yml ps

# Stop services
docker-compose -f deployments/docker/docker-compose.yml down
```

### Individual Service Deployment

```bash
# Build image
docker build -f deployments/docker/Dockerfile.orchestrator -t repograph-orchestrator:latest .

# Run container
docker run -d \
  --name orchestrator \
  -p 8088:8088 \
  --env-file .env \
  repograph-orchestrator:latest
```

### Docker Compose Production

```bash
# Production compose file
docker-compose -f deployments/docker/docker-compose.yml \
  -f deployments/docker/docker-compose.prod.yml up -d
```

---

## Kubernetes Deployment

### 1. Create Namespace

```bash
kubectl create namespace repograph
```

### 2. Create Secrets

```bash
# Azure OpenAI credentials
kubectl create secret generic azure-credentials \
  --from-literal=api-key=$AZURE_OPENAI_API_KEY \
  --from-literal=endpoint=$AZURE_OPENAI_ENDPOINT \
  -n repograph

# Pinecone credentials
kubectl create secret generic pinecone-credentials \
  --from-literal=api-key=$PINECONE_API_KEY \
  -n repograph

# Google Vision credentials
kubectl create secret generic google-credentials \
  --from-file=service-account.json=credentials/service-account.json \
  -n repograph
```

### 3. Deploy PostgreSQL

```bash
kubectl apply -f deployments/kubernetes/postgres-statefulset.yaml
kubectl apply -f deployments/kubernetes/postgres-service.yaml
```

### 4. Deploy Redis

```bash
kubectl apply -f deployments/kubernetes/redis-deployment.yaml
kubectl apply -f deployments/kubernetes/redis-service.yaml
```

### 5. Deploy Microservices

```bash
# Deploy all services
kubectl apply -f deployments/kubernetes/deployments/

# Deploy services
kubectl apply -f deployments/kubernetes/services/

# Deploy ingress
kubectl apply -f deployments/kubernetes/ingress/
```

### 6. Verify Deployment

```bash
# Check pods
kubectl get pods -n repograph

# Check services
kubectl get svc -n repograph

# Check logs
kubectl logs -f deployment/orchestrator -n repograph
```

### Using Helm (Recommended for Production)

```bash
# Add repo (if published)
helm repo add repograph https://charts.repograph.ai

# Install
helm install repograph repograph/repograph \
  --namespace repograph \
  --create-namespace \
  --set azure.apiKey=$AZURE_OPENAI_API_KEY \
  --set pinecone.apiKey=$PINECONE_API_KEY

# Upgrade
helm upgrade repograph repograph/repograph -n repograph

# Uninstall
helm uninstall repograph -n repograph
```

---

## Cloud Deployment

### Google Cloud Run

#### 1. Build and Push Images

```bash
# Set project
gcloud config set project YOUR_PROJECT_ID

# Enable services
gcloud services enable run.googleapis.com
gcloud services enable containerregistry.googleapis.com

# Build and push
for service in orchestrator document-scanner content-extractor vision-service summarization-service embedding-service vector-store query-service; do
  docker build -f deployments/docker/Dockerfile.$service -t gcr.io/YOUR_PROJECT_ID/repograph-$service:latest .
  docker push gcr.io/YOUR_PROJECT_ID/repograph-$service:latest
done
```

#### 2. Deploy Services

```bash
# Deploy orchestrator
gcloud run deploy repograph-orchestrator \
  --image gcr.io/YOUR_PROJECT_ID/repograph-orchestrator:latest \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated \
  --set-env-vars AZURE_OPENAI_API_KEY=$AZURE_OPENAI_API_KEY \
  --set-env-vars PINECONE_API_KEY=$PINECONE_API_KEY

# Repeat for other services
```

### AWS ECS

```bash
# Create ECR repositories
aws ecr create-repository --repository-name repograph-orchestrator

# Build and push
docker build -f deployments/docker/Dockerfile.orchestrator -t repograph-orchestrator:latest .
docker tag repograph-orchestrator:latest $AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com/repograph-orchestrator:latest
docker push $AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com/repograph-orchestrator:latest

# Create task definition and service using AWS Console or CLI
```

### Azure Container Instances

```bash
# Create resource group
az group create --name repograph-rg --location eastus

# Create container instance
az container create \
  --resource-group repograph-rg \
  --name repograph-orchestrator \
  --image your-registry.azurecr.io/repograph-orchestrator:latest \
  --dns-name-label repograph-orchestrator \
  --ports 8088 \
  --environment-variables \
    AZURE_OPENAI_API_KEY=$AZURE_OPENAI_API_KEY \
    PINECONE_API_KEY=$PINECONE_API_KEY
```

---

## Configuration

### Environment Variables

Create `.env` file:

```bash
# Azure OpenAI
AZURE_OPENAI_API_KEY=your_key
AZURE_OPENAI_ENDPOINT=https://your-resource.openai.azure.com/
AZURE_OPENAI_EMBEDDINGS_DEPLOYMENT=text-embedding-ada-002
AZURE_OPENAI_CHAT_DEPLOYMENT=gpt-4
AZURE_OPENAI_API_VERSION=2024-02-01

# Pinecone
PINECONE_API_KEY=your_key
PINECONE_INDEX_NAME=repograph-ai-index
PINECONE_DIMENSION=1536
PINECONE_CLOUD=aws
PINECONE_REGION=us-east-1

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

# Application
DATA_DIRECTORY=./data/diagrams
LOG_LEVEL=info
CHUNK_SIZE=1000
CHUNK_OVERLAP=200
```

### Configuration File

Create `configs/config.yaml`:

```yaml
server:
  port: 8080
  read_timeout: 30s
  write_timeout: 30s

app:
  data_directory: ./data/diagrams
  log_level: info
  chunk_size: 1000
  chunk_overlap: 200
  skip_existing_documents: true

services:
  document_scanner_url: http://document-scanner:8081
  content_extractor_url: http://content-extractor:8082
  vision_service_url: http://vision-service:8083
  summarization_service_url: http://summarization-service:8084
  embedding_service_url: http://embedding-service:8085
  vector_store_service_url: http://vector-store:8086
  query_service_url: http://query-service:8087
  orchestrator_service_url: http://orchestrator:8088
```

---

## Monitoring

### Health Checks

```bash
# Check all services
for port in 8081 8082 8083 8084 8085 8086 8087 8088; do
  echo "Checking port $port..."
  curl -s http://localhost:$port/health | jq
done
```

### Prometheus Metrics

Add to each service:

```go
import "github.com/prometheus/client_golang/prometheus/promhttp"

http.Handle("/metrics", promhttp.Handler())
```

Access metrics: `http://localhost:8088/metrics`

### Logging

Logs are in JSON format:

```bash
# View logs
docker-compose logs -f orchestrator

# Filter by level
docker-compose logs orchestrator | grep '"level":"error"'
```

### Grafana Dashboard

Import dashboard from `deployments/monitoring/grafana-dashboard.json`

---

## Troubleshooting

### Service Won't Start

```bash
# Check logs
docker logs repograph-orchestrator

# Check configuration
./bin/repograph-cli health

# Verify connectivity
curl http://localhost:8088/health
```

### Database Connection Issues

```bash
# Test PostgreSQL connection
psql -h localhost -U repograph -d repograph_db

# Check if running
docker ps | grep postgres

# Restart
docker restart postgres
```

### Redis Connection Issues

```bash
# Test Redis
redis-cli ping

# Check if running
docker ps | grep redis

# Restart
docker restart redis
```

### External API Errors

```bash
# Test Azure OpenAI
curl -X POST "$AZURE_OPENAI_ENDPOINT/openai/deployments/gpt-4/chat/completions?api-version=2024-02-01" \
  -H "api-key: $AZURE_OPENAI_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"messages":[{"role":"user","content":"test"}]}'

# Test Pinecone
curl -X GET "https://api.pinecone.io/indexes" \
  -H "Api-Key: $PINECONE_API_KEY"
```

### Performance Issues

```bash
# Check resource usage
docker stats

# Scale embedding service
docker-compose up -d --scale embedding-service=3

# Check Kubernetes resources
kubectl top pods -n repograph
```

---

## Backup and Recovery

### Database Backup

```bash
# Backup PostgreSQL
docker exec postgres pg_dump -U repograph repograph_db > backup.sql

# Restore
docker exec -i postgres psql -U repograph repograph_db < backup.sql
```

### Configuration Backup

```bash
# Backup all configs
tar -czf config-backup.tar.gz .env configs/ credentials/
```

---

## Security Best Practices

1. **Never commit credentials** to git
2. **Use secrets management** (HashiCorp Vault, AWS Secrets Manager)
3. **Enable TLS** for all services
4. **Implement authentication** for production
5. **Regular security audits**
6. **Update dependencies** regularly
7. **Use read-only containers** where possible
8. **Network segmentation** in Kubernetes
9. **Rate limiting** on all endpoints
10. **Audit logging** enabled

---

*Last Updated: February 2026*
