#!/usr/bin/env bash
set -euo pipefail

APP_ROOT="${APP_ROOT:-/opt/shieldpanel}"
ENV_DIR="/etc/shieldpanel"
DROP_DATABASE="${DROP_DATABASE:-false}"
DB_NAME="${DB_NAME:-shieldpanel}"
DB_USER="${DB_USER:-shieldpanel}"

if [[ "${EUID}" -ne 0 ]]; then
  echo "Run uninstall.sh as root." >&2
  exit 1
fi

systemctl disable --now shieldpanel.service || true
rm -f /etc/systemd/system/shieldpanel.service
systemctl daemon-reload

rm -rf "${APP_ROOT}" "${ENV_DIR}"
rm -f /etc/nginx/conf.d/shieldpanel-http.conf /etc/nginx/conf.d/shieldpanel-zones.conf
rm -rf /etc/nginx/shieldpanel /etc/nginx/sites-enabled/shieldpanel
systemctl reload nginx || true

if [[ "${DROP_DATABASE}" == "true" ]]; then
  runuser -u postgres -- psql -c "drop database if exists ${DB_NAME};" || true
  runuser -u postgres -- psql -c "drop user if exists ${DB_USER};" || true
fi

echo "ShieldPanel removed."
