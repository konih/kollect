# Performance and scalability

kollect is designed for large clusters with many watched resources. This guide summarizes
tuning knobs introduced in [ADR-0026](adr/0026-performance-scalability.md).

## Controller parallelism

| Flag | Default | Controller |
| --- | --- | --- |
| `--max-concurrent-reconciles-target` | `5` | `KollectTarget` |
| `--max-concurrent-reconciles-inventory` | `3` | `KollectInventory` |
| `--max-concurrent-reconciles-hub` | `2` | `KollectHub` |

Raise concurrency when reconcile latency grows while CPU is underutilized. Lower it when
API server throttling or etcd watch pressure appears.

## Workqueue rate limiting

When `--reconcile-rate-limit` is **unset** (`0`), controller-runtime uses its default
exponential failure rate limiter (5ms base, 1000s cap). Set a positive duration (for example
`100ms`) to slow retries on terminal errors without changing success-path throughput.

`kollect_workqueue_depth` approximates queue pressure as **in-flight reconciles** per controller
(not the internal client-go queue length).

## Export debouncing

`--export-debounce` (default `30s`) coalesces identical inventory payloads per `KollectInventory`
before hitting external sinks. Lower for fresher Git/DB exports; raise to reduce sink API load.

## Collection engine

- **In-memory store:** O(n) memory in collected object count; one `RWMutex` guards nested maps.
- **Informers:** When all active targets for a GVK resolve to **one** namespace via
  `spec.namespaceSelector`, the dynamic informer is scoped to that namespace. Otherwise the
  informer watches all namespaces and filters events at dispatch time (correctness over cache size).
- **Resync:** 12h informer resync is a correctness backstop, not a freshness driver.

## Metrics

| Metric | Type | Notes |
| --- | --- | --- |
| `kollect_workqueue_depth` | Gauge | In-flight reconciles per controller |
| `kollect_reconcile_duration_seconds` | Histogram | Per-controller reconcile latency |
| `kollect_informer_objects` | Gauge | Indexer size by GVR |
| `kollect_export_bytes_total` | Counter | Payload bytes per sink type |

## Profiling

`--enable-pprof` serves standard Go profiles on `--pprof-bind-address` (default `:6060`),
separate from Prometheus metrics (`:8080` / `:8443`). Helm sets `pprof.enabled: false` by default;
enable in dev overlays only.

## Benchmarks and load tests

```bash
task bench
KOLECT_LOAD_TEST=1 task load-test
```

Default `go test ./...` excludes `load`-tagged tests.
