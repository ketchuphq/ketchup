# shell is set so new path is picked up
SHELL := /bin/bash
export PATH := $(GOPATH):$(PATH)

dev:
	make -C admin
	go build .

osx: admin build-darwin
linux: admin build-linux
admin:
	make -C admin production
build-%:
	GOOS=$* GOARCH=amd64 go build \
	     -ldflags="-s -w" \
	     -gcflags=-trimpath=$$GOPATH \
	     -asmflags=-trimpath=$$GOPATH \
	     .

release: goreleaser.yml
	goreleaser
release-nr: goreleaser.yml
	goreleaser -nr
goreleaser.yml:
	@cat ./goreleaser.yml.tmpl \
		| sed -e "s,{GOPATH},$$GOPATH,g" \
		> goreleaser.yml

include make/*.mk

goimports:
	GO_DIRS=$$(find . -name "*.go" -exec dirname {} \; | sort -u); \
	  goimports -w -local github.com/octavore/ketchup $$GO_DIRS

.PHONY: admin goreleaser.yml

