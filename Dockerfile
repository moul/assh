FROM golang:1.12-alpine as build
RUN             apk add --update --no-cache git gcc musl-dev make
COPY            contrib go.* /go/src/moul.io/assh/
WORKDIR         /go/src/moul.io/assh
RUN             GO111MODULE=on go get .
COPY            . /go/src/moul.io/assh
RUN             make install

FROM            alpine
COPY            --from=build /go/bin/assh /bin/assh
ENTRYPOINT      ["/bin/assh"]
