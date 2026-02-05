# CodeQL Upload Fix - Summary

## Date
February 4, 2026

## Issue
CodeQL workflow was failing with the following error during the upload phase:
```
Warning: Request body length does not match content-length header
Error: Request body length does not match content-length header
Warning: An unexpected error occurred when sending a status report: Request body length does not match content-length header
```

## Root Cause
The error occurs when the CodeQL action attempts to upload SARIF results to GitHub's API, but the HTTP content-length header doesn't match the actual request body length. This can happen due to:

1. **Compression issues**: The SARIF file is compressed differently than expected
2. **Database clustering**: The database finalization process may cause size mismatches
3. **Network issues**: Intermediate proxies or network equipment modifying content
4. **Large result files**: Size exceeds expected limits causing chunking problems
5. **Duplicate upload**: Both analyze and upload-sarif actions try to upload with the same category

## Solution Applied

### Modified File: `.github/workflows/codeql.yml`

**Changes Made:**

1. **Split Analysis and Upload Steps**
   - Changed from single `analyze` action to separate `analyze` and `upload-sarif` steps
   - This gives better control over the upload process

2. **Added Output Directory and Disabled Upload**
   ```yaml
   output: sarif-results
   upload: never
   ```
   - Saves SARIF results to a directory instead of automatic upload
   - **`upload: never`** - CRITICAL: Prevents analyze action from uploading SARIF
   - This avoids the "only one run per category" error

3. **Added Environment Options**
   ```yaml
   env:
     CODEQL_ACTION_EXTRA_OPTIONS: '{"database": {"finalize": ["--no-db-cluster"]}}'
   ```
   - Disables database clustering which can cause content-length mismatches
   - Helps with compression issues

4. **Separate Upload Step**
   ```yaml
   - name: Filter and Upload SARIF
     uses: github/codeql-action/upload-sarif@v4
     if: always()
     with:
       sarif_file: sarif-results/${{ matrix.language }}.sarif
       category: "/language:${{ matrix.language }}"
       wait-for-processing: true
   ```
   - Uses dedicated upload action for better reliability
   - `if: always()` ensures upload happens even if analysis has warnings
   - `wait-for-processing: true` ensures proper handling of upload

5. **Added Go Cache**
   ```yaml
   cache: true
   ```
   - Speeds up builds by caching Go modules

## Files Modified

1. ✅ `.github/workflows/codeql.yml` - Updated workflow with fix
2. ✅ `docs/CODEQL_TROUBLESHOOTING.md` - Created comprehensive troubleshooting guide
3. ✅ `README.md` - Added link to troubleshooting guide

## Testing

To verify the fix:

1. **Commit and push changes:**
   ```bash
   cd /home/nadeeshame/go/rag_knowledge_service
   git add .github/workflows/codeql.yml docs/CODEQL_TROUBLESHOOTING.md README.md
   git commit -m "fix: resolve CodeQL upload content-length mismatch error"
   git push origin master
   ```

2. **Monitor workflow:**
   - Go to GitHub Actions tab
   - Watch the CodeQL workflow execution
   - Verify successful upload in the "Filter and Upload SARIF" step

3. **Verify results:**
   - Navigate to Security → Code scanning
   - Confirm alerts are visible and up to date

## Expected Behavior After Fix

### Before (Error):
```
Uploading code scanning results
  Uploading results
  Warning: Request body length does not match content-length header
  Error: Request body length does not match content-length header
```

### After (Success):
```
Perform CodeQL Analysis
  ✓ Analysis completed successfully
  ✓ SARIF file saved to sarif-results/go.sarif

Filter and Upload SARIF
  ✓ Reading SARIF file
  ✓ Uploading to GitHub
  ✓ Upload successful
  ✓ Results processed
```

## Alternative Solutions (If This Doesn't Work)

If the primary fix doesn't resolve the issue, try these alternatives (documented in `docs/CODEQL_TROUBLESHOOTING.md`):

1. **Update CodeQL tools version**
2. **Reduce query scope** (use `security-extended` instead of `security-and-quality`)
3. **Create custom query configuration**
4. **Manual SARIF compression**
5. **Check GitHub Status** (might be API issue)
6. **Use larger GitHub runner**
7. **Enable debug logging**

## Benefits

1. ✅ **More reliable uploads**: Separation of concerns reduces failure points
2. ✅ **Better error handling**: `if: always()` ensures upload attempts even with warnings
3. ✅ **Easier debugging**: Separate steps make it clear where failures occur
4. ✅ **Future-proof**: Configuration options prevent similar issues
5. ✅ **Documented**: Comprehensive troubleshooting guide for future reference

## References

- [CodeQL Action Documentation](https://github.com/github/codeql-action)
- [SARIF Support for Code Scanning](https://docs.github.com/en/code-security/code-scanning/integrating-with-code-scanning/sarif-support-for-code-scanning)
- [Troubleshooting Code Scanning](https://docs.github.com/en/code-security/code-scanning/troubleshooting-code-scanning)
- [GitHub CodeQL Action Issues](https://github.com/github/codeql-action/issues)

## Status

✅ **Fix Applied** - Ready for testing
📝 **Documentation Updated** - Troubleshooting guide created
🔄 **Awaiting Verification** - Push changes and monitor workflow

## Next Steps

1. Commit and push the changes
2. Monitor the next CodeQL workflow run
3. If successful, close any related issues
4. If unsuccessful, consult `docs/CODEQL_TROUBLESHOOTING.md` for alternative solutions
