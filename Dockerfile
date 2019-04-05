# Dockerfile for test image
FROM golang:1.12.1-alpine as build-env
MAINTAINER Sergey Gladkovskiy <smgladkovskiy@gmail.com>

ENV GO111MODULE=on

RUN apk update \
 && apk add --no-cache \
    curl \
    git \
    make \
    openssh-client \
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