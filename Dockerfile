# Dockerfile for test image
FROM golang:alpine as build-env
MAINTAINER Sergey Gladkovskiy <smgladkovskiy@gmail.com>

ARG DEP_VERSION="0.5.0"
ENV IN_CONTAINER="true"

RUN apk update \
 && apk add --no-cache \
    ca-certificates \
    curl \
    git \
    make \
    openssl \
    openssh-client \
    gcc \
    musl-dev \
 && curl -L -s https://github.com/golang/dep/releases/download/v${DEP_VERSION}/dep-linux-amd64 -o /bin/dep \
 && chmod +x /bin/dep \
 && rm -rf /var/cache/apk/* \
 && rm -rf /tmp/*

ARG SRC="/go/src/github.com/microparts/configuration-golang"

COPY . ${SRC}
WORKDIR ${SRC}

RUN mkdir -p /cfgs/defaults \
 && mkdir -p /cfgs/development \
 && ln -s ${SRC}/test/configuration/defaults/* /cfgs/defaults/ \
 && ln -s ${SRC}/test/configuration/development/* /cfgs/development/

RUN make deps