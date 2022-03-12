# dynamic config
ARG             BUILD_DATE
ARG             VCS_REF
ARG             VERSION

# build
FROM            golang:1.17.8-alpine as builder
RUN             apk add --no-cache git gcc musl-dev make
ENV             GO111MODULE=on
WORKDIR         /go/src/moul.io/assh
COPY            go.* ./
RUN             go mod download
COPY            . ./
RUN             make install

# minimalist runtime
FROM alpine:3.15.0
LABEL           org.label-schema.build-date=$BUILD_DATE \
                org.label-schema.name="assh" \
                org.label-schema.description="" \
                org.label-schema.url="https://moul.io/assh/" \
                org.label-schema.vcs-ref=$VCS_REF \
                org.label-schema.vcs-url="https://github.com/moul/assh" \
                org.label-schema.vendor="Manfred Touron" \
                org.label-schema.version=$VERSION \
                org.label-schema.schema-version="1.0" \
                org.label-schema.cmd="docker run -i -t --rm moul/assh" \
                org.label-schema.help="docker exec -it $CONTAINER assh --help"
COPY            --from=builder /go/bin/assh /bin/
ENTRYPOINT      ["/bin/assh"]
#CMD             []
