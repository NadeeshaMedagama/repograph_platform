# Tests

This directory contains tests for RepoGraph Platform.

## Structure (To Be Implemented)

```
tests/
├── unit/              # Unit tests (alongside code)
├── integration/       # Integration tests
│   ├── services_test.go
│   └── workflow_test.go
├── e2e/              # End-to-end tests
│   └── full_pipeline_test.go
├── fixtures/         # Test data
│   ├── documents/
│   └── expected/
└── testutils/        # Test utilities
    └── helpers.go
```

## Running Tests

```bash
# Unit tests (in each package)
go test ./...

# Integration tests
go test -tags=integration ./tests/integration/...

# E2E tests
go test -tags=e2e ./tests/e2e/...

# With coverage
make test-coverage
```

## Test Coverage Goal

Target: **80%+** coverage for all packages

Current status: To be implemented

## Guidelines

1. Write tests alongside implementation
2. Use table-driven tests
3. Mock external dependencies
4. Use testcontainers for integration tests
5. Follow Go testing best practices

See `docs/DEVELOPMENT.md` for testing guidelines.
