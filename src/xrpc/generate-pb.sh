#!/bin/bash
set -ex

BASEPATH=$(pwd)
GOPATH=$(go env GOPATH)
OUT=$BASEPATH
repos=$(find . -name "*.proto")
for FILE in $repos; do
  protoc -I . -I `go env GOPATH`/src --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. --grpc-gateway_out=paths=source_relative:. $FILE
done
