APP_NAME=go-webserver
IMAGE ?= local/go-webserver
TAG ?= local

.PHONY: help run test test-unit test-integration docker-build docker-run docker-stop

help:
	@echo "make run              - run server locally"
	@echo "make test-unit        - run unit tests"
	@echo "make test-integration - run integration tests"
	@echo "make test             - run all tests"
	@echo "make docker-build     - build docker image"
	@echo "make docker-run       - run docker container"
	@echo "make docker-stop      - stop docker container"

run:
	go run ./cmd/server

test-unit:
	go test ./... -run '^TestUnit' -v

test-integration:
	go test ./... -run '^TestIntegration' -v

test:
	go test ./... -v

docker-build:
	docker build -t $(IMAGE):$(TAG) .

docker-run:
	docker run --rm -d --name $(APP_NAME) -p 8080:8080 $(IMAGE):$(TAG)

docker-stop:
	docker stop $(APP_NAME)
