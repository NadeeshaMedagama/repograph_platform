# Kubernetes Deployment Manifests

This directory will contain Kubernetes deployment manifests for RAG Knowledge Service.

## Structure (To Be Created)

```
kubernetes/
├── deployments/
│   ├── orchestrator.yaml
│   ├── document-scanner.yaml
│   ├── content-extractor.yaml
│   ├── vision-service.yaml
│   ├── summarization-service.yaml
│   ├── embedding-service.yaml
│   ├── vector-store.yaml
│   └── query-service.yaml
├── services/
│   └── *.yaml (service definitions)
├── ingress/
│   └── ingress.yaml
├── configmaps/
│   └── app-config.yaml
├── secrets/
│   └── credentials.yaml (encrypted)
└── statefulsets/
    ├── postgres.yaml
    └── redis.yaml
```

## Priority

This is planned for **Phase 2** (production deployment).

For now, use Docker Compose for local development and testing.

See `docs/DEPLOYMENT.md` for current deployment options.
