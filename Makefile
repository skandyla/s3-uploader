# Makefile
.PHONY: build
.DEFAULT_GOAL := all
all: lint test

IMAGENAME=microtester
PROJECT ?= github.com/skandyla/${IMAGENAME}
BUILD_VERSION=$(shell git describe --tags --always)
BUILD_COMMIT?=$(shell git rev-parse --short HEAD)
BUILD_TIME?=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
BUILD_BRANCH=$(shell git rev-parse --symbolic-full-name --abbrev-ref HEAD)

lint:
	golangci-lint run

install_linter:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.43.0

get_deps:
	go mod download

test:
	go test -v ./...
	go test -race ./...

cover:
	go test -v ./... -coverprofile=/tmp/c.out
	go tool cover -func=/tmp/c.out


GO_VERSION_FLAGS= -X ${PROJECT}/internal/version.BuildVersion=${BUILD_VERSION} -X ${PROJECT}/internal/version.BuildCommit=${BUILD_COMMIT} -X ${PROJECT}/internal/version.BuildTime=${BUILD_TIME} -X ${PROJECT}/internal/version.BuildBranch=${BUILD_BRANCH}

build:
	@echo MARK: build go code
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -v --ldflags "-extldflags '-static' ${GO_VERSION_FLAGS}" -o main


docker_image:
	@echo MARK: package build code inside docker image
	docker version
	docker build -t ${IMAGENAME} --build-arg GIT_COMMIT=${BUILD_COMMIT}  -f build/Dockerfile .

docker_compose_start:
	@echo "MARK: Starting docker-compose"
	@docker compose up -d
	@echo "MARK: Sleeping 5 seconds until dependencies start."
	@sleep 5

docker_compose_stop:
	@echo MARK: Stopping cluster
	@docker compose down --remove-orphans

# example. go integration tests TBD
docker_compose_run_tests:
	@echo MARK: Testing via docker-compose 
	curl localhost:8080/liveness
	curl localhost:8080/__service/info
