.PHONY: start
start:
	go run cmd/server/main.go

.PHONY: proto
proto:
	protoc --proto_path=api/proto --proto_path=third_party --experimental_allow_proto3_optional  --go_out=plugins=grpc:pkg/api accident.proto
	protoc --proto_path=api/proto --proto_path=third_party --experimental_allow_proto3_optional  --go_out=plugins=grpc:pkg/api drowsiness.proto