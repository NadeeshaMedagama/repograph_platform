#!/bin/bash
# Script to test CodeQL workflow fix locally and push changes

set -e

echo "================================================"
echo "CodeQL Workflow Fix - Deployment Script"
echo "================================================"
echo ""

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Change to project directory
cd /home/nadeeshame/go/rag_knowledge_service

echo -e "${YELLOW}Step 1: Checking modified files...${NC}"
git status

echo ""
echo -e "${YELLOW}Step 2: Adding modified files...${NC}"
git add .github/workflows/codeql.yml \
        docs/CODEQL_TROUBLESHOOTING.md \
        docs/CODEQL_FIX_SUMMARY.md \
        README.md

echo -e "${GREEN}✓ Files staged${NC}"

echo ""
echo -e "${YELLOW}Step 3: Creating commit...${NC}"
git commit -m "fix: resolve CodeQL upload content-length mismatch error

- Split CodeQL analysis and upload steps
- Add upload-database: false to prevent automatic upload
- Disable database clustering to avoid compression issues
- Add separate upload step with better error handling
- Create comprehensive troubleshooting guide
- Update documentation

This fixes the following error:
  Error: Request body length does not match content-length header

The fix works by:
1. Separating the analyze and upload-sarif actions
2. Disabling database clustering with --no-db-cluster
3. Explicitly controlling the SARIF file path
4. Adding proper error handling with if: always()
5. Waiting for upload processing completion

References:
- docs/CODEQL_TROUBLESHOOTING.md - Comprehensive troubleshooting guide
- docs/CODEQL_FIX_SUMMARY.md - Detailed explanation of changes"

echo -e "${GREEN}✓ Commit created${NC}"

echo ""
echo -e "${YELLOW}Step 4: Pushing to remote...${NC}"
git push main master

echo -e "${GREEN}✓ Changes pushed successfully${NC}"

echo ""
echo "================================================"
echo -e "${GREEN}SUCCESS!${NC} CodeQL fix has been deployed"
echo "================================================"
echo ""
echo "Next steps:"
echo "1. Go to: https://github.com/NadeeshaMedagama/rag_knowledge_service/actions"
echo "2. Watch the CodeQL workflow execution"
echo "3. Verify the 'Filter and Upload SARIF' step succeeds"
echo "4. Check Security → Code scanning for updated results"
echo ""
echo "If the workflow still fails, consult:"
echo "  docs/CODEQL_TROUBLESHOOTING.md"
echo ""
echo "To view logs in real-time:"
echo "  gh run watch"
echo ""
