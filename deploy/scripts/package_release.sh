#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
PROJECT_DIR="$(dirname "${ROOT_DIR}")"
PROJECT_NAME="$(basename "${ROOT_DIR}")"
OUT_DIR="${1:-${ROOT_DIR}/dist-release}"
ARCHIVE_NAME="${2:-shieldpanel.tar.gz}"
ARCHIVE_PATH="${OUT_DIR}/${ARCHIVE_NAME}"

mkdir -p "${OUT_DIR}"

tar \
  --exclude='.git' \
  --exclude='frontend/node_modules' \
  --exclude='frontend/dist' \
  --exclude='dist-release' \
  --exclude='backend/*.exe' \
  -czf "${ARCHIVE_PATH}" \
  -C "${PROJECT_DIR}" \
  "${PROJECT_NAME}"

printf '%s\n' "${ARCHIVE_PATH}"
