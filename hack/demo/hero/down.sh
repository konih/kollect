#!/usr/bin/env bash
# Tear down the hero demo kind cluster and local port-forward.
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=lib.sh
source "${SCRIPT_DIR}/lib.sh"

_hero_require kind
_hero_stop_port_forward
kind_delete_cluster "$HERO_CLUSTER"
rm -f "$HERO_STATE_FILE" /tmp/kollect-hero-forgejo-pf.log
rm -rf "$HERO_INVENTORY_CLONE_DIR"
_hero_log "Hero demo torn down."
