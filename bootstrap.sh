#!/usr/bin/env bash
set -euo pipefail

PACKAGE_URL="${SHIELDPANEL_PACKAGE_URL:-}"
GITHUB_REPO="${SHIELDPANEL_GITHUB_REPO:-}"
GITHUB_REF="${SHIELDPANEL_GITHUB_REF:-main}"
GITHUB_RELEASE_TAG="${SHIELDPANEL_GITHUB_RELEASE_TAG:-}"
GITHUB_ASSET_NAME="${SHIELDPANEL_GITHUB_ASSET_NAME:-shieldpanel.tar.gz}"
WORK_DIR="${SHIELDPANEL_BOOTSTRAP_DIR:-/usr/local/src/shieldpanel-bootstrap}"
ARCHIVE_PATH="${WORK_DIR}/shieldpanel.tar.gz"
EXTRACT_DIR="${WORK_DIR}/package"

require_root() {
  if [[ "${EUID}" -ne 0 ]]; then
    echo "Run bootstrap.sh as root." >&2
    exit 1
  fi
}

resolve_package_url() {
  if [[ -n "${PACKAGE_URL}" ]]; then
    return
  fi

  if [[ -z "${GITHUB_REPO}" ]]; then
    cat >&2 <<'EOF'
Provide one of these variables before running bootstrap.sh:

  1. SHIELDPANEL_PACKAGE_URL
  2. SHIELDPANEL_GITHUB_REPO + SHIELDPANEL_GITHUB_REF
  3. SHIELDPANEL_GITHUB_REPO + SHIELDPANEL_GITHUB_RELEASE_TAG

Examples:
  SHIELDPANEL_PACKAGE_URL=https://downloads.example.com/shieldpanel.tar.gz bash bootstrap.sh
  SHIELDPANEL_GITHUB_REPO=myorg/shieldpanel SHIELDPANEL_GITHUB_REF=main bash bootstrap.sh
  SHIELDPANEL_GITHUB_REPO=myorg/shieldpanel SHIELDPANEL_GITHUB_RELEASE_TAG=v1.0.0 bash bootstrap.sh
EOF
    exit 1
  fi

  if [[ -n "${GITHUB_RELEASE_TAG}" ]]; then
    PACKAGE_URL="https://github.com/${GITHUB_REPO}/releases/download/${GITHUB_RELEASE_TAG}/${GITHUB_ASSET_NAME}"
    return
  fi

  PACKAGE_URL="https://github.com/${GITHUB_REPO}/archive/refs/heads/${GITHUB_REF}.tar.gz"
}

download_file() {
  local url="$1"
  local output="$2"
  if command -v curl >/dev/null 2>&1; then
    curl -fsSL "${url}" -o "${output}"
  elif command -v wget >/dev/null 2>&1; then
    wget --no-check-certificate -O "${output}" "${url}"
  else
    apt-get update
    apt-get install -y curl
    curl -fsSL "${url}" -o "${output}"
  fi
}

prepare_dirs() {
  mkdir -p "${WORK_DIR}"
  rm -rf "${EXTRACT_DIR}"
  mkdir -p "${EXTRACT_DIR}"
}

extract_package() {
  tar -xzf "${ARCHIVE_PATH}" -C "${EXTRACT_DIR}" --strip-components=1
}

run_install() {
  cd "${EXTRACT_DIR}"
  bash install.sh
}

require_root
resolve_package_url
prepare_dirs
download_file "${PACKAGE_URL}" "${ARCHIVE_PATH}"
extract_package
run_install
