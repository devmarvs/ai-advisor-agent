
package main

import (
  "context"
  "database/sql"
  "log"
)

type WaitEmailReplyPayload struct{ ThreadID string }

func handleWaitEmailReply(ctx context.Context, tx *sql.Tx, p WaitEmailReplyPayload) error {
  log.Printf("[worker] wait_email_reply thread=%s", p.ThreadID)
  return nil
}
