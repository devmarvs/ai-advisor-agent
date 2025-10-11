
package main

import (
  "context"
  "database/sql"
  "log"
)

type SendEmailPayload struct{ To, Subject, Body, ThreadID string }

func handleSendEmail(ctx context.Context, tx *sql.Tx, p SendEmailPayload) error {
  log.Printf("[worker] send_email to=%s subject=%s", p.To, p.Subject)
  return nil
}
