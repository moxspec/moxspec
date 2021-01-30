#!/usr/bin/make -f

GO                  := go

DOCKER_WORKDIR      := /go/src/github.com/actapio/moxspec
DOCKER_RUN          := sudo docker run --rm -v $(CURDIR):$(DOCKER_WORKDIR) --workdir=$(DOCKER_WORKDIR)

CENTOS_CONTAINER    := takaswat/moxspec-centos:7

