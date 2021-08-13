<img align="left" src="https://i.imgur.com/fquebpE.png" width="90" height="90">

### Hedwig

> Standalone server for sending out notifications with pluggable service providers
---

[![Go Report Card](https://goreportcard.com/badge/github.com/kamikazechaser/hedwig)](https://goreportcard.com/report/github.com/kamikazechaser/hedwig)
[![license: Unlicense](https://img.shields.io/badge/license-Unlicense-brightgreen)](https://opensource.org/licenses/Unlicense)


### Background

This is a replacement for one of my project's notification server which is crudely based on serverless functions. The older implementation suffered from poor resilliancy and zero persistance. This project is heavily inspired by [shove](https://github.com/pennersr/shove). Unfortunately shove is limited to FCM, Telegram and Email which makes it unsuitable for my use case where I change service providers regularly.

The current design makes it as easy as writing a plugin following the common [`plugin spec`](https://github.com/kamikazechaser/hedwig/blob/master/internal/svcplugin/svcplugin.go), compiling it and loading it at run-time. Managing the plugins is entirely done inside a `config.json` file.

All incoming messages are required to follow a common `message spec`.

**Status**: This project is currently undergoing **development**. Expect breaking changes (especially to plugins).

### Features

- Service providers are loaded as compiled [Go plugins](https://pkg.go.dev/plugin)
- Underlying resiliancy and persistance provided by [asynq](https://github.com/hibiken/asynq)
- Scheduled notifications with optional recall option
- Queue from any source that talks HTTP
- Fanout to multiple service providers from a single request

### Plugins

Hedwig comes with a couple of sample plugins to give you an idea of how to write your own:

- Telegram
- Reddit
- Mailgun
- FCM
- Discord
