BINARIES ?=	assh webapp
GODIR ?=	github.com/moul/advanced-ssh-config

PKG_BASE_DIR ?=	./pkg
CONVEY_PORT ?=	9042
SOURCES :=	$(shell find . -type f -name "*.go")
GOLIST :=	$(shell go list ./cmd/... ./pkg/...)
COMMANDS :=	$(shell go list ./cmd/...)
PACKAGES :=	$(shell go list ./pkg/...)
REL_COMMANDS := $(subst $(GODIR),./,$(COMMANDS))
REL_PACKAGES := $(subst $(GODIR),./,$(PACKAGES))
GOENV ?=	GO15VENDOREXPERIMENT=1
GO ?=		$(GOENV) go
USER ?=		$(shell whoami)


all:	build


.PHONY: build
build:	$(BINARIES)


.PHONY: docker
docker:
	docker build -t moul/assh .


$(BINARIES):	$(SOURCES)
	$(GO) build -ldflags=-s -i -v -o $@ ./cmd/$@


.PHONY: test
test:
	#$(GO) get -t ./...
	$(GO) test -i $(PACKAGES) $(COMMANDS)
	$(GO) test -v $(PACKAGES) $(COMMANDS)


.PHONY: examples
examples:
	@for example in $(dir $(wildcard examples/*/assh.yml)); do                \
	  set -xe;                                                                \
	  ./assh -c $$example/assh.yml config build > $$example/ssh_config;       \
	  ./assh -c $$example/assh.yml config graphviz > $$example/graphviz.dot;  \
	  dot -Tpng $$example/graphviz.dot > $$example/graphviz.png;              \
	done

.PHONY: install
install:
	$(GO) install $(COMMANDS)


.PHONY: clean
clean:
	rm -f $(BINARIES)


.PHONY: re
re:	clean all


.PHONY: convey
convey:
	$(GO) get github.com/smartystreets/goconvey
	goconvey -cover -port=$(CONVEY_PORT) -workDir="$(realpath $(PKG_BASE_DIR))" -depth=1


.PHONY:	cover
cover:	profile.out


profile.out:	$(SOURCES)
	rm -f $@
	find . -name profile.out -delete
	for package in $(REL_PACKAGES); do \
	  rm -f $$package/profile.out; \
	  $(GO) test -covermode=count -coverpkg=$(PKG_BASE_DIR)/... -coverprofile=$$package/profile.out $$package; \
	done
	echo "mode: count" > profile.out.tmp
	cat `find . -name profile.out` | grep -v mode: | sort -r | awk '{if($$1 != last) {print $$0;last=$$1}}' >> profile.out.tmp
	mv profile.out.tmp profile.out
