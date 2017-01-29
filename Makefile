include make/proto.mk

all:
	make -C admin sources

goimports:
	GO_DIRS=$$(find . -name "*.go" -exec dirname {} \; | sort -u); \
		$$GOPATH/bin/goimports -w -local github.com/octavore/ketchup $$GO_DIRS
