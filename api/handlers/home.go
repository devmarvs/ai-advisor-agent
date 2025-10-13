package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Home serves the chat UI with History and New Thread actions.
// - History button loads /messages (no hardcoded bubbles)
// - +New thread clears the current conversation view (client-side)
func Home() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(`<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width,initial-scale=1" />
  <title>AI Advisor Agent</title>
  <style>
    :root { --bg:#0b0c10; --fg:#e5e7eb; --muted:#9ca3af; --card:#111827; --line:#1f2937; --accent:#0ea5e9; }
    * { box-sizing: border-box; }
    body { margin:0; font-family: Inter, system-ui, -apple-system, Segoe UI, Roboto, sans-serif; background: var(--bg); color: var(--fg); }
    header { padding: 16px 20px; border-bottom: 1px solid var(--line); display:flex; align-items:center; gap:16px; }
    header .tabs { display:flex; gap:8px; }
    header .tab { padding:8px 12px; border-radius:10px; border:1px solid var(--line); cursor:pointer; color: var(--muted); background:#0d1117; }
    header .tab.active { background: var(--accent); color:white; border-color: transparent; }
    .wrap { display:flex; max-width: 1100px; margin: 0 auto; }
    .sidebar { width: 260px; border-right: 1px solid var(--line); padding: 16px; min-height: calc(100vh - 58px); }
    .main { flex:1; padding: 16px; }
    .history { display:none; gap:8px; flex-direction: column; }
    .history.show { display:flex; }
    .msglist { display:flex; flex-direction: column; gap: 8px; padding-bottom: 88px; }
    .bubble { padding:12px 14px; border-radius:12px; line-height:1.5; white-space:pre-wrap; background: var(--card); border:1px solid var(--line); }
    .bubble.user { background:#1f2937; }
    .bubble.assistant { background:#0e7490; }
    .composer { position: fixed; left: 280px; right: 20px; bottom: 20px; display:flex; gap:8px; }
    .composer input { flex:1; padding:12px 14px; border-radius:10px; border:1px solid var(--line); background: var(--card); color: var(--fg); }
    .composer button { padding:12px 16px; border-radius:10px; border:0; background: var(--accent); color:white; cursor:pointer; }
    @media (max-width: 880px) { .sidebar { display:none; } .composer { left: 20px; } }
  </style>
</head>
<body>
  <header>
    <strong>AI Advisor Agent</strong>
    <div class="tabs">
      <button class="tab active" id="tab-chat">Chat</button>
      <button class="tab" id="tab-history">History</button>
      <button class="tab" id="tab-new">+ New thread</button>
    </div>
  </header>
  <div class="wrap">
    <aside class="sidebar">
      <div class="history" id="history"></div>
    </aside>
    <main class="main">
      <div id="msgs" class="msglist"></div>
    </main>
  </div>
  <div class="composer">
    <input id="msg" type="text" placeholder="Ask about clients or say e.g. 'Schedule an appointment with Sara Smith next week'"/>
    <button id="sendBtn">Send</button>
  </div>

<script>
const $ = s => document.querySelector(s);
const msgs = $("#msgs");
const historyBox = $("#history");
const tabChat = $("#tab-chat");
const tabHistory = $("#tab-history");
const tabNew = $("#tab-new");
const input = $("#msg");
const sendBtn = $("#sendBtn");

function addBubble(role, text){
  const d = document.createElement("div");
  d.className = "bubble " + role;
  d.textContent = text;
  msgs.appendChild(d);
  window.scrollTo(0, document.body.scrollHeight);
}

async function loadHistory(){
  try{
    const r = await fetch("/messages");
    if(!r.ok){ throw new Error("messages fetch failed"); }
    const j = await r.json();
    historyBox.innerHTML = "";
    (j.messages || []).forEach(m => {
      const d = document.createElement("div");
      const role = (m.role || m.Role || "assistant");
      const content = (m.content || m.Content || "");
      d.className = "bubble " + role;
      d.textContent = content;
      historyBox.appendChild(d);
    });
  }catch(e){
    historyBox.innerHTML = '<div class="bubble">Failed to load history.</div>';
  }
}

async function send(){
  const v = input.value.trim();
  if(!v) return;
  sendBtn.disabled = true;
  addBubble("user", v);
  input.value = "";
  try{
    const r = await fetch("/chat", { method:"POST", headers:{ "Content-Type":"application/json" }, body: JSON.stringify({ message: v })});
    const j = await r.json();
    if(j.error){ addBubble("assistant", "Error: " + j.error); return; }
    // ensure plain text, not JSON blob
    const txt = typeof j.reply === "string" ? j.reply : JSON.stringify(j.reply);
    addBubble("assistant", txt);
  }catch(e){
    addBubble("assistant", "Error: " + e.message);
  }finally{
    sendBtn.disabled = false;
  }
}

// Tabs behavior
tabChat.addEventListener("click", () => {
  tabChat.classList.add("active");
  tabHistory.classList.remove("active");
  $("#history").classList.remove("show");
});
tabHistory.addEventListener("click", async () => {
  tabHistory.classList.add("active");
  tabChat.classList.remove("active");
  await loadHistory();
  $("#history").classList.add("show");
  window.scrollTo({ top: 0, behavior: "smooth" });
});
tabNew.addEventListener("click", () => {
  // Clear current conversation on the UI (client-side new thread)
  msgs.innerHTML = "";
  tabChat.classList.add("active");
  tabHistory.classList.remove("active");
  $("#history").classList.remove("show");
  addBubble("assistant", "Started a new thread. How can I help?");
});

// Enter to send
input.addEventListener("keydown", (e) => { if(e.key === "Enter"){ e.preventDefault(); send(); } });
sendBtn.addEventListener("click", send);
</script>
</body>
</html>`))
	}
}
