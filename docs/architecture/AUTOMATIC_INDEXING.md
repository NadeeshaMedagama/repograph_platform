# Automatic Document Indexing

## Overview

RAG Knowledge Service automatically indexes all documents when you start the services with `docker-compose up -d`. No manual steps required!

## How It Works

### 1. Startup Process

```bash
# Start services
docker-compose up -d
```

When the orchestrator service starts, it automatically:

1. **Waits for services** - Ensures all microservices are healthy (5 seconds)
2. **Scans directory** - Reads all files from `DATA_DIRECTORY`
3. **Processes documents** - Extracts, analyzes, embeds, and stores
4. **Runs in background** - Doesn't block the API server
5. **Deduplicates** - Skips already-indexed documents

### 2. Indexing Pipeline

```
┌──────────────────────────────────────────────────────────────────┐
│              Automatic Indexing on Startup                       │
├──────────────────────────────────────────────────────────────────┤
│                                                                  │
│  docker-compose up -d                                            │
│         │                                                        │
│         ▼                                                        │
│  Orchestrator Starts                                             │
│         │                                                        │
│         ├──▶ HTTP Server (port 8088) ────▶ Ready for API calls  │
│         │                                                        │
│         └──▶ Background Indexer                                 │
│                     │                                            │
│                     ▼                                            │
│              Wait 5 seconds (services ready)                     │
│                     │                                            │
│                     ▼                                            │
│              Scan DATA_DIRECTORY                                 │
│                     │                                            │
│                     ▼                                            │
│         ┌───────────┴───────────┐                               │
│         │                       │                               │
│    For Each File           Check if indexed                     │
│         │                   (by file hash)                      │
│         │                       │                               │
│         ├──────┬────────────────┤                               │
│         │      │                │                               │
│    Not Indexed │           Already Indexed                      │
│         │      │                │                               │
│         ▼      │                └──▶ Skip                        │
│                │                                                 │
│    1. Extract Content (via Content Extractor)                   │
│    2. Analyze Images (via Vision Service)                       │
│    3. Generate Summary (via Summarization Service)              │
│    4. Create Embeddings (via Embedding Service)                 │
│    5. Store in Pinecone (via Vector Store)                      │
│                │                                                 │
│                ▼                                                 │
│         Log Progress                                             │
│                │                                                 │
│                ▼                                                 │
│         Next File                                                │
│                                                                  │
└──────────────────────────────────────────────────────────────────┘
```

### 3. What Gets Indexed

**All file types** in `DATA_DIRECTORY`:
- 📄 Documents (PDF, DOCX, PPTX, TXT, MD)
- 🖼️ Images (PNG, JPG, SVG)
- 📊 Diagrams (DrawIO, Excalidraw)
- 📈 Spreadsheets (XLSX, CSV)
- 💻 Code files (Go, Python, JS, etc.)
- 🔧 Structured data (JSON, YAML, XML)

### 4. Deduplication

Each file gets a **SHA-256 hash** of its content:
- ✅ **First time**: File is processed and indexed
- ⚡ **Already indexed**: File is skipped (saves time & API costs)
- 🔄 **File changed**: New hash triggers re-indexing

## Configuration

### Environment Variables

```bash
# Data directory to index
DATA_DIRECTORY=./data/diagrams

# Skip already-indexed documents (recommended)
SKIP_EXISTING_DOCUMENTS=true

# Logging level (see indexing progress)
LOG_LEVEL=info
```

### Custom Directory

Change the directory in `.env`:

```bash
# Index different directory
DATA_DIRECTORY=./my-documents

# Or multiple directories (comma-separated)
DATA_DIRECTORY=./docs,./diagrams,./reports
```

## Monitoring Indexing Progress

### View Logs

```bash
# Watch orchestrator logs in real-time
docker-compose logs -f orchestrator

# You'll see:
# - "Starting automatic document indexing"
# - "Processing file: document.pdf"
# - "Extracted 1500 words"
# - "Generated embedding (1536 dimensions)"
# - "Stored in Pinecone: doc-abc123"
# - "Skipping already-indexed: diagram.png"
```

### Check Progress via API

```bash
# Get indexing statistics
curl http://localhost:8088/api/v1/status

# Response:
{
  "indexing": {
    "status": "running",
    "total_files": 150,
    "processed": 45,
    "skipped": 100,
    "failed": 5,
    "progress_percent": 96
  }
}
```

## Manual Control

### Disable Automatic Indexing

Set `DATA_DIRECTORY` to empty:

```bash
# In .env
DATA_DIRECTORY=

# Or start orchestrator without indexing
docker-compose run --rm orchestrator --no-auto-index
```

### Trigger Re-indexing

```bash
# Force reindex all documents (ignores deduplication)
curl -X POST http://localhost:8088/api/v1/process/directory \
  -H "Content-Type: application/json" \
  -d '{
    "directory": "./data/diagrams",
    "force_reprocess": true
  }'

# Or use CLI
./bin/repograph-cli index --force
```

### Index Specific Files

```bash
# Index a single file
curl -X POST http://localhost:8088/api/v1/process/document \
  -H "Content-Type: application/json" \
  -d '{
    "file_path": "./data/diagrams/architecture.pdf"
  }'
```

## Performance

### Indexing Speed

| File Type | Avg Time | Notes |
|-----------|----------|-------|
| Text (TXT, MD) | 1-2 sec | Fast extraction |
| Documents (PDF, DOCX) | 3-5 sec | OCR if needed |
| Images (PNG, JPG) | 5-8 sec | Vision API call |
| Diagrams (DrawIO) | 2-4 sec | XML parsing + Vision |
| Code Files | 1-2 sec | Syntax highlighting |
| Spreadsheets | 3-6 sec | Multiple sheets |

**Typical performance**: 
- 100 documents in ~5-10 minutes
- 1000 documents in ~1-2 hours
- Parallel processing (can be scaled)

### Optimization Tips

1. **Skip existing documents**: Set `SKIP_EXISTING_DOCUMENTS=true`
2. **Filter file types**: Only index relevant files
3. **Scale services**: Increase replicas for embedding-service
4. **Batch processing**: Group small files together

```bash
# Scale embedding service for faster processing
docker-compose up -d --scale embedding-service=3
```

## Data Storage

### Pinecone Structure

Each document chunk is stored as:

```json
{
  "id": "doc-abc123-chunk-0",
  "values": [0.123, -0.456, ...],  // 1536-dimension embedding
  "metadata": {
    "document_id": "doc-abc123",
    "file_name": "architecture.pdf",
    "file_path": "./data/diagrams/architecture.pdf",
    "file_type": "pdf",
    "file_size": 1024000,
    "file_hash": "sha256:abc123...",
    "chunk_index": 0,
    "chunk_total": 10,
    "content": "This document describes the system architecture...",
    "summary": "Architecture overview with microservices design...",
    "indexed_at": "2026-02-02T10:05:00Z",
    "created_at": "2026-02-01T15:30:00Z"
  }
}
```

### Metadata Benefits

- **Fast filtering**: Query by file type, date, size
- **No SQL needed**: All data in Pinecone
- **Rich context**: Full document info with vectors
- **Deduplication**: Check by file_hash

## Troubleshooting

### Indexing Not Starting

**Check**:
1. `DATA_DIRECTORY` is set in `.env`
2. Directory exists and has files
3. Orchestrator service is running

```bash
# Verify config
docker-compose exec orchestrator printenv | grep DATA_DIRECTORY

# Check orchestrator logs
docker-compose logs orchestrator | grep "automatic indexing"
```

### Slow Indexing

**Solutions**:
1. Scale embedding service: `docker-compose up -d --scale embedding-service=3`
2. Check API rate limits (Azure OpenAI, Google Vision)
3. Verify network connectivity to external services

### Files Not Being Indexed

**Check**:
1. File format is supported
2. File is not corrupted
3. Sufficient API quota (Azure, Google)
4. Check error logs: `docker-compose logs orchestrator | grep ERROR`

### Out of Memory

**Solutions**:
1. Increase Docker memory limit
2. Process smaller batches
3. Reduce chunk size in configuration

```bash
# In .env
CHUNK_SIZE=500  # Default: 1000
```

## Best Practices

1. **✅ Use deduplication**: Set `SKIP_EXISTING_DOCUMENTS=true`
2. **✅ Monitor logs**: Watch for errors during indexing
3. **✅ Verify Pinecone**: Check index stats after indexing
4. **✅ Test with small dataset**: Verify before large-scale indexing
5. **✅ Backup data**: Keep original documents safe
6. **✅ Set up alerts**: Monitor indexing failures

## Example Workflow

### First Time Setup

```bash
# 1. Configure environment
cp .env.example .env
nano .env  # Add API keys

# 2. Add documents to data directory
cp my-docs/* ./data/diagrams/

# 3. Start services (auto-indexes)
docker-compose up -d

# 4. Monitor progress
docker-compose logs -f orchestrator

# 5. Wait for completion
# "Automatic indexing completed: 150 files processed"

# 6. Query your documents
./bin/repograph-cli query ask "What's in my documents?"
```

### Adding New Documents

```bash
# 1. Add files to data directory
cp new-docs/* ./data/diagrams/

# 2. Restart orchestrator (or wait for next startup)
docker-compose restart orchestrator

# New files will be automatically indexed
# Already-indexed files will be skipped
```

### Regular Updates

```bash
# Option 1: Restart periodically
docker-compose restart orchestrator

# Option 2: Manual trigger
curl -X POST http://localhost:8088/api/v1/process/directory \
  -d '{"directory": "./data/diagrams"}'

# Option 3: Watch directory (future feature)
# Orchestrator will detect new files automatically
```

## Security Considerations

1. **API Keys**: Never commit `.env` file
2. **Data Privacy**: Documents are sent to Azure and Google APIs
3. **Access Control**: Implement authentication for production
4. **Audit Logs**: Track what was indexed and when
5. **Encryption**: Use HTTPS for API calls

## Summary

✅ **Automatic**: No manual steps required  
✅ **Smart**: Deduplicates to save time and costs  
✅ **Fast**: Parallel processing with scalable services  
✅ **Monitored**: Full logging and progress tracking  
✅ **Reliable**: Error handling and retry logic  
✅ **Flexible**: API and CLI for manual control  

**Just run `docker-compose up -d` and your documents are automatically indexed!** 🚀

---

**Last Updated**: February 2, 2026  
**Version**: 1.0.0
