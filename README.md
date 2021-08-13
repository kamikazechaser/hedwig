<img align="left" src="https://i.imgur.com/fquebpE.png" width="90" height="90">

### Hedwig

> Standalone server for sending out notifications with pluggable service providers

---

[![Go Report Card](https://goreportcard.com/badge/github.com/kamikazechaser/hedwig)](https://goreportcard.com/report/github.com/kamikazechaser/hedwig)
[![license: Unlicense](https://img.shields.io/badge/license-Unlicense-brightgreen)](https://opensource.org/licenses/Unlicense)

### Background

This is a replacement for one of my project's notification server which was crudely based on serverless functions. The older implementation suffered from poor resilliancy and zero persistance. The other alternatives ([shove](https://github.com/pennersr/shove)) are not suitable for my use case where I change service providers regularly and want finer control over how messages are created.

The current design is based on Go plugins and follow a specific [`plugin spec`](https://github.com/kamikazechaser/hedwig/blob/master/internal/svcplugin/svcplugin.go). This allows for high flexibility and control over individual service clients. The compiled plugins are loaded during runtime and individual plugin configuration pulled from the [`config file`](https://github.com/kamikazechaser/hedwig/blob/master/config.example.json).

Hedwig exposes a protected HTTP endpoint to enqueue incoming messages which are required to follow a common [`message spec`](https://github.com/kamikazechaser/hedwig/blob/master/internal/message/message.go).

**Status**: This project is currently undergoing **development**. Expect breaking changes (especially to plugins).

### Features

- Service providers are loaded as compiled [Go plugins](https://pkg.go.dev/plugin)
- Underlying resiliancy and persistance provided by [asynq](https://github.com/hibiken/asynq)
- Scheduled notifications
- Queue from any source that talks HTTP

### Plugins

Hedwig comes with a couple of sample plugins to give you an idea of how to write your own:

- Telegram
- Reddit
- Mailgun
