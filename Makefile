# Build vars
export COMMIT_HASH := $(shell git rev-parse --short HEAD)
export VERSION := $(shell git rev-parse --abbrev-ref HEAD)
export BUILD_DATE := $(shell date +%Y%m%d%H%M%S)

## install project tools
install:
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.58.2
	@go install go.uber.org/mock/mockgen@v0.4.0

## test: run unit/mock tests
test: generate
	go test -v ./...

## unit-test-only: run unit tests without any dependent step
unit-test-only:
	go test -failfast ./...

## generate: runs go generate
generate:
	go generate ./...

## clean-mock: removes all generated mocks
clean-mock:
	find . -iname '*_mock.go' -exec rm {} \;

## update: runs go mod vendor and tidy
update: tidy

## tidy: runs go mod tidy
tidy:
	go mod tidy -v

lint:
	golangci-lint run -v ./...

standalone:
	CGO_ENABLED=0 GOOS=linux go build -o load ./loader/load/load.go

