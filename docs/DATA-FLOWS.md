# kollect data flows

Visual walkthroughs of how data moves through the operator. For CRD roles see
[ARCHITECTURE.md](ARCHITECTURE.md); for locked decisions see
[PLATFORM-DECISIONS.md](PLATFORM-DECISIONS.md).

---

## 1. Export debouncing

**Problem:** Event-driven informers can fire hundreds of updates per minute. Without coalescing,
every watch event would trigger a Postgres upsert or Git commit.

**Design:** The in-memory collect store updates **immediately** on every target reconcile. Only the
**sink export** step is debounced per `KollectInventory`.

### Per-inventory state machine

```mermaid
flowchart TD
  Start([Inventory reconcile]) --> Scope{In KollectScope?}
  Scope -->|no / not enforced| SinkOK{Sink reachable?}
  Scope -->|violation| Degrade[Degraded — no export]
  SinkOK -->|no| Degrade
  SinkOK -->|yes| Marshal[Marshal namespace payload]
  Marshal --> Size{Within maxExportBytes?}
  Size -->|no| Degrade
  Size -->|yes| Hash[Compute payload checksum]
  Hash --> Gen{generation changed?}
  Gen -->|yes| Export[Export to sinks now]
  Gen -->|no| Chk{checksum changed?}
  Chk -->|yes| Export
  Chk -->|no| Interval{elapsed ≥ exportMinInterval?}
  Interval -->|yes| Export
  Interval -->|no| Wait[RequeueAfter remainder]
  Export --> Status[Update status: lastExportTime, conditions]
  Wait --> Start
```

### Timing example (default `exportMinInterval: 30s`)

```mermaid
sequenceDiagram
  participant API as Kubernetes API
  participant Store as Collect store
  participant Inv as KollectInventory
  participant PG as Postgres sink

  API->>Store: Deployment image patch (t=0s)
  Note over Store: row updated immediately
  Store->>Inv: trigger reconcile
  Inv->>PG: export (checksum changed)

  API->>Store: unrelated resync (t=5s)
  Store->>Inv: trigger reconcile
  Note over Inv: same checksum — debounced
  Inv-->>Inv: requeue ~25s

  API->>Store: second patch (t=12s)
  Store->>Inv: trigger reconcile
  Inv->>PG: export (checksum changed — bypass interval)

  Note over Inv,PG: identical payload at t=40s → export allowed (30s elapsed)
```

### Configuration

| Field | Default | Effect |
| --- | --- | --- |
| `KollectInventory.spec.exportMinInterval` | **30s** (CRD default) | Min gap between exports of **identical** payload |
| `metadata.generation` bump | — | Immediate export (spec edit) |
| Payload checksum change | — | Immediate export (material inventory change) |

Operator `--export-debounce` is a **deprecated fallback** when the spec field is unset on legacy
manifests.

---

## 2. Collection pipeline

How a watched object becomes an inventory row.

```mermaid
flowchart LR
  subgraph api [Kubernetes API]
    Obj[Deployment / CRD / …]
  end

  subgraph operator [kollect operator]
    Inf[Shared informer<br/>one per GVK]
    Tgt[KollectTarget<br/>reconciler]
    Prof[KollectProfile<br/>schema]
    Ext[Extractor<br/>JSONPath / CEL]
    Store[(Collect store<br/>per namespace)]
    Inv[KollectInventory]
  end

  Obj -->|watch| Inf
  Inf -->|enqueue| Tgt
  Prof -.->|profileRef| Tgt
  Tgt -->|label/NS filter| Ext
  Obj -.->|cached object| Ext
  Ext -->|attribute map| Store
  Store --> Inv
```

**Key properties:**

- **One informer per GVK** across all targets ([ADR-0014](adr/0014-event-driven-informers.md)).
- Targets only differ by **namespace/label selectors** and **profileRef**.
- Extraction runs on the **cached unstructured object** — no per-target API list calls.

---

## 3. Attribute extraction (JSONPath arrays)

`KollectProfile` attributes are evaluated per object. Single-index paths return a scalar; wildcard
paths return a **JSON array** in the export row.

```mermaid
flowchart TD
  Obj[Unstructured object] --> Path{Path type?}
  Path -->|cel:…| CEL[CEL evaluator]
  Path -->|$.… or {.…}| JP[kubectl JSONPath]
  CEL --> Val[Go value]
  JP --> Matches{match count}
  Matches -->|1| Scalar[scalar in row]
  Matches -->|>1| List[array in row]
  Matches -->|0| Opt{optional?}
  Opt -->|yes| Skip[omit attribute]
  Opt -->|no| Null[null in row]
```

**Deployment containers example:**

| Path | Result for 2-container pod |
| --- | --- |
| `$.spec.template.spec.containers[0].image` | `"app:1.0"` (string) |
| `$.spec.template.spec.containers[*].image` | `["app:1.0", "sidecar:2.0"]` (list) |

See [ADR-0003](adr/0003-cel-jsonpath-extraction.md) for syntax rules.

---

## 4. `KollectScope` enforcement gate

Static scope object; enforced at **target** and **inventory** reconcile time (hard degrade).

```mermaid
flowchart TD
  Tgt[KollectTarget reconcile] --> Load[Load KollectScope in namespace]
  Load --> Enforced{scope exists?}
  Enforced -->|no| Collect[Register watch + collect]
  Enforced -->|yes| GVK{profile GVK in allowedGVKs?}
  GVK -->|no| DenyT[Degraded ScopeGVKDenied]
  GVK -->|yes| NS{workload NS in allowedNamespaces?}
  NS -->|no| DenyT
  NS -->|yes| Collect

  Inv[KollectInventory reconcile] --> Load2[Load KollectScope]
  Load2 --> Sinks{sinkRefs ⊆ scope.sinkRefs?}
  Sinks -->|no| DenyI[Degraded ScopeSinkDenied]
  Sinks -->|yes| Export[Continue export path]
```

Example: [ADR-0016](adr/0016-namespaced-multi-tenancy.md#enforcement-example-gvk-denied).

---

## 5. `KollectConnectionTest` lifecycle

One-shot probe CR for audited connectivity checks.

```mermaid
stateDiagram-v2
  [*] --> Pending: CR created
  Pending --> Probing: reconciler starts
  Probing --> Succeeded: sink OK
  Probing --> Failed: sink error
  Succeeded --> TTL: status.completed
  Failed --> TTL
  TTL --> Deleted: after ttlSecondsAfterFinished
  Pending --> Probing: spec change (generation bump)
  Succeeded --> Probing: spec change re-probe
  Failed --> Probing: spec change re-probe
```

Default TTL: **300s**. Patch `spec.sinkRef` to force a fresh probe.

---

## See also

- [ARCHITECTURE.md](ARCHITECTURE.md) — CRD model and deployment defaults
- [PERFORMANCE.md](PERFORMANCE.md) — metrics and tuning
- [examples/deployment-inventory.md](examples/deployment-inventory.md) — end-to-end walkthrough
