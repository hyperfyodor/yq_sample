GIT_TAG := 0.0.1
#GIT_TAG := $(shell git describe --tags --abbrev=0)
MODULE_NAME := $(shell go list -m)
LDFLAGS := -X $(MODULE_NAME)/pkg.Version=$(GIT_TAG)

sqlc_gen:
	cd db && sqlc generate

proto_gen:
	cd ./proto/consumer && protoc --go_out=./gen --go_opt=paths=source_relative \
                               --go-grpc_out=./gen --go-grpc_opt=paths=source_relative \
                           	consumer.proto

producer:
	go run -ldflags "$(LDFLAGS)" ./cmd/producer

consumer:
	go run -ldflags  "$(LDFLAGS)" ./cmd/consumer

migrate:
	go run -ldflags  "$(LDFLAGS)" ./cmd/migrator

compose:
	cd docker && docker-compose up -d

compose_clean:
	cd docker && docker-compose down -v

explain:
	go run ./cmd/consumer -cfg_explain
	go run ./cmd/producer -cfg_explain
	go run ./cmd/migrator -cfg_explain

version:
	go run -ldflags "$(LDFLAGS)" ./cmd/consumer -version
	go run -ldflags "$(LDFLAGS)" ./cmd/producer -version
	go run -ldflags "$(LDFLAGS)" ./cmd/migrator -version