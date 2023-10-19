#!/usr/bin/env bash

function build() {
    docker run --rm -i --name checker-build \
      -v "$(pwd)":/go/src/work \
      -v "$GOPATH/pkg/mod":/go/src/mod \
      -w /go/src/work \
      golang:1.20-alpine sh -c "go mod tidy && go build -o checker cmd/main.go"
}

build && docker build -f manifest/docker/Dockerfile -t "web-checker" .