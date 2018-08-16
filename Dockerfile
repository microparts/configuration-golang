# Dockerfile for test image
FROM golang:alpine as build-env
MAINTAINER Sergey Gladkovskiy <smgladkovskiy@gmail.com>

RUN apk update \
 && apk add --no-cache \
    git \
    make \
 && rm -rf /var/cache/apk/* \
 && rm -rf /tmp/*

COPY . /go/src/gitlab.teamc.io/teamc.io/microservice/configuration/golang-pkg.git
WORKDIR /go/src/gitlab.teamc.io/teamc.io/microservice/configuration/golang-pkg.git

RUN make deps