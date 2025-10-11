
package main

import (
  "context"
  "database/sql"
  "log"
)

type CreateEventPayload struct{ Title, Start, End string; Attendees []string }

func handleCreateEvent(ctx context.Context, tx *sql.Tx, p CreateEventPayload) error {
  log.Printf("[worker] create_event title=%s", p.Title)
  return nil
}
