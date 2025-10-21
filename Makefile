PROTOBUF_FILES = $(shell find ./proto/ -name '*.proto')
PROTOBUF_FILES_COMPILED = $(patsubst %.proto,%.pb.go,$(PROTOBUF_FILES))
PROTOBUF_FILES_COMPILED += $(patsubst %.proto,%_grpc.pb.go,$(PROTOBUF_FILES))

.PHONY: clean protos

clean:
	-rm $(PROTOBUF_FILES_COMPILED)

protos: $(PROTOBUF_FILES_COMPILED)

%.pb.go %_grpc.pb.go: %.proto
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative $^
