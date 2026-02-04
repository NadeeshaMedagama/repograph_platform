#!/bin/bash

# Verification Script for rag-knowledge-service Rename
# This script verifies all changes were applied correctly

set -e

GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "========================================"
echo "RAG Knowledge Service - Verification"
echo "========================================"
echo ""

ERRORS=0
WARNINGS=0

# Function to check and report
check_pass() {
    echo -e "${GREEN}✓${NC} $1"
}

check_fail() {
    echo -e "${RED}✗${NC} $1"
    ERRORS=$((ERRORS + 1))
}

check_warn() {
    echo -e "${YELLOW}⚠${NC} $1"
    WARNINGS=$((WARNINGS + 1))
}

# 1. Check go.mod
echo "1. Checking go.mod..."
if grep -q "module github.com/nadeeshame/rag-knowledge-service" go.mod; then
    check_pass "Module path updated in go.mod"
else
    check_fail "Module path NOT updated in go.mod"
fi

if grep -q "go 1.23" go.mod; then
    check_pass "Go version is 1.23"
else
    check_warn "Go version might not be 1.23"
fi

# 2. Check imports in Go files
echo ""
echo "2. Checking Go file imports..."
OLD_IMPORTS=$(grep -r "github.com/nadeeshame/repograph_platform" --include="*.go" . 2>/dev/null | wc -l)
if [ "$OLD_IMPORTS" -eq 0 ]; then
    check_pass "No old import paths found"
else
    check_fail "Found $OLD_IMPORTS files with old import paths"
fi

NEW_IMPORTS=$(grep -r "github.com/nadeeshame/rag-knowledge-service" --include="*.go" . 2>/dev/null | wc -l)
if [ "$NEW_IMPORTS" -gt 0 ]; then
    check_pass "New import paths found in $NEW_IMPORTS locations"
else
    check_fail "No new import paths found"
fi

# 3. Check CLI directory
echo ""
echo "3. Checking CLI directory..."
if [ -d "cmd/rag-cli" ]; then
    check_pass "CLI directory renamed to rag-cli"
else
    check_fail "CLI directory not found at cmd/rag-cli"
fi

if [ -d "cmd/repograph-cli" ]; then
    check_warn "Old CLI directory still exists at cmd/repograph-cli"
fi

# 4. Check Dockerfiles
echo ""
echo "4. Checking Dockerfiles..."
OLD_GO_VERSION=$(grep -r "FROM golang:1.21" deployments/docker/Dockerfile.* 2>/dev/null | wc -l)
if [ "$OLD_GO_VERSION" -eq 0 ]; then
    check_pass "No Dockerfiles using Go 1.21"
else
    check_fail "Found $OLD_GO_VERSION Dockerfiles still using Go 1.21"
fi

NEW_GO_VERSION=$(grep -r "FROM golang:1.23" deployments/docker/Dockerfile.* 2>/dev/null | wc -l)
if [ "$NEW_GO_VERSION" -gt 0 ]; then
    check_pass "Found $NEW_GO_VERSION Dockerfiles using Go 1.23"
else
    check_fail "No Dockerfiles using Go 1.23"
fi

# 5. Check docker-compose.yml
echo ""
echo "5. Checking docker-compose.yml..."
if grep -q "rag-knowledge-" docker-compose.yml; then
    check_pass "Container names updated in docker-compose.yml"
else
    check_fail "Container names not updated in docker-compose.yml"
fi

if grep -q "repograph-" docker-compose.yml; then
    check_warn "Old container name prefixes still found"
fi

# 6. Check documentation
echo ""
echo "6. Checking documentation..."
if grep -q "RAG Knowledge Service" README.md; then
    check_pass "README.md updated with new name"
else
    check_fail "README.md not updated"
fi

if grep -q "RepoGraph Platform" README.md; then
    check_warn "Old project name still in README.md"
fi

# 7. Try to build
echo ""
echo "7. Testing build..."
if go build ./... 2>/dev/null; then
    check_pass "Project builds successfully"
else
    check_fail "Project build failed"
fi

# 8. Check Makefile
echo ""
echo "8. Checking Makefile..."
if grep -q "rag-cli" Makefile; then
    check_pass "Makefile updated with new CLI name"
else
    check_fail "Makefile not updated"
fi

# 9. Check for old binary names
echo ""
echo "9. Checking binary directory..."
if [ -d "bin" ]; then
    if [ -f "bin/repograph-cli" ]; then
        check_warn "Old CLI binary still exists: bin/repograph-cli"
    fi

    if [ -f "bin/rag-cli" ]; then
        check_pass "New CLI binary exists: bin/rag-cli"
    fi
fi

# 10. Check golangci-lint config
echo ""
echo "10. Checking linter configuration..."
if [ -f ".golangci.yml" ]; then
    check_pass "Linter configuration exists"
else
    check_warn "No .golangci.yml found"
fi

# Summary
echo ""
echo "========================================"
echo "VERIFICATION SUMMARY"
echo "========================================"
echo ""

if [ $ERRORS -eq 0 ]; then
    echo -e "${GREEN}✓ All checks passed!${NC}"
    echo ""
    echo "The project has been successfully renamed to rag-knowledge-service."
    echo ""
    echo "Next steps:"
    echo "1. Update GitHub repository name (Settings → Repository name)"
    echo "2. Update GitHub Actions permissions (Settings → Actions → General)"
    echo "3. Verify GitHub secrets are configured"
    echo ""
    exit 0
else
    echo -e "${RED}✗ $ERRORS error(s) found${NC}"
    if [ $WARNINGS -gt 0 ]; then
        echo -e "${YELLOW}⚠ $WARNINGS warning(s)${NC}"
    fi
    echo ""
    echo "Please review the errors above and fix them."
    echo ""
    exit 1
fi
