include make/proto.mk

goimports:
	GO_DIRS=$$(find . -name "*.go" -exec dirname {} \; | sort -u); \
		$$GOPATH/bin/goimports -w -local github.com/octavore/press $$GO_DIRS
