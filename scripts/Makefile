.PHONY: init run clean build-linux build-mac run-test run-test-html
GIT_COMMIT = $(shell git rev-parse --short HEAD)
BUILD_LDFLAGS = "-s -w -X main.GitCommit=$(GIT_COMMIT)"

deps:
	@go mod download

run: deps
	@go run main.go

clean:
	@rm -rf build

build-linux: deps
	@GOOS=linux GOARCH=amd64 go build -ldflags=$(BUILD_LDFLAGS) -o ./app -v

# go test may execute parallel. if parallel, tests used database fixture would be failed.
run-test:
	@go test ./... -cover -p 1
run-test-html:
	@go test ./... -coverprofile=cover.out -p 1 && go tool cover -html=cover.out -o cover.html && rm cover.out
