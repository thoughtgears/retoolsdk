#!make
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

.PHONY: test example

example:
	@go run example/main.go

test:
	@go test -v ./...
