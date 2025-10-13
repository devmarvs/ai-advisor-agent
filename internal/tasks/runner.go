package tasks

import "time"

type Task struct {
	ID        string
	UserID    string
	Type      string
	State     string
	Payload   string
	Result    string
	LastError string
	UpdatedAt time.Time
	CreatedAt time.Time
}

func Tick() error {
	return nil
}
