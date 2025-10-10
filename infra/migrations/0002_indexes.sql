-- SQL migration indexes
-- Vector indexes
CREATE INDEX ON email_embedding USING ivfflat (embedding vector_l2_ops) WITH (lists=100);
CREATE INDEX ON note_embedding USING ivfflat (embedding vector_l2_ops) WITH (lists=100);
CREATE INDEX ON instruction_embedding USING ivfflat (embedding vector_l2_ops) WITH (lists=50);

-- Useful btrees
CREATE INDEX ON email (user_id, sent_at DESC);
CREATE INDEX ON contact (user_id, email);
CREATE INDEX ON task (user_id, status, run_at);