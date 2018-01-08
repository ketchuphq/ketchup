# shell is set so new path is picked up
SHELL := /bin/bash
export PATH := $(GOPATH)/bin:$(PATH)

GO_DIRS := $(shell find . -name "*.go" -exec dirname {} \; | grep -v vendor | grep -v proto | sort -u)

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
	goreleaser --skip-publish
goreleaser.yml:
	@cat ./goreleaser.yml.tmpl \
		| sed -e "s,{GOPATH},$$GOPATH,g" \
		> goreleaser.yml

include make/*.mk

prepare-admin:
	make -C admin prepare

prepare-vendor:
	go get -u github.com/kardianos/govendor
	govendor sync

goimports:
	@goimports -w -local github.com/ketchuphq/ketchup $(GO_DIRS)

test:
	@go test $(GO_DIRS)

cover:
	@go get -u github.com/go-playground/overalls
	overalls -project github.com/ketchuphq/ketchup
	@find . -name 'profile.coverprofile' | xargs rm

circle-cover: cover
	@go get -u github.com/mattn/goveralls
	@goveralls \
		-coverprofile=overalls.coverprofile \
		-service=circle-ci \
		-repotoken=$$COVERALLS_TOKEN


.PHONY: admin goreleaser.yml

