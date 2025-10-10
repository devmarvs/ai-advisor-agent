package main

import (
	"context"

	"github.com/hibiken/asynq"
)

const (
	TaskWaitForEmailReply = "wait:email-reply"
)

func RegisterHandlers(mux *asynq.ServeMux) {
	mux.HandleFunc(TaskWaitForEmailReply, handleWaitForEmailReply)
}

func handleWaitForEmailReply(c context.Context, t *asynq.Task) error {
	// TODO: poll Gmail history until reply in thread; then continue workflow
	return nil
}
