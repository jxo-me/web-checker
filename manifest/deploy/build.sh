#!/usr/bin/env bash

function build() {
    docker run --rm -i --name checker-build \
      -v "$(pwd)":/go/src/work \
      -v "$GOPATH/pkg/mod":/go/src/mod \
      -w /go/src/work \
      golang:1.18-alpine sh -c "go mod tidy && go build -y -o checker cmd/main.go"
}

build