-- SQL migration init
CREATE EXTENSION IF NOT EXISTS vector;

-- Gmail email storage
CREATE TABLE email (
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
    embedding VECTOR(3072)
);

-- HubSpot contacts & notes
CREATE TABLE contact (
    id BIGSERIAL PRIMARY KEY,
    user_id UUID REFERENCES app_user(id) ON DELETE CASCADE,
    hubspot_id TEXT,
    email TEXT,
    first_name TEXT,
    last_name TEXT,
    metadata JSONB DEFAULT '{}'
);

CREATE TABLE note (
    id BIGSERIAL PRIMARY KEY,
    user_id UUID REFERENCES app_user(id) ON DELETE CASCADE,
    hubspot_id TEXT,
    contact_id BIGINT REFERENCES contact(id) ON DELETE SET NULL,
    body TEXT,
    created_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE note_embedding (
    note_id BIGINT REFERENCES note(id) ON DELETE CASCADE,
    embedding VECTOR(3072)
);

-- Calendar events we create/track
CREATE TABLE meeting (
    id BIGSERIAL PRIMARY KEY,
    user_id UUID REFERENCES app_user(id) ON DELETE CASCADE,
    gcal_event_id TEXT,
    title TEXT,
    start_time TIMESTAMPTZ,
    end_time TIMESTAMPTZ,
    attendees JSONB
);

-- Agent messages (for chat history & audits)
CREATE TABLE agent_message (
    id BIGSERIAL PRIMARY KEY,
    user_id UUID REFERENCES app_user(id) ON DELETE CASCADE,
    role TEXT CHECK (role IN ('user','assistant','tool')),
    content TEXT,
    tool_calls JSONB,
    created_at TIMESTAMPTZ DEFAULT now(),
    thread_id TEXT
);

-- Durable tasks / state machines
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
