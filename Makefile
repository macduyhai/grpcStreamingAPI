gen-proto:
	protoc streamproto/streampb.proto --go_out=plugins=grpc:.
	// protoc --java_out=. example.proto
        // protoc --python_out=. example.proto
	// protoc --js_out=import_style=commonjs,binary:. example.proto
	// protoc --ruby_out=. example.proto
