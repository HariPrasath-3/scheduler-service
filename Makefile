PWD := $(shell pwd)
PROTO_DIR := $(PWD)/proto
PROTO_OUT := $(PWD)/client/golang/proto

.PHONY: run build test tidy docker-up docker-down docker-logs proto proto-clean

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

proto:
	mkdir -p $(PROTO_OUT)
	protoc \
		--proto_path=$(PROTO_DIR) \
		--go_out=$(PROTO_OUT) \
		--go-grpc_out=$(PROTO_OUT) \
		$(shell find $(PROTO_DIR) -name "*.proto")

proto-clean:
	rm -rf $(PROTO_OUT)