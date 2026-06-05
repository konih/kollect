# ADR-0026: Performance and scalability

## Status

Accepted

## Context

kollect watches arbitrary GVKs, aggregates attributes in memory, and exports on inventory
reconcile. Large clusters need tunable controller parallelism, observable queue pressure,
bounded sink churn, and optional profiling without coupling to Prometheus scrape paths.

## Decision

1. **Controller options:** Expose `MaxConcurrentReconciles` per reconciler
   (`KollectTarget`, `KollectInventory`, `KollectHub`) via operator flags with documented defaults.
2. **Workqueue:** Use controller-runtime default exponential failure rate limiting unless
   `--reconcile-rate-limit` overrides the base delay. Approximate queue depth with an in-flight
   reconcile gauge (`kollect_workqueue_depth`).
3. **Metrics:** Add reconcile duration histogram, informer indexer size gauge, and export byte
   counter alongside existing export latency histogram.
4. **Export debounce:** Make inventory export debounce configurable via `--export-debounce`.
5. **Informers:** Scope dynamic informers to a single namespace when all targets for a GVR agree;
   otherwise watch all namespaces and filter by `namespaceSelector` at dispatch.
6. **Profiling:** Optional `--enable-pprof` on `:6060`; disabled in production Helm values.
7. **Tests:** `go test -bench` for extraction; optional `load`-tagged test gated by
   `KOLECT_LOAD_TEST=1`.

## Consequences

- Operators can scale reconcile throughput without rebuilding images.
- In-flight gauge is an approximation, not a substitute for controller-runtime's internal queue metrics.
- Multi-namespace targets still use cluster-wide informer caches when scopes differ.

## References

- [ADR-0014](0014-event-driven-informers.md) — event-driven collection
- [ADR-0020](0020-error-taxonomy.md) — error classes and requeue behavior
