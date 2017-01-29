all:
	make -C admin sources

include make/proto.mk

goimports:
	GO_DIRS=$$(find . -name "*.go" -exec dirname {} \; | sort -u); \
	  $$GOPATH/bin/goimports -w -local github.com/octavore/ketchup $$GO_DIRS

linux:
	GOOS=linux GOARCH=amd64 go build .

ship:
	scp press ketchup1:p2
