# Manual Steps Required

After all the automated fixes, you need to complete these manual steps:

## 1. 🔧 Update GitHub Repository Name

### Steps:
1. Go to GitHub repository: https://github.com/NadeeshaMedagama/repograph_platform
2. Click **Settings**
3. Scroll to **Repository name**
4. Change from `repograph_platform` to `rag-knowledge-service`
5. Click **Rename**

### After Renaming:
```bash
# Update your local git remote
cd /home/nadeeshame/go/rag_knowledge_service
git remote set-url origin https://github.com/NadeeshaMedagama/rag-knowledge-service.git

# Optional: Rename local directory
cd /home/nadeeshame/go
mv rag_knowledge_service rag-knowledge-service
cd rag-knowledge-service
```

## 2. 🔐 Fix GitHub Actions Permissions

### Problem:
```
Error: denied: installation not allowed to Create organization package
```

### Solution:
1. Go to: **Settings** → **Actions** → **General**
2. Scroll to **Workflow permissions**
3. Select: ✅ **Read and write permissions**
4. ✅ Check: **Allow GitHub Actions to create and approve pull requests**
5. Click **Save**

## 3. 🔑 Verify GitHub Secrets

Ensure these secrets are configured:

### Repository Secrets:
1. Go to: **Settings** → **Secrets and variables** → **Actions**
2. Verify/Add these secrets:
   - ✅ `AZURE_OPENAI_API_KEY`
   - ✅ `PINECONE_API_KEY`
   - ✅ `GOOGLE_VISION_API_KEY` (optional)

### How to Add:
1. Click **New repository secret**
2. Enter **Name** and **Value**
3. Click **Add secret**

## 4. 📝 Update Pinecone Index Name (Optional)

If you want to match the new project name:

### Option A: Create New Index
```bash
# In Pinecone Console, create new index:
# Name: rag-knowledge-service
# Dimensions: 1536
# Metric: cosine
```

### Option B: Keep Existing Index
Just update `.env`:
```bash
PINECONE_INDEX_NAME=repograph-platform  # or your current index name
```

## 5. 🧪 Test the Changes

### Local Testing:
```bash
cd /home/nadeeshame/go/rag-knowledge-service  # after rename

# 1. Build all services
make build

# 2. Run tests
make test

# 3. Start with Docker
docker compose up -d

# 4. Check logs
docker compose logs -f orchestrator

# 5. Test CLI
./bin/rag-cli --help
```

### CI/CD Testing:
1. Push changes to trigger workflow:
```bash
git add .
git commit -m "chore: rename project to rag-knowledge-service"
git push origin master
```

2. Check GitHub Actions:
   - Go to **Actions** tab
   - Verify all workflows pass ✅

## 6. 📦 Update Package Registry (If Using GHCR)

After fixing permissions, images will be pushed to:
```
ghcr.io/nadeeshamedagama/rag-knowledge-service-orchestrator:master
ghcr.io/nadeeshamedagama/rag-knowledge-service-document-scanner:master
ghcr.io/nadeeshamedagama/rag-knowledge-service-content-extractor:master
# ... etc for all services
```

### Make Images Public (Optional):
1. Go to package in GHCR
2. Click **Package settings**
3. Change visibility to **Public**

## 7. 🔄 Update External References

Update any external systems that reference the old name:

### Documentation Sites:
- [ ] Update links in external docs
- [ ] Update README badges
- [ ] Update API documentation URLs

### Cloud Deployments:
- [ ] Update Kubernetes manifests
- [ ] Update Helm charts
- [ ] Update terraform configs
- [ ] Update environment variables

### Monitoring/Logging:
- [ ] Update service names in dashboards
- [ ] Update alert configurations
- [ ] Update log aggregation queries

## 8. 📢 Communicate Changes

### Team Notification:
```markdown
🚀 Project Renamed: repograph_platform → rag-knowledge-service

What changed:
- Repository: github.com/NadeeshaMedagama/rag-knowledge-service
- Go module: github.com/nadeeshame/rag-knowledge-service
- CLI: repograph-cli → rag-cli
- Docker images: ghcr.io/.../rag-knowledge-service-*

Action required:
1. Update your git remote:
   git remote set-url origin https://github.com/NadeeshaMedagama/rag-knowledge-service.git

2. Rebuild local environment:
   go mod download
   make build

3. Update any scripts/configs that reference old name
```

## ✅ Checklist

- [ ] Rename GitHub repository
- [ ] Update git remote locally
- [ ] Fix GitHub Actions permissions
- [ ] Verify GitHub secrets
- [ ] Test local build (`make build`)
- [ ] Test Docker build (`docker compose build`)
- [ ] Push and verify CI/CD pipeline
- [ ] Update external references
- [ ] Notify team
- [ ] Update documentation

## 🎉 Done!

Once these steps are complete, your project will be fully renamed and operational as **rag-knowledge-service**.

---
*Last updated: February 3, 2026*
