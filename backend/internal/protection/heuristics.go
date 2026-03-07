package protection

import (
	"regexp"
	"strings"

	"shieldpanel/backend/internal/domain"
)

var (
	traversalPattern   = regexp.MustCompile(`(?i)(\.\./|%2e%2e|%252e%252e)`)
	exploitPattern     = regexp.MustCompile(`(?i)(union\s+select|sleep\(|benchmark\(|or\s+1=1|<\?php|cmd=|wget\s|curl\s|etc/passwd|/bin/sh|base64_,|select.+from)`)
	suspiciousUAPattern = regexp.MustCompile(`(?i)(curl|python-requests|go-http-client|wget|masscan|sqlmap|nikto|zgrab|httpclient|libwww-perl)`)
)

func detectBadBot(userAgent string) bool {
	return suspiciousUAPattern.MatchString(strings.TrimSpace(userAgent))
}

func scoreHeaders(headers map[string]string, userAgent string) int {
	score := 0
	if strings.TrimSpace(userAgent) == "" {
		score += 3
	}
	if headers["accept"] == "" {
		score++
	}
	if headers["accept-encoding"] == "" {
		score++
	}
	if strings.Contains(userAgent, "Mozilla") && headers["accept-language"] == "" {
		score++
	}
	if strings.Contains(userAgent, "Mozilla") && headers["sec-fetch-site"] == "" {
		score++
	}
	return score
}

func detectSuspiciousPath(path string, query string) (bool, string) {
	combined := path + "?" + query
	if traversalPattern.MatchString(combined) {
		return true, "path_traversal"
	}
	if exploitPattern.MatchString(combined) {
		return true, "exploit_probe"
	}
	return false, ""
}

func matchCustomRules(rules []domain.DomainRule, ctx RequestContext) (string, string) {
	pathQuery := strings.ToLower(ctx.Path + "?" + ctx.Query)
	ua := strings.ToLower(ctx.UserAgent)
	for _, rule := range rules {
		if !rule.Enabled {
			continue
		}
		pattern := strings.ToLower(rule.Pattern)
		switch rule.Type {
		case "path", "query":
			if strings.Contains(pathQuery, pattern) {
				return rule.Action, "rule:" + rule.Name
			}
		case "ua":
			if strings.Contains(ua, pattern) {
				return rule.Action, "rule:" + rule.Name
			}
		case "ip":
			if strings.Contains(strings.ToLower(ctx.RemoteIP), pattern) {
				return rule.Action, "rule:" + rule.Name
			}
		}
	}
	return "", ""
}
