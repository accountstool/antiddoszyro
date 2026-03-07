package protection

import (
	"fmt"
	"html/template"
	"strings"
)

func ChallengePage(language string, host string, originalURI string, mode string, reason string) string {
	title := "ShieldPanel Verification"
	description := "Complete the browser check to continue."
	buttonLabel := "Continue"
	waitText := "Verifying your browser..."
	if strings.HasPrefix(strings.ToLower(language), "vi") {
		title = "Xac thuc ShieldPanel"
		description = "Hoan tat kiem tra trinh duyet de tiep tuc."
		buttonLabel = "Tiep tuc"
		waitText = "Dang xac minh trinh duyet..."
	}
	autoSubmit := ""
	if mode == "js" {
		autoSubmit = `<script>setTimeout(function(){ document.getElementById('challenge-form').submit(); }, 1800);</script>`
	}
	return fmt.Sprintf(`<!doctype html>
<html lang="%s">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>%s</title>
  <style>
    body{font-family:system-ui,sans-serif;margin:0;min-height:100vh;display:grid;place-items:center;background:linear-gradient(180deg,#f8fafc,#dbeafe);color:#0f172a}
    .card{width:min(480px,92vw);background:#fff;border-radius:20px;padding:32px;box-shadow:0 24px 70px rgba(15,23,42,.12)}
    .badge{display:inline-block;padding:6px 10px;border-radius:999px;background:#e0f2fe;color:#0369a1;font-size:12px;font-weight:700;letter-spacing:.08em;text-transform:uppercase}
    h1{font-size:28px;margin:16px 0 12px}
    p{line-height:1.6;color:#334155}
    button{margin-top:20px;border:none;border-radius:12px;padding:12px 18px;background:#0f766e;color:#fff;font-weight:700;cursor:pointer}
  </style>
</head>
<body>
  <main class="card">
    <span class="badge">ShieldPanel</span>
    <h1>%s</h1>
    <p>%s</p>
    <p><strong>%s</strong></p>
    <form id="challenge-form" method="post" action="/__shieldpanel_verify">
      <input type="hidden" name="host" value="%s">
      <input type="hidden" name="redirect_uri" value="%s">
      <button type="submit">%s</button>
    </form>
  </main>
  %s
</body>
</html>`,
		template.HTMLEscapeString(language),
		template.HTMLEscapeString(title),
		template.HTMLEscapeString(title),
		template.HTMLEscapeString(description),
		template.HTMLEscapeString(waitText+" "+reason),
		template.HTMLEscapeString(host),
		template.HTMLEscapeString(originalURI),
		template.HTMLEscapeString(buttonLabel),
		autoSubmit,
	)
}

func BlockPage(language string, reason string) string {
	title := "Request blocked"
	description := "Your request triggered ShieldPanel protection."
	if strings.HasPrefix(strings.ToLower(language), "vi") {
		title = "Yeu cau bi chan"
		description = "Yeu cau cua ban da kich hoat co che bao ve ShieldPanel."
	}
	return fmt.Sprintf(`<!doctype html>
<html lang="%s"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width, initial-scale=1">
<title>%s</title>
<style>
body{font-family:system-ui,sans-serif;margin:0;min-height:100vh;display:grid;place-items:center;background:linear-gradient(180deg,#111827,#1f2937);color:#f9fafb}
.card{width:min(460px,92vw);background:#111827;border:1px solid rgba(255,255,255,.1);border-radius:20px;padding:32px}
.reason{margin-top:16px;padding:14px 16px;border-radius:12px;background:rgba(239,68,68,.12);color:#fecaca}
</style></head>
<body><main class="card"><h1>%s</h1><p>%s</p><div class="reason">%s</div></main></body></html>`,
		template.HTMLEscapeString(language),
		template.HTMLEscapeString(title),
		template.HTMLEscapeString(title),
		template.HTMLEscapeString(description),
		template.HTMLEscapeString(reason),
	)
}
