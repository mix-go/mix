#!/bin/bash
set -ex

BASEPATH=$(pwd)
GOPATH=$(go env GOPATH)
repos=$(find . -name "*.proto")
for FILE in $repos; do
  protoc -I . -I $GOPATH/src --go_out=paths=source_relative:$BASEPATH --go-grpc_out=paths=source_relative:$BASEPATH --grpc-gateway_out=paths=source_relative:$BASEPATH $FILE
done
