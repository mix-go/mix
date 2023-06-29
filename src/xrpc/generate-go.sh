#!/bin/bash
set -ex

cd proto
protoc -I . -I `go env GOPATH`/src --go_out=../testdata --go-grpc_out=../testdata --grpc-gateway_out=../testdata *.proto
