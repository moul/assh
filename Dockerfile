FROM golang:1.8.1
COPY . /go/src/github.com/moul/advanced-ssh-config
WORKDIR /go/src/github.com/moul/advanced-ssh-config
RUN make
ENTRYPOINT ["/go/src/github.com/moul/advanced-ssh-config/assh"]
