# Architecture and operations

The token bucket stores a floating-point token balance plus last-refill time. A decision first adds elapsed refill, caps at capacity, then consumes only when enough tokens exist. Local buckets are protected by a mutex. The supplied Redis script executes the equivalent transition atomically on the server.

## Distributed design

Run stateless API replicas close to callers and attach regional Redis. Namespace keys by policy version and tenant. Set expiry long enough to preserve low-rate limits but finite enough to reclaim abandoned tenants. Decide explicitly whether Redis failure should fail open or closed per endpoint risk; do not silently mix policies.

Metrics use decision labels only. Tenant IDs belong in sampled structured logs or traces—not Prometheus labels—because unbounded cardinality can destabilize monitoring.

## Twelve-factor mapping

| Factor | Implementation |
|---|---|
| Codebase | One quota-service codebase |
| Dependencies | Go modules plus a versioned Lua artifact |
| Config | Capacity, refill, Redis URL, and port come from environment |
| Backing services | Redis is an attached resource, not embedded infrastructure |
| Build, release, run | Tested source -> immutable image -> configured runtime |
| Processes | API replicas keep no durable distributed quota state |
| Port binding | Service exports HTTP on `$PORT` |
| Concurrency | Horizontal processes share atomic Redis state |
| Disposability | Atomic checks and short requests make restarts safe |
| Dev/prod parity | Same algorithm in local tests and Redis Lua deployment |
| Logs | Structured events and metrics stream externally |
| Admin processes | Policy migration or key cleanup uses the same release |
