# ADR-0012: Prometheus metrics and export sink (stub)

## Status

Accepted (stub, 2026-06-05)

## Context

Phase 4 targets kube-state-metrics-style custom resource metrics and a Prometheus export sink.
Phase 1 ships operator metrics per [ADR-0020](0020-error-taxonomy.md) instead of a full CR metrics
pipeline.

## Decision

1. **Operator metrics (Phase 1):** expose cardinality-safe gauges and histograms on the controller
   `/metrics` endpoint — including `kollect_collect_items_total`, `kollect_collected_objects`,
   `kollect_export_duration_seconds`, and reconcile counters.
2. **Prometheus sink (`spec.type: prometheus`):** register in the sink factory but return a clear
   error on export until Phase 4 implements remote-write or scrape-friendly export.
3. **Full kube-state-metrics-style `CustomResourceStateMetrics` config:** deferred to Phase 4; optional
   `KollectProfile.spec.metrics` field reserved in design docs only.

## Consequences

- Portals can scrape operator metrics today; sink-type `prometheus` is not a supported export path yet.
- Phase 4 work must extend this ADR rather than invent a parallel metrics model.

## See also

- [ADR-0020: Error taxonomy and metrics](0020-error-taxonomy.md)
- [ADR-0024: Inventory HTTP auth](0024-inventory-api-auth.md)
