include .env

GIT_TAG := 0.0.1
#GIT_TAG := $(shell git describe --tags --abbrev=0)
MODULE_NAME := $(shell go list -m)
LDFLAGS := -X $(MODULE_NAME)/pkg.Version=$(GIT_TAG) -s -w

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

test_unit:
	go test ./internal...

test_integration:
	go test ./test -count=1 -v

explain:
	go run ./cmd/consumer -cfg_explain
	go run ./cmd/producer -cfg_explain
	go run ./cmd/migrator -cfg_explain

version:
	go run -ldflags "$(LDFLAGS)" ./cmd/consumer -version
	go run -ldflags "$(LDFLAGS)" ./cmd/producer -version
	go run -ldflags "$(LDFLAGS)" ./cmd/migrator -version

consumer_cpu:
	go tool pprof -http=127.0.0.1:9999 http://localhost:${CSM_PROFILING_PORT}/debug/pprof/profile
consumer_goroutine:
	go tool pprof -http=127.0.0.1:9998 http://localhost:${CSM_PROFILING_PORT}/debug/pprof/goroutine
consumer_heap:
	go tool pprof -http=127.0.0.1:9997 http://localhost:${CSM_PROFILING_PORT}/debug/pprof/heap

consumer_profiles:
	curl -o ./profiles/csm_cpu.pprof http://localhost:${CSM_PROFILING_PORT}/debug/pprof/profile
	curl -o ./profiles/csm_goroutine.pprof http://localhost:${CSM_PROFILING_PORT}/debug/pprof/goroutine
	curl -o ./profiles/csm_heap.pprof http://localhost:${CSM_PROFILING_PORT}/debug/pprof/heap

producer_cpu:
	go tool pprof -http=127.0.0.1:9999 http://localhost:${PRD_PROFILING_PORT}/debug/pprof/profile
consumer_goroutine:
	go tool pprof -http=127.0.0.1:9998 http://localhost:${PRD_PROFILING_PORT}/debug/pprof/goroutine
consumer_heap:
	go tool pprof -http=127.0.0.1:9997 http://localhost:${PRD_PROFILING_PORT}/debug/pprof/heap

producer_profiles:
	curl -o ./profiles/prd_cpu.pprof http://localhost:${PRD_PROFILING_PORT}/debug/pprof/profile
	curl -o ./profiles/prd_goroutine.pprof http://localhost:${PRD_PROFILING_PORT}/debug/pprof/goroutine
	curl -o ./profiles/prd_heap.pprof http://localhost:${PRD_PROFILING_PORT}/debug/pprof/heap
