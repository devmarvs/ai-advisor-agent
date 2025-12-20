package storage

import (
	"context"
	"database/sql"
)

// ListRecentMessages returns the most recent messages (newest last).
// This powers the "History" view and lets the chat preload context.
func ListRecentMessages(ctx context.Context, db *sql.DB, limit int) ([]Message, error) {
	if limit <= 0 {
		limit = 20
	}

	items, err := LoadMessages(ctx, db, limit)
	if err != nil {
		return nil, err
	}

	// Reverse so the caller gets oldest → newest
	for i, j := 0, len(items)-1; i < j; i, j = i+1, j-1 {
		items[i], items[j] = items[j], items[i]
	}
	return items, nil
}
