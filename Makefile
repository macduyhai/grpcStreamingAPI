gen-proto:
	protoc streamproto/streampb.proto --go_out=plugins=grpc:.