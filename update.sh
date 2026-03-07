#!/usr/bin/env bash
set -euo pipefail

APP_ROOT="${APP_ROOT:-/opt/shieldpanel}"
ENV_FILE="${ENV_FILE:-/etc/shieldpanel/shieldpanel.env}"
REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

if [[ "${EUID}" -ne 0 ]]; then
  echo "Run update.sh as root." >&2
  exit 1
fi

rsync -a --delete \
  --exclude ".git" \
  --exclude "frontend/node_modules" \
  --exclude "frontend/dist" \
  "${REPO_ROOT}/" "${APP_ROOT}/"

pushd "${APP_ROOT}/frontend" >/dev/null
npm install
npm run build
popd >/dev/null

export PATH="/usr/local/go/bin:${PATH}"
pushd "${APP_ROOT}/backend" >/dev/null
go mod tidy
go build -o "${APP_ROOT}/bin/shieldpanel" ./cmd/server
go build -o "${APP_ROOT}/bin/shieldpanel-migrate" ./cmd/migrate
popd >/dev/null

set -a
source "${ENV_FILE}"
set +a
"${APP_ROOT}/bin/shieldpanel-migrate"
systemctl restart shieldpanel.service
echo "ShieldPanel updated."
