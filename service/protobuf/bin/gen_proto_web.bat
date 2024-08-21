@echo off

@setlocal
pushd "%~dp0"

set path=%cd%
set outdir=%path%\..\..\assets\proto
set srcdir=%path%\..\proto

if not exist %outdir% (
    mkdir %outdir%
)

%path%\protoc.exe %srcdir%\*.proto --proto_path="%srcdir%" --plugin=protoc-gen-grpc=%path%\protoc-gen-grpc-web.exe --js_out=import_style=commonjs,binary:"%outdir%" --grpc-web_out=import_style=typescript,mode=grpcweb:"%outdir%"

