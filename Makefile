.DEFAULT_GOAL := test

.PHONY: fmt vet build clean test examples

fmt:
	go fmt ./...

vet: fmt
	go vet ./...

build: vet
	go build

clean: vet
	go clean

test: vet
	go test

examples:
	go test -tags=examples -v
