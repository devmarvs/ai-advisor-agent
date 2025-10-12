
ALTER TABLE task
  ADD COLUMN IF NOT EXISTS claimed_at timestamptz,
  ADD COLUMN IF NOT EXISTS priority int DEFAULT 100,
  ADD COLUMN IF NOT EXISTS dedupe_key text;

CREATE INDEX IF NOT EXISTS task_ready_idx
  ON task (status, priority, run_at NULLS FIRST, id);

CREATE UNIQUE INDEX IF NOT EXISTS task_dedupe_unique
  ON task (user_id, dedupe_key)
  WHERE dedupe_key IS NOT NULL;

CREATE INDEX IF NOT EXISTS task_waiting_idx
  ON task (run_at)
  WHERE status IN ('pending','waiting');
