#!/usr/bin/env bash
# Pre-recording checks: inventory Ready, Forgejo reachable, first export landed.
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=lib.sh
source "${SCRIPT_DIR}/lib.sh"

_hero_require_tools
_hero_source_state

if ! kind_cluster_exists "$HERO_CLUSTER"; then
  echo "Cluster ${HERO_CLUSTER} not found — run: task demo-hero-up" >&2
  exit 1
fi

kind_use_context "$HERO_CLUSTER"

_hero_log "Checking KollectInventory Ready..."
kubectl wait --for=condition=Ready kollectinventory/demo-inventory -n default --timeout=30s

_hero_log "Checking Git sink ConnectionVerified..."
kubectl wait --for=condition=ConnectionVerified kollectsnapshotsink/hero-git-sink \
  -n default --timeout=30s

_hero_start_port_forward
_hero_log "Checking Forgejo API..."
curl -fsS "http://127.0.0.1:${HERO_FORGEJO_PF_PORT}/api/v1/version" >/dev/null

if [[ ! -d "$HERO_INVENTORY_CLONE_DIR/.git" ]]; then
  _hero_clone_inventory_repo
fi

if ! find "$HERO_INVENTORY_CLONE_DIR" -type f \( -name '*.yaml' -o -name '*.yml' -o -name '*.json' \) \
  ! -path '*/.git/*' | grep -q .; then
  echo "No exported inventory files in ${HERO_INVENTORY_CLONE_DIR} — wait for debounce/export or re-run task demo-hero-up" >&2
  exit 1
fi

_hero_log "Pre-flight OK — safe to record."
