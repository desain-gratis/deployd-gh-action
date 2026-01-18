ifneq (,$(wildcard .env))
    include .env
    export $(shell sed 's/=.*//' .env)
endif

build:
	GOOS=linux GOARCH=amd64 go build -o dist/action-linux-amd64 ./cmd/action/main.go

run-local:
	GOOS=linux GOARCH=amd64 go run ./cmd/action/main.go

run:
	dist/action-linux-amd64
