package config

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	App      AppConfig
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	Auth     AuthConfig
	Security SecurityConfig
	Nginx    NginxConfig
	Defaults DefaultsConfig
	Paths    PathsConfig
}

type AppConfig struct {
	Name    string
	Version string
	Env     string
}

type ServerConfig struct {
	Host            string
	Port            int
	PublicURL       string
	PanelPort       int
	FrontendDistDir string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
}

type DatabaseConfig struct {
	URL          string
	MaxConns     int32
	MinConns     int32
	MaxConnIdle  time.Duration
	MaxConnLife  time.Duration
	HealthPeriod time.Duration
}

type RedisConfig struct {
	Addr         string
	Password     string
	DB           int
	KeyPrefix    string
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type AuthConfig struct {
	SessionCookieName string
	CSRFCookieName    string
	SessionTTL        time.Duration
	RememberMeTTL     time.Duration
	CookieSecure      bool
	CookieDomain      string
	BcryptCost        int
}

type SecurityConfig struct {
	SessionSecret       string
	ChallengeSecret     string
	AllowedMethods      []string
	ChallengeTTL        time.Duration
	ChallengeGraceTTL   time.Duration
	AutoBanThreshold    int
	AutoBanWindow       time.Duration
	AutoBanTTL          time.Duration
	TrustedProxyMode    bool
	TrustedProxyHeaders []string
}

type NginxConfig struct {
	Enabled          bool
	BinaryPath       string
	ConfigPath       string
	SitesAvailable   string
	SitesEnabled     string
	ZonesPath        string
	TemplatesDir     string
	ACMEWebroot      string
	PanelUpstreamURL string
	ReloadOnChange   bool
}

type DefaultsConfig struct {
	Language             string
	DefaultProtection    string
	DefaultRateLimitRPS  int
	DefaultRateLimitBurst int
	DefaultChallengeMode string
}

type PathsConfig struct {
	MigrationsDir string
	DataDir       string
	LogDir        string
	IssueCertScript string
	RenewCertScript string
}

func Load() Config {
	root := env("SHIELDPANEL_ROOT", "/opt/shieldpanel")
	return Config{
		App: AppConfig{
			Name:    env("APP_NAME", "ShieldPanel"),
			Version: env("APP_VERSION", "0.1.0"),
			Env:     env("APP_ENV", "production"),
		},
		Server: ServerConfig{
			Host:            env("SERVER_HOST", "0.0.0.0"),
			Port:            envInt("SERVER_PORT", 8080),
			PanelPort:       envInt("PANEL_PORT", 8080),
			PublicURL:       env("PUBLIC_URL", "http://127.0.0.1:8080"),
			FrontendDistDir: env("FRONTEND_DIST_DIR", filepath.Join(root, "frontend", "dist")),
			ReadTimeout:     envDuration("SERVER_READ_TIMEOUT", 15*time.Second),
			WriteTimeout:    envDuration("SERVER_WRITE_TIMEOUT", 30*time.Second),
		},
		Database: DatabaseConfig{
			URL:          env("DATABASE_URL", "postgres://shieldpanel:shieldpanel@127.0.0.1:5432/shieldpanel?sslmode=disable"),
			MaxConns:     int32(envInt("DB_MAX_CONNS", 20)),
			MinConns:     int32(envInt("DB_MIN_CONNS", 2)),
			MaxConnIdle:  envDuration("DB_MAX_CONN_IDLE", 10*time.Minute),
			MaxConnLife:  envDuration("DB_MAX_CONN_LIFE", 2*time.Hour),
			HealthPeriod: envDuration("DB_HEALTH_PERIOD", 30*time.Second),
		},
		Redis: RedisConfig{
			Addr:         env("REDIS_ADDR", "127.0.0.1:6379"),
			Password:     env("REDIS_PASSWORD", ""),
			DB:           envInt("REDIS_DB", 0),
			KeyPrefix:    env("REDIS_KEY_PREFIX", "shieldpanel:"),
			DialTimeout:  envDuration("REDIS_DIAL_TIMEOUT", 5*time.Second),
			ReadTimeout:  envDuration("REDIS_READ_TIMEOUT", 3*time.Second),
			WriteTimeout: envDuration("REDIS_WRITE_TIMEOUT", 3*time.Second),
		},
		Auth: AuthConfig{
			SessionCookieName: env("AUTH_SESSION_COOKIE", "shieldpanel_session"),
			CSRFCookieName:    env("AUTH_CSRF_COOKIE", "shieldpanel_csrf"),
			SessionTTL:        envDuration("AUTH_SESSION_TTL", 12*time.Hour),
			RememberMeTTL:     envDuration("AUTH_REMEMBER_TTL", 30*24*time.Hour),
			CookieSecure:      envBool("AUTH_COOKIE_SECURE", false),
			CookieDomain:      env("AUTH_COOKIE_DOMAIN", ""),
			BcryptCost:        envInt("AUTH_BCRYPT_COST", 12),
		},
		Security: SecurityConfig{
			SessionSecret:       env("SESSION_SECRET", "change-me-session-secret"),
			ChallengeSecret:     env("CHALLENGE_SECRET", "change-me-challenge-secret"),
			AllowedMethods:      strings.Split(env("ALLOWED_METHODS", "GET,POST,PUT,PATCH,DELETE,HEAD,OPTIONS"), ","),
			ChallengeTTL:        envDuration("CHALLENGE_TTL", 15*time.Minute),
			ChallengeGraceTTL:   envDuration("CHALLENGE_GRACE_TTL", 12*time.Hour),
			AutoBanThreshold:    envInt("AUTO_BAN_THRESHOLD", 8),
			AutoBanWindow:       envDuration("AUTO_BAN_WINDOW", 15*time.Minute),
			AutoBanTTL:          envDuration("AUTO_BAN_TTL", 30*time.Minute),
			TrustedProxyMode:    envBool("TRUSTED_PROXY_MODE", false),
			TrustedProxyHeaders: strings.Split(env("TRUSTED_PROXY_HEADERS", "CF-Connecting-IP,X-Forwarded-For"), ","),
		},
		Nginx: NginxConfig{
			Enabled:          envBool("NGINX_ENABLED", true),
			BinaryPath:       env("NGINX_BIN", "/usr/sbin/nginx"),
			ConfigPath:       env("NGINX_CONFIG_PATH", "/etc/nginx/nginx.conf"),
			SitesAvailable:   env("NGINX_SITES_AVAILABLE", "/etc/nginx/shieldpanel/sites-available"),
			SitesEnabled:     env("NGINX_SITES_ENABLED", "/etc/nginx/shieldpanel/sites-enabled"),
			ZonesPath:        env("NGINX_ZONES_PATH", "/etc/nginx/shieldpanel/zones.conf"),
			TemplatesDir:     env("NGINX_TEMPLATES_DIR", filepath.Join(root, "deploy", "nginx", "templates")),
			ACMEWebroot:      env("NGINX_ACME_WEBROOT", "/var/www/shieldpanel/acme"),
			PanelUpstreamURL: env("PANEL_UPSTREAM_URL", "http://127.0.0.1:8080"),
			ReloadOnChange:   envBool("NGINX_RELOAD_ON_CHANGE", true),
		},
		Defaults: DefaultsConfig{
			Language:              env("DEFAULT_LANGUAGE", "en"),
			DefaultProtection:     env("DEFAULT_PROTECTION_MODE", "basic"),
			DefaultRateLimitRPS:   envInt("DEFAULT_RATE_LIMIT_RPS", 20),
			DefaultRateLimitBurst: envInt("DEFAULT_RATE_LIMIT_BURST", 40),
			DefaultChallengeMode:  env("DEFAULT_CHALLENGE_MODE", "cookie"),
		},
		Paths: PathsConfig{
			MigrationsDir: env("MIGRATIONS_DIR", filepath.Join(root, "backend", "migrations")),
			DataDir:       env("DATA_DIR", filepath.Join(root, "data")),
			LogDir:        env("LOG_DIR", "/var/log/shieldpanel"),
			IssueCertScript: env("ISSUE_CERT_SCRIPT", filepath.Join(root, "deploy", "scripts", "issue_cert.sh")),
			RenewCertScript: env("RENEW_CERT_SCRIPT", filepath.Join(root, "deploy", "scripts", "renew_cert.sh")),
		},
	}
}

func env(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}
	return fallback
}

func envInt(key string, fallback int) int {
	value := env(key, "")
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func envBool(key string, fallback bool) bool {
	value := strings.ToLower(env(key, ""))
	if value == "" {
		return fallback
	}
	switch value {
	case "1", "true", "yes", "on":
		return true
	case "0", "false", "no", "off":
		return false
	default:
		return fallback
	}
}

func envDuration(key string, fallback time.Duration) time.Duration {
	value := env(key, "")
	if value == "" {
		return fallback
	}
	parsed, err := time.ParseDuration(value)
	if err != nil {
		return fallback
	}
	return parsed
}
