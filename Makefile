.PHONY: run
run:
	go run cmd/bot/main.go

build:
	go build -o bin/bot cmd/bot/main.go

LOCAL_BIN:=$CURDIR/bin
.PHONY: .deps
.deps:
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway && \
	go install google.golang.org/protobuf/cmd/protoc-gen-go && \
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc

.PHONY: generate
generate:
	buf mod update
	buf generate


MIGRATIONS_DIR=./migrations
.PHONY: migration
migration:
	goose -dir=${MIGRATIONS_DIR} create $(NAME) sql