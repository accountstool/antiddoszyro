package domain

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	UserRoleOwner  UserRole = "owner"
	UserRoleAdmin  UserRole = "admin"
	UserRoleViewer UserRole = "viewer"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	DisplayName  string    `json:"displayName"`
	Role         UserRole  `json:"role"`
	Language     string    `json:"language"`
	LastLoginAt  *time.Time `json:"lastLoginAt"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type Session struct {
	ID            uuid.UUID  `json:"id"`
	UserID        uuid.UUID  `json:"userId"`
	TokenHash     string     `json:"-"`
	CSRFToken     string     `json:"-"`
	IPAddress     string     `json:"ipAddress"`
	UserAgent     string     `json:"userAgent"`
	RememberMe    bool       `json:"rememberMe"`
	LastSeenAt    time.Time  `json:"lastSeenAt"`
	ExpiresAt     time.Time  `json:"expiresAt"`
	RevokedAt     *time.Time `json:"revokedAt"`
	CreatedAt     time.Time  `json:"createdAt"`
}

type ProtectionMode string

const (
	ProtectionOff         ProtectionMode = "off"
	ProtectionBasic       ProtectionMode = "basic"
	ProtectionAggressive  ProtectionMode = "aggressive"
	ProtectionUnderAttack ProtectionMode = "under_attack"
)

type ChallengeMode string

const (
	ChallengeOff    ChallengeMode = "off"
	ChallengeCookie ChallengeMode = "cookie"
	ChallengeJS     ChallengeMode = "js"
)

type Domain struct {
	ID                  uuid.UUID      `json:"id"`
	Name                string         `json:"name"`
	OriginHost          string         `json:"originHost"`
	OriginPort          int            `json:"originPort"`
	OriginProtocol      string         `json:"originProtocol"`
	OriginServerName    string         `json:"originServerName"`
	Enabled             bool           `json:"enabled"`
	ProtectionEnabled   bool           `json:"protectionEnabled"`
	ProtectionMode      ProtectionMode `json:"protectionMode"`
	ChallengeMode       ChallengeMode  `json:"challengeMode"`
	CloudflareMode      bool           `json:"cloudflareMode"`
	SSLAutoIssue        bool           `json:"sslAutoIssue"`
	SSLEnabled          bool           `json:"sslEnabled"`
	ForceHTTPS          bool           `json:"forceHttps"`
	RateLimitRPS        int            `json:"rateLimitRps"`
	RateLimitBurst      int            `json:"rateLimitBurst"`
	BadBotMode          bool           `json:"badBotMode"`
	HeaderValidation    bool           `json:"headerValidation"`
	JSChallengeEnabled  bool           `json:"jsChallengeEnabled"`
	AllowedMethods      []string       `json:"allowedMethods"`
	Notes               string         `json:"notes"`
	CreatedAt           time.Time      `json:"createdAt"`
	UpdatedAt           time.Time      `json:"updatedAt"`
}

type DomainRule struct {
	ID         uuid.UUID `json:"id"`
	DomainID   uuid.UUID `json:"domainId"`
	Name       string    `json:"name"`
	Type       string    `json:"type"`
	Pattern    string    `json:"pattern"`
	Action     string    `json:"action"`
	Enabled    bool      `json:"enabled"`
	Priority   int       `json:"priority"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type ListType string

const (
	ListTypeWhitelist ListType = "whitelist"
	ListTypeBlacklist ListType = "blacklist"
)

type IPEntry struct {
	ID         uuid.UUID  `json:"id"`
	DomainID   *uuid.UUID `json:"domainId"`
	ListType   ListType   `json:"listType"`
	IP         string     `json:"ip"`
	CIDR       string     `json:"cidr"`
	Reason     string     `json:"reason"`
	ExpiresAt  *time.Time `json:"expiresAt"`
	CreatedBy  *uuid.UUID `json:"createdBy"`
	CreatedAt  time.Time  `json:"createdAt"`
}

type TemporaryBan struct {
	ID         uuid.UUID  `json:"id"`
	DomainID   *uuid.UUID `json:"domainId"`
	IP         string     `json:"ip"`
	Reason     string     `json:"reason"`
	Source     string     `json:"source"`
	ExpiresAt  time.Time  `json:"expiresAt"`
	CreatedAt  time.Time  `json:"createdAt"`
}

type RequestDecision string

const (
	RequestDecisionAllowed         RequestDecision = "allowed"
	RequestDecisionBlocked         RequestDecision = "blocked"
	RequestDecisionChallenged      RequestDecision = "challenged"
	RequestDecisionChallengePassed RequestDecision = "challenge_passed"
)

type RequestLog struct {
	ID              uuid.UUID       `json:"id"`
	DomainID        *uuid.UUID      `json:"domainId"`
	DomainName      string          `json:"domainName"`
	ClientIP        string          `json:"clientIp"`
	CountryCode     string          `json:"countryCode"`
	Method          string          `json:"method"`
	Path            string          `json:"path"`
	QueryString     string          `json:"queryString"`
	UserAgent       string          `json:"userAgent"`
	RequestID       string          `json:"requestId"`
	Decision        RequestDecision `json:"decision"`
	StatusCode      int             `json:"statusCode"`
	BlockReason     string          `json:"blockReason"`
	ResponseTimeMS  int             `json:"responseTimeMs"`
	Score           int             `json:"score"`
	ChallengeType   string          `json:"challengeType"`
	CreatedAt       time.Time       `json:"createdAt"`
}

type RequestLogInput struct {
	DomainID        *uuid.UUID
	DomainName      string
	ClientIP        string
	CountryCode     string
	Method          string
	Path            string
	QueryString     string
	UserAgent       string
	RequestID       string
	Decision        RequestDecision
	StatusCode      int
	BlockReason     string
	ResponseTimeMS  int
	Score           int
	ChallengeType   string
	OccurredAt      time.Time
}

type TrafficRollup struct {
	WindowStart time.Time       `json:"windowStart"`
	Granularity string          `json:"granularity"`
	DomainID    *uuid.UUID      `json:"domainId"`
	DomainName  string          `json:"domainName"`
	Decision    RequestDecision `json:"decision"`
	BlockReason string          `json:"blockReason"`
	RequestCount int64          `json:"requestCount"`
	UniqueIPs   int64           `json:"uniqueIps"`
}

type SSLCertificate struct {
	ID          uuid.UUID  `json:"id"`
	DomainID     uuid.UUID `json:"domainId"`
	Issuer      string     `json:"issuer"`
	Status      string     `json:"status"`
	CertPath    string     `json:"certPath"`
	KeyPath     string     `json:"keyPath"`
	ExpiresAt   *time.Time `json:"expiresAt"`
	LastRenewAt *time.Time `json:"lastRenewAt"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}

type SystemSetting struct {
	Key        string    `json:"key"`
	Value      string    `json:"value"`
	Type       string    `json:"type"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type AuditLog struct {
	ID          uuid.UUID  `json:"id"`
	UserID      *uuid.UUID `json:"userId"`
	Username    string     `json:"username"`
	Action      string     `json:"action"`
	EntityType  string     `json:"entityType"`
	EntityID    string     `json:"entityId"`
	IPAddress   string     `json:"ipAddress"`
	UserAgent   string     `json:"userAgent"`
	Details     string     `json:"details"`
	CreatedAt   time.Time  `json:"createdAt"`
}

type DashboardSummary struct {
	Healthy           bool   `json:"healthy"`
	TotalDomains      int64  `json:"totalDomains"`
	BlockedToday      int64  `json:"blockedToday"`
	CurrentRPS        int64  `json:"currentRps"`
	CurrentBlockedPS  int64  `json:"currentBlockedPerSecond"`
	TopAttackedDomain string `json:"topAttackedDomain"`
	TopAttackingIP    string `json:"topAttackingIp"`
	TotalRequests24h  int64  `json:"totalRequests24h"`
	Allowed24h        int64  `json:"allowed24h"`
	Blocked24h        int64  `json:"blocked24h"`
	Challenged24h     int64  `json:"challenged24h"`
	ChallengePassRate float64 `json:"challengePassRate"`
}

type TimePoint struct {
	Label   string `json:"label"`
	Allowed int64  `json:"allowed"`
	Blocked int64  `json:"blocked"`
	Challenge int64 `json:"challenge"`
}

type RankedMetric struct {
	Name  string `json:"name"`
	Value int64  `json:"value"`
}

type StatsOverview struct {
	IncomingRequests   int64          `json:"incomingRequests"`
	AllowedRequests    int64          `json:"allowedRequests"`
	BlockRequests      int64          `json:"blockRequests"`
	ChallengedRequests int64          `json:"challengedRequests"`
	ChallengePassRate  float64        `json:"challengePassRate"`
	UniqueIPs          int64          `json:"uniqueIps"`
	PeakRPS            int64          `json:"peakRps"`
	PeakTime           string         `json:"peakTime"`
	TopIPs             []RankedMetric `json:"topIps"`
	TopUserAgents      []RankedMetric `json:"topUserAgents"`
	TopDomains         []RankedMetric `json:"topDomains"`
	TopReasons         []RankedMetric `json:"topReasons"`
	TopCountries       []RankedMetric `json:"topCountries"`
	RequestSeries      []TimePoint    `json:"requestSeries"`
}

type Pagination struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"pageSize"`
	TotalItems int64 `json:"totalItems"`
	TotalPages int64 `json:"totalPages"`
}
