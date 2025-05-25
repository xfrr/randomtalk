#!/usr/bin/env just --justfile
default: clean go-format go-lint

# start the docker-compose stack. supports [nats,]
up stream_system:
  @docker-compose \
    -f ./deployments/docker/docker-compose.yml \
    -f ./deployments/docker/{{stream_system}}/docker-compose.yml \
    up -d --build --force-recreate

up-xray stream_system:
  @docker-compose \
    -f ./deployments/docker/docker-compose.yml \
    -f ./deployments/docker/{{stream_system}}/docker-compose.yml \
    -f ./deployments/docker/observability.docker-compose.yml \
    up -d --build --force-recreate

# stop the docker-compose stack. supports [nats,]
down stream_system="nats":
  @docker-compose \
    -f ./deployments/docker/docker-compose.yml \
    -f ./deployments/docker/{{stream_system}}/docker-compose.yml \
    down --remove-orphans --volumes --rmi all --timeout 0

down-xray stream_system:
  @docker-compose \
    -f ./deployments/docker/docker-compose.yml \
    -f ./deployments/docker/{{stream_system}}/docker-compose.yml \
    -f ./deployments/docker/observability.docker-compose.yml \
    down --remove-orphans

# build Go applications
go-build:
  @go build ./...

# generate Go code
go-generate:
  @go generate ./...

# generate protobuf stubs
go-generate-protobuf:
  @cd proto && buf dep update && buf generate

# test Go codebase
go-tests:
  @go test ./... -v -race \
    -timeout=30s -count=1 \
    -cover -coverpkg=./... -covermode=atomic \
    -coverprofile=.coverage

# run Go integration tests
go-integration-tests:
  @go test ./... -v -tags=integration \
    -race -timeout=30s -count=1 \
    -cover -coverpkg=./... -covermode=atomic \
    -coverprofile=.integration.coverage

go-load-tests:
  @ghz --insecure --async \
    --count-errors \
    --config ./tests/load/ghz-config.json \
    -c 20 -n 1000

# lint Go codebase
go-lint:
  @golangci-lint run

# format the Go codebase
go-format:
  # format Go code
  @gofmt -s -w .
  # format swag docs
  @swag fmt --dir ./

# clean up the project's build and test artifacts
clean:
  @rm -rf bin
  @rm -rf tmp
  @rm -rf vendor
  @rm -rf .env
  @rm -rf .cache
  @rm -rf .coverage

# display all available commands
help:
  @just --list

