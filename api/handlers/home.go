package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Home() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(`<!doctype html>
<html>
<head>
  <meta charset="utf-8">
  <title>AI Advisor Agent</title>
  <style>
    body{font-family:system-ui,Segoe UI,Arial;max-width:680px;margin:40px auto;padding:0 16px}
    h1{margin:0 0 16px}
    #msg{width:100%;padding:10px;border:1px solid #ccc;border-radius:8px}
    button{padding:8px 14px;border:0;border-radius:8px;cursor:pointer}
    pre{background:#f6f7f9;border:1px solid #e3e5e8;border-radius:8px;padding:12px;white-space:pre-wrap}
    .muted{color:#666}
  </style>
</head>
<body>
  <h1>AI Advisor Agent</h1>
  <p class="muted">Type a message to enqueue a demo task.</p>

  <div style="display:flex;gap:8px;align-items:center">
    <input id="msg" placeholder="hello">
    <button id="sendBtn" type="button">Send</button>
  </div>

  <p id="status" class="muted" style="margin-top:10px"></p>
  <pre id="out"></pre>

<script>
const $ = (s)=>document.querySelector(s);
const api = window.location.origin;

async function callChat() {
  const btn = $("#sendBtn");
  const msg = $("#msg").value || "hello";
  $("#status").textContent = "Sending...";
  $("#out").textContent = "";
  btn.disabled = true;

  try {
    const res = await fetch(api + "/chat", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ message: msg })
    });

    const text = await res.text();
    $("#out").textContent = text;
    $("#status").textContent = res.ok ? "✅ OK" : "❌ Error " + res.status;
  } catch (e) {
    $("#status").textContent = "❌ Network/JS error";
    $("#out").textContent = String(e);
  } finally {
    btn.disabled = false;
  }
}

$("#sendBtn").addEventListener("click", callChat);
// also allow Enter key
$("#msg").addEventListener("keydown", (e)=>{ if(e.key==="Enter"){ e.preventDefault(); callChat(); }});
</script>
</body>
</html>`))
	}
}
