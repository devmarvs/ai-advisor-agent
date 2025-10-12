package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Home() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(`<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width,initial-scale=1" />
  <title>AI Advisor Agent</title>
  <style>
    :root{
      --bg:#f7f8fb; --panel:#fff; --muted:#6b7280; --text:#0f172a; --border:#e5e7eb;
      --primary:#0ea5e9; --chip:#f1f5f9; --chipText:#0f172a;
    }
    *{box-sizing:border-box}
    html,body{height:100%}
    body{margin:0;background:var(--bg);color:var(--text);font:16px/1.4 system-ui,Segoe UI,Arial}
    .app{max-width:980px;margin:0 auto;min-height:100%;display:grid;grid-template-columns:1fr;gap:16px;padding:16px}
    @media (min-width:900px){ .app{grid-template-columns:360px minmax(0,1fr)} }

    /* Left chat panel */
    .panel{background:var(--panel);border:1px solid var(--border);border-radius:16px;display:flex;flex-direction:column;min-height:70vh;overflow:hidden}
    .panel-header{padding:20px 20px 8px;border-bottom:1px solid var(--border)}
    .title{font-size:24px;font-weight:700;margin:0}
    .tabs{display:flex;gap:8px;margin-top:12px}
    .tab{padding:6px 10px;border-radius:999px;border:1px solid var(--border);background:#fff;font-weight:600}
    .tab.active{background:#eef6ff;border-color:#bfdbfe;color:#1d4ed8}

    .section{padding:16px 20px;border-bottom:1px solid var(--border)}
    .dim{color:var(--muted);font-size:14px}

    .chip-row{display:flex;flex-wrap:wrap;gap:8px;margin-top:12px}
    .chip{display:inline-flex;align-items:center;gap:6px;background:var(--chip);color:var(--chipText);
          padding:10px 12px;border-radius:12px;font-weight:600;border:1px solid var(--border)}
    .chip .face{width:18px;height:18px;border-radius:50%;background:#d1d5db;display:inline-block}

    .thread{padding:10px 20px 0 20px;display:flex;flex-direction:column;gap:14px;overflow:auto}
    .card{border:1px solid var(--border);border-radius:16px;background:#fff;padding:16px}
    .card .time{color:var(--muted);font-weight:600;margin-bottom:6px}
    .card .title{font-size:20px;margin:0 0 8px 0}
    .avatars{display:flex;gap:4px}
    .avatars .a{width:22px;height:22px;border-radius:50%;background:#cbd5e1;border:2px solid #fff}

    /* Composer */
    .composer{margin-top:auto;border-top:1px solid var(--border);padding:12px;background:linear-gradient(#fff,#fff)}
    .composer-inner{display:flex;gap:10px;align-items:center}
    .input{flex:1;border:1px solid var(--border);border-radius:14px;padding:12px 14px;font-size:16px}
    .btn{border:0;background:var(--primary);color:#fff;padding:10px 14px;border-radius:12px;font-weight:700;cursor:pointer}
    .btn:disabled{opacity:.6;cursor:not-allowed}

    /* Right column (optional preview stream) */
    .right{display:flex;flex-direction:column;gap:12px}
    .bubble{background:#f1f5f9;border:1px solid var(--border);border-radius:12px;padding:12px;white-space:pre-wrap}
    .muted{color:var(--muted)}
    .hint{font-size:13px;color:var(--muted)}
    .pill{display:inline-block;border:1px solid var(--border);border-radius:999px;padding:6px 10px;background:#fff;font-weight:600}
  </style>
</head>
<body>
  <main class="app">
    <!-- LEFT: Chat panel, matches the product screenshot closely -->
    <section class="panel">
      <div class="panel-header">
        <h1 class="title">Ask Anything</h1>
        <div class="tabs">
          <span class="tab active">Chat</span>
          <span class="tab">History</span>
          <span class="tab">+ New thread</span>
        </div>
      </div>

      <div class="section">
        <div class="dim">Context set to all meetings</div>
        <div class="dim" style="margin-top:4px;">11:17am – May 13, 2025</div>
      </div>

      <div class="section">
        <p style="margin:0 0 8px 0">I can answer questions about any Jump meeting. What do you want to know?</p>
        <div class="chip-row">
          <span class="chip">
            Find meetings I’ve had with
            <span class="face"></span> Bill
            and <span class="face"></span> Tim
            this month
          </span>
        </div>
      </div>

      <div class="section">
        <p style="margin:0">
          Sure, here are some recent meetings that you, Bill, and Tim all attended. I found 2 in May. <span class="pill">⎯⎯⎯</span>
        </p>
      </div>

      <div class="thread" id="thread">
        <div class="card">
          <div class="time">8 Thursday</div>
          <h3 class="title">Quarterly All Team Meeting</h3>
          <div class="avatars">
            <span class="a"></span><span class="a"></span><span class="a"></span><span class="a"></span>
          </div>
        </div>

        <div class="card">
          <div class="time">16 Friday</div>
          <h3 class="title">Strategy review</h3>
          <div class="avatars">
            <span class="a"></span><span class="a"></span>
          </div>
        </div>

        <div class="hint">I can summarize these meetings, schedule a follow up, and more!</div>
      </div>

      <div class="composer">
        <div class="composer-inner">
          <input id="msg" class="input" placeholder="Ask anything about your meetings..." />
          <button id="sendBtn" class="btn">Send</button>
        </div>
      </div>
    </section>

    <!-- RIGHT: reply area (shows what the backend returns) -->
    <section class="right">
      <div class="bubble muted">Replies</div>
      <div id="reply" class="bubble">No messages yet.</div>
      <div class="hint">The input sends a POST to <code>/chat</code>. The agent’s JSON reply is rendered here.</div>
    </section>
  </main>

<script>
const $ = (s)=>document.querySelector(s);
const thread = $("#thread");
const reply = $("#reply");
const api = window.location.origin;

function addUserBubble(text){
  const el = document.createElement("div");
  el.className = "bubble";
  el.textContent = text;
  thread.appendChild(el);
  thread.scrollTop = thread.scrollHeight;
}

function addAssistantBubble(text){
  const el = document.createElement("div");
  el.className = "bubble";
  el.textContent = text;
  thread.appendChild(el);
  thread.scrollTop = thread.scrollHeight;
}

async function send(){
  const btn = $("#sendBtn");
  const input = $("#msg");
  const value = input.value.trim();
  if(!value){ input.focus(); return }
  addUserBubble(value);
  btn.disabled = true;
  reply.textContent = "Sending...";

  try{
    const res = await fetch(api + "/chat", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ message: value })
    });
    const text = await res.text();
    reply.textContent = text;
    addAssistantBubble(text);
    input.value = "";
    input.focus();
  }catch(e){
    reply.textContent = "Error: " + e;
    addAssistantBubble("Error: " + e);
  }finally{
    btn.disabled = false;
  }
}

$("#sendBtn").addEventListener("click", send);
$("#msg").addEventListener("keydown", e=>{ if(e.key==="Enter"){ e.preventDefault(); send(); } });
</script>
</body>
</html>`))
	}
}
