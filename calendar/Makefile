SHALL=/bin/bash

export CGO_ENABLED=0
export DSN=psql://gouser:gopassword@localhost:5432/gotest

default: build
.PHONY: default

build:
	@ echo "-> build binary ..."
	@ go build -ldflags "-X main.HashCommit=`git rev-parse HEAD` -X main.BuildStamp=`date -u '+%Y-%m-%d_%I:%M:%S%p'`" -o ./calendar .
.PHONY: build

test:
	@ echo "-> running tests ..."
	@ CGO_ENABLED=1 go test -race ./...
.PHONY: test

lint:
	@ echo "-> running linters ..."
	@ golangci-lint run ./...
.PHONY: lint