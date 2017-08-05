#!/usr/bin/env bash
set -xe
mkdir -p $GOPATH/src/githbub.com/goccmack
git clone https://github.com/goccmack/gocc $GOPATH/src/github.com/goccmack/gocc
go get golang.org/x/tools/cmd/goimports
go get github.com/kisielk/errcheck
go get -u github.com/golang/lint/golint