APP_NAME=go-webserver
IMAGE ?= local/go-webserver
TAG ?= local

.PHONY: help run test test-unit test-integration docker-build docker-run docker-stop docker-build-dev docker-run-dev

help:
	@echo "make run              - run server locally"
	@echo "make test-unit        - run unit tests"
	@echo "make test-integration - run integration tests"
	@echo "make test             - run all tests"
	@echo "make docker-build     - build production image (Dockerfile)"
	@echo "make docker-run       - run production container"
	@echo "make docker-build-dev - build development image (Dockerfile.dev)"
	@echo "make docker-run-dev   - run development container"
	@echo "make docker-stop      - stop container"

run:
	go run ./cmd/server

test-unit:
	go test ./... -run '^TestUnit' -v

test-integration:
	go test ./... -run '^TestIntegration' -v

test:
	go test ./... -v

docker-build:
	docker build -f Dockerfile -t $(IMAGE):$(TAG) .

docker-run:
	docker run --rm -d --name $(APP_NAME) -p 8080:8080 $(IMAGE):$(TAG)

docker-build-dev:
	docker build -f Dockerfile.dev -t $(IMAGE):dev .

docker-run-dev:
	docker run --rm -d --name $(APP_NAME)-dev -p 8080:8080 $(IMAGE):dev

docker-stop:
	-docker stop $(APP_NAME)
	-docker stop $(APP_NAME)-dev
