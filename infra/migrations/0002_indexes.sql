
CREATE INDEX IF NOT EXISTS email_sent_idx ON email (user_id, sent_at DESC);
CREATE INDEX IF NOT EXISTS contact_user_idx ON contact (user_id, email);
CREATE INDEX IF NOT EXISTS task_status_run_idx ON task (user_id, status, run_at);
CREATE INDEX IF NOT EXISTS email_embed_idx ON email_embedding USING ivfflat (embedding vector_l2_ops) WITH (lists=100);
CREATE INDEX IF NOT EXISTS note_embed_idx ON note_embedding USING ivfflat (embedding vector_l2_ops) WITH (lists=100);
CREATE INDEX IF NOT EXISTS instr_embed_idx ON instruction_embedding USING ivfflat (embedding vector_l2_ops) WITH (lists=50);
