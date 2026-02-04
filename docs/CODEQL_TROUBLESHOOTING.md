# CodeQL Troubleshooting Guide

## Issue: Request body length does not match content-length header

### Problem Description
When running CodeQL analysis in GitHub Actions, you may encounter this error during the upload phase:

```
Warning: Request body length does not match content-length header
Error: Request body length does not match content-length header
Warning: An unexpected error occurred when sending a status report: Request body length does not match content-length header
```

### Root Causes

1. **HTTP Compression Mismatch**: The SARIF file is compressed differently than expected by GitHub's API
2. **Large SARIF Files**: Results file exceeds size limits or causes chunking issues
3. **Network/Proxy Issues**: Intermediate proxies or network equipment modifying content
4. **CodeQL Action Version**: Bugs in specific versions of the action
5. **Database Clustering**: The database finalization process may cause issues

## Solution 1: Split Analysis and Upload (Applied)

The workflow has been updated to separate the analysis and upload steps:

```yaml
- name: Perform CodeQL Analysis
  uses: github/codeql-action/analyze@v4
  with:
    category: "/language:${{ matrix.language }}"
    output: sarif-results
    upload-database: false
  env:
    CODEQL_ACTION_EXTRA_OPTIONS: '{"database": {"finalize": ["--no-db-cluster"]}}'

- name: Filter and Upload SARIF
  uses: github/codeql-action/upload-sarif@v4
  if: always()
  with:
    sarif_file: sarif-results/${{ matrix.language }}.sarif
    category: "/language:${{ matrix.language }}"
    wait-for-processing: true
```

**Key changes:**
- `output: sarif-results` - Saves results to a directory
- `upload-database: false` - Prevents automatic upload during analysis
- `--no-db-cluster` - Disables database clustering which can cause issues
- Separate `upload-sarif` step with `if: always()` ensures upload happens even if analysis has warnings

## Solution 2: Use Latest CodeQL CLI

If the issue persists, try using a specific CodeQL CLI version:

```yaml
- name: Initialize CodeQL
  uses: github/codeql-action/init@v4
  with:
    languages: ${{ matrix.language }}
    queries: +security-and-quality
    tools: latest  # or specify a version like '2.15.4'
```

## Solution 3: Reduce Query Scope

If the SARIF file is too large, reduce the queries:

```yaml
- name: Initialize CodeQL
  uses: github/codeql-action/init@v4
  with:
    languages: ${{ matrix.language }}
    queries: +security-extended  # Instead of security-and-quality
```

Or use only security queries:

```yaml
queries: security-extended
```

## Solution 4: Disable Specific Queries

Create a custom query suite to exclude problematic queries:

1. Create `.github/codeql/codeql-config.yml`:

```yaml
name: "CodeQL Config"
queries:
  - uses: security-extended
  - uses: security-and-quality
    
paths-ignore:
  - "**/test/**"
  - "**/vendor/**"
  - "**/node_modules/**"
```

2. Update workflow:

```yaml
- name: Initialize CodeQL
  uses: github/codeql-action/init@v4
  with:
    languages: ${{ matrix.language }}
    config-file: ./.github/codeql/codeql-config.yml
```

## Solution 5: Manual SARIF Processing

If all else fails, process SARIF manually:

```yaml
- name: Perform CodeQL Analysis
  uses: github/codeql-action/analyze@v4
  with:
    category: "/language:${{ matrix.language }}"
    upload: Failed  # Don't upload
    output: sarif-results

- name: Process SARIF
  if: always()
  run: |
    # Compress SARIF file to reduce size
    gzip -c sarif-results/${{ matrix.language }}.sarif > results.sarif.gz
    gunzip results.sarif.gz
    mv results.sarif processed-results.sarif

- name: Upload Processed SARIF
  uses: github/codeql-action/upload-sarif@v4
  if: always()
  with:
    sarif_file: processed-results.sarif
```

## Solution 6: Check GitHub Status

Sometimes this is a transient GitHub API issue:

1. Check [GitHub Status](https://www.githubstatus.com/)
2. Check [CodeQL Action Issues](https://github.com/github/codeql-action/issues)
3. Wait and retry the workflow

## Solution 7: Update Actions Versions

Ensure all actions are using latest versions:

```yaml
- uses: actions/checkout@v4
- uses: actions/setup-go@v5
- uses: github/codeql-action/init@v4
- uses: github/codeql-action/autobuild@v4
- uses: github/codeql-action/analyze@v4
- uses: github/codeql-action/upload-sarif@v4
```

## Solution 8: Increase Runner Resources

For large codebases, use a larger runner:

```yaml
jobs:
  analyze:
    runs-on: ubuntu-latest-4-cores  # or ubuntu-latest-8-cores
```

## Verification

After applying fixes, verify the workflow:

1. **Push a test commit:**
   ```bash
   git add .github/workflows/codeql.yml
   git commit -m "fix: resolve CodeQL upload content-length mismatch"
   git push
   ```

2. **Monitor the workflow:**
   - Go to **Actions** tab in GitHub
   - Watch the CodeQL workflow
   - Check for successful upload

3. **Verify results:**
   - Go to **Security** → **Code scanning**
   - Verify alerts are visible
   - Check the timestamp is recent

## Debug Information

If the issue persists, enable debug logging:

1. Add repository secret `ACTIONS_STEP_DEBUG` with value `true`
2. Add repository secret `ACTIONS_RUNNER_DEBUG` with value `true`
3. Re-run the workflow
4. Check logs for detailed output

## Common Error Messages

| Error | Likely Cause | Solution |
|-------|--------------|----------|
| Request body length mismatch | Compression/size issue | Solution 1 (applied) |
| Upload failed | Network/API issue | Retry or Solution 6 |
| SARIF file not found | Wrong path | Check `sarif_file` path |
| Category mismatch | Incorrect category format | Use `/language:go` format |
| Permission denied | Missing permissions | Check workflow permissions |

## Additional Resources

- [CodeQL Documentation](https://codeql.github.com/docs/)
- [CodeQL Action Repository](https://github.com/github/codeql-action)
- [SARIF Support](https://docs.github.com/en/code-security/code-scanning/integrating-with-code-scanning/sarif-support-for-code-scanning)
- [Troubleshooting Code Scanning](https://docs.github.com/en/code-security/code-scanning/troubleshooting-code-scanning)

## Contact

If none of these solutions work:
1. Open an issue on [github/codeql-action](https://github.com/github/codeql-action/issues)
2. Include workflow logs with debug enabled
3. Include SARIF file size and runner information
