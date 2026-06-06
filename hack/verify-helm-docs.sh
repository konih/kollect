#!/usr/bin/env bash
# Fail if charts/kollect/README.md is stale relative to values.yaml and README.md.gotmpl.
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "${ROOT}"

HELM_DOCS="${HELM_DOCS_BIN:-${ROOT}/bin/helm-docs}"
if [[ ! -x "${HELM_DOCS}" ]]; then
  echo "verify-helm-docs: installing helm-docs into bin/" >&2
  bash hack/install-helm-docs.sh v1.14.2 bin/helm-docs
  HELM_DOCS="${ROOT}/bin/helm-docs"
fi

scratch="$(mktemp -d)"
trap 'rm -rf "${scratch}"' EXIT

cp -a charts/kollect/README.md "${scratch}/README.md"

"${HELM_DOCS}" --chart-search-root charts/kollect

if ! diff -u "${scratch}/README.md" charts/kollect/README.md; then
  echo "verify-helm-docs: charts/kollect/README.md drift — run 'task helm-docs' and commit" >&2
  exit 1
fi

echo "verify-helm-docs: ok"
