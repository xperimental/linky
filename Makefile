.PHONY: all test build-binary install clean

GO ?= go
GO_CMD := CGO_ENABLED=0 $(GO)
GIT_VERSION := $(shell git describe --tags --dirty)
VERSION := $(GIT_VERSION:v%=%)
GIT_COMMIT := $(shell git rev-parse HEAD)

all: test build-binary

test:
	$(GO_CMD) test -cover ./...

build-binary:
	$(GO_CMD) build -tags netgo -ldflags "-w -X main.Version=$(VERSION) -X main.GitCommit=$(GIT_COMMIT)" -o linky .

build-all:
	./build-all.sh

install:
	install -D -t $(DESTDIR)/usr/bin/ linky

docker-image:
	docker build -t xperimental/linky .

clean:
	rm -f linky
