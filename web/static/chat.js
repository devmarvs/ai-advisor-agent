const elMsgs=document.getElementById('msgs');const elForm=document.getElementById('f');const elText=document.getElementById('t');
function render(messages){elMsgs.innerHTML='';messages.forEach(m=>{const div=document.createElement('div');div.className='bubble '+(m.role||m.Role);div.textContent=m.content||m.Content;elMsgs.appendChild(div);});window.scrollTo(0,document.body.scrollHeight);}
async function load(){const r=await fetch('/messages');const j=await r.json();render(j.messages||[]);}
elForm.addEventListener('submit',async e=>{e.preventDefault();const text=elText.value.trim();if(!text)return;elText.value='';await fetch('/chat',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({message:text})});await load();});
setInterval(load,2500);load();
