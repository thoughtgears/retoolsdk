#!make
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

.PHONY: lint test

lint:
	@golangci-lint run

test:
	@go test -v ./...
