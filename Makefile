LAST_COMMIT := $(shell git rev-parse --short HEAD)
BUILDSTR := ${LAST_COMMIT}
GIN := jsoniter
SERVICES := reddit telegram mailgun
BIN := hedwig

.PHONY: build
build:
	$(foreach service,$(SERVICES),go build -ldflags="-s -w" -buildmode=plugin -o $(service).svc services/$(service)/$(service).go;)
	go build -tags=${GIN} -o ${BIN} -ldflags="-s -w -X 'main.version=${BUILDSTR}'" core/*.go

clean:
	go clean
	$(foreach service,$(SERVICES),rm -f $(service))
	rm -f ${BIN}
