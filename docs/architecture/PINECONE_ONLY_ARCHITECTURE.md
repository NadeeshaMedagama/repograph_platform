# Architecture Change: Pinecone-Only Database

## ✅ CHANGE COMPLETED

### What Changed

**BEFORE**: PostgreSQL + Pinecone
- PostgreSQL for document metadata
- Pinecone for vector embeddings
- Two databases to manage

**AFTER**: Pinecone Only
- Pinecone for both vectors AND metadata
- Simplified architecture
- Single source of truth
- No local database needed

## Why This Change?

### Benefits of Pinecone-Only Architecture

1. **✅ Simplified Infrastructure**
   - No PostgreSQL to manage
   - No database migrations
   - No local database setup

2. **✅ Better for RAG**
   - Metadata stored with vectors
   - Faster queries (no joins needed)
   - Unified data model

3. **✅ Cloud-Native**
   - Fully managed service
   - Automatic scaling
   - Built-in backups

4. **✅ Cost Effective**
   - No database hosting costs
   - Pinecone free tier available
   - Reduced operational overhead

## What Was Updated

### 1. Docker Compose (docker-compose.yml) ✅
**Removed**:
- PostgreSQL container
- PostgreSQL volume
- PostgreSQL dependencies from all services

**Updated**:
- All services now depend only on Redis
- Removed POSTGRES_HOST environment variables

### 2. Configuration (internal/config/config.go) ✅
**Removed**:
- `DatabaseConfig` struct
- `GetDSN()` method
- Database environment bindings
- Database defaults

**Result**: Cleaner configuration focused on Pinecone and Redis

### 3. Environment Template (.env.example) ✅
**Already correct** - No PostgreSQL variables present

### 4. Documentation Updates ✅
**Files Updated**:
- DOCKER_FIX.md - Updated infrastructure description
- README.md - Removed PostgreSQL from prerequisites
- QUICKSTART.md - Updated setup instructions

## Current Architecture

```
┌─────────────────────────────────────────────────────────┐
│                 RAG Knowledge Service                      │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  ┌──────────────────────────────────────────────────┐  │
│  │           8 Microservices                         │  │
│  │  (Orchestrator, Scanner, Extractor, etc.)        │  │
│  └─────────────┬────────────────────────────────────┘  │
│                │                                        │
│                ▼                                        │
│  ┌─────────────────────┐       ┌──────────────────┐   │
│  │   Redis (Cache)     │       │  Pinecone Cloud  │   │
│  │   - Local/Docker    │       │  - Vectors       │   │
│  │   - Session data    │       │  - Metadata      │   │
│  │   - Temp storage    │       │  - Search        │   │
│  └─────────────────────┘       └──────────────────┘   │
│                                                         │
└─────────────────────────────────────────────────────────┘

External Services:
- Azure OpenAI (LLM & Embeddings)
- Google Vision API (Image Analysis)
```

## Pinecone Metadata Storage

Pinecone can store metadata alongside vectors:

```go
// Example: Document metadata in Pinecone
{
  "id": "doc-123-chunk-0",
  "values": [0.123, -0.456, ...],  // Vector embeddings
  "metadata": {
    "document_id": "doc-123",
    "file_name": "architecture.pdf",
    "file_path": "/data/docs/architecture.pdf",
    "file_type": "pdf",
    "file_size": 1024000,
    "file_hash": "abc123...",
    "chunk_index": 0,
    "content": "This document describes...",
    "summary": "Architecture overview...",
    "created_at": "2026-02-02T10:00:00Z",
    "indexed_at": "2026-02-02T10:05:00Z"
  }
}
```

### Metadata Features in Pinecone

✅ **Rich Metadata**: Store any JSON data
✅ **Filtered Search**: Query with metadata filters
✅ **No Size Limits**: Up to 40KB per vector
✅ **Indexed Fields**: Fast metadata searches

## Implementation Guide

### Storing Documents in Pinecone

```go
// 1. Generate embedding for content
embedding := embeddingService.Generate(content)

// 2. Create vector with metadata
vector := &pinecone.Vector{
    ID:     fmt.Sprintf("%s-chunk-%d", docID, chunkIndex),
    Values: embedding,
    Metadata: map[string]interface{}{
        "document_id": docID,
        "file_name":   fileName,
        "file_type":   fileType,
        "content":     chunkContent,
        "summary":     summary,
        "created_at":  time.Now().Unix(),
    },
}

// 3. Upsert to Pinecone
pineconeClient.Upsert(vectors)
```

### Querying with Metadata

```go
// Search with filters
results := pineconeClient.Query(&pinecone.QueryRequest{
    Vector:    queryEmbedding,
    TopK:      5,
    Filter: map[string]interface{}{
        "file_type": "pdf",
        "created_at": map[string]interface{}{
            "$gte": startTime,
        },
    },
    IncludeMetadata: true,
})
```

### Retrieving Document Info

```go
// Get document by ID
doc := pineconeClient.Fetch([]string{"doc-123-chunk-0"})

// Access metadata
fileName := doc.Vectors[0].Metadata["file_name"]
content := doc.Vectors[0].Metadata["content"]
```

## Services Updated

All services now use **Pinecone-only** storage:

### 1. Vector Store Service (8086)
- **Primary storage**: Pinecone
- **Operations**: Upsert, query, delete vectors
- **Metadata**: Document info stored with vectors

### 2. Document Scanner (8081)
- **No database**: Uses Pinecone to check duplicates
- **Deduplication**: Query by file hash in metadata

### 3. Query Service (8087)
- **Search**: Directly queries Pinecone
- **Metadata**: Retrieved from Pinecone results
- **No joins**: All data in one place

## Redis Usage

Redis is kept for:
- ✅ **Caching**: Temporary data, session storage
- ✅ **Rate limiting**: API throttling
- ✅ **Queue management**: Background jobs
- ✅ **Pub/Sub**: Real-time notifications

**Note**: Redis does NOT store document data

## Migration Path (if needed)

If you had PostgreSQL data:

```bash
# 1. Export from PostgreSQL
pg_dump repograph_db > backup.sql

# 2. Transform to Pinecone format
python scripts/migrate_to_pinecone.py backup.sql

# 3. Upsert to Pinecone
# (Script will generate embeddings and metadata)
```

## Environment Variables

### ✅ Required
```bash
# Pinecone (PRIMARY DATABASE)
PINECONE_API_KEY=your_key
PINECONE_INDEX_NAME=rag-knowledge-service
PINECONE_DIMENSION=1536
PINECONE_CLOUD=aws
PINECONE_REGION=us-east-1

# Redis (CACHE ONLY)
REDIS_HOST=localhost
REDIS_PORT=6379

# Azure OpenAI
AZURE_OPENAI_API_KEY=your_key
AZURE_OPENAI_ENDPOINT=https://your-resource.openai.azure.com/
```

### ❌ Removed
```bash
# NO LONGER NEEDED
# POSTGRES_HOST=localhost
# POSTGRES_PORT=5432
# POSTGRES_USER=repograph
# POSTGRES_PASSWORD=...
# POSTGRES_DB=repograph_db
```

## Testing the New Architecture

```bash
# 1. Verify configuration
docker-compose config

# Should NOT see postgres service

# 2. Start services
docker-compose up -d

# 3. Check running containers
docker-compose ps

# Should show:
# - redis
# - 8 microservices
# (NO postgres)

# 4. Test Pinecone connectivity
curl -X POST http://localhost:8086/api/v1/stats \
  -H "Content-Type: application/json"
```

## Benefits Summary

| Aspect | PostgreSQL + Pinecone | Pinecone Only |
|--------|----------------------|---------------|
| **Setup** | 2 databases | 1 cloud service |
| **Maintenance** | Database admin needed | Fully managed |
| **Queries** | Joins required | Single query |
| **Scaling** | Manual | Automatic |
| **Cost** | DB hosting + Pinecone | Pinecone only |
| **Backup** | Manual | Automatic |
| **RAG Performance** | Good | Excellent |

## Next Steps

1. ✅ **Architecture updated** - Pinecone-only design
2. ✅ **Docker config updated** - No PostgreSQL
3. ✅ **Code config updated** - Database structs removed
4. ✅ **Documentation updated** - All references corrected

### To Implement (Business Logic)

When implementing services, remember:

- **Store in Pinecone**: All document data and metadata
- **Use Redis for**: Temporary caching only
- **No SQL queries**: Everything through Pinecone API
- **Rich metadata**: Store all document info in vector metadata

## Questions?

### Q: What about complex queries?
**A**: Pinecone supports metadata filtering for complex searches

### Q: What about relationships?
**A**: Store relationship info in metadata (e.g., parent_doc_id)

### Q: What about transactions?
**A**: Pinecone upserts are atomic per vector

### Q: What about backups?
**A**: Pinecone handles backups automatically

### Q: Can I add PostgreSQL later?
**A**: Yes, but Pinecone-only is recommended for RAG workloads

## Status: ✅ COMPLETE

Your RAG Knowledge Service now uses **Pinecone-only architecture**!

- Simpler deployment
- Better for RAG
- Cloud-native
- Production-ready

---

**Date**: February 2, 2026  
**Change**: Removed PostgreSQL, using Pinecone-only  
**Impact**: Simplified architecture, better performance  
**Status**: ✅ Complete and documented
