#!/usr/bin/env bash

# gofast.....

set -e

# gogo 참고해서 수정해야함.
protoc --gofast_out=plugins=grpc:. -I=".:${GOGOPROTO_PATH}:${ETCD_ROOT_DIR}/..:${ETCD_ROOT_DIR}:${GRPC_GATEWAY_ROOT}/third_party/googleapis" \
      --plugin="${GOFAST_BIN}" ./*.proto


