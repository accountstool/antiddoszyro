#!/usr/bin/env bash
set -euo pipefail

APP_ROOT="${APP_ROOT:-/opt/shieldpanel}"
APP_USER="${APP_USER:-shieldpanel}"
APP_GROUP="${APP_GROUP:-shieldpanel}"
ENV_DIR="/etc/shieldpanel"
ENV_FILE="${ENV_FILE:-${ENV_DIR}/shieldpanel.env}"
REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
GO_VERSION="${GO_VERSION:-1.22.12}"
NODE_MAJOR="${NODE_MAJOR:-22}"
DB_NAME="${DB_NAME:-shieldpanel}"
DB_USER="${DB_USER:-shieldpanel}"
DB_PASSWORD="${DB_PASSWORD:-change-me-db-password}"

require_root() {
  if [[ "${EUID}" -ne 0 ]]; then
    echo "Run install.sh as root." >&2
    exit 1
  fi
}

install_base_packages() {
  apt-get update
  apt-get install -y ca-certificates curl gnupg lsb-release rsync unzip build-essential nginx redis-server postgresql postgresql-contrib certbot
}

install_node() {
  if command -v node >/dev/null 2>&1; then
    local major
    major="$(node -p 'process.versions.node.split(`.`)[0]')"
    if [[ "${major}" -ge 18 ]]; then
      return
    fi
  fi
  curl -fsSL "https://deb.nodesource.com/setup_${NODE_MAJOR}.x" | bash -
  apt-get install -y nodejs
}

install_go() {
  if command -v go >/dev/null 2>&1; then
    return
  fi
  curl -fsSL "https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz" -o /tmp/go.tar.gz
  rm -rf /usr/local/go
  tar -C /usr/local -xzf /tmp/go.tar.gz
  ln -sf /usr/local/go/bin/go /usr/local/bin/go
}

ensure_user() {
  if ! getent group "${APP_GROUP}" >/dev/null 2>&1; then
    groupadd --system "${APP_GROUP}"
  fi
  if ! id -u "${APP_USER}" >/dev/null 2>&1; then
    useradd --system --home "${APP_ROOT}" --shell /usr/sbin/nologin --gid "${APP_GROUP}" "${APP_USER}"
  fi
}

prepare_paths() {
  mkdir -p "${APP_ROOT}" "${APP_ROOT}/bin" "${ENV_DIR}" /var/log/shieldpanel /var/www/shieldpanel/acme /etc/nginx/shieldpanel/sites-available /etc/nginx/shieldpanel/sites-enabled
  rm -rf /etc/nginx/sites-enabled/shieldpanel
  chown -R "${APP_USER}:${APP_GROUP}" /var/log/shieldpanel /var/www/shieldpanel
}

sync_project() {
  rsync -a --delete \
    --exclude ".git" \
    --exclude "frontend/node_modules" \
    --exclude "frontend/dist" \
    "${REPO_ROOT}/" "${APP_ROOT}/"
}

create_env_file() {
  if [[ ! -f "${ENV_FILE}" ]]; then
    install -m 640 "${APP_ROOT}/.env.example" "${ENV_FILE}"
    sed -i "s|change-me-db-password|${DB_PASSWORD}|g" "${ENV_FILE}"
    sed -i "s|http://your-vps-ip:8080|http://$(hostname -I | awk '{print $1}'):8080|g" "${ENV_FILE}"
  fi
}

migrate_env_file() {
  sed -i 's|^NGINX_SITES_AVAILABLE=/etc/nginx/sites-available/shieldpanel$|NGINX_SITES_AVAILABLE=/etc/nginx/shieldpanel/sites-available|' "${ENV_FILE}"
  sed -i 's|^NGINX_SITES_ENABLED=/etc/nginx/sites-enabled/shieldpanel$|NGINX_SITES_ENABLED=/etc/nginx/shieldpanel/sites-enabled|' "${ENV_FILE}"
}

configure_database() {
  runuser -u postgres -- psql -tc "select 1 from pg_roles where rolname='${DB_USER}'" | grep -q 1 || runuser -u postgres -- psql -c "create user ${DB_USER} with password '${DB_PASSWORD}';"
  runuser -u postgres -- psql -tc "select 1 from pg_database where datname='${DB_NAME}'" | grep -q 1 || runuser -u postgres -- psql -c "create database ${DB_NAME} owner ${DB_USER};"
}

build_frontend() {
  pushd "${APP_ROOT}/frontend" >/dev/null
  npm install
  npm run build
  popd >/dev/null
}

build_backend() {
  export PATH="/usr/local/go/bin:${PATH}"
  pushd "${APP_ROOT}/backend" >/dev/null
  go mod tidy
  go build -o "${APP_ROOT}/bin/shieldpanel" ./cmd/server
  go build -o "${APP_ROOT}/bin/shieldpanel-migrate" ./cmd/migrate
  go build -o "${APP_ROOT}/bin/shieldpanel-seed" ./cmd/seed
  popd >/dev/null
}

run_database_tasks() {
  set -a
  source "${ENV_FILE}"
  set +a
  "${APP_ROOT}/bin/shieldpanel-migrate"
  "${APP_ROOT}/bin/shieldpanel-seed"
}

install_systemd() {
  install -m 644 "${APP_ROOT}/deploy/systemd/shieldpanel.service" /etc/systemd/system/shieldpanel.service
  systemctl daemon-reload
  systemctl enable shieldpanel.service
}

install_nginx_base() {
  install -m 644 "${APP_ROOT}/deploy/nginx/includes/shieldpanel-http.conf" /etc/nginx/conf.d/shieldpanel-http.conf
  touch /etc/nginx/conf.d/shieldpanel-zones.conf
  nginx -t
  systemctl enable nginx redis-server postgresql
  systemctl restart nginx redis-server postgresql
}

start_services() {
  chown -R "${APP_USER}:${APP_GROUP}" "${APP_ROOT}"
  systemctl restart shieldpanel.service
}

print_summary() {
  echo
  echo "ShieldPanel installed."
  echo "Panel URL: $(grep '^PUBLIC_URL=' "${ENV_FILE}" | cut -d= -f2-)"
  echo "Seeded admin: admin@shieldpanel.local"
  echo "Seeded password: ChangeMe123!"
  echo
  echo "Recommended firewall rules:"
  echo "  ufw allow 22/tcp"
  echo "  ufw allow 80/tcp"
  echo "  ufw allow 443/tcp"
  echo "  ufw allow 8080/tcp"
  echo
}

require_root
install_base_packages
install_node
install_go
ensure_user
prepare_paths
sync_project
create_env_file
migrate_env_file
configure_database
build_frontend
build_backend
run_database_tasks
install_systemd
install_nginx_base
start_services
print_summary
