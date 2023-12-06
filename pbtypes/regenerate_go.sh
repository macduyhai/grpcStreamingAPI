#!/usr/bin/env bash

set -e

PROTOC=$(which protoc)

# Generate protobufs types all the dir. Exclude hidden directories such as .git
shopt -u dotglob

find service/* -prune -type d ! -name third_party | while read -r d; do
  ${PROTOC} \
  -I=. \
  -I=${GOPATH}/src \
  -I=${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis/ \
  -I=${GOPATH}/src/github.com/gogo/protobuf/protobuf \
  -I=${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/ \
  --gogo_out=plugins=grpc,\
Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types,\
Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types,\
Mgoogle/protobuf/empty.proto=github.com/gogo/protobuf/types,\
Mgoogle/api/annotations.proto=github.com/gogo/googleapis/google/api,\
Mgoogle/protobuf/field_mask.proto=github.com/gogo/protobuf/types,\
Mgoogle/protobuf/struct.proto=github.com/gogo/protobuf/types,\
Mgoogle/protobuf/any.proto=github.com/gogo/protobuf/types,\
Mgoogle/protobuf/wrappers.proto=github.com/gogo/protobuf/types:\
${GOPATH}/src/grpc-appota/grpcStreamingAPI/pbtypes \
${d}/*.proto \
    --swagger_out=third_party/OpenAPI/

done
