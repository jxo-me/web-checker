.EXPORT_ALL_VARIABLES:

CONFIG=./example/config_local.yaml

run:
	go run cmd/main.go

build:
	chmod a+x manifest/deploy/build.sh && ./manifest/deploy/build.sh

update:
	docker rm -f web-checker && docker run -d --restart always --name web-checker --net=host web-checker:latest