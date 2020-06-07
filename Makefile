GOPKG ?=	moul.io/assh
DOCKER_IMAGE ?=	moul/assh
GOBINS ?=	.

PRE_INSTALL_STEPS += generate
PRE_UNITTEST_STEPS += generate
PRE_TEST_STEPS += generate
PRE_BUILD_STEPS += generate
PRE_LINT_STEPS += generate
PRE_TIDY_STEPS += generate
PRE_BUMPDEPS_STEPS += generate

VERSION ?= `git describe --tags --always`
VCS_REF ?= `git rev-parse --short HEAD`

GO_INSTALL_OPTS = -ldflags="-X 'moul.io/assh/v2/pkg/version.Version=$(VERSION)' -X 'moul.io/assh/v2/pkg/version.VcsRef=$(VCS_REF)' "

include rules.mk

.PHONY: generate
generate:
	go generate

.PHONY: examples
examples: $(TARGET)
	@for example in $(dir $(wildcard examples/*/assh.yml)); do                    \
	  set -xe;                                                                    \
	  $(TARGET) -c $$example/assh.yml config build > $$example/ssh_config;           \
	  $(TARGET) -c $$example/assh.yml config graphviz > $$example/graphviz.dot;      \
	  dot -Tpng $$example/graphviz.dot > $$example/graphviz.png;                  \
	  if [ -x $$example/test.sh ]; then (cd $$example; ./test.sh || exit 1); fi;  \
	done

.PHONY: gen-release
gen-release: generate
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
