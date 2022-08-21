SHALL=/bin/bash

export CGO_ENABLED=0
export DSN=psql://gouser:gopassword@localhost:5432/gotest

default: build
.PHONY: default

build:
	@ echo "-> build binary ..."
	@ go build -ldflags "-X main.HashCommit=`git rev-parse HEAD` -X main.BuildStamp=`date -u '+%Y-%m-%d_%I:%M:%S%p'`" -o ./calendar/cmd/server .
.PHONY: build

proto:
	@echo "-> generate protobuff and grpc resources ..."
	@ protoc --proto_path=proto proto/*.proto --go_out=. --go-grpc_out=.
.PHONY: proto

test:
	@ echo "-> running unit tests ..."
	@ CGO_ENABLED=1 go test -race ./...
.PHONY: test

test-it:
	@ echo "-> running all tests ..."
	@ CGO_ENABLED=1 go test -race ./... --tags=integration
.PHONY: test-it

lint:
	@ echo "-> running linters ..."
	@ golangci-lint run ./...
.PHONY: lint