#!/bin/bash
set -ex

cd proto
protoc -I . -I `go env GOPATH`/src --go_out=../openmix --go-grpc_out=../openmix --grpc-gateway_out=../openmix *.proto
