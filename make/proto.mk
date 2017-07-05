export PATH := $(GOPATH)/bin:$(PATH)

# all proto files in subfolders
PROTO_FILES=$(shell find . -name "*.proto" | grep -v node_modules)

# includes for all subfolders containing proto files, e.g. `-I./path/to/proto/folder`
INCL_PROTO_DIR=$(shell find . -name "*.proto" -exec dirname {} \; | grep -v node_modules | sort -u | sed -e 's/^/-I/')
INCL_WKT=-I $$GOPATH/src/github.com/golang/protobuf/ptypes/struct

PROTO_PREFIX=import_prefix_proto=github.com/ketchuphq/ketchup/proto/

protos:
	@mkdir -p proto
	protoc $(INCL_PROTO_DIR) $(INCL_WKT) $(PROTO_FILES) $$GOPATH/src/github.com/golang/protobuf/ptypes/struct/struct.proto --go-2_out=$(PROTO_PREFIX),plugins=setter:./proto
	go run util/gots/main.go ./admin/src/js/lib/api.ts

protos-list:
	@echo $(PROTO_FILES) | tr " " "\n"

prepare-protos:
	brew install grpc/grpc/google-protobuf
	go get -u github.com/octavore/protobuf/protoc-gen-go-2
