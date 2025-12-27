.PHONY: run build test tidy

run:
	go run ./cmd

build:
	go build ./cmd

tidy:
	go mod tidy
