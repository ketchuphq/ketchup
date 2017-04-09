dev:
	make -C admin
	go build .

osx: admin build-darwin

linux: admin build-linux

build-%:
	GOOS=$* GOARCH=amd64 go build .

admin:
	make -C admin production

include make/*.mk

goimports:
	GO_DIRS=$$(find . -name "*.go" -exec dirname {} \; | sort -u); \
	  $$GOPATH/bin/goimports -w -local github.com/octavore/ketchup $$GO_DIRS

.PHONY: admin

