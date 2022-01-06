LAST_COMMIT := $(shell git rev-parse --short HEAD)
BUILDSTR := ${LAST_COMMIT}
GIN := jsoniter
BIN := hedwig

.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux go build -tags=${GIN} -o build/${BIN} -ldflags="-s -w -X 'main.version=${BUILDSTR}'" cmd/*.go
