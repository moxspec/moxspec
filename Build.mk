#!/usr/bin/make -f
.PHONY: mox test link vet goimports 

include $(CURDIR)/Version.mk
include $(CURDIR)/Config.mk

mox: test ## build just mox 
	mkdir -p bin
	CGO_ENABLED=0 $(GO) build \
		-a -tags netgo -installsuffix netgo \
		-ldflags "-s -w -X main.version=$(PKGVER) -X main.revision=`git rev-parse --short HEAD` -extldflags '-static'" \
		-o bin/$@ cmd/$@/*.go
	-cd bin && ln -s $@ lsdiag && ln -s $@ lsraid && ln -s $@ lssn

test: goimports lint vet ## run unit tests
	$(GO) test -coverprofile c.out ./...

lint: ## run golint
	golint -set_exit_status ./...

vet: ## run go vet
	$(GO) vet ./...

goimports: ## run goimports
	goimports -l ./ | xargs -r false

