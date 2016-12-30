# all proto files in subfolders
PROTO_FILES=$(shell find . -name "*.proto")

# includes for all subfolders containing proto files, e.g. `-I./path/to/proto/folder`
INCL_PROTO_DIR=$(shell find . -name "*.proto" -exec dirname {} \; | sort -u | sed -e 's/^/-I/')

PROTO_PREFIX=import_prefix_proto=github.com/octavore/press/proto/

protos:
	@mkdir -p proto
	protoc $(INCL_PROTO_DIR) $(PROTO_FILES) --go_out=$(PROTO_PREFIX),plugins=setter:./proto
	go run util/gots/main.go ./admin/src/js/lib/api.ts

protos_list:
	@echo $(PROTO_FILES) | tr " " "\n"

prepare:
	brew install grpc/grpc/google-protobuf
	go get -u github.com/golang/protobuf/protoc-gen-go