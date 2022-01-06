package main

import (
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"
	"github.com/kamikazechaser/hedwig/internal/message"
)

var (
	qClient *asynq.Client
)

func genericHandler(ctx context.Context, t *asynq.Task) error {
	var msg message.Message

	if err := json.Unmarshal(t.Payload(), &msg); err != nil {
		return err
	}

	if err := hedwig[msg.Service].Push(msg); err != nil {
		return err
	}

	return nil
}

func initQueue() {
	redisConnection := asynq.RedisClientOpt{Addr: conf.String("queue.redis")}
	qClient = asynq.NewClient(redisConnection)

	queueProcessor := asynq.NewServer(redisConnection, asynq.Config{
		Concurrency: len(hedwig) * conf.Int("queue.concurrency"),
	})

	mux := asynq.NewServeMux()

	for plugin := range hedwig {
		mux.HandleFunc(plugin, genericHandler)
	}

	queueProcessor.Run(mux)
}
