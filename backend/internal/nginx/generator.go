package nginx

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"shieldpanel/backend/internal/config"
	"shieldpanel/backend/internal/domain"
)

var safeNamePattern = regexp.MustCompile(`[^a-zA-Z0-9]+`)

type Manager struct {
	cfg       config.NginxConfig
	logger    *slog.Logger
	siteTmpl  *template.Template
	zonesTmpl *template.Template
}

type siteTemplateData struct {
	Domain         domain.Domain
	ZoneName       string
	ConnZoneName   string
	OriginURL      string
	PanelUpstream  string
	ACMEWebroot    string
	VerifyRoute    string
	ChallengeRoute string
	BlockRoute     string
	AccessLog      string
	ErrorLog       string
	HasCertificate bool
	CertPath       string
	KeyPath        string
}

func NewManager(cfg config.NginxConfig, logger *slog.Logger) *Manager {
	return &Manager{
		cfg:       cfg,
		logger:    logger,
		siteTmpl:  template.Must(template.New("site").Parse(siteTemplate)),
		zonesTmpl: template.Must(template.New("zones").Funcs(template.FuncMap{"safeName": safeName}).Parse(zonesTemplate)),
	}
}

func (m *Manager) Sync(ctx context.Context, domains []domain.Domain) error {
	if !m.cfg.Enabled {
		return nil
	}
	if err := os.MkdirAll(m.cfg.SitesAvailable, 0o755); err != nil {
		return err
	}
	if err := os.MkdirAll(m.cfg.SitesEnabled, 0o755); err != nil {
		return err
	}

	previousFiles := map[string][]byte{}
	if err := m.writeZones(domains, previousFiles); err != nil {
		return err
	}
	expected := map[string]struct{}{}
	for _, item := range domains {
		path := filepath.Join(m.cfg.SitesAvailable, item.Name+".conf")
		expected[path] = struct{}{}
		if err := m.writeSite(item, previousFiles); err != nil {
			m.rollback(previousFiles)
			return err
		}
		enabledPath := filepath.Join(m.cfg.SitesEnabled, item.Name+".conf")
		expected[enabledPath] = struct{}{}
		if err := m.copyFile(path, enabledPath, previousFiles); err != nil {
			m.rollback(previousFiles)
			return err
		}
	}

	if err := m.removeStaleFiles(expected, previousFiles); err != nil {
		m.rollback(previousFiles)
		return err
	}

	if err := m.validate(ctx); err != nil {
		m.rollback(previousFiles)
		return err
	}
	if m.cfg.ReloadOnChange {
		if err := m.reload(ctx); err != nil {
			m.rollback(previousFiles)
			return err
		}
	}
	return nil
}

func (m *Manager) writeZones(domains []domain.Domain, previous map[string][]byte) error {
	path := m.cfg.ZonesPath
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	var buf bytes.Buffer
	if err := m.zonesTmpl.Execute(&buf, map[string]any{"Domains": domains}); err != nil {
		return err
	}
	return m.writeFile(path, buf.Bytes(), previous)
}

func (m *Manager) writeSite(item domain.Domain, previous map[string][]byte) error {
	path := filepath.Join(m.cfg.SitesAvailable, item.Name+".conf")
	var buf bytes.Buffer
	certPath := filepath.Join("/etc/letsencrypt/live", item.Name, "fullchain.pem")
	keyPath := filepath.Join("/etc/letsencrypt/live", item.Name, "privkey.pem")
	hasCertificate := fileExists(certPath) && fileExists(keyPath)
	if item.SSLEnabled && !hasCertificate {
		m.logger.Warn("ssl enabled for domain without certificate, falling back to http only", "domain", item.Name, "cert_path", certPath)
	}
	data := siteTemplateData{
		Domain:         item,
		ZoneName:       "sp_" + safeName(item.Name),
		ConnZoneName:   "spc_" + safeName(item.Name),
		OriginURL:      fmt.Sprintf("%s://%s:%d", item.OriginProtocol, item.OriginHost, item.OriginPort),
		PanelUpstream:  strings.TrimSuffix(m.cfg.PanelUpstreamURL, "/"),
		ACMEWebroot:    m.cfg.ACMEWebroot,
		VerifyRoute:    "/__shieldpanel_verify",
		ChallengeRoute: "/__shieldpanel_challenge",
		BlockRoute:     "/__shieldpanel_block",
		AccessLog:      fmt.Sprintf("/var/log/nginx/shieldpanel-%s.access.log", safeName(item.Name)),
		ErrorLog:       fmt.Sprintf("/var/log/nginx/shieldpanel-%s.error.log", safeName(item.Name)),
		HasCertificate: hasCertificate,
		CertPath:       certPath,
		KeyPath:        keyPath,
	}
	if err := m.siteTmpl.Execute(&buf, data); err != nil {
		return err
	}
	return m.writeFile(path, buf.Bytes(), previous)
}

func (m *Manager) copyFile(source string, target string, previous map[string][]byte) error {
	data, err := os.ReadFile(source)
	if err != nil {
		return err
	}
	return m.writeFile(target, data, previous)
}

func (m *Manager) writeFile(path string, data []byte, previous map[string][]byte) error {
	if _, ok := previous[path]; !ok {
		if existing, err := os.ReadFile(path); err == nil {
			previous[path] = existing
		} else if os.IsNotExist(err) {
			previous[path] = nil
		} else {
			return err
		}
	}
	return os.WriteFile(path, data, 0o644)
}

func (m *Manager) removeStaleFiles(expected map[string]struct{}, previous map[string][]byte) error {
	paths := []string{m.cfg.SitesAvailable, m.cfg.SitesEnabled}
	for _, base := range paths {
		entries, err := os.ReadDir(base)
		if err != nil {
			return err
		}
		for _, entry := range entries {
			if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".conf") {
				continue
			}
			fullPath := filepath.Join(base, entry.Name())
			if _, ok := expected[fullPath]; ok {
				continue
			}
			if _, tracked := previous[fullPath]; !tracked {
				existing, _ := os.ReadFile(fullPath)
				previous[fullPath] = existing
			}
			if err := os.Remove(fullPath); err != nil {
				return err
			}
		}
	}
	return nil
}

func (m *Manager) rollback(previous map[string][]byte) {
	for path, data := range previous {
		if data == nil {
			_ = os.Remove(path)
			continue
		}
		_ = os.WriteFile(path, data, 0o644)
	}
}

func (m *Manager) validate(ctx context.Context) error {
	cmd, err := m.execWithPrivileges(ctx, "-t", "-c", m.cfg.ConfigPath)
	if err != nil {
		return err
	}
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("nginx validation failed: %w: %s", err, string(output))
	}
	return nil
}

func (m *Manager) reload(ctx context.Context) error {
	cmd, err := m.execWithPrivileges(ctx, "-s", "reload")
	if err != nil {
		return err
	}
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("nginx reload failed: %w: %s", err, string(output))
	}
	return nil
}

func (m *Manager) execWithPrivileges(ctx context.Context, args ...string) (*exec.Cmd, error) {
	if os.Geteuid() == 0 {
		return exec.CommandContext(ctx, m.cfg.BinaryPath, args...), nil
	}
	if _, err := exec.LookPath("sudo"); err != nil {
		return nil, fmt.Errorf("nginx command requires elevated privileges and sudo is not available: %w", err)
	}
	sudoArgs := append([]string{m.cfg.BinaryPath}, args...)
	return exec.CommandContext(ctx, "sudo", sudoArgs...), nil
}

func safeName(input string) string {
	safe := strings.Trim(safeNamePattern.ReplaceAllString(strings.ToLower(input), "_"), "_")
	if safe == "" {
		return "default"
	}
	return safe
}

func fileExists(path string) bool {
	info, err := os.Lstat(path)
	return err == nil && !info.IsDir()
}

const zonesTemplate = `
{{- range .Domains }}
limit_req_zone $binary_remote_addr zone=sp_{{ safeName .Name }}:20m rate={{ .RateLimitRPS }}r/s;
limit_conn_zone $binary_remote_addr zone=spc_{{ safeName .Name }}:20m;
{{- end }}
`

const siteTemplate = `
server {
    listen 80;
    server_name {{ .Domain.Name }};
    access_log {{ .AccessLog }};
    error_log {{ .ErrorLog }} warn;
    client_max_body_size 64m;

    location ^~ /.well-known/acme-challenge/ {
        root {{ .ACMEWebroot }};
    }

    location = {{ .VerifyRoute }} {
        proxy_pass {{ .PanelUpstream }}/internal/protection/verify;
        proxy_set_header X-Original-Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header CF-Connecting-IP $http_cf_connecting_ip;
        proxy_set_header CF-IPCountry $http_cf_ipcountry;
    }

{{- if and .HasCertificate .Domain.ForceHTTPS }}
    location / {
        return 301 https://$host$request_uri;
    }
{{- else }}
    {{ template "protected" . }}
{{- end }}
}

{{- if .HasCertificate }}
server {
    listen 443 ssl http2;
    server_name {{ .Domain.Name }};
    access_log {{ .AccessLog }};
    error_log {{ .ErrorLog }} warn;
    client_max_body_size 64m;

    ssl_certificate {{ .CertPath }};
    ssl_certificate_key {{ .KeyPath }};
    ssl_session_timeout 1d;
    ssl_session_cache shared:SSL:10m;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;

    location = {{ .VerifyRoute }} {
        proxy_pass {{ .PanelUpstream }}/internal/protection/verify;
        proxy_set_header X-Original-Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header CF-Connecting-IP $http_cf_connecting_ip;
        proxy_set_header CF-IPCountry $http_cf_ipcountry;
    }

    {{ template "protected" . }}
}
{{- end }}

{{- define "protected" }}
location / {
    limit_req zone={{ .ZoneName }} burst={{ .Domain.RateLimitBurst }} nodelay;
    limit_conn {{ .ConnZoneName }} 80;

    auth_request /__shieldpanel_check;
    auth_request_set $shieldpanel_set_cookie $upstream_http_set_cookie;
    auth_request_set $shieldpanel_reason $upstream_http_x_shieldpanel_reason;
    auth_request_set $shieldpanel_challenge $upstream_http_x_shieldpanel_challenge;
    add_header Set-Cookie $shieldpanel_set_cookie always;
    error_page 401 = {{ .ChallengeRoute }};
    error_page 403 = {{ .BlockRoute }};

    proxy_http_version 1.1;
    proxy_set_header Host $host;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_read_timeout 60s;
    proxy_connect_timeout 10s;
    proxy_pass {{ .OriginURL }};
}

location = /__shieldpanel_check {
    internal;
    proxy_pass {{ .PanelUpstream }}/internal/protection/check;
    proxy_pass_request_body off;
    proxy_set_header Content-Length "";
    proxy_set_header X-Original-Host $host;
    proxy_set_header X-Original-URI $request_uri;
    proxy_set_header X-Original-Path $uri;
    proxy_set_header X-Original-Query $query_string;
    proxy_set_header X-Original-Method $request_method;
    proxy_set_header X-Forwarded-Proto $scheme;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header User-Agent $http_user_agent;
    proxy_set_header Accept $http_accept;
    proxy_set_header Accept-Language $http_accept_language;
    proxy_set_header Accept-Encoding $http_accept_encoding;
    proxy_set_header Sec-Fetch-Site $http_sec_fetch_site;
    proxy_set_header Sec-Fetch-Mode $http_sec_fetch_mode;
    proxy_set_header Sec-Ch-Ua $http_sec_ch_ua;
    proxy_set_header CF-Connecting-IP $http_cf_connecting_ip;
    proxy_set_header CF-IPCountry $http_cf_ipcountry;
}

location = {{ .ChallengeRoute }} {
    internal;
    proxy_pass {{ .PanelUpstream }}/internal/protection/challenge;
    proxy_set_header X-Original-Host $host;
    proxy_set_header X-Original-URI $request_uri;
    proxy_set_header X-Original-Method $request_method;
    proxy_set_header X-Forwarded-Proto $scheme;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Challenge-Mode $shieldpanel_challenge;
    proxy_set_header X-Reason $shieldpanel_reason;
    proxy_set_header Accept-Language $http_accept_language;
}

location = {{ .BlockRoute }} {
    internal;
    proxy_pass {{ .PanelUpstream }}/internal/protection/block;
    proxy_set_header X-Original-Host $host;
    proxy_set_header X-Original-URI $request_uri;
    proxy_set_header X-Original-Method $request_method;
    proxy_set_header X-Forwarded-Proto $scheme;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Reason $shieldpanel_reason;
    proxy_set_header Accept-Language $http_accept_language;
}
{{- end }}
`
