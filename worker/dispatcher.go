
package main

import (
  "context"
  "database/sql"
  "encoding/json"
  "fmt"
)

func dispatch(ctx context.Context, tx *sql.Tx, id int64, kind, payload string) error {
  switch kind {
  case "send_email":
    var p struct{ To, Subject, Body, ThreadID string }
    if err := json.Unmarshal([]byte(payload), &p); err != nil { return err }
    return handleSendEmail(ctx, tx, p)
  case "create_calendar_event":
    var p struct{ Title, Start, End string; Attendees []string }
    if err := json.Unmarshal([]byte(payload), &p); err != nil { return err }
    return handleCreateEvent(ctx, tx, p)
  case "wait_email_reply":
    var p struct{ ThreadID string }
    if err := json.Unmarshal([]byte(payload), &p); err != nil { return err }
    return handleWaitEmailReply(ctx, tx, p)
  default:
    return fmt.Errorf("unknown task kind: %s", kind)
  }
}
