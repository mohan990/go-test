# Go Webserver + Docker + GitHub Actions CI

This repo is focused on:
1. Local development
2. Unit and integration testing
3. Docker build/run
4. GitHub Actions CI

## Project structure

- `cmd/server/main.go` - Go HTTP server
- `cmd/server/main_test.go` - unit tests
- `cmd/server/integration_test.go` - integration tests
- `Dockerfile` - standard image
- `Dockerfile.dev` - development image
- `Makefile` - local run/test/docker commands
- `.github/workflows/ci.yml` - GitHub Actions CI pipeline

## API endpoints

- `GET /`
- `GET /hello?name=YourName`
- `GET /healthz`

## Prerequisites

- Go 1.22+
- Docker
- GitHub account

## Local run

```bash
make run
```

Test quickly:

```bash
curl http://localhost:8080/
curl "http://localhost:8080/hello?name=HP"
curl http://localhost:8080/healthz
```

## Tests

Run unit tests:

```bash
make test-unit
```

Run integration tests:

```bash
make test-integration
```

Run all tests:

```bash
make test
```

## Docker

Build/run standard image (`Dockerfile`):

```bash
make docker-build
make docker-run
curl http://localhost:8080/healthz
```

Build/run development image (`Dockerfile.dev`):

```bash
make docker-build-dev
make docker-run-dev
curl http://localhost:8080/healthz
```

Stop containers:

```bash
make docker-stop
```

Build with your Docker Hub repo:

```bash
make docker-build IMAGE=<dockerhub-username>/go-webserver TAG=local
```

## GitHub setup

```bash
git init
git add .
git commit -m "Initial Go server with tests and CI"
git branch -M main
git remote add origin https://github.com/<your-github-username>/<your-repo>.git
git push -u origin main
```

## GitHub Actions CI

Workflow file:
- `.github/workflows/ci.yml`

On every push and pull request to every branch, CI runs:
1. `make test-unit`
2. `make test-integration`
3. `make docker-build IMAGE=local/go-webserver TAG=ci`
