
PROTOC_CMD = protoc
PROTOC_GRPC_CMD = protoc --plugin=$(shell go env GOPATH)/bin/protoc-gen-go-grpc

all: minialertAisprid

minialertAisprid: messages/serverProtocol_grpc.pb.go messages/serverProtocol.pb.go
	mkdir -p log
	go build
	# mkdir log

messages/serverProtocol_grpc.pb.go: messages/serverProtocol.proto
	$(PROTOC_GRPC_CMD) --go-grpc_out=. --go-grpc_opt=paths=source_relative  messages/serverProtocol.proto

messages/serverProtocol.pb.go: messages/serverProtocol.proto
	$(PROTOC_CMD) --go_out=. --go_opt=paths=source_relative messages/serverProtocol.proto

dependencies:
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

clean:
	rm -f minialertAisprid
	rm -f messages/*.go

