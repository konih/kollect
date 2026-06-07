#!/usr/bin/env bash
# Build or push a maintainer-only controller image (local / test-* tags — not for production).
set -euo pipefail

MODE="${1:?usage: $0 build|push}"

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "${REPO_ROOT}"

IMAGE_REPO="${IMAGE_REPO:-ghcr.io/konih/kollect}"
PLATFORMS="${PLATFORMS:-linux/amd64}"
CONTAINER_TOOL="${CONTAINER_TOOL:-docker}"

IMAGE_TAG="${IMAGE_TAG:-}"
if [[ -z "${IMAGE_TAG}" ]]; then
  if [[ "${MODE}" == "push" ]]; then
    IMAGE_TAG="test-$(git rev-parse --short HEAD)"
  else
    IMAGE_TAG="local"
  fi
fi

if [[ "${IMAGE_TAG}" =~ ^v?[0-9]+\.[0-9]+\.[0-9]+(-.+)?$ ]]; then
  echo "Refusing semver-like tag '${IMAGE_TAG}' — use maintainer tags (local, test-*, dev-*)" >&2
  exit 1
fi

IMAGE="${IMAGE_REPO}:${IMAGE_TAG}"

echo "[docker-local] ${MODE} ${IMAGE} (platform=${PLATFORMS}; maintainer-only — not for production)"

if [[ "${PLATFORMS}" == *,* ]]; then
  if [[ "${MODE}" == "build" ]]; then
    echo "Multi-arch PLATFORMS requires push mode; use PLATFORMS=linux/amd64 for kind/minikube load." >&2
    exit 1
  fi
  "${CONTAINER_TOOL}" buildx build --push --platform "${PLATFORMS}" -t "${IMAGE}" .
  exit 0
fi

"${CONTAINER_TOOL}" build --platform "${PLATFORMS}" -t "${IMAGE}" .
if [[ "${MODE}" == "push" ]]; then
  "${CONTAINER_TOOL}" push "${IMAGE}"
fi
