LAST_COMMIT := $(shell git rev-parse --short HEAD)
BUILDSTR := ${LAST_COMMIT}
GIN_MODE := release

BIN := hedwig

REDDIT_BIN := reddit.svc
TELEGRAM_BIN := telegram.svc

.PHONY: build
build:
	go build -ldflags="-s -w" -buildmode=plugin -o ${REDDIT_BIN} services/reddit/reddit.go
	go build -ldflags="-s -w" -buildmode=plugin -o ${TELEGRAM_BIN} services/telegram/telegram.go

	go build -o ${BIN} -ldflags="-s -w -X 'main.version=${BUILDSTR}'" core/*.go

clean:
	go clean
	- rm -f ${BIN} ${REDDIT_BIN} ${TELEGRAM_BIN}
