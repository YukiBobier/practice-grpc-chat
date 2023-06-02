#!/bin/bash

generate_code() {
    IN=$1
    OUT=$2
    PACKAGE=$3

    protoc \
    --go_out=${OUT} \
    --go_opt=paths=import \
    --go_opt=M${IN}=./${PACKAGE} \
    --go-grpc_out=${OUT} \
    --go-grpc_opt=paths=import \
    --go-grpc_opt=M${IN}=./${PACKAGE} \
    ${IN}
}

generate_code proto/chat.proto . internal/chat