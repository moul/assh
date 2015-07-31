PACKAGES := $(shell go list ./... | grep -v vendor)

build:
	go get ./...
	go build -o assh ./cmd/assh

test:
	go test $(PACKAGES)

install:
	go install $(PACKAGES)
