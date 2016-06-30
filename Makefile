BINARIES ?=	assh
GODIR ?=	github.com/moul/advanced-ssh-config

PKG_BASE_DIR ?=	./pkg
CONVEY_PORT ?=	9042
SOURCES :=	$(shell find . -type f -name "*.go")
COMMANDS :=	$(shell go list ./... | grep -v /vendor/ | grep /cmd/)
PACKAGES :=	$(shell go list ./... | grep -v /vendor/ | grep -v /cmd/)
REL_COMMANDS := $(subst $(GODIR),./,$(COMMANDS))
REL_PACKAGES := $(subst $(GODIR),./,$(PACKAGES))
GOENV ?=	GO15VENDOREXPERIMENT=1
GO ?=		$(GOENV) go
USER ?=		$(shell whoami)



all:	build


.PHONY: build
build:	$(BINARIES)


$(BINARIES):	$(SOURCES)
	$(GO) build -o $@ ./cmd/$@


.PHONY: test
test:
	#$(GO) get -t ./...
	$(GO) test -i $(PACKAGES) $(COMMANDS)
	$(GO) test -v $(PACKAGES) $(COMMANDS)


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


.PHONY: docker-build
docker-build:
	go get github.com/laher/goxc
	rm -rf contrib/docker/linux_386
	for binary in $(BINARIES); do                                             \
	  goxc -bc="linux,386" -d . -pv contrib/docker -n $$binary xc;            \
	  mv contrib/docker/linux_386/$$binary contrib/docker/entrypoint;         \
	  docker build -t $(USER)/$$binary contrib/docker;                        \
	  docker run -it --rm $(USER)/$$binary || true;                           \
	  docker inspect --type=image --format="{{ .Id }}" moul/$$binary || true; \
	  echo "Now you can run 'docker push $(USER)/$$binary'";                  \
	done
