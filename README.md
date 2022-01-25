<img align="left" src="https://i.imgur.com/fquebpE.png" width="90" height="90">

### Hedwig

> Standalone server for sending out notifications with pluggable service providers

---

[![Go Report Card](https://goreportcard.com/badge/github.com/kamikazechaser/hedwig)](https://goreportcard.com/report/github.com/kamikazechaser/hedwig)
[![license: Unlicense](https://img.shields.io/badge/license-Unlicense-brightgreen)](https://opensource.org/licenses/Unlicense)

### Background

The current design is based on Go plugins and follow a specific [`plugin spec`](https://github.com/kamikazechaser/hedwig/blob/master/internal/svcplugin/svcplugin.go). The compiled plugins are loaded during runtime. Any plugin specific configuration is baked into the plugin itself.

Hedwig exposes a protected HTTP endpoint to enqueue incoming messages which are required to follow a common [`message spec`](https://github.com/kamikazechaser/hedwig/blob/master/internal/message/message.go).

**Status**: This project is currently undergoing **development**. Expect breaking changes (especially to plugins).

### Features

- Service providers are loaded as compiled [Go plugins](https://pkg.go.dev/plugin)
- Underlying resiliancy and persistance provided by [asynq](https://github.com/hibiken/asynq)
- Scheduled notifications
- Queue from any source that talks HTTP

### Plugins

A sample plugin would look like:

```go
package main

import (
	"context"
	"time"

	"github.com/kamikazechaser/hedwig/internal/message"
	"github.com/mailgun/mailgun-go/v4"
)

const (
	pluginName string = "mailgun"
)

type client struct {
	mailgunClient *mailgun.MailgunImpl
}

func New() (interface{}, error) {
	mg := mailgun.NewMailgun("", "")

	return &client{
		mailgunClient: mg,
	}, nil
}

func (c *client) PluginName() string {
	return pluginName
}

func (c *client) HealthCheck() bool {
    // check if you have enough credits e.t.c.
    // TODO: not implemnted in Hedwig yet
	return true
}

func (c *client) Push(msg message.Message) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	mail := c.mailgunClient.NewMessage("Domain <noreply@your.domain>", msg.Title, "", msg.To)
	mail.SetTemplate("master")

    // msg.Params for additional properties that you want to include
	mail.AddVariable("link", msg.Params[0])
	mail.AddVariable("action", msg.Params[1])
	mail.AddVariable("message", msg.Content)

	_, _, err := c.mailgunClient.Send(ctx, mail)

	if err != nil {
        // return errors only if you want to keep retrying
        // return nil in cases where the service is blocked by the receipient e.t.c.
		return err
	}

	return nil
}

```

### Usage

**prerequisites**

- Go 1.6
- Make
- Redis server (add to `config.json`)

```bash
$ cp config.example.json config.json
# Edit config.json
$ make build
$ ./hedwig
# Hedwig will look for plugins in the root folder and load them at run time
```

Queue notifications by sending a request in the format:

Request:

```
> POST /push?key=test HTTP/1.1
> Host: localhost:3000
> Content-Type: application/json
> Accept: */*

| {
| 	"to": "135207785",
| 	"service": "telegram",
| 	"title": "hello from hedwig",
| 	"content": "test123",
| 	"delay": 1,
| }
```

Response:

```
< HTTP/1.1 200 OK
< Content-Type: application/json; charset=utf-8

| {
|   "message": "6c9d0793-41e8-438c-9181-170af6a6da7f",
|   "ok": true
| }
```
