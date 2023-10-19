.EXPORT_ALL_VARIABLES:

CONFIG=./example/config_local.yaml

run:
	go run cmd/main.go

build:
	go build -o checker cmd/main.go