package worker

import (
  "context"
  "database/sql"
  "log"
  "time"
)

func Start(db *sql.DB) {
  go func() {
    log.Println("[worker] loop started")
    for {
      if err := claimAndRunOne(db); err != nil {
        log.Printf("[worker] error: %v", err)
        time.Sleep(2 * time.Second)
      }
    }
  }()
}

func claimAndRunOne(db *sql.DB) error {
  ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
  defer cancel()

  tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
  if err != nil { return err }

  var id int64; var kind, payload string
  err = tx.QueryRowContext(ctx, `
    UPDATE task SET status='running', claimed_at=now(), updated_at=now()
     WHERE id = (
       SELECT id FROM task
        WHERE status IN ('pending')
          AND (run_at IS NULL OR run_at <= now())
        ORDER BY priority ASC, run_at NULLS FIRST, id
        FOR UPDATE SKIP LOCKED
        LIMIT 1)
    RETURNING id, kind, payload::text
  `).Scan(&id, &kind, &payload)

  if err == sql.ErrNoRows {
    tx.Commit()
    time.Sleep(1500 * time.Millisecond)
    return nil
  }
  if err != nil { tx.Rollback(); return err }

  if err := dispatch(ctx, tx, id, kind, payload); err != nil {
    _, _ = tx.ExecContext(ctx, `UPDATE task SET status='failed', last_error=$2, updated_at=now() WHERE id=$1`, id, err.Error())
    tx.Commit()
    return err
  }

  _, err = tx.ExecContext(ctx, `UPDATE task SET status='done', updated_at=now() WHERE id=$1`, id)
  if err != nil { tx.Rollback(); return err }
  return tx.Commit()
}
