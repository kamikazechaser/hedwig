package main

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/kamikazechaser/hedwig/internal/message"
	"github.com/mailgun/mailgun-go/v4"
)

const (
	serviceName string = "mailgun"
)

type config struct {
	Domain     string `json:"domain"`
	PrivateKey string `json:"private_key"`
	From       string `json:"from"`
}

type client struct {
	config        *config
	mailgunClient *mailgun.MailgunImpl
}

// New initializes the service client
func New(jsonConfig []byte) (interface{}, error) {
	var conf *config

	if err := json.Unmarshal(jsonConfig, &conf); err != nil {
		return nil, err
	}

	// Ensure all configs are present in the config file before attempting to push notifications
	if conf.Domain == "" || conf.PrivateKey == "" || conf.From == "" {
		return nil, errors.New("mailgun config is incomplete")
	}

	mg := mailgun.NewMailgun(conf.Domain, conf.PrivateKey)

	return &client{
		config:        conf,
		mailgunClient: mg,
	}, nil
}

func (c *client) ServiceName() string {
	return serviceName
}

func (c *client) Push(msg message.Message) error {
	// very specific push configuration, won't work for most
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	mail := c.mailgunClient.NewMessage(c.config.From, msg.Title, "", msg.To)
	mail.SetTemplate("master")

	mail.AddVariable("link", msg.Params[0])
	mail.AddVariable("action", msg.Params[1])
	mail.AddVariable("message", msg.Content)

	_, _, err := c.mailgunClient.Send(ctx, mail)

	if err != nil {
		return err
	}

	return nil
}
