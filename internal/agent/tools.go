package agent

import (
	"context"
	"time"
)

type Toolset interface {
	SearchContext(ctx context.Context, userID, query string, limit int) ([]ContextDoc, error)
	SendEmail(ctx context.Context, userID, to, subject, text string) error
	FindSlots(ctx context.Context, userID string, from, to time.Time, attendees []string) ([]TimeSlot, error)
	CreateEvent(ctx context.Context, userID, title string, when time.Time, attendees []string, description string) (string, error)
}

type DefaultToolset struct{}

type ContextDoc struct {
	Kind    string
	Snippet string
	Source  string
	When    time.Time
}

type TimeSlot struct {
	Start time.Time
	End   time.Time
}

func (DefaultToolset) SearchContext(ctx context.Context, userID, query string, limit int) ([]ContextDoc, error) { return []ContextDoc{}, nil }
func (DefaultToolset) SendEmail(ctx context.Context, userID, to, subject, text string) error { return nil }
func (DefaultToolset) FindSlots(ctx context.Context, userID string, from, to time.Time, attendees []string) ([]TimeSlot, error) { return []TimeSlot{}, nil }
func (DefaultToolset) CreateEvent(ctx context.Context, userID, title string, when time.Time, attendees []string, description string) (string, error) { return "event_123", nil }
