# ShieldPanel

ShieldPanel is a production-style MVP for self-hosted anti-DDoS, anti-bot, and reverse proxy management on Ubuntu VPS. One VPS can front multiple domains, enforce lightweight WAF rules, serve challenge pages, collect traffic statistics, and expose everything through a bilingual React admin panel.

## Architecture

- Edge: Nginx receives traffic for protected domains and calls the backend through `auth_request` before traffic reaches the origin.
- Backend: Go + Gin exposes the admin API, protection engine, log ingestion, Nginx config generation, migrations, and health checks.
- Frontend: React + Vite + Tailwind renders the admin panel with English and Vietnamese translations.
- Data: PostgreSQL stores configuration, audit data, logs, and rollups. Redis stores rate-limit state, temporary bans, login throttles, and realtime counters.

## Key features

- Multi-domain reverse proxy management with per-domain origin, SSL, HTTPS redirect, challenge mode, and protection mode.
- Score-based anti-bot decisions with bad user-agent detection, header validation, request heuristics, cookie challenges, optional JS challenges, and temporary auto-ban.
- Dashboard, domain detail, statistics, IP blacklist/whitelist, users, settings, and audit logs.
- Ubuntu deployment assets: `install.sh`, `update.sh`, `uninstall.sh`, systemd service, Nginx include files, and Let's Encrypt helpers.

## Folder layout

```text
shieldpanel/
  backend/
    cmd/
    internal/
    migrations/
  frontend/
    src/
  deploy/
    nginx/
    scripts/
    systemd/
  .env.example
  install.sh
  update.sh
  uninstall.sh
```

## Ubuntu install

1. Copy the project to the target VPS.
2. Review `.env.example` and export overrides if you want custom database credentials.
3. Run:

```bash
sudo bash install.sh
```

## One-line bootstrap install

`bootstrap.sh` supports three remote install modes:

- direct package URL with `SHIELDPANEL_PACKAGE_URL`
- GitHub branch archive with `SHIELDPANEL_GITHUB_REPO` + `SHIELDPANEL_GITHUB_REF`
- GitHub release asset with `SHIELDPANEL_GITHUB_REPO` + `SHIELDPANEL_GITHUB_RELEASE_TAG`

If you want the aaPanel-style command using `raw.githubusercontent.com`, put `bootstrap.sh` in your repo and run one of these.

Recommended for production, using a GitHub Release asset:

```bash
URL=https://raw.githubusercontent.com/<owner>/<repo>/main/bootstrap.sh && \
if command -v curl >/dev/null 2>&1; then curl -fsSLo shieldpanel-bootstrap.sh "$URL"; else wget --no-check-certificate -O shieldpanel-bootstrap.sh "$URL"; fi && \
sudo SHIELDPANEL_GITHUB_REPO=<owner>/<repo> SHIELDPANEL_GITHUB_RELEASE_TAG=v0.1.0 bash shieldpanel-bootstrap.sh
```

Simpler for testing, using the current branch archive directly:

```bash
URL=https://raw.githubusercontent.com/<owner>/<repo>/main/bootstrap.sh && \
if command -v curl >/dev/null 2>&1; then curl -fsSLo shieldpanel-bootstrap.sh "$URL"; else wget --no-check-certificate -O shieldpanel-bootstrap.sh "$URL"; fi && \
sudo SHIELDPANEL_GITHUB_REPO=<owner>/<repo> SHIELDPANEL_GITHUB_REF=main bash shieldpanel-bootstrap.sh
```

If you want a custom archive URL, this still works:

```bash
URL=https://raw.githubusercontent.com/<owner>/<repo>/main/bootstrap.sh && \
PACKAGE_URL=https://downloads.example.com/shieldpanel/shieldpanel.tar.gz && \
if command -v curl >/dev/null 2>&1; then curl -fsSLo shieldpanel-bootstrap.sh "$URL"; else wget --no-check-certificate -O shieldpanel-bootstrap.sh "$URL"; fi && \
sudo SHIELDPANEL_PACKAGE_URL="$PACKAGE_URL" bash shieldpanel-bootstrap.sh
```

To build a release archive locally:

```bash
bash deploy/scripts/package_release.sh
```

This behaves like a remote setup entrypoint: it downloads the package and then runs `install.sh`.

The installer will:

- install Nginx, Redis, PostgreSQL, Certbot, Node.js, and Go
- copy the project to `/opt/shieldpanel`
- build the frontend and backend
- run migrations and seed demo data
- install the systemd service
- install the base Nginx include

Default seeded admin credentials:

- Email: `admin@shieldpanel.local`
- Password: `ChangeMe123!`

## Operations

- Health check: `GET /healthz`
- Backend service: `systemctl status shieldpanel`
- Update: `sudo bash update.sh`
- Uninstall: `sudo bash uninstall.sh`
- Migrate manually: `sudo bash /opt/shieldpanel/deploy/scripts/migrate.sh`
- Seed manually: `sudo bash /opt/shieldpanel/deploy/scripts/seed_admin.sh`

## Nginx integration

- Base include: `/etc/nginx/conf.d/shieldpanel-http.conf`
- Generated zones: `/etc/nginx/conf.d/shieldpanel-zones.conf`
- Generated domain files: `/etc/nginx/sites-available/shieldpanel/*.conf`
- Enabled domain files: `/etc/nginx/sites-enabled/shieldpanel/*.conf`

ShieldPanel writes configs, runs `nginx -t`, and only reloads Nginx when validation passes.

## SSL

- Configure `LETSENCRYPT_EMAIL` in `/etc/shieldpanel/shieldpanel.env`
- Issue a certificate from the panel or run:

```bash
sudo bash /opt/shieldpanel/deploy/scripts/issue_cert.sh example.com /var/www/shieldpanel/acme
```

- Renew:

```bash
sudo bash /opt/shieldpanel/deploy/scripts/renew_cert.sh example.com /var/www/shieldpanel/acme
```

## Notes

- Country statistics use `CF-IPCountry` when traffic comes through Cloudflare or another trusted proxy that sets compatible headers.
- The backend entrypoints were validated with `go build`, and the frontend build was validated with `npm run build`.
