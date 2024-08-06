# Unit test coverage floor percentage (CI should fail if coverage is below this)
COV_CUTOVER := "50"

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
	go build -o shared_library/ir-standings.so -buildmode=c-shared shared_library/main.go

## clean-mock: removes all generated mocks
clean-mock:
	find . -iname '*_mock.go' -exec rm {} \;

## update: runs go mod vendor and tidy
update: tidy

## tidy: runs go mod tidy
tidy:
	go mod tidy -v

lint: generate
	golangci-lint run -v ./...

## cover: run unit/mock tests with coverage report. Generated mocks are filtered out of the report
cover: generate
	go test -failfast -count=2 --race -coverprofile=coverage.out -coverpkg=./... ./...
	cat coverage.out | grep -v "_mock.go" | grep -v redact.go > coverage.nomocks.out
	go tool cover -func coverage.nomocks.out

## cover-check: checks the code coverage to be beyond a certain threshold
cover-check: cover
	COV_CUTOVER=${COV_CUTOVER} ./.github/cover-check.sh

standalone:
	CGO_ENABLED=0 GOOS=linux go build -o load ./loader/load/load.go

## Shared library to statically link
shared:
	go build -o shared_library/libgoir.a -buildmode=c-archive shared_library/main.go
	gcc -g -o shared_library/client shared_library/client.c shared_library/latest-standings.c  shared_library/libgoir.a

## Windows

#  go build -o libgoir.lib -buildmode=c-archive .\main.go
#  gcc -g -c .\latest-standings.c
#  lib libgoir.lib latest-standings.o
#  cp libgoir.lib windows/
#  dumpbin /symbols libgoir.lib | findstr "Stand"
# PS C:\Users\Ian\ir-standings\shared_library> dumpbin /symbols .\libgoir.lib | findstr "Stand"
# 002 00000000 SECT1  notype ()    External     | LiveStandings
# 005 00000000 SECT10 notype       Static       | .rdata$.refptr._cgoexp_fc2f6fae3e8a_LiveStandings
# 026 00000000 SECT10 notype       External     | .refptr._cgoexp_fc2f6fae3e8a_LiveStandings
# 02B 00000000 UNDEF  notype       External     | _cgoexp_fc2f6fae3e8a_LiveStandings
# 4F0 00066000 SECT1  notype ()    External     | _cgoexp_fc2f6fae3e8a_LiveStandings
# 002 00000000 SECT1  notype ()    External     | GoLatestStandings
# 020 00000000 UNDEF  notype ()    External     | LiveStandings
#           10 .rdata$.refptr._cgoexp_fc2f6fae3e8a_LiveStandings
# PS C:\Users\Ian\ir-standings\shared_library>
#
# Add libgoir.lib to linker options
# Needed to add msvcrt.lib and remove libcmt.lib  (d suffixes for debug)