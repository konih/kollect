# ADR-0003: CEL and JSONPath attribute extraction

## Status

Accepted

## Context

`KollectProfile` stores named attribute paths that the collection engine evaluates against
unstructured Kubernetes objects. Profiles must support both **JSONPath** (kubectl-style field
selection) and **CEL** (computed predicates and string formatting) without ambiguity.

Prior art: kube-state-metrics custom-resource-state metrics, Flux CEL filtering, and the legacy
batch collector's fixed Go schema (rejected — see ADR-0004).

## Decision

- Implement extraction in `internal/collect/extractor.go`.
- **JSONPath** paths use kubectl syntax (`{.metadata.name}`) or `$`-prefixed JSONPath
  (`$.metadata.name`). No prefix required.
- **CEL** expressions are prefixed with `cel:` (e.g. `cel:object.metadata.name`). The CEL
  environment exposes the full unstructured object as `object`.
- Empty paths are terminal errors (`ErrTerminal` per ADR-0020).
- Optional attributes (`spec.attributes[].optional: true`) skip extraction failures and absent
  values without failing the whole object.
- Required attributes that resolve to `null`/missing still record `null` in the result map for
  JSONPath; CEL `null` is treated the same.

## Consequences

### Positive

- Unambiguous path language selection without a separate CRD field.
- CEL prefix is grep-friendly in profiles and samples.
- Table-driven unit tests cover both engines without a cluster.

### Negative

- Authors must remember the `cel:` prefix; webhook validation could enforce it later.
- CEL and JSONPath type coercion differs — `type` on `AttributeSpec` is documentary until export
  validation lands.

## Open questions

- **OPEN:** Add OpenAPI/webhook validation that paths are either valid JSONPath or `cel:` prefixed?
- **OPEN:** Support JSONPath filters (`[?(@.label)]`) in Phase 1 or defer to property tests?
