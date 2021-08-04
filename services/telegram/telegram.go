package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/kamikazechaser/hedwig/internal/message"
	"github.com/yi-jiayu/ted"
)

const (
	serviceName string = "telegram"
)

type config struct {
	Token string `json:"token"`
}

type client struct {
	config     *config
	httpClient *http.Client
}

// New initializes the service client
func New(jsonConfig []byte) (interface{}, error) {
	var conf *config

	if err := json.Unmarshal(jsonConfig, &conf); err != nil {
		return nil, err
	}

	// Ensure all api keys are present in the config file before attempting to push notifications
	if conf.Token == "" {
		return nil, errors.New("telegram token not provided")
	}

	httpClient := &http.Client{
		Timeout: time.Second * 15,
	}

	return &client{
		config:     conf,
		httpClient: httpClient,
	}, nil
}

func (c *client) ServiceName() string {
	return serviceName
}

func (c *client) Push(msg message.Message) error {
	tg := ted.Bot{
		Token:      c.config.Token,
		HTTPClient: c.httpClient,
	}

	pushMessage := ted.SendMessageRequest{
		ChatID:    msg.To,
		Text:      fmt.Sprintf("*%s*\n\n%s", msg.Title, msg.Content),
		ParseMode: "Markdown",
	}

	_, err := tg.Do(pushMessage)
	if err != nil {
		if tgError, ok := err.(ted.Response); ok {
			return tgError
		}

		return err
	}

	return nil
}
