# AI Advisor Agent

AI Advisor Agent integrates with Google (Gmail + Calendar) to help financial advisors or client-facing professionals manage communications, schedule meetings, and maintain proactive client engagement.
> Status: Work in progress - expect active development and changes.

The agent uses LLM reasoning and account context to answer questions such as:

- "Who mentioned their kid plays baseball?"
- "Why did Greg say he wanted to sell AAPL stock?"
- "Schedule an appointment with Sara Smith next week."
- "When I add an event to my calendar, send a reminder email to attendees."

---

## Features

- Google OAuth integration to read/write Gmail and Calendar data
- Chat-based interface
- Persistent chat memory stored in PostgreSQL
- Automatic syncing of emails and calendar data
- Responses powered by Groq
- Optional vector search via pgvector
- Proactive automation based on Gmail or Calendar events

---

## Tech Stack

| Layer | Technology |
|-------|-------------|
| Frontend (UI) | HTML, TailwindCSS, Vanilla JS (with minimal React-like structure) |
| Backend (API) | Go (Gin framework) |
| Database | PostgreSQL |
| ORM / Data Access | native SQL via `database/sql` |
| Authentication | OAuth 2.0 (Google) |
| AI Integration | `go-openai` (Groq OpenAI-compatible API) |
| Deployment | Vercel (planned) |
| Storage | `agent_message` table for conversation history |
| Vector Search (optional) | pgvector / embeddings |

---

## Folder Structure

```
ai-advisor-agent-scaffold/
│
├── api/                  # Vercel serverless entrypoint
│   └── index.go
│
├── server/
│   ├── handlers/         # HTTP endpoints (chat, auth, google)
│   ├── storage/          # Database helper functions
│   ├── router.go         # Route definitions
│   └── main.go           # Local entry point
│
├── web/
│   ├── templates/        # HTML UI templates
│   └── static/           # JS/CSS assets
│
├── api/migrations/       # Postgres migrations (bundled for Vercel)
├── .env.example          # Environment variable template
└── README.md             # This file
```

---

## Setup Instructions

### 1. Environment Setup

Create a `.env` file based on `.env.example`:

```bash
PORT=8080

DB_HOST=<host>
DB_PORT=5432
DB_NAME=<dbname>
DB_USER=<user>
DB_PASSWORD=<password>
DB_SSLMODE=require
DB_CHANNEL_BINDING=

GROQ_API_KEY=gsk_xxxxx
GROQ_MODEL=llama-3.1-8b-instant
GROQ_BASE_URL=https://api.groq.com/openai/v1

OAUTH_REDIRECT_BASE_URL=https://your-app.vercel.app
APP_BASE_URL=https://your-app.vercel.app
POST_CONNECT_REDIRECT=/

GOOGLE_CLIENT_ID=xxxxxx
GOOGLE_CLIENT_SECRET=xxxxxx

CRON_TOKEN=change-me
```

### 2. Database Schema

Run the migrations in `api/migrations` (recommended). For a quick local setup, create the minimum tables:

```sql
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS app_user (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  email TEXT UNIQUE NOT NULL,
  google_refresh_token TEXT,
  created_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE IF NOT EXISTS agent_message (
  id BIGSERIAL PRIMARY KEY,
  user_id UUID REFERENCES app_user(id) ON DELETE CASCADE,
  role TEXT NOT NULL,
  content TEXT NOT NULL,
  created_at TIMESTAMPTZ DEFAULT now(),
  thread_id TEXT
);
```

### 3. Install Dependencies

```bash
cd server
go mod tidy
```

### 4. Run Locally

```bash
cd server
go run .
```

Access the app locally at http://localhost:8080

---

## OAuth Setup

### Google Cloud Console
- Enable Gmail API and Calendar API
- Add your redirect URI:
  `https://your-app.vercel.app/oauth/google/callback`
- Add your test user (e.g. marvin.dev.ph@gmail.com)

---

## Screenshots

| Page | Preview |
|------|----------|
| 1. Connect Page | <img width="800" height="605" alt="Screenshot 2025-10-13 07-03-10" src="https://github.com/user-attachments/assets/ed255c0b-f51b-43b9-a263-75ead1c5c27f" /> |
| 2. Chat Page (Initial Message) | <img width="1599" height="769" alt="Screenshot 2025-10-13 07-01-25" src="https://github.com/user-attachments/assets/370995f5-e353-48c3-9bb0-55f843c3030b" /> |
| 3. Chat History View | <img width="502" height="711" alt="Screenshot 2025-10-13 07-01-45" src="https://github.com/user-attachments/assets/70da548e-a939-40ca-a379-d2642c96dd99" /> |
| 4. New Thread Example | <img width="1606" height="778" alt="Screenshot 2025-10-13 07-01-53" src="https://github.com/user-attachments/assets/1411a2f4-2160-4000-915f-c0fb6144179f" /> |

Place your screenshots in a `screenshots/` folder and rename them to match the filenames above.

---

## Example Questions to Try

- "Who mentioned their kid plays baseball?"
- "Why did Greg say he wanted to sell AAPL stock?"
- "Schedule an appointment with Sara Smith next week."

---

## Deployment

- Set environment variables in your deployment platform.
- Build and run the API from `server/` with `go build -o server .` and `./server`.
- Set `OAUTH_REDIRECT_BASE_URL` to your public base URL (for example, `https://your-app.vercel.app`).
- Migrations are applied automatically on startup. Set `MIGRATIONS_DIR` if the default path isn't found.

---

## Future Improvements

- [ ] Full vector memory with embeddings (pgvector)
- [ ] Real-time webhook ingestion for Gmail
- [ ] Multi-user chat threads with context persistence
- [ ] Dashboard for managing assistant instructions
- [ ] Conversation analytics and insights

---

## Author

Marvin (marvin.dev.ph@gmail.com)
Built using Go, PostgreSQL, and Groq.
