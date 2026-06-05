#!/usr/bin/env bash
# Post-export integration asserts for git (bare repo) and postgres (testcontainers).
# Complements in-cluster git-export-assert.sh when sinks are not wired in kind smoke.
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$ROOT"

_log() { echo "[export-assert] $*"; }

_log "Git bare-repo export round-trip..."
go test -count=1 -timeout=3m -tags=integration \
  ./internal/sink/git/ -run '^TestExportBareRepoIntegration$'

_log "Postgres upsert export round-trip..."
go test -count=1 -timeout=10m -tags=integration \
  ./internal/sink/postgres/ -run '^TestExportPostgres$'

_log "Export integration asserts passed."
