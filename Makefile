
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

build:
	make gen
	go build -o ./bin/app ./cmd/app

run:
	make build
	./bin/app

PB=pb
PROTO=proto/.proto
protobuf:
	protoc  \
		--go_out=. \
		--go_opt=M$(PROTO)=$(PB)/userspb \
		--go-grpc_out=. \
		--go-grpc_opt=M$(PROTO)=$(PB)/userspb \
		$(PROTO)

wire-gen:
	wire ./internal/app/

gen:
	make protobuf
	make wire-gen

migrate.up:
	migrate -path ./migrations -database 'postgres://$(PG_USER):$(PG_PASS)@$(PG_HOST):$(PG_PORT)/$(PG_NAME)?sslmode=disable' up

migrate.down:
	migrate -path ./migrations -database 'postgres://$(PG_USER):$(PG_PASS)@$(PG_HOST):$(PG_PORT)/$(PG_NAME)?sslmode=disable' down