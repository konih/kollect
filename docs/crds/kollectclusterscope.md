# KollectClusterScope

Cluster-scoped tenancy **ceiling** for platform operators ([ADR-0207](../adr/0207-target-collection-filtering.md)).

## Spec

| Field | Role |
| --- | --- |
| `allowedGVKs` | Cap on GVKs cluster targets may collect |
| `allowedNamespaces` | Cap on workload namespaces |
| `deniedNamespaces` | Platform blacklist — not overridable by Targets |
| `sinkRefs` | Permitted cluster sink names for export |

Static config only — no status subresource ([ADR-0202](../adr/0202-static-vs-reconciled.md)).

## See also

- [KollectScope](kollectscope.md) — namespaced ceiling
- [KollectClusterTarget](kollectclustertarget.md) — collection intent
