#!/usr/bin/env bash
set -euo pipefail

APP_ROOT="${APP_ROOT:-/opt/shieldpanel}"
ENV_FILE="${ENV_FILE:-/etc/shieldpanel/shieldpanel.env}"
REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

migrate_env_file() {
  if [[ -f "${ENV_FILE}" ]]; then
    sed -i 's|^NGINX_SITES_AVAILABLE=/etc/nginx/sites-available/shieldpanel$|NGINX_SITES_AVAILABLE=/etc/nginx/shieldpanel/sites-available|' "${ENV_FILE}"
    sed -i 's|^NGINX_SITES_ENABLED=/etc/nginx/sites-enabled/shieldpanel$|NGINX_SITES_ENABLED=/etc/nginx/shieldpanel/sites-enabled|' "${ENV_FILE}"
  fi
}

prepare_paths() {
  mkdir -p /var/log/shieldpanel /var/www/shieldpanel/acme /etc/nginx/shieldpanel/sites-available /etc/nginx/shieldpanel/sites-enabled
  rm -rf /etc/nginx/sites-enabled/shieldpanel
  touch /etc/nginx/conf.d/shieldpanel-zones.conf
  chown -R shieldpanel:shieldpanel /var/log/shieldpanel /var/www/shieldpanel /etc/nginx/shieldpanel
  chown shieldpanel:shieldpanel /etc/nginx/conf.d/shieldpanel-zones.conf
}

install_sudoers() {
  cat >/etc/sudoers.d/shieldpanel-nginx <<'EOF'
shieldpanel ALL=(root) NOPASSWD: /usr/sbin/nginx -t -c /etc/nginx/nginx.conf, /usr/sbin/nginx -s reload
EOF
  chmod 440 /etc/sudoers.d/shieldpanel-nginx
}

install_nginx_base() {
  install -m 644 "${APP_ROOT}/deploy/nginx/includes/shieldpanel-http.conf" /etc/nginx/conf.d/shieldpanel-http.conf
  touch /etc/nginx/conf.d/shieldpanel-zones.conf
  chown shieldpanel:shieldpanel /etc/nginx/conf.d/shieldpanel-zones.conf
  nginx -t
  systemctl restart nginx
}

if [[ "${EUID}" -ne 0 ]]; then
  echo "Run update.sh as root." >&2
  exit 1
fi

rsync -a --delete \
  --exclude ".git" \
  --exclude "frontend/node_modules" \
  --exclude "frontend/dist" \
  "${REPO_ROOT}/" "${APP_ROOT}/"

migrate_env_file
prepare_paths
install_sudoers

pushd "${APP_ROOT}/frontend" >/dev/null
echo "[ShieldPanel] Installing frontend dependencies..."
npm install
echo "[ShieldPanel] Building frontend bundle..."
npm run build
popd >/dev/null

export PATH="/usr/local/go/bin:${PATH}"
pushd "${APP_ROOT}/backend" >/dev/null
echo "[ShieldPanel] Building backend binaries..."
go mod tidy
go build -o "${APP_ROOT}/bin/shieldpanel" ./cmd/server
go build -o "${APP_ROOT}/bin/shieldpanel-migrate" ./cmd/migrate
popd >/dev/null

set -a
source "${ENV_FILE}"
set +a
echo "[ShieldPanel] Running database migrations..."
"${APP_ROOT}/bin/shieldpanel-migrate"
install_nginx_base
systemctl restart shieldpanel.service
echo "ShieldPanel updated."
