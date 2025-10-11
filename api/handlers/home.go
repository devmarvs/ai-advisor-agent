package handlers

import (
  "net/http"
  "github.com/gin-gonic/gin"
)

func Home() gin.HandlerFunc {
  return func(c *gin.Context){
    c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(`<!doctype html>
<html>
<head><meta charset="utf-8"><title>AI Advisor Agent</title>
<style>body{font-family:system-ui;max-width:680px;margin:40px auto}</style></head>
<body>
  <h1>AI Advisor Agent</h1>
  <p>Type a message to enqueue a demo task.</p>
  <form onsubmit="event.preventDefault(); send();">
    <input id="msg" placeholder="hello" style="width:100%;padding:8px">
    <button style="margin-top:10px">Send</button>
  </form>
  <pre id="out"></pre>
<script>
async function send(){
  const msg = document.getElementById('msg').value || "hello";
  const res = await fetch('/chat',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({message:msg})});
  document.getElementById('out').textContent = await res.text();
}
</script>
</body></html>`))
  }
}
