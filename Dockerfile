FROM golang:1.16-alpine AS builder

RUN apk update && apk add gcc libc-dev make git

WORKDIR /hedwig-build/
COPY . .
ENV CGO_ENABLED=1 GOOS=linux
RUN --mount=type=cache,target=/root/.cache/go-build \
    make build && \
    go build -ldflags="-s -w" -buildmode=plugin -o build/telegram.plugin _plugins/telegram/telegram.go

FROM alpine:3.15

WORKDIR /hedwig

COPY --from=builder /hedwig-build/config.toml /hedwig-build/hedwig /hedwig-build/build/ /hedwig/

EXPOSE 3000
CMD ["/hedwig/hedwig"]
