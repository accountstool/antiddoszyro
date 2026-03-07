#!/usr/bin/env bash
set -euo pipefail

DOMAIN="${1:?domain is required}"
WEBROOT="${2:-/var/www/shieldpanel/acme}"
EMAIL="${LETSENCRYPT_EMAIL:?LETSENCRYPT_EMAIL must be set in the environment}"

mkdir -p "${WEBROOT}"
certbot certonly \
  --non-interactive \
  --agree-tos \
  --email "${EMAIL}" \
  --webroot \
  -w "${WEBROOT}" \
  -d "${DOMAIN}" \
  --force-renewal
systemctl reload nginx
