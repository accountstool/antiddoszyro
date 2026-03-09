package protection

import (
	"fmt"
	"html/template"
	"strings"
)

func ChallengePage(language string, host string, originalURI string, mode string, reason string, nonce string) string {
	title := "ShieldPanel Verification"
	description := "Complete the browser verification to continue to the protected site."
	eyebrow := "ShieldPanel Security"
	sliderLabel := "Drag the slider to complete verification"
	sliderReady := "Verification complete. Redirecting..."
	sliderDragging := "Keep dragging to the end..."
	fallbackLabel := "Continue manually"
	waitText := "Browser signals are being evaluated."
	featureA := "Browser telemetry and interaction timing"
	featureB := "Cookie clearance and IP-bound verification"
	featureC := "Bot-resistant slider confirmation"
	noScript := "JavaScript is required to complete this verification."
	if strings.HasPrefix(strings.ToLower(language), "vi") {
		title = "Xác thực ShieldPanel"
		description = "Hoàn tất xác minh trình duyệt để tiếp tục vào website đang được bảo vệ."
		eyebrow = "Bảo vệ ShieldPanel"
		sliderLabel = "Kéo thanh trượt sang hết bên phải để xác thực"
		sliderReady = "Xác thực thành công. Đang chuyển tiếp..."
		sliderDragging = "Tiếp tục kéo đến hết thanh..."
		fallbackLabel = "Tiếp tục thủ công"
		waitText = "Đang đánh giá tín hiệu trình duyệt."
		featureA = "Đo tín hiệu trình duyệt và thời gian tương tác"
		featureB = "Cookie clearance gắn với IP truy cập"
		featureC = "Xác thực slider chống bot"
		noScript = "Trang xác thực này yêu cầu JavaScript để hoạt động."
	}

	statusText := challengeReasonText(language, reason)

	return fmt.Sprintf(`<!doctype html>
<html lang="%s">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>%s</title>
  <style>
    :root{
      color-scheme:dark;
      --bg:#050816;
      --panel:#0c1224;
      --panel-soft:rgba(255,255,255,.04);
      --line:rgba(255,255,255,.12);
      --text:#f8fafc;
      --muted:#9fb0c9;
      --accent:#8be9fd;
      --accent-strong:#24c8db;
      --success:#34d399;
      --danger:#fb7185;
      --shadow:0 30px 80px rgba(2,6,23,.45);
    }
    *{box-sizing:border-box}
    body{
      margin:0;
      min-height:100vh;
      font-family:"Segoe UI",system-ui,sans-serif;
      color:var(--text);
      background:
        radial-gradient(circle at 15%% 20%%, rgba(36,200,219,.22), transparent 26%%),
        radial-gradient(circle at 85%% 12%%, rgba(148,163,184,.14), transparent 24%%),
        linear-gradient(180deg,#030712 0%%,#071128 52%%,#09152e 100%%);
      overflow:hidden;
    }
    .grid{
      position:relative;
      min-height:100vh;
      display:grid;
      place-items:center;
      padding:32px 18px;
    }
    .shell{
      width:min(1020px,100%%);
      display:grid;
      grid-template-columns:minmax(0,1.1fr) minmax(0,.9fr);
      gap:24px;
      align-items:stretch;
    }
    .panel,.info{
      position:relative;
      overflow:hidden;
      border:1px solid var(--line);
      border-radius:28px;
      background:linear-gradient(180deg,rgba(12,18,36,.96),rgba(7,11,24,.92));
      box-shadow:var(--shadow);
      backdrop-filter:blur(18px);
    }
    .panel{padding:34px}
    .info{padding:30px}
    .eyebrow{
      display:inline-flex;
      align-items:center;
      gap:8px;
      padding:8px 12px;
      border-radius:999px;
      background:rgba(139,233,253,.12);
      border:1px solid rgba(139,233,253,.22);
      color:var(--accent);
      font-size:12px;
      font-weight:700;
      letter-spacing:.16em;
      text-transform:uppercase;
    }
    h1{
      margin:18px 0 10px;
      font-size:clamp(34px,5vw,56px);
      line-height:.98;
      letter-spacing:-.05em;
    }
    p{
      margin:0;
      color:var(--muted);
      line-height:1.7;
      font-size:15px;
    }
    .status{
      margin-top:22px;
      padding:16px 18px;
      border-radius:18px;
      background:rgba(15,23,42,.78);
      border:1px solid rgba(139,233,253,.16);
      color:#dbeafe;
      font-weight:600;
    }
    .slider-shell{
      margin-top:26px;
      padding:22px;
      border-radius:22px;
      border:1px solid rgba(255,255,255,.08);
      background:linear-gradient(180deg,rgba(255,255,255,.04),rgba(255,255,255,.02));
    }
    .slider-head{
      display:flex;
      justify-content:space-between;
      gap:16px;
      align-items:center;
      margin-bottom:14px;
    }
    .slider-head strong{font-size:15px}
    .slider-hint{font-size:13px;color:var(--muted)}
    .slider{
      position:relative;
      width:100%%;
      height:74px;
      border-radius:999px;
      border:1px solid rgba(255,255,255,.08);
      background:rgba(3,7,18,.78);
      overflow:hidden;
      user-select:none;
      -webkit-user-select:none;
    }
    .slider-progress{
      position:absolute;
      inset:0 auto 0 0;
      width:74px;
      border-radius:999px;
      background:linear-gradient(90deg,rgba(36,200,219,.22),rgba(36,200,219,.55));
      transition:width .14s ease;
    }
    .slider-label{
      position:absolute;
      inset:0;
      display:grid;
      place-items:center;
      padding:0 96px 0 112px;
      text-align:center;
      color:#dbeafe;
      font-weight:700;
      letter-spacing:.01em;
      font-size:14px;
      pointer-events:none;
    }
    .slider-knob{
      position:absolute;
      left:8px;
      top:8px;
      width:58px;
      height:58px;
      border:none;
      border-radius:999px;
      background:linear-gradient(180deg,#f8fafc,#cbd5e1);
      color:#0f172a;
      display:grid;
      place-items:center;
      font-size:24px;
      font-weight:900;
      cursor:grab;
      box-shadow:0 12px 28px rgba(15,23,42,.32);
      transition:transform .14s ease, box-shadow .14s ease;
    }
    .slider-knob:active{cursor:grabbing}
    .slider-complete .slider-progress{background:linear-gradient(90deg,rgba(52,211,153,.35),rgba(52,211,153,.68))}
    .slider-complete .slider-knob{background:linear-gradient(180deg,#ecfdf5,#bbf7d0)}
    .actions{
      display:flex;
      gap:12px;
      align-items:center;
      margin-top:18px;
    }
    .button{
      appearance:none;
      border:none;
      border-radius:14px;
      padding:13px 18px;
      font-size:14px;
      font-weight:700;
      cursor:pointer;
      transition:transform .16s ease, opacity .16s ease, background .16s ease;
    }
    .button[disabled]{opacity:.45;cursor:not-allowed}
    .button-primary{
      background:linear-gradient(135deg,var(--accent-strong),#67e8f9);
      color:#062033;
    }
    .button-secondary{
      background:transparent;
      border:1px solid rgba(255,255,255,.16);
      color:var(--text);
    }
    .meta{
      margin-top:22px;
      display:grid;
      gap:12px;
    }
    .meta-card{
      border-radius:18px;
      padding:16px 18px;
      border:1px solid rgba(255,255,255,.08);
      background:rgba(255,255,255,.03);
    }
    .meta-card span{
      display:block;
      font-size:12px;
      letter-spacing:.12em;
      text-transform:uppercase;
      color:var(--muted);
      margin-bottom:6px;
    }
    .meta-card strong{
      font-size:16px;
      line-height:1.5;
    }
    .list{
      margin:22px 0 0;
      padding:0;
      list-style:none;
      display:grid;
      gap:14px;
    }
    .list li{
      display:flex;
      gap:12px;
      align-items:flex-start;
      padding:14px 16px;
      border-radius:18px;
      background:rgba(255,255,255,.03);
      border:1px solid rgba(255,255,255,.07);
      color:#dbeafe;
    }
    .bullet{
      flex:0 0 auto;
      width:30px;
      height:30px;
      border-radius:12px;
      display:grid;
      place-items:center;
      background:rgba(36,200,219,.16);
      color:var(--accent);
      font-size:14px;
      font-weight:800;
    }
    .footer{
      margin-top:18px;
      font-size:12px;
      color:var(--muted);
    }
    .sr-only{
      position:absolute;
      width:1px;
      height:1px;
      padding:0;
      margin:-1px;
      overflow:hidden;
      clip:rect(0,0,0,0);
      border:0;
    }
    @media (max-width: 900px){
      body{overflow:auto}
      .shell{grid-template-columns:1fr}
      .panel,.info{padding:24px}
      .slider-label{padding:0 80px}
    }
  </style>
</head>
<body>
  <div class="grid">
    <div class="shell">
      <main class="panel">
        <span class="eyebrow">%s</span>
        <h1>%s</h1>
        <p>%s</p>
        <div class="status" id="challenge-status">%s</div>

        <div class="slider-shell">
          <div class="slider-head">
            <strong id="slider-title">%s</strong>
            <span class="slider-hint">%s</span>
          </div>
          <div class="slider" id="slider">
            <div class="slider-progress" id="slider-progress"></div>
            <div class="slider-label" id="slider-label">%s</div>
            <button class="slider-knob" id="slider-knob" type="button" aria-label="%s">></button>
          </div>
          <div class="actions">
            <button class="button button-primary" id="submit-button" type="submit" form="challenge-form" disabled>%s</button>
            <button class="button button-secondary" id="retry-button" type="button">%s</button>
          </div>
          <div class="footer">%s</div>
        </div>

        <form id="challenge-form" method="post" action="/__shieldpanel_verify">
          <input type="hidden" name="host" value="%s">
          <input type="hidden" name="redirect_uri" value="%s">
          <input type="hidden" name="challenge_mode" value="%s">
          <input type="hidden" name="challenge_nonce" value="%s">
          <input type="hidden" name="slider_completed" id="slider_completed" value="0">
          <input type="hidden" name="interaction_ms" id="interaction_ms" value="0">
          <input type="hidden" name="pointer_moves" id="pointer_moves" value="0">
          <input type="hidden" name="webdriver" id="webdriver" value="0">
          <input type="hidden" name="language" id="language" value="">
          <input type="hidden" name="timezone" id="timezone" value="">
          <input type="hidden" name="screen" id="screen" value="">
          <input type="hidden" name="viewport" id="viewport" value="">
          <input type="hidden" name="platform" id="platform" value="">
          <input type="hidden" name="device_memory" id="device_memory" value="">
          <input type="hidden" name="hardware_concurrency" id="hardware_concurrency" value="">
          <input class="sr-only" type="text" name="website" id="website" autocomplete="off" tabindex="-1">
        </form>
      </main>

      <aside class="info">
        <span class="eyebrow">%s</span>
        <div class="meta">
          <div class="meta-card">
            <span>Host</span>
            <strong>%s</strong>
          </div>
          <div class="meta-card">
            <span>Request</span>
            <strong>%s</strong>
          </div>
        </div>
        <ul class="list">
          <li><span class="bullet">1</span><div>%s</div></li>
          <li><span class="bullet">2</span><div>%s</div></li>
          <li><span class="bullet">3</span><div>%s</div></li>
        </ul>
        <div class="footer"><noscript>%s</noscript></div>
      </aside>
    </div>
  </div>
  <script>
  (function () {
    var slider = document.getElementById("slider");
    var knob = document.getElementById("slider-knob");
    var progress = document.getElementById("slider-progress");
    var label = document.getElementById("slider-label");
    var submitButton = document.getElementById("submit-button");
    var retryButton = document.getElementById("retry-button");
    var form = document.getElementById("challenge-form");
    var status = document.getElementById("challenge-status");
    var pageStart = performance.now();
    var dragStart = 0;
    var dragging = false;
    var pointerMoves = 0;
    var startOffset = 0;
    var knobStartX = 8;
    var completed = false;

    function setField(id, value) {
      var el = document.getElementById(id);
      if (el) {
        el.value = value;
      }
    }

    function syncBrowserSignals() {
      setField("language", navigator.language || "");
      setField("timezone", (Intl.DateTimeFormat().resolvedOptions().timeZone || ""));
      setField("screen", window.screen.width + "x" + window.screen.height);
      setField("viewport", window.innerWidth + "x" + window.innerHeight);
      setField("platform", navigator.platform || "");
      setField("device_memory", navigator.deviceMemory || "");
      setField("hardware_concurrency", navigator.hardwareConcurrency || "");
      setField("webdriver", navigator.webdriver ? "1" : "0");
    }

    function maxOffset() {
      return slider.clientWidth - knob.offsetWidth - 16;
    }

    function setOffset(offset) {
      var bounded = Math.max(0, Math.min(offset, maxOffset()));
      knob.style.transform = "translateX(" + bounded + "px)";
      progress.style.width = (bounded + knob.offsetWidth + 8) + "px";
      return bounded;
    }

    function resetSlider() {
      dragging = false;
      pointerMoves = 0;
      setOffset(0);
      label.textContent = %q;
    }

    function finishSlider() {
      completed = true;
      dragging = false;
      setField("slider_completed", "1");
      setField("interaction_ms", String(Math.round(performance.now() - pageStart)));
      setField("pointer_moves", String(pointerMoves));
      slider.classList.add("slider-complete");
      label.textContent = %q;
      if (status) {
        status.textContent = %q;
      }
      submitButton.disabled = false;
      window.setTimeout(function () {
        try {
          form.requestSubmit();
        } catch (error) {
          form.submit();
        }
      }, 280);
    }

    syncBrowserSignals();
    window.addEventListener("resize", syncBrowserSignals);
    submitButton.disabled = true;
    retryButton.addEventListener("click", function () {
      window.location.reload();
    });

    knob.addEventListener("pointerdown", function (event) {
      if (completed) {
        return;
      }
      dragging = true;
      pointerMoves = 0;
      dragStart = performance.now();
      knobStartX = event.clientX;
      startOffset = parseFloat((knob.style.transform.match(/translateX\(([0-9.]+)px\)/) || [0, 0])[1]);
      label.textContent = %q;
      knob.setPointerCapture(event.pointerId);
      event.preventDefault();
    });

    window.addEventListener("pointermove", function (event) {
      if (!dragging || completed) {
        return;
      }
      pointerMoves += 1;
      var nextOffset = startOffset + (event.clientX - knobStartX);
      var bounded = setOffset(nextOffset);
      if (bounded >= maxOffset() * 0.94 && performance.now() - dragStart > 350) {
        finishSlider();
      }
    });

    function stopDragging() {
      if (!dragging || completed) {
        return;
      }
      dragging = false;
      if (parseFloat((knob.style.transform.match(/translateX\(([0-9.]+)px\)/) || [0, 0])[1]) < maxOffset() * 0.94) {
        resetSlider();
      }
    }

    window.addEventListener("pointerup", stopDragging);
    window.addEventListener("pointercancel", stopDragging);
    resetSlider();
  })();
  </script>
</body>
</html>`,
		template.HTMLEscapeString(language),
		template.HTMLEscapeString(title),
		template.HTMLEscapeString(eyebrow),
		template.HTMLEscapeString(title),
		template.HTMLEscapeString(description),
		template.HTMLEscapeString(waitText + " " + statusText),
		template.HTMLEscapeString(sliderLabel),
		template.HTMLEscapeString(statusText),
		template.HTMLEscapeString(sliderLabel),
		template.HTMLEscapeString(sliderLabel),
		template.HTMLEscapeString(fallbackLabel),
		template.HTMLEscapeString(fallbackLabel),
		template.HTMLEscapeString(waitText),
		template.HTMLEscapeString(host),
		template.HTMLEscapeString(originalURI),
		template.HTMLEscapeString(mode),
		template.HTMLEscapeString(nonce),
		template.HTMLEscapeString(eyebrow),
		template.HTMLEscapeString(host),
		template.HTMLEscapeString(originalURI),
		template.HTMLEscapeString(featureA),
		template.HTMLEscapeString(featureB),
		template.HTMLEscapeString(featureC),
		template.HTMLEscapeString(noScript),
		template.JSEscapeString(sliderLabel),
		template.JSEscapeString(sliderReady),
		template.JSEscapeString(waitText),
		template.JSEscapeString(sliderDragging),
	)
}

func BlockPage(language string, reason string) string {
	title := "Request blocked"
	description := "Your request triggered ShieldPanel protection."
	if strings.HasPrefix(strings.ToLower(language), "vi") {
		title = "Yêu cầu bị chặn"
		description = "Yêu cầu của bạn đã kích hoạt cơ chế bảo vệ ShieldPanel."
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

func challengeReasonText(language string, reason string) string {
	key := strings.TrimSpace(strings.ToLower(reason))
	if strings.HasPrefix(strings.ToLower(language), "vi") {
		switch key {
		case "challenge_required", "rate_limit":
			return "Lưu lượng này cần xác minh thêm trước khi được chuyển tiếp."
		case "under_attack":
			return "Website đang ở chế độ bảo vệ cao, cần xác minh đầy đủ."
		case "webdriver_detected":
			return "Trình duyệt tự động hóa đã bị phát hiện."
		case "slider_incomplete":
			return "Bạn cần kéo thanh trượt đến cuối để tiếp tục."
		case "browser_signals_missing":
			return "Thiếu tín hiệu trình duyệt cần thiết cho bước xác thực."
		case "challenge_too_fast", "interaction_too_low":
			return "Tương tác xác minh chưa đủ tự nhiên. Hãy thử lại chậm hơn."
		case "challenge_expired", "challenge_mismatch":
			return "Phiên xác thực đã hết hạn hoặc không còn hợp lệ."
		case "bot_signal_detected":
			return "Tín hiệu bot đã bị phát hiện trong bước xác thực."
		default:
			return "Hệ thống đang yêu cầu xác minh trình duyệt trước khi cho phép truy cập."
		}
	}

	switch key {
	case "challenge_required", "rate_limit":
		return "This request needs additional browser verification before it can continue."
	case "under_attack":
		return "The site is currently in a high-protection mode."
	case "webdriver_detected":
		return "Automation tooling was detected during verification."
	case "slider_incomplete":
		return "Complete the slider to continue."
	case "browser_signals_missing":
		return "Required browser signals are missing."
	case "challenge_too_fast", "interaction_too_low":
		return "Verification interaction was too fast. Please try again more naturally."
	case "challenge_expired", "challenge_mismatch":
		return "The verification session expired or no longer matches this request."
	case "bot_signal_detected":
		return "Bot-like verification signals were detected."
	default:
		return "ShieldPanel is asking for a browser verification before forwarding this request."
	}
}
