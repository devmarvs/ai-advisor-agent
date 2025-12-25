package storage

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"time"
)

type Message struct {
	ID        int64     `json:"id"`
	Role      string    `json:"role"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

var (
	ensureOnce sync.Once
	ensureErr  error
)

// EnsureSchema creates/updates the agent_message table if it doesn't exist.
// Safe to call multiple times; it will run only once.
func EnsureSchema(db *sql.DB) error {
	ensureOnce.Do(func() {
		const q = `
CREATE TABLE IF NOT EXISTS agent_message (
  id BIGSERIAL PRIMARY KEY,
  user_id UUID,
  role TEXT NOT NULL,
  content TEXT NOT NULL,
  tool_calls JSONB,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  thread_id TEXT
);`
		_, ensureErr = db.Exec(q)
		if ensureErr != nil {
			return
		}
		_, _ = db.Exec(`ALTER TABLE agent_message ADD COLUMN IF NOT EXISTS user_id UUID`)
		_, _ = db.Exec(`ALTER TABLE agent_message ADD COLUMN IF NOT EXISTS tool_calls JSONB`)
		_, _ = db.Exec(`ALTER TABLE agent_message ADD COLUMN IF NOT EXISTS thread_id TEXT`)
	})
	return ensureErr
}

func SaveMessage(ctx context.Context, db *sql.DB, userID, role, content string) (int64, error) {
	if err := EnsureSchema(db); err != nil {
		return 0, fmt.Errorf("ensure schema: %w", err)
	}

	const q = `
INSERT INTO agent_message (user_id, role, content)
VALUES ($1, $2, $3)
RETURNING id;`
	var id int64
	var userArg any = userID
	if strings.TrimSpace(userID) == "" {
		userArg = nil
	}
	if err := db.QueryRowContext(ctx, q, userArg, role, content).Scan(&id); err != nil {
		return 0, fmt.Errorf("insert message: %w", err)
	}
	return id, nil
}

func LoadMessages(ctx context.Context, db *sql.DB, userID string, limit int) ([]Message, error) {
	if err := EnsureSchema(db); err != nil {
		return nil, fmt.Errorf("ensure schema: %w", err)
	}

	if limit <= 0 || limit > 500 {
		limit = 200
	}

	const q = `
SELECT id, role, content, created_at
FROM agent_message
WHERE user_id IS NOT DISTINCT FROM $1
ORDER BY created_at DESC, id DESC
LIMIT $2;`

	var userArg any = userID
	if strings.TrimSpace(userID) == "" {
		userArg = nil
	}
	rows, err := db.QueryContext(ctx, q, userArg, limit)
	if err != nil {
		return nil, fmt.Errorf("select messages: %w", err)
	}
	defer rows.Close()

	var out []Message
	for rows.Next() {
		var m Message
		if err := rows.Scan(&m.ID, &m.Role, &m.Content, &m.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		out = append(out, m)
	}
	return out, rows.Err()
}
