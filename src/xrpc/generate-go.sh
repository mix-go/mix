#!/bin/bash
set -ex

cd proto
protoc -I . -I `go env GOPATH`/src --go_out=../protobuf --go-grpc_out=../protobuf --grpc-gateway_out=../protobuf *.proto
