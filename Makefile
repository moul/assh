rwildcard =	$(foreach d,$(wildcard $1*),$(call rwildcard,$d/,$2) $(filter $(subst *,%,$2),$d))
uniq =		$(if $1,$(firstword $1) $(call uniq,$(filter-out $(firstword $1),$1)))
SOURCES :=	$(call rwildcard,./cmd/ ./pkg/,*.go) go.*
PACKAGES :=	$(call uniq,$(dir $(call rwildcard,./pkg/,*.go)))
GOPATH ?=	$(HOME)/go
GO ?=		GO111MODULE=on go
TARGET ?=	$(GOPATH)/bin/assh


all:	install


.PHONY: install
install: $(TARGET)
$(TARGET): $(SOURCES)
	$(GO) install -v


.PHONY: docker.build
docker.build:
	docker build -t moul/assh .


.PHONY: test
test:
	$(GO) test -v ./...


.PHONY: examples
examples: $(TARGET)
	@for example in $(dir $(wildcard examples/*/assh.yml)); do                    \
	  set -xe;                                                                    \
	  $(TARGET) -c $$example/assh.yml config build > $$example/ssh_config;           \
	  $(TARGET) -c $$example/assh.yml config graphviz > $$example/graphviz.dot;      \
	  dot -Tpng $$example/graphviz.dot > $$example/graphviz.png;                  \
	  if [ -x $$example/test.sh ]; then (cd $$example; ./test.sh || exit 1); fi;  \
	done


.PHONY: clean
clean:
	rm -f $(TARGET) $(call rwildcard,./,*profile.out)
	rm -rf .release


.PHONY: re
re:	clean all


.PHONY:	cover
cover:	profile.out


profile.out: $(SOURCES)
	rm -f $@
	find . -name profile.out -delete
	for package in $(PACKAGES); do \
	  rm -f $$package/profile.out; \
	  $(GO) test -covermode=count -coverpkg=./... -coverprofile=$$package/profile.out $$package; \
	done
	echo "mode: count" > profile.out.tmp
	cat `find . -name profile.out` | grep -v mode: | sort -r | awk '{if($$1 != last) {print $$0;last=$$1}}' >> profile.out.tmp
	mv profile.out.tmp profile.out


.PHONY: lint
lint:
	golangci-lint run --verbose ./...


.PHONY: release
release:
	mkdir -p .release
	GOOS=linux   GOARCH=amd64 go build -i -v -o .release/assh_linux_amd64   .
	GOOS=linux   GOARCH=386   go build -i -v -o .release/assh_linux_386     .
	GOOS=linux   GOARCH=arm   go build -i -v -o .release/assh_linux_arm     .
	GOOS=openbsd GOARCH=amd64 go build -i -v -o .release/assh_openbsd_amd64 .
	GOOS=openbsd GOARCH=386   go build -i -v -o .release/assh_openbsd_386   .
	GOOS=openbsd GOARCH=arm   go build -i -v -o .release/assh_openbsd_arm   .
	GOOS=darwin  GOARCH=amd64 go build -i -v -o .release/assh_darwin_amd64  .
	GOOS=darwin  GOARCH=386   go build -i -v -o .release/assh_darwin_386    .
	#GOOS=darwin  GOARCH=arm   go build -i -v -o .release/assh_darwin_arm    .
	GOOS=netbsd  GOARCH=amd64 go build -i -v -o .release/assh_netbsd_amd64  .
	GOOS=netbsd  GOARCH=386   go build -i -v -o .release/assh_netbsd_386    .
	GOOS=netbsd  GOARCH=arm   go build -i -v -o .release/assh_netbsd_arm    .
	GOOS=freebsd GOARCH=amd64 go build -i -v -o .release/assh_freebsd_amd64 .
	GOOS=freebsd GOARCH=386   go build -i -v -o .release/assh_freebsd_386   .
	GOOS=freebsd GOARCH=arm   go build -i -v -o .release/assh_freebsd_arm   .
	GOOS=windows GOARCH=amd64 go build -i -v -o .release/assh_windows_amd64.exe .
	GOOS=windows GOARCH=386   go build -i -v -o .release/assh_windows_386.exe   .
	#GOOS=windows GOARCH=arm   go build -i -v -o .release/assh_windows_arm.exe   .
