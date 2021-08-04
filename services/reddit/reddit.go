package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/kamikazechaser/hedwig/internal/message"
	"github.com/vartanbeno/go-reddit/reddit"
)

const (
	serviceName string = "reddit"
	version     string = "v1"
)

type config struct {
	ID       string `json:"app_id"`
	Secret   string `json:"app_secret"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type client struct {
	config       *config
	redditClient *reddit.Client
}

// New initializes the service client
func New(jsonConfig []byte) (interface{}, error) {
	var conf *config

	if err := json.Unmarshal(jsonConfig, &conf); err != nil {
		return nil, err
	}

	// Ensure all configs are present in the config file before attempting to push notifications
	if conf.ID == "" || conf.Secret == "" || conf.Username == "" || conf.Password == "" {
		return nil, errors.New("reddit config is incomplete")
	}

	httpClient := &http.Client{Timeout: time.Second * 15}
	credentials := reddit.Credentials(*conf)
	newRedditClient, _ := reddit.NewClient(httpClient, &credentials)

	return &client{
		config:       conf,
		redditClient: newRedditClient,
	}, nil
}

func (c *client) ServiceName() string {
	return serviceName
}

func (c *client) Push(msg message.Message) error {
	_, err := c.redditClient.Message.Send(context.Background(), &reddit.SendMessageRequest{
		To:      msg.To,
		Subject: msg.Title,
		Text:    msg.Content,
	})

	if err != nil {
		if redditError, ok := err.(*reddit.ErrorResponse); ok {
			return redditError
		}

		return err
	}

	return nil
}
