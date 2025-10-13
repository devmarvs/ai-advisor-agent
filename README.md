# 🤖 AI Advisor Agent

AI Advisor Agent is an intelligent assistant that integrates with **Google (Gmail + Calendar)** and **HubSpot CRM** to help financial advisors or client-facing professionals manage communications, schedule meetings, and maintain proactive client engagement.

The agent combines **LLM reasoning**, **tool-calling**, and **contextual memory** from your connected accounts to understand clients, automate follow-ups, and answer natural language questions such as:

- “Who mentioned their kid plays baseball?”
- “Why did Greg say he wanted to sell AAPL stock?”
- “Schedule an appointment with Sara Smith next week.”
- “When I add an event to my calendar, send a reminder email to attendees.”

---

## 🧠 Features

- **Google OAuth integration** — read/write Gmail and Calendar data  
- **HubSpot OAuth integration** — sync contacts, notes, and CRM data  
- **Chat-based interface** (ChatGPT-style UI)  
- **Persistent chat memory** stored in PostgreSQL  
- **Automatic syncing** of emails and CRM data  
- **Context-aware responses** powered by OpenAI  
- **RAG-ready** backend architecture (vector storage optional)  
- **Proactive automation** based on Gmail, Calendar, or HubSpot events  

---

## 🧰 Tech Stack

| Layer | Technology |
|-------|-------------|
| **Frontend (UI)** | HTML, TailwindCSS, Vanilla JS (with minimal React-like structure) |
| **Backend (API)** | Go (Gin framework) |
| **Database** | PostgreSQL |
| **ORM / Data Access** | native SQL via `database/sql` |
| **Authentication** | OAuth 2.0 (Google & HubSpot) |
| **AI Integration** | `go-openai` (OpenAI API) |
| **Deployment** | Render |
| **Storage** | `agent_message` table for conversation history |
| **RAG / Memory (optional)** | pgvector / embeddings (future-ready) |

---

## 🧩 Folder Structure

```
ai-advisor-agent/
│
├── api/
│   ├── handlers/         # HTTP endpoints (chat, auth, hubspot, google)
│   ├── storage/          # Database helper functions
│   ├── router.go         # Route definitions
│   └── main.go           # Entry point
│
├── static/
│   ├── connect.html      # Account connection UI
│   ├── chat.html         # Chat UI
│   └── assets/           # CSS, JS, icons
│
├── .env.example          # Environment variable template
└── README.md             # This file
```

---

## ⚙️ Setup Instructions

### 1️⃣ Environment Setup

Create a `.env` file based on `.env.example`:

```bash
PORT=8080
DATABASE_URL=postgres://<user>:<password>@<host>:5432/<dbname>?sslmode=require
OPENAI_API_KEY=sk-xxxxx
HUBSPOT_CLIENT_ID=xxxxxx
HUBSPOT_CLIENT_SECRET=xxxxxx
HUBSPOT_REDIRECT_URI=https://your-app.onrender.com/oauth/hubspot/callback
GOOGLE_CLIENT_ID=xxxxxx
GOOGLE_CLIENT_SECRET=xxxxxx
GOOGLE_REDIRECT_URI=https://your-app.onrender.com/oauth/google/callback
```

### 2️⃣ Database Schema

Run this SQL snippet to create the message store:

```sql
CREATE TABLE IF NOT EXISTS agent_message (
  id BIGSERIAL PRIMARY KEY,
  role TEXT NOT NULL,
  content TEXT NOT NULL,
  created_at TIMESTAMPTZ DEFAULT now()
);
```

### 3️⃣ Install Dependencies

```bash
go mod tidy
```

### 4️⃣ Run Locally

```bash
go run main.go
```

Access the app locally at  
👉 http://localhost:8080

---

## 🧭 OAuth Setup

### Google Cloud Console
- Enable **Gmail API** and **Calendar API**
- Add your redirect URI  
  → `https://your-app.onrender.com/oauth/google/callback`
- Add your test user (e.g. marvin.dev.ph@gmail.com)

### HubSpot Developer Portal
- Create a new private app
- Add scopes:
  ```
  crm.objects.contacts.read
  crm.objects.contacts.write
  crm.objects.owners.read
  crm.schemas.contacts.read
  crm.objects.deals.read
  crm.objects.deals.write
  crm.lists.read
  crm.objects.notes.read
  crm.objects.notes.write
  ```
- Redirect URI:  
  → `https://your-app.onrender.com/oauth/hubspot/callback`

---

## 💬 Screenshots

| Page | Preview |
|------|----------|
| **1. Connect Page** | <img width="800" height="605" alt="Screenshot 2025-10-13 at 7 03 10 AM" src="https://github.com/user-attachments/assets/ed255c0b-f51b-43b9-a263-75ead1c5c27f" />
 |
| **2. Chat Page (Initial Message)** | <img width="1599" height="769" alt="Screenshot 2025-10-13 at 7 01 25 AM" src="https://github.com/user-attachments/assets/370995f5-e353-48c3-9bb0-55f843c3030b" />
|
| **3. Chat History View** | <img width="502" height="711" alt="Screenshot 2025-10-13 at 7 01 45 AM" src="https://github.com/user-attachments/assets/70da548e-a939-40ca-a379-d2642c96dd99" />
 |
| **4. New Thread Example** | <img width="1606" height="778" alt="Screenshot 2025-10-13 at 7 01 53 AM" src="https://github.com/user-attachments/assets/1411a2f4-2160-4000-915f-c0fb6144179f" />
|

> 💡 Place your screenshots in a `/screenshots` folder and rename them to match the filenames above.

---

## 🧪 Example Questions to Try

- “Who mentioned their kid plays baseball?”
- “Why did Greg say he wanted to sell AAPL stock?”
- “Schedule an appointment with Sara Smith next week.”
- “When someone emails me that isn’t in HubSpot, create a contact.”

---

## 🚀 Deployment

Deployed easily on **Render** or **Fly.io**.

For Render:
- Create a **Web Service**
- Set build command: `go build -o server .`
- Set start command: `./server`
- Add your `.env` variables under **Environment**

---

## 🧱 Future Improvements

- [ ] Full vector memory with embeddings (pgvector)
- [ ] Real-time webhook ingestion for Gmail/HubSpot
- [ ] Multi-user chat threads with context persistence
- [ ] Dashboard for managing AI instructions
- [ ] Conversation analytics and insights

---

## 👨‍💻 Author

**Marvin (marvin.dev.ph@gmail.com)**  
Built using Go, PostgreSQL, and OpenAI.
