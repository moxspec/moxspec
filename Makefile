#!/usr/bin/make -f
.PHONY: help mox test clean init

include $(CURDIR)/Config.mk

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' Makefile | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-10s\033[0m %s\n", $$1, $$2}'

mox: ## build moxspec
ifdef CI
	$(MAKE) -f Build.mk mox
else
	$(DOCKER_RUN) $(CENTOS_CONTAINER) make -f Build.mk mox
	sudo chown -R $(shell id -u):$(shell id -g) bin
endif

test: ## run tests
ifdef CI
	$(MAKE) -f Build.mk test
else
	$(DOCKER_RUN) $(CENTOS_CONTAINER) make -f Build.mk test
	sudo chown -R $(shell id -u):$(shell id -g) c.out
endif

clean: ## clean all artifacts
	-rm -rf bin/
	-rm -rf c.out

init: ## install requirements 
	go get -u github.com/digitalocean/go-smbios/smbios
	go get -u golang.org/x/crypto/ssh/terminal
	go get -u github.com/kylelemons/godebug/{pretty,diff}
