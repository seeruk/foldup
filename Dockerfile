# Build
FROM golang:1.8
MAINTAINER Elliot Wright <hello@elliotdwright.com>

WORKDIR /go/src/github.com/SeerUK/foldup

COPY . /go/src/github.com/SeerUK/foldup

RUN set -x \
    && CGO_ENABLED=0 GOOS=linux go build -a ./cmd/...

# Foldup packaging
FROM alpine:latest
MAINTAINER Elliot Wright <hello@elliotdwright.com>

WORKDIR /root/

COPY --from=0 /go/src/github.com/SeerUK/foldup/foldup .

RUN set -x \
    && apk add --update \
        ca-certificates \
    && rm -rf /var/cache/apk/*

ENTRYPOINT ["/root/foldup"]
