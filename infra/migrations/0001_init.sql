
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE app_user (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  email TEXT UNIQUE NOT NULL,
  google_refresh_token TEXT,
  created_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE instruction (
  id BIGSERIAL PRIMARY KEY,
  user_id UUID REFERENCES app_user(id) ON DELETE CASCADE,
  text TEXT NOT NULL,
  active BOOLEAN DEFAULT TRUE,
  created_at TIMESTAMPTZ DEFAULT now()
);
CREATE TABLE instruction_embedding (
  instruction_id BIGINT REFERENCES instruction(id) ON DELETE CASCADE,
  embedding JSONB
);

CREATE TABLE email (
  id BIGSERIAL PRIMARY KEY,
  user_id UUID REFERENCES app_user(id) ON DELETE CASCADE,
  gmail_message_id TEXT UNIQUE,
  thread_id TEXT,
  sender TEXT,
  recipients TEXT[],
  subject TEXT,
  snippet TEXT,
  body_text TEXT,
  body_html TEXT,
  sent_at TIMESTAMPTZ,
  history_id BIGINT
);
CREATE TABLE email_embedding (
  email_id BIGINT REFERENCES email(id) ON DELETE CASCADE,
  embedding JSONB
);

CREATE TABLE contact (
  id BIGSERIAL PRIMARY KEY,
  user_id UUID REFERENCES app_user(id) ON DELETE CASCADE,
  email TEXT,
  first_name TEXT,
  last_name TEXT,
  metadata JSONB DEFAULT '{}'::jsonb
);

CREATE TABLE note (
  id BIGSERIAL PRIMARY KEY,
  user_id UUID REFERENCES app_user(id) ON DELETE CASCADE,
  contact_id BIGINT REFERENCES contact(id) ON DELETE SET NULL,
  body TEXT,
  created_at TIMESTAMPTZ DEFAULT now()
);
CREATE TABLE note_embedding (
  note_id BIGINT REFERENCES note(id) ON DELETE CASCADE,
  embedding JSONB
);

CREATE TABLE meeting (
  id BIGSERIAL PRIMARY KEY,
  user_id UUID REFERENCES app_user(id) ON DELETE CASCADE,
  gcal_event_id TEXT,
  title TEXT,
  start_time TIMESTAMPTZ,
  end_time TIMESTAMPTZ,
  attendees JSONB
);

CREATE TABLE agent_message (
  id BIGSERIAL PRIMARY KEY,
  user_id UUID REFERENCES app_user(id) ON DELETE CASCADE,
  role TEXT CHECK (role IN ('user','assistant','tool')),
  content TEXT,
  tool_calls JSONB,
  created_at TIMESTAMPTZ DEFAULT now(),
  thread_id TEXT
);

CREATE TYPE task_status AS ENUM ('pending','waiting','running','done','failed');

CREATE TABLE task (
  id BIGSERIAL PRIMARY KEY,
  user_id UUID REFERENCES app_user(id) ON DELETE CASCADE,
  kind TEXT,
  status task_status DEFAULT 'pending',
  payload JSONB,
  result JSONB,
  parent_task_id BIGINT REFERENCES task(id) ON DELETE SET NULL,
  run_at TIMESTAMPTZ,
  retries INT DEFAULT 0,
  last_error TEXT,
  created_at TIMESTAMPTZ DEFAULT now(),
  updated_at TIMESTAMPTZ DEFAULT now()
);
