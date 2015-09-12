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
	find . -name profile.out -delete
	for package in $(PACKAGES); do \
	  rm -f $$package/profile.out; \
	  go test -covermode=count -coverpkg=./pkg/... -coverprofile=$$package/profile.out $$package; \
	done
	echo "mode: count" > profile.out.tmp
	cat `find . -name profile.out` | grep -v mode: | sort -r | awk '{if($$1 != last) {print $$0;last=$$1}}' >> profile.out.tmp
	mv profile.out.tmp profile.out


.PHONY: convey
convey:
	go get github.com/smartystreets/goconvey
	goconvey -cover -port=9042 -workDir="$(realpath .)/pkg" -depth=-1


.PHONY: docker/assh
docker/assh:
	goxc -bc=linux,386 -d=docker -o="{{.Dest}}{{.PS}}{{.ExeName}}{{.Ext}}" -include="" compile


.PHONY: docker
docker: docker/assh
	docker build -t moul/assh:latest docker
	docker run -it --rm moul/assh --version
