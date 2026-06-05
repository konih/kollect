# ADR-0004: CRD model — prefixed kinds, static vs reconciled split

## Status

Accepted

## Context

kollect replaces a hardcoded batch collector schema with CRD-driven configuration. Operators in
this space use different tenancy and config patterns:

- **external-secrets** splits `SecretStore` (namespaced) vs `ClusterSecretStore` (cluster) with
  namespace `conditions` on the cluster variant; provider config is a discriminated union in spec.
- **Flux notification-controller** keeps `Provider` and `Alert` as **static** objects (no status
  subresource, no dedicated reconciler) while `Receiver` is reconciled.
- **Argo CD** uses `AppProject` as a tenancy boundary and `Application` as the reconciled unit.

kollect must combine generic attribute selection, resource selection, aggregation, multi-backend
export, and (later) doc-sync — no single OSS project covers all of this.

Platform users commonly terminate TLS to internal Git/GitLab with **custom CAs**; sink configuration
must support that from early phases, not as a later bolt-on.

## Decision

API group `kollect.dev/v1alpha1`. All kinds are **prefixed** (`Kollect*`) to avoid collisions.

### Static config (validated, no controller)

| Kind | Scope | Role |
| --- | --- | --- |
| `KollectProfile` | Cluster | Reusable extraction schema: GVK + named CEL/JSONPath attributes |
| `KollectSink` | Cluster | Backend config: `type` + endpoint + `secretRef` + TLS trust; resolved via Go registry |
| `KollectScope` | Cluster | Tenancy boundary (Phase 3): allowed GVKs, namespaces, sinks |

### Reconciled (controller + dynamic informers)

| Kind | Scope | Role |
| --- | --- | --- |
| `KollectTarget` | Namespaced | `profileRef` + selectors + optional CEL predicate; drives collection |
| `KollectInventory` | Cluster | Aggregates targets; dispatches to sinks |
| `KollectPublication` | Namespaced | **Deferred** until collection is mature — template render + doc-backend sync |

### Reserved (designed, not built yet)

- `KollectReceiver` — inbound webhook → trigger (Flux Receiver pattern).
- `KollectTargetSet` — generator templating many Targets (ApplicationSet pattern).

Short names: `kprof`, `ksink`, `kscope`, `ktgt`, `kinv`, `kpub`.

All reconciled kinds support `spec.suspend` and `kollect.dev/requestedAt` manual-trigger annotation.

### Validating admission (early)

**Validating webhooks** (Phase 0/1, not post-hoc workarounds) must reject at apply time:

- Invalid or non-compilable **CEL** and **JSONPath** expressions on `KollectProfile` attributes
- Unknown `KollectSink.type` values (enum aligned with [ADR-0020](0020-error-taxonomy.md))
- Cross-field constraints already expressed as CEL `x-kubernetes-validations` in OpenAPI

Runtime collection may still surface `ErrTerminal` for GVK/API discovery failures; expression syntax
is webhook-validated first.

### `KollectSink` TLS trust (custom CA)

Git and GitLab sink specs include explicit trust material (at least one of):

- `tls.caBundle` — inline PEM (base64) for org-internal CA
- `tls.caSecretRef` — namespaced or cluster secret key containing CA bundle

Default remains **verify TLS**; `insecureSkipVerify` is opt-in only for dev and must be flagged in
status when used.

Credential material stays in `secretRef` (tokens, SSH keys) — never in spec/status.

### JSONPath filters on targets

**Deferred:** sink-side or target-side JSONPath *filters* (post-extraction row filtering) are not
Phase 1 API. Schema clarity and aggregation matter more than where filtering runs ([REQUIREMENTS.md](../REQUIREMENTS.md)).

## Consequences

### Positive

- Static Profile/Sink cuts moving parts (validated at admission, read at reconcile time).
- Namespaced Targets align with team ownership; cluster Profiles/Sinks enable shared platform config.
- Prefix naming is grep-friendly and avoids generic kind collisions in multi-operator clusters.
- Early webhooks prevent bad profiles from wedging reconcilers.

### Negative

- Cluster-scoped `KollectInventory` may be awkward for strict multi-tenant isolation (see open questions).
- `KollectSink` is cluster-scoped only today — no namespaced sink variant unlike ESO's Store split.
- Webhook + CEL maintenance cost on every new attribute type rule.

## Open questions

- **OPEN:** Should `KollectInventory` be **namespaced** (team-owned rollup) with optional cluster
  aggregator, mirroring `SecretStore` vs `ClusterSecretStore`? Cluster scope simplifies a single
  portal view but widens RBAC blast radius.
- **OPEN:** Add `KollectClusterSink` + namespaced `KollectSink` split in Phase 1 or defer to Phase 3
  with `KollectScope`?
- **OPEN:** Single `caBundle` field vs only `secretRef` for CA — size limits on CRD spec?
