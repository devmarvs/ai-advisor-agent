package tools

type SendEmailArgs struct {
	To       string
	Subject  string
	Body     string
	ThreadID *string
}

func SendEmail(userID string, a SendEmailArgs) error { /* TODO: use Gmail API */ return nil }
