#!/bin/bash
# vim: tabstop=4:softtabstop=2:shiftwidth=4:noexpandtab

DESTDIR=/
INSTALL_LOCATION=$(DESTDIR)/opt/harr
PROJECT=github.com/blendle/harr

.PHONY: all clean install deps cli worker legacy test watch

all: install

install: harr

harr: deps
	$(GOPATH)/bin/gom build -o $(INSTALL_LOCATION)/bin/harr src/harr.go

deps:
	GOPATH=$(GOPATH) go get github.com/mattn/gom
	GOPATH=$(GOPATH) $(GOPATH)/bin/gom install

clean:
	rm -rf $(GOPATH)/pkg/linux_amd64/$(PROJECT)
	rm -rf $(GOPATH)/pkg/darwin_amd64/$(PROJECT)

test: clean
	gom exec ginkgo -r --randomizeAllSpecs --randomizeSuites --succinct -cover

watch: test
	gom exec ginkgo watch -r
