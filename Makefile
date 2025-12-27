.PHONY: run build test tidy docker-up docker-down docker-logs

run:
	go run ./cmd/main

build:
	go build ./cmd/main

test:
	go test ./...

tidy:
	go mod tidy

docker-up:
	docker compose up --build

docker-down:
	docker compose down

docker-logs:
	docker compose logs -f