PACKAGES := $(shell go list ./... | grep -v vendor)

build:
	go get ./...
	go build -o assh ./cmd/assh

test:
	go test -v $(PACKAGES)

install:
	go install $(PACKAGES)

cover:
	rm -f profile.out file-profile.out
	for package in $(PACKAGES); do \
	  go test -coverprofile=file-profile.out $$package; \
	  if [ -f file-profile.out ]; then cat file-profile.out | grep -v "mode: set" >> profile.out || true; rm -f file-profile.out; fi \
	done
	echo "mode: set" | cat - profile.out > profile.out.tmp && mv profile.out.tmp profile.out
