package worker

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
    fmt.Printf("[worker] send_email to=%s subject=%s\n", p.To, p.Subject)
    return nil
  case "create_calendar_event":
    var p struct{ Title, Start, End string; Attendees []string }
    if err := json.Unmarshal([]byte(payload), &p); err != nil { return err }
    fmt.Printf("[worker] create_event title=%s\n", p.Title)
    return nil
  case "wait_email_reply":
    var p struct{ ThreadID string }
    if err := json.Unmarshal([]byte(payload), &p); err != nil { return err }
    return nil
  default:
    return fmt.Errorf("unknown task kind: %s", kind)
  }
}
