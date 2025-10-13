package storage

import (
	"context"
	"database/sql"
)

// Message represents a chat message row.
// type Message struct {
// 	ID        int64     `json:"id"`
// 	Role      string    `json:"role"`
// 	Content   string    `json:"content"`
// 	CreatedAt time.Time `json:"created_at"`
// }

// // SaveMessage inserts a message into the `agent_message` table.
// // Assumes the table has columns: role TEXT, content TEXT, created_at TIMESTAMPTZ DEFAULT now().
// func SaveMessage(ctx context.Context, db *sql.DB, role, content string) (int64, error) {
// 	const q = `
// 		INSERT INTO agent_message (role, content)
// 		VALUES ($1, $2)
// 		RETURNING id;
// 	`
// 	var id int64
// 	if err := db.QueryRowContext(ctx, q, role, content).Scan(&id); err != nil {
// 		return 0, err
// 	}
// 	return id, nil
// }

// ListRecentMessages returns the most recent messages (newest last).
// This powers the "History" view and lets the chat preload context.
func ListRecentMessages(ctx context.Context, db *sql.DB, limit int) ([]Message, error) {
	if limit <= 0 {
		limit = 20
	}

	const q = `
		SELECT id, role, content, COALESCE(created_at, now()) AS created_at
		FROM agent_message
		ORDER BY created_at DESC
		LIMIT $1;
	`

	rows, err := db.QueryContext(ctx, q, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []Message
	for rows.Next() {
		var m Message
		if err := rows.Scan(&m.ID, &m.Role, &m.Content, &m.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, m)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Reverse so the caller gets oldest → newest
	for i, j := 0, len(items)-1; i < j; i, j = i+1, j-1 {
		items[i], items[j] = items[j], items[i]
	}
	return items, nil
}
