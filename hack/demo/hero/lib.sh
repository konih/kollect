#!/usr/bin/env bash
# Shared helpers for the hero demo harness (Forgejo in kind + golden Git sample).
set -euo pipefail

HERO_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${HERO_DIR}/../../.." && pwd)"
# shellcheck source=../../kind/common.sh
source "${REPO_ROOT}/hack/kind/common.sh"

readonly HERO_CLUSTER="${HERO_CLUSTER:-kollect-hero}"
readonly HERO_CLUSTER_CONFIG="${HERO_CLUSTER_CONFIG:-${HERO_DIR}/cluster.yaml}"
readonly HERO_DEV_VALUES="${HERO_DEV_VALUES:-${REPO_ROOT}/charts/kollect/ci/dev-values.yaml}"
readonly HERO_STATE_FILE="${HERO_STATE_FILE:-/tmp/kollect-hero-state.env}"
readonly HERO_PF_PID_FILE="${HERO_PF_PID_FILE:-/tmp/kollect-hero-forgejo-pf.pid}"
readonly HERO_FORGEJO_NS="${HERO_FORGEJO_NS:-forgejo}"
readonly HERO_FORGEJO_USER="${HERO_FORGEJO_USER:-kollect}"
readonly HERO_FORGEJO_PASS="${HERO_FORGEJO_PASS:-kollect-demo}"
readonly HERO_FORGEJO_REPO="${HERO_FORGEJO_REPO:-inventory-demo}"
readonly HERO_FORGEJO_PF_PORT="${HERO_FORGEJO_PF_PORT:-13000}"
readonly HERO_INVENTORY_CLONE_DIR="${HERO_INVENTORY_CLONE_DIR:-/tmp/kollect-hero-inventory}"
readonly HERO_GIT_SECRET="${HERO_GIT_SECRET:-hero-git-credentials}"

_hero_log() { echo "[hero] $*"; }

_hero_require_tools() {
  _kind_require_tools
  _kind_require docker "https://docs.docker.com/get-docker/"
  _kind_require git "https://git-scm.com/"
  _kind_require curl "https://curl.se/"
  _kind_require task "https://taskfile.dev/"
}

_hero_write_state() {
  cat >"$HERO_STATE_FILE" <<EOF
HERO_CLUSTER=${HERO_CLUSTER}
HERO_FORGEJO_NS=${HERO_FORGEJO_NS}
HERO_FORGEJO_USER=${HERO_FORGEJO_USER}
HERO_FORGEJO_PASS=${HERO_FORGEJO_PASS}
HERO_FORGEJO_REPO=${HERO_FORGEJO_REPO}
HERO_FORGEJO_PF_PORT=${HERO_FORGEJO_PF_PORT}
HERO_INVENTORY_CLONE_DIR=${HERO_INVENTORY_CLONE_DIR}
HERO_GIT_SECRET=${HERO_GIT_SECRET}
FORGEJO_TOKEN=${FORGEJO_TOKEN:-}
FORGEJO_INTERNAL_URL=http://forgejo.${HERO_FORGEJO_NS}.svc.cluster.local:3000
FORGEJO_HOST_URL=http://127.0.0.1:${HERO_FORGEJO_PF_PORT}
GIT_CLONE_URL=http://${HERO_FORGEJO_USER}:\${FORGEJO_TOKEN}@127.0.0.1:${HERO_FORGEJO_PF_PORT}/${HERO_FORGEJO_USER}/${HERO_FORGEJO_REPO}.git
EOF
}

_hero_source_state() {
  # shellcheck disable=SC1090
  [[ -f "$HERO_STATE_FILE" ]] && source "$HERO_STATE_FILE"
}

_hero_start_port_forward() {
  if [[ -f "$HERO_PF_PID_FILE" ]] && kill -0 "$(cat "$HERO_PF_PID_FILE")" 2>/dev/null; then
    _hero_log "Forgejo port-forward already running (pid $(cat "$HERO_PF_PID_FILE"))."
    return 0
  fi

  _hero_log "Starting Forgejo port-forward localhost:${HERO_FORGEJO_PF_PORT}..."
  kubectl port-forward -n "$HERO_FORGEJO_NS" "svc/forgejo" "${HERO_FORGEJO_PF_PORT}:3000" \
    >/tmp/kollect-hero-forgejo-pf.log 2>&1 &
  echo $! >"$HERO_PF_PID_FILE"
  local deadline=$((SECONDS + 30))
  while (( SECONDS < deadline )); do
    if curl -fsS "http://127.0.0.1:${HERO_FORGEJO_PF_PORT}/api/v1/version" >/dev/null 2>&1; then
      _hero_log "Forgejo reachable on localhost:${HERO_FORGEJO_PF_PORT}."
      return 0
    fi
    sleep 1
  done
  _hero_log "Forgejo port-forward failed — see /tmp/kollect-hero-forgejo-pf.log"
  return 1
}

_hero_stop_port_forward() {
  if [[ -f "$HERO_PF_PID_FILE" ]]; then
    kill "$(cat "$HERO_PF_PID_FILE")" 2>/dev/null || true
    rm -f "$HERO_PF_PID_FILE"
  fi
}

_hero_clone_inventory_repo() {
  _hero_source_state
  local clone_url="http://${HERO_FORGEJO_USER}:${FORGEJO_TOKEN}@127.0.0.1:${HERO_FORGEJO_PF_PORT}/${HERO_FORGEJO_USER}/${HERO_FORGEJO_REPO}.git"
  rm -rf "$HERO_INVENTORY_CLONE_DIR"
  git clone --branch main --single-branch "$clone_url" "$HERO_INVENTORY_CLONE_DIR"
}

_hero_wait_inventory_ready() {
  local name="${1:-demo-inventory}"
  _hero_log "Waiting for KollectInventory/${name} Ready..."
  kubectl wait --for=condition=Ready "kollectinventory/${name}" -n default --timeout=180s
}

_hero_wait_git_export() {
  _hero_source_state
  _hero_start_port_forward
  local deadline=$((SECONDS + 240))
  _hero_log "Waiting for first Git export in ${HERO_FORGEJO_REPO}..."
  while (( SECONDS < deadline )); do
    if git -C "$HERO_INVENTORY_CLONE_DIR" pull -q 2>/dev/null; then
      if git -C "$HERO_INVENTORY_CLONE_DIR" rev-parse HEAD >/dev/null 2>&1 \
        && find "$HERO_INVENTORY_CLONE_DIR" -type f \( -name '*.yaml' -o -name '*.yml' -o -name '*.json' \) \
          ! -path '*/.git/*' | grep -q .; then
        _hero_log "First export detected in ${HERO_INVENTORY_CLONE_DIR}."
        return 0
      fi
    else
      _hero_clone_inventory_repo 2>/dev/null || true
    fi
    sleep 5
  done
  _hero_log "Timed out waiting for Git export."
  kubectl get kinv,ksnap -A 2>/dev/null || true
  kubectl logs -n kollect-system -l app.kubernetes.io/name=kollect --tail=60 2>/dev/null || true
  return 1
}
