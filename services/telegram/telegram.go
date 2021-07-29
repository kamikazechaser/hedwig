package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/kamikazechaser/hoodwink/internal/message"
)

const (
	serviceName string = "telegram"
)

type config struct {
	Token string `json:"token"`
}

type client struct {
	config *config
	request *http.Client
}

type telegramResponse struct {
	Ok bool `json:"ok"`
	Result interface{} `json:"result"`
}

func New(jsonConfig map[string]interface{}) (interface{}, error) {
	var conf *config

	jsonString, err := json.Marshal(jsonConfig)

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(jsonString, &conf); err != nil {
		return nil, err
	}

	if conf.Token == "" {
		return nil, errors.New("telegram token not provided")
	}

	// replace with asynq queue builder
	httpClient := &http.Client{
		Timeout: time.Duration(5) * time.Second,
	}

	return &client{
		config: conf,
		request: httpClient,
	}, nil
}

func (c *client) ServiceName() string {
	return serviceName
}

func (c *client) Push(msg message.Message) error {
	fmt.Println(c.config.Token)

	var tgEndpoint = "https://api.telegram.org/bot" + c.config.Token + "/sendMessage"

	resp, err := c.request.PostForm(
		tgEndpoint,
		url.Values{
			"chat_id": {msg.To},
			"text": {msg.Title + "\n\n" + msg.Content},
		})

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	tgRes := telegramResponse{}

	if err := json.Unmarshal(b, &tgRes); err != nil {
		return err
	}

	fmt.Printf("%v", tgRes)

	if !tgRes.Ok {
		return errors.New("telegram side error")
	}

	return nil
}