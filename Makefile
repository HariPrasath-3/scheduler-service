PWD := $(shell pwd)
PROTO_DIR := $(PWD)/proto
PROTO_OUT := $(PWD)

.PHONY: run build test tidy docker-up docker-down docker-logs dynamo-create proto proto-clean

run:
	go run ./cmd/main

build:
	go build ./cmd/main

docker-up:
	docker compose up --build

docker-down:
	docker compose down

docker-logs:
	docker compose logs -f

dynamo-create:
	./scripts/create-dynamo-table.sh

proto:
	protoc \
		--proto_path=$(PWD) \
		--go_out=$(PROTO_OUT) \
		--go_opt=paths=source_relative \
		--go-grpc_out=$(PROTO_OUT) \
		--go-grpc_opt=paths=source_relative \
		$(shell find $(PROTO_DIR) -name "*.proto")

proto-clean:
	find $(PWD)/proto -name "*.pb.go" -delete