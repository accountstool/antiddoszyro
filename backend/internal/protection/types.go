package protection

import (
	"time"

	"shieldpanel/backend/internal/domain"
)

type RequestContext struct {
	Host      string
	URI       string
	Path      string
	Query     string
	Method    string
	Scheme    string
	RemoteIP  string
	UserAgent string
	Headers   map[string]string
	Cookies   map[string]string
	RequestID string
}

type DecisionAction string

const (
	ActionAllow     DecisionAction = "allow"
	ActionChallenge DecisionAction = "challenge"
	ActionBlock     DecisionAction = "block"
)

type Decision struct {
	Action         DecisionAction
	Reason         string
	StatusCode     int
	ChallengeMode  string
	SetCookie      string
	Domain         domain.Domain
	ClientIP       string
	CountryCode    string
	Score          int
	ResponseTimeMS int
	OccurredAt     time.Time
}
