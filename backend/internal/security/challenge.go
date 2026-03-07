package security

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type ChallengeSigner struct {
	secret []byte
}

func NewChallengeSigner(secret string) *ChallengeSigner {
	return &ChallengeSigner{secret: []byte(secret)}
}

func (s *ChallengeSigner) Sign(domain string, clientIP string, expiresAt time.Time) string {
	payload := fmt.Sprintf("%s|%s|%d", domain, clientIP, expiresAt.Unix())
	signature := s.sign(payload)
	return base64.RawURLEncoding.EncodeToString([]byte(payload + "|" + signature))
}

func (s *ChallengeSigner) Verify(token string, domain string, clientIP string, now time.Time) bool {
	raw, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil {
		return false
	}
	parts := strings.Split(string(raw), "|")
	if len(parts) != 4 {
		return false
	}
	payload := strings.Join(parts[:3], "|")
	if !hmac.Equal([]byte(parts[3]), []byte(s.sign(payload))) {
		return false
	}
	expiresAt, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		return false
	}
	return parts[0] == domain && parts[1] == clientIP && now.Unix() <= expiresAt
}

func (s *ChallengeSigner) sign(payload string) string {
	mac := hmac.New(sha256.New, s.secret)
	_, _ = mac.Write([]byte(payload))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}
