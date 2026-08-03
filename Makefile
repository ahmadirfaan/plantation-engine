export PATH := $(PATH):$(shell go env GOPATH)/bin

.PHONY: clean all init generate test test_api integration-test

all: build/main

build/main: cmd/main.go generated
	@echo "Building..."
	go build -o $@ $<

clean:
	rm -rf generated

init: clean generate
	go mod tidy
	go mod vendor

test:
	go clean -testcache
	go test -short -coverprofile coverage.out -v $(shell go list ./... | grep -v /generated)
	go tool cover -func=coverage.out

test_api:
	go clean -testcache
	go test ./tests/...

integration-test:
	go clean -testcache
	TESTCONTAINERS=1 go test -v ./integration/...

generate: generated

generated: api.yml
	@echo "Generating files..."
	mkdir generated || true
	go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --package generated -generate types,server,spec $< > generated/api.gen.go
