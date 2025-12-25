const $ = (s) => document.querySelector(s);

const msgs = $("#msgs");
const historyBox = $("#history");
const tabChat = $("#tab-chat");
const tabHistory = $("#tab-history");
const tabNew = $("#tab-new");
const input = $("#msg");
const sendBtn = $("#sendBtn");

function addBubble(role, text) {
  const d = document.createElement("div");
  d.className = "bubble " + role;
  d.textContent = text;
  msgs.appendChild(d);
  window.scrollTo(0, document.body.scrollHeight);
}

function renderHistory(groups, fallbackMessages) {
  historyBox.innerHTML = "";
  if (groups.length === 0 && Array.isArray(fallbackMessages)) {
    const h = document.createElement("div");
    h.className = "bubble";
    h.textContent = "History";
    historyBox.appendChild(h);
    fallbackMessages.forEach((m) => {
      const d = document.createElement("div");
      const role = m.role || m.Role || "assistant";
      d.className = "bubble " + role;
      d.textContent = m.content || m.Content || "";
      historyBox.appendChild(d);
    });
    return;
  }

  groups.forEach((g) => {
    const h = document.createElement("div");
    h.className = "bubble";
    h.textContent = g.date;
    historyBox.appendChild(h);

    (g.items || []).forEach((m) => {
      const d = document.createElement("div");
      const role = m.role || "assistant";
      d.className = "bubble " + role;
      d.textContent = m.content || "";
      historyBox.appendChild(d);
    });
  });
}

async function loadHistory() {
  try {
    const r = await fetch("/messages");
    if (!r.ok) {
      throw new Error("messages fetch failed");
    }
    const j = await r.json();
    renderHistory(j.groups || [], j.messages);
  } catch (e) {
    historyBox.innerHTML = '<div class="bubble">Failed to load history.</div>';
  }
}

async function send() {
  const v = input.value.trim();
  if (!v) {
    return;
  }
  sendBtn.disabled = true;
  addBubble("user", v);
  input.value = "";
  try {
    const r = await fetch("/chat", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ message: v }),
    });
    const j = await r.json();
    if (j.error) {
      addBubble("assistant", "Error: " + j.error);
      return;
    }
    const txt = typeof j.reply === "string" ? j.reply : JSON.stringify(j.reply);
    addBubble("assistant", txt);
  } catch (e) {
    addBubble("assistant", "Error: " + e.message);
  } finally {
    sendBtn.disabled = false;
  }
}

tabChat.addEventListener("click", () => {
  tabChat.classList.add("active");
  tabHistory.classList.remove("active");
  historyBox.classList.remove("show");
});

tabHistory.addEventListener("click", async () => {
  tabHistory.classList.add("active");
  tabChat.classList.remove("active");
  await loadHistory();
  historyBox.classList.add("show");
  window.scrollTo({ top: 0, behavior: "smooth" });
});

tabNew.addEventListener("click", () => {
  msgs.innerHTML = "";
  tabChat.classList.add("active");
  tabHistory.classList.remove("active");
  historyBox.classList.remove("show");
  addBubble("assistant", "Started a new thread. How can I help?");
});

input.addEventListener("keydown", (e) => {
  if (e.key === "Enter") {
    e.preventDefault();
    send();
  }
});

sendBtn.addEventListener("click", send);
