# all proto files in subfolders
PROTO_FILES=$(shell find . -name "*.proto")

# includes for all subfolders containing proto files, e.g. `-I./path/to/proto/folder`
INCL_PROTO_DIR=$(shell find . -name "*.proto" -exec dirname {} \; | sort -u | sed -e 's/^/-I/')

protos:
	@mkdir -p proto
	protoc $(INCL_PROTO_DIR) $(PROTO_FILES) --go_out=plugins=setter+grpc:./proto

prepare:
	brew install grpc/grpc/google-protobuf
	go get -u github.com/golang/protobuf/protoc-gen-go