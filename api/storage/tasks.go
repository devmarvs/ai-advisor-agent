package storage

import (
  "context"
  "database/sql"
  "encoding/json"
)

func Enqueue(ctx context.Context, db *sql.DB, userID string, kind string, payload any, runAt *string, dedupeKey *string) (int64, error) {
  b, _ := json.Marshal(payload)
  q := `INSERT INTO task (user_id, kind, status, payload, run_at, dedupe_key)
        VALUES ($1,$2,'pending',$3, $4, $5)
        ON CONFLICT (user_id, dedupe_key) WHERE $5 IS NOT NULL DO NOTHING
        RETURNING id`
  var id int64
  err := db.QueryRowContext(ctx, q, userID, kind, string(b), runAt, dedupeKey).Scan(&id)
  if err == sql.ErrNoRows { return 0, nil }
  return id, err
}

func WakeTask(ctx context.Context, db *sql.DB, taskID int64) error {
  _, err := db.ExecContext(ctx, `UPDATE task SET status='pending', updated_at=now() WHERE id=$1`, taskID)
  return err
}
