# Testing strategy

Kollect is **TDD-first**. Quality gates follow a six-tier test pyramid (L0–L5) defined in
[ADR-0706: Testing and merge-gate architecture](../adr/0706-testing-merge-gate-architecture.md).

!!! tip "Quick local loop"
    Before opening a PR: `task lint` · `task coverage` · `task coverage:race` (recommended) ·
    `task verify` · `task scrub`. See [Coding standards](coding-standards.md) and
    [CONTRIBUTING.md](../../CONTRIBUTING.md) for the full checklist.

## Test pyramid (L0–L5)

| Tier | Scope | Blocks merge? | Typical command |
| --- | --- | --- | --- |
| **L0 — Unit** | Pure packages, table-driven tests, mocks | Yes | `task test` |
| **L1 — Controller / API** | envtest reconcilers, webhooks | Yes | `task coverage` (no `-race` in CI) |
| **L2 — Golden / contract** | OpenAPI fragments, sample YAML, extractor goldens | Yes | `task test` |
| **L3 — Integration** | Real Postgres, Kafka, Git, S3, GCS, Redis, NATS (testcontainers) | Yes | `task test-integration` |
| **L4 — E2E** | Kind cluster: Helm install, smoke, export asserts | **PR smoke (required)** + nightly / extended | `task test:e2e` |
| **UI — Playwright** | React SPA smoke (MSW dev server); not backend Kind L4 | No (optional locally) | `task ui-e2e` · `task ui-e2e:docker` |
| **L5 — Load / perf** | Bounded synthetic scale (≤2000 objects), micro-benchmarks | Opt-in | `task load-test` · `task perf-report` |

**Direction:** Most tests live at L0–L2. Every new sink backend must reach **L3** before merge
([NFR-EXT-3](../REQUIREMENTS.md)). L4 catches wiring regressions that unit tests miss. L5 stays
opt-in so default CI stays fast.

## Coverage target

Statement coverage on `./internal/...` is enforced by `hack/coverage.sh` via `task coverage`:

| Setting | Value |
| --- | --- |
| **Target (pre-v0.1.0)** | **80%** — ratchet `COVERAGE_MIN` when measured coverage is sustained |
| **Current CI floor** | 65% (`COVERAGE_MIN` in `.github/workflows/ci.yaml`) |
| **Codecov project target** | 70% (`codecov.yml`) |

Regressions below the enforced floor fail CI. Raise the floor only after coverage has grown
sustainably — see ADR-0706 for the ratchet policy.

```sh
task coverage          # unit + envtest + floor check → coverage.out (CI path; no -race)
task coverage:race     # local-only: COVERAGE_RACE=1 + CGO_ENABLED=1
task coverage:report   # go tool cover -func summary
task coverage:html     # coverage.html for browser review
```

Integration-tagged tests (`-tags=integration`) and e2e packages are excluded from the default
coverage profile.

## What CI runs on every PR

**Path filters:** GitHub skips **CI**, **Preflight**, and **CodeQL** when a PR or push to `main`
changes *only* documentation paths (`docs/**`, `mkdocs.yml`, root prose `*.md`, `CHANGELOG.md`,
`LICENSE`, issue templates). The **Docs** workflow then runs `task lint:markdown`, `mkdocs build`,
and (on `main` push) deploys to [konih.github.io/kollect](https://konih.github.io/kollect/). Any
change under `api/`, `internal/`, `charts/`, `cmd/`, `config/`, `hack/`, `ui/`, `test/`, `go.mod`,
or `.github/workflows/` — or a mixed docs+code PR — runs the full gate below. Release tags no
longer trigger docs deploys; the site tracks `main` only.

Binding jobs in `.github/workflows/ci.yaml` (see ADR-0706 for the full matrix):

- Secret scan (`gitleaks`), codegen drift (`task verify`), vulncheck, lint/format
- **Architecture fitness:** `go-arch-lint` via `task arch-lint` (import boundaries in
  `.go-arch-lint.yml`)
- **Dependency policy:** golangci-lint `depguard` + `gomodguard` (same `task lint` job)
- **L0–L2:** `task coverage` with coverage floor
- **L3:** `task test-integration` (Docker required)
- Helm packaging (`task helm-test`), image build (`task docker:build`)
- Native Go fuzz (CEL/JSONPath extractors, content hash)
- RBAC audit (`hack/audit-rbac.sh`)

**E2E smoke (L4 Tier 0):** `.github/workflows/e2e-smoke.yaml` job **`kind-smoke`** on every
non-docs PR and push to `main` (same `paths-ignore` as CI). Required for branch protection — see
[coding-standards.md](coding-standards.md).

**Non-blocking on PR:** `e2e-extended.yaml` (Tier 1 matrix + webhook profile; label `e2e/full` or
path-filtered), `task perf-report` (promoted to blocking at **v0.2.0** per ADR-0706);
**SonarCloud** scan (`sonarcloud` job — needs `SONAR_TOKEN`; see
[tooling-setup.md](tooling-setup.md)).

### Holistic maintainability (SonarCloud)

SonarCloud mirrors coverage trends and surfaces duplication / technical-debt ratios over time —
complementing point-in-time `dupl` and Codecov. Configured in `sonar-project.properties`; optional
until the maintainer adds `SONAR_TOKEN`. Does not replace `task lint` or arch-lint.

## Scheduled and manual tiers

| Workflow | Tier | Purpose |
| --- | --- | --- |
| `docs.yaml` | Docs | Markdown lint, MkDocs build, GitHub Pages deploy (`main` only) |
| `e2e-smoke.yaml` | L4 Tier 0 | **Mandatory** kind smoke on PR + `main` (job `kind-smoke`) |
| `e2e-extended.yaml` | L4 Tier 1 | Optional git-export, multitenant, tenant-mode, webhook profile |
| `e2e-nightly.yaml` | L4 Tier 2 | Full Kind matrix + bench/perf (deduped L3) |
| `test-e2e.yaml` | L4 Tier 3 | Manual full matrix (`workflow_dispatch`) |
| `release.yaml` | Supply chain | Image signing, SBOM, chart publish |

Set repository variable **`GIT_EXPORT_TEST_REPO`** (Settings → Actions → Variables) to enable full
remote git SHA assert in git-export scenarios. Without it, jobs verify inventory HTTP hash only.

## Local development commands

| Task | Purpose |
| --- | --- |
| `task test` | Unit + envtest (no floor check; no race detector) |
| `task coverage` | Unit + envtest + 65% floor (CI; CGO off, no `-race`) |
| `task coverage:race` | Same as coverage with race detector (local pre-PR) |
| `task test-integration` | L3 sink/transport integration (Docker) |
| `task test:e2e` | L4 kind smoke (setup → smoke → teardown) |
| `task ui-ci` | UI PR gate — Vitest, lint, build, mock drift (no Playwright) |
| `task ui-e2e` · `task ui-e2e:docker` | UI Playwright smoke — optional; distinct from backend L4 |
| `task bench` | Micro-benchmarks on hot paths |
| `KOLECT_LOAD_TEST=1 task load-test` | L5 bounded load (≤2000 objects, opt-in) |
| `task perf-report` | Benchmark + unit pass summary (local only, gitignored output) |

Full local setup: [DEVELOPMENT.md](../DEVELOPMENT.md).

## Definition of done

Per-change checklist: [guidelines § 6](guidelines.md#6-definition-of-done-per-change).
PR workflow: [CONTRIBUTING.md § Pull request process](../../CONTRIBUTING.md#pull-request-process).

## Further reading

- [ADR-0706: Testing and merge-gate architecture](../adr/0706-testing-merge-gate-architecture.md)
- [Engineering guidelines](https://github.com/konih/kollect/blob/main/docs/development/guidelines.md) §4 (testing rules)
- [REQUIREMENTS.md](../REQUIREMENTS.md) — NFR-TEST-* priorities
- [PERFORMANCE.md](../PERFORMANCE.md) — scale bounds and perf-report workflow
- [tooling-setup.md](tooling-setup.md) — arch-lint, depguard, SonarCloud maintainer steps
