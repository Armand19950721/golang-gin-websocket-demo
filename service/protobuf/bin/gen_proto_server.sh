#!/bin/bash -x
set -e

SRC_DIR=${PWD}/../proto
OUT_DIR=${PWD}/../../src/protos
COMPILER=${PWD}/../../node_modules/.bin/protoc-gen-ts

mkdir -p ${OUT_DIR}

pushd ${SRC_DIR}

grpc_tools_node_protoc ${SRC_DIR}/*.proto \
--proto_path=${SRC_DIR} \
--plugin="protoc-gen-ts=${COMPILER}" \
--js_out=import_style=commonjs,binary:${OUT_DIR} \
--ts_out=grpc_js:${OUT_DIR} \
--grpc_out=grpc_js:${OUT_DIR}

echo "finish"
