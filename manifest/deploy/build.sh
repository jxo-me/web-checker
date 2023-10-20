#!/usr/bin/env bash

function build() {
docker run --rm -i --name checker-build \
  -v "$(pwd)":/go/src/work \
  -v "$GOPATH/pkg/mod":/go/src/mod \
  --net=host \
  -w /go/src/work \
  golang:1.20-alpine sh -c "go mod tidy && go build -o checker cmd/main.go"
}

# docker run -d --restart always --name web-checker --net=host web-checker:latest
function start() {
docker run -d --restart always \
 --name web-checker \
 --net=host \
 web-checker:latest
}

function cleanImages() {
  if [ $(docker images -q --filter "dangling=true") ]; then
    # clean untagged images
    docker rmi -f $(docker images -q --filter dangling=true)
  fi
}

build && docker build -f manifest/docker/Dockerfile -t "web-checker" .