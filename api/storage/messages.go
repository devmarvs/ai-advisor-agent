package storage

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type ChatMessage struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Role      string
	Content   string
	CreatedAt time.Time
}

// SaveMessage inserts a chat message.
func SaveMessage(ctx context.Context, db *sql.DB, userID string, role, content string) error {
	_, err := db.ExecContext(ctx, `
		INSERT INTO agent_message (id, user_id, role, content, created_at)
		VALUES (gen_random_uuid(), $1::uuid, $2, $3, now())
	`, userID, role, content)
	return err
}

// ListRecentMessages returns the most recent N messages for a user (ascending by time).
func ListRecentMessages(ctx context.Context, db *sql.DB, userID string, limit int) ([]ChatMessage, error) {
	if limit <= 0 {
		limit = 100
	}
	rows, err := db.QueryContext(ctx, `
		SELECT id, user_id, role, content, created_at
		FROM agent_message
		WHERE user_id = $1::uuid
		ORDER BY created_at DESC
		LIMIT $2
	`, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tmp []ChatMessage
	for rows.Next() {
		var m ChatMessage
		if err := rows.Scan(&m.ID, &m.UserID, &m.Role, &m.Content, &m.CreatedAt); err != nil {
			return nil, err
		}
		tmp = append(tmp, m)
	}
	// reverse to ascending for UI
	n := len(tmp)
	out := make([]ChatMessage, n)
	for i := 0; i < n; i++ {
		out[i] = tmp[n-1-i]
	}
	return out, nil
}
