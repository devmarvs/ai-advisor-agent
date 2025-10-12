
CREATE INDEX IF NOT EXISTS email_sent_idx ON email (user_id, sent_at DESC);
CREATE INDEX IF NOT EXISTS contact_user_idx ON contact (user_id, email);
CREATE INDEX IF NOT EXISTS task_status_run_idx ON task (user_id, status, run_at);
CREATE INDEX IF NOT EXISTS email_embedding_gin ON email_embedding USING GIN (embedding);
CREATE INDEX IF NOT EXISTS note_embedding_gin ON note_embedding USING GIN (embedding);
CREATE INDEX IF NOT EXISTS instruction_embedding_gin ON instruction_embedding USING GIN (embedding);
