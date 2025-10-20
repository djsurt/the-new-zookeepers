PROTOBUF_FILES = $(shell find ./proto/ -name '*.proto')
PROTOBUF_FILES_COMPILED = $(patsubst %.proto,%.go,$(PROTOBUF_FILES))

.PHONY: protos

protos: $(PROTOBUF_FILES_COMPILED)

%.go: %.proto
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative $^
