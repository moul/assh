PACKAGES := $(addprefix ./,$(wildcard pkg/*))
COMMANDS := $(addprefix ./,$(wildcard cmd/*))


all: build


build:
	go get ./...
	gofmt -w $(PACKAGES) $(COMMANDS)
	go test -i $(PACKAGES) $(COMMANDS)
	for command in $(COMMANDS); do \
	  go build -o `basename $$command` $$command; \
	done


test:
	go get ./...
	go test -i $(PACKAGES) $(COMMANDS)
	go test -v $(PACKAGES) $(COMMANDS)


install:
	go install $(COMMANDS)


cover:
	rm -f profile.out file-profile.out
	for package in $(PACKAGES); do \
	  go test -coverprofile=file-profile.out $$package; \
	  if [ -f file-profile.out ]; then cat file-profile.out | grep -v "mode: set" >> profile.out || true; rm -f file-profile.out; fi \
	done
	echo "mode: set" | cat - profile.out > profile.out.tmp && mv profile.out.tmp profile.out


.PHONY: convey
convey:
	go get github.com/smartystreets/goconvey
	goconvey -cover -port=9042 -workDir="$(realpath .)/pkg" -depth=-1
