# Quota Lab

[![CI](https://github.com/borikhan/quota-lab/actions/workflows/ci.yml/badge.svg)](https://github.com/borikhan/quota-lab/actions/workflows/ci.yml)
[![Go](https://img.shields.io/badge/Go-1.26-00ADD8?logo=go)](https://go.dev/)

A tenant-aware token-bucket quota service in Go, with both a concurrency-safe local implementation and an atomic Redis Lua algorithm for distributed deployments.

## Engineering highlights

- Continuous token refill allows bursts without losing long-term rate guarantees
- Decisions return remaining capacity and an actionable retry delay
- Tenant isolation is explicit in every check
- Redis Lua performs read/refill/consume/write atomically across API replicas
- Prometheus-format counters expose allowed and rejected traffic

```mermaid
flowchart LR
  Gateway --> API[Quota API]
  API --> Lua[Atomic Lua]
  Lua --> Redis[(Redis buckets)]
  API --> Metrics[Prometheus metrics]
```

## Quick start

```bash
go test -race ./...
go run ./cmd/server
curl -i -X POST localhost:8080/v1/check -H 'X-Tenant-ID: acme'
```

Capacity and refill rate are externalized in [`.env.example`](.env.example). `GET /metrics` exposes low-cardinality decision counters.

Read [architecture and operations](docs/architecture.md) for distributed correctness, failure policy, cardinality controls, and twelve-factor design.

## License

MIT
