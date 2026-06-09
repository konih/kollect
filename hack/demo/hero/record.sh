#!/usr/bin/env bash
# Record a VHS tape into docs/assets/demo/.
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=lib.sh
source "${SCRIPT_DIR}/lib.sh"

TAPE="${1:?tape file name required (e.g. demo-git-only.tape)}"
TAPE_PATH="${SCRIPT_DIR}/${TAPE}"

if [[ ! -f "$TAPE_PATH" ]]; then
  echo "Tape not found: ${TAPE_PATH}" >&2
  exit 1
fi

if ! command -v vhs >/dev/null 2>&1; then
  cat >&2 <<'EOF'
vhs is not installed — install before recording:

  go install github.com/charmbracelet/vhs@latest

Or: brew install vhs

The checked-in .tape files are the source of truth; run this script locally to generate GIF/MP4.
EOF
  exit 1
fi

_hero_require_tools
bash "${SCRIPT_DIR}/preflight.sh"

_hero_log "Recording ${TAPE} (cwd=${REPO_ROOT})..."
cd "$REPO_ROOT"
vhs "$TAPE_PATH"
_hero_log "Artifacts written under docs/assets/demo/ (see tape Output directives)."
