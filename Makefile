build:
	go get ./...
	go build

test:
	go test $(go list ./... | grep -v vendor)
