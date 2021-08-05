package main

import (
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"
	"github.com/kamikazechaser/hedwig/internal/message"
)

func genericHandler(ctx context.Context, t *asynq.Task) error {
	var msg message.Message

	if err := json.Unmarshal(t.Payload(), &msg); err != nil {
		return err
	}

	if err := app.services[msg.Service].Push(msg); err != nil {
		return err
	}

	return nil
}

func initQueue() (*asynq.Server, *asynq.ServeMux) {
	redisConnection := asynq.RedisClientOpt{Addr: conf.String("asynq.redis")}
	q = asynq.NewClient(redisConnection)

	queueProcessor := asynq.NewServer(redisConnection, asynq.Config{
		Concurrency: len(app.enabledServices) * 2,
	})

	mux := asynq.NewServeMux()

	for _, service := range app.enabledServices {
		mux.HandleFunc(service, genericHandler)
	}

	return queueProcessor, mux
}
