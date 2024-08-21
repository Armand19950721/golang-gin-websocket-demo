#!/bin/bash

SRC_DIR=${PWD}/../proto
OUT_DIR=${PWD}/../../assets/protos
COMPILER=${PWD}/../grpc_plugin/protoc-gen-grpc-web

mkdir -p ${OUT_DIR}

pushd ${SRC_DIR}

protoc ${SRC_DIR}/*.proto \
--proto_path=${SRC_DIR} \
--plugin="protoc-gen-grpc=${COMPILER}" \
--js_out=import_style=commonjs,binary:${OUT_DIR} \
--grpc-web_out=import_style=typescript,mode=grpcweb:${OUT_DIR}

echo "finish"
