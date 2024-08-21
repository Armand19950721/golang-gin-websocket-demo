#!/bin/bash

SRC_DIR=${PWD}/../proto
OUT_DIR=${PWD}/../../..

mkdir -p ${OUT_DIR}

pushd ${SRC_DIR}

protoc ${SRC_DIR}/*.proto \
--proto_path=${SRC_DIR} \
--go_out=${OUT_DIR} \
--go-grpc_out=${OUT_DIR}