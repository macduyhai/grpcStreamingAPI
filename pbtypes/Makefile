gen-proto:
	protoc --go_out=. --go_opt=paths=source_relative \
            --go-grpc_out=. --go-grpc_opt=paths=source_relative \
              service/ewallet/*.proto \
#	 protoc streamproto/streampb.proto --go_out=plugins=grpc:.
#	 protoc --java_out=. example.proto
#    protoc --python_out=. example.proto
#	 protoc --js_out=import_style=commonjs,binary:. example.proto
#	 protoc --ruby_out=. example.proto
