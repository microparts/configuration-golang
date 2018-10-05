# Dockerfile for test image
FROM golang:alpine as build-env
MAINTAINER Sergey Gladkovskiy <smgladkovskiy@gmail.com>

ARG DEP_VERSION="0.4.1"

RUN apk update \
 && apk add --no-cache \
    ca-certificates \
    curl \
    git \
    make \
    openssl \
    openssh-client \
 && curl -L -s https://github.com/golang/dep/releases/download/v${DEP_VERSION}/dep-linux-amd64 -o /bin/dep \
 && chmod +x /bin/dep \
 && rm -rf /var/cache/apk/* \
 && rm -rf /tmp/*

COPY . /go/src/gitlab.teamc.io/teamc.io/microservice/configuration/golang-pkg.git
WORKDIR /go/src/gitlab.teamc.io/teamc.io/microservice/configuration/golang-pkg.git

RUN mkdir -p /cfgs/defaults \
 && mkdir -p /cfgs/test \
 && ln -s /go/src/gitlab.teamc.io/teamc.io/microservice/configuration/golang-pkg.git/test/configuration/defaults/* /cfgs/defaults/ \
 && ln -s /go/src/gitlab.teamc.io/teamc.io/microservice/configuration/golang-pkg.git/test/configuration/test/* /cfgs/test/

RUN make deps