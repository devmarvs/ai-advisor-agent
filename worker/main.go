package main

import (
	"log"
	"os"

	"github.com/hibiken/asynq"
)

func main() {
	r := asynq.NewRedisClientOpt(asynq.RedisClientOpt{Addr: os.Getenv("REDIS_URL")})
	srv := asynq.NewServer(r, asynq.Config{Concurrency: 10})
	mux := asynq.NewServeMux()
	RegisterHandlers(mux)
	if err := srv.Run(mux); err != nil {
		log.Fatal(err)
	}
}
