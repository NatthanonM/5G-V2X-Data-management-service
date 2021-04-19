# Go related commands
PROJECTNAME := $(shell basename "$(PWD)")
GOCMD = go
OS := $(shell uname -s | awk '{print tolower($$0)}')
GOARCH := amd64
GOBIN = bin

## bin: build go server/main to binary
.PHONY: bin
bin:
	env CGO_ENABLED=0 GOOS=$(OS) GOARCH=${GOARCH} go build -a -installsuffix cgo -o ${GOBIN}/server-$(OS)-${GOARCH} -mod=vendor cmd/server/main.go

.PHONY: start
start:
	go run cmd/server/main.go

.PHONY: proto
proto:
	protoc --proto_path=api/proto --proto_path=third_party --experimental_allow_proto3_optional  --go_out=plugins=grpc:pkg/api accident.proto
	protoc --proto_path=api/proto --proto_path=third_party --experimental_allow_proto3_optional  --go_out=plugins=grpc:pkg/api drowsiness.proto
	protoc --proto_path=api/proto --proto_path=third_party --experimental_allow_proto3_optional  --go_out=plugins=grpc:pkg/api car.proto